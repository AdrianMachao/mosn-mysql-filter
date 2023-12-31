/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tunnel

import (
	"encoding/json"
	"net"
	"sync"
	"time"

	"mosn.io/api"
	v2 "mosn.io/mosn/pkg/config/v2"
	"mosn.io/mosn/pkg/filter/network/tunnel/ext"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/server"
	"mosn.io/mosn/pkg/stagemanager"
	"mosn.io/mosn/pkg/types"
	"mosn.io/pkg/utils"
)

type AgentBootstrapConfig struct {
	Enable bool `json:"enable"`
	// The number of connections established between the agent and each server
	ConnectionNum int `json:"connection_num"`
	// The cluster of remote server
	Cluster string `json:"cluster"`
	// After the connection is established, the data transmission is processed by this listener
	HostingListener string `json:"hosting_listener"`
	// Static remote server list
	StaticServerList []string `json:"server_list"`

	// DynamicServerListConfig is used to specify dynamic server configuration
	DynamicServerListConfig struct {
		DynamicServerLister string `json:"dynamic_server_lister"`
	} `json:"dynamic_server_list_config"`

	// ConnectRetryTimes
	ConnectRetryTimes int `json:"connect_retry_times"`
	// ReconnectBaseDuration
	ReconnectBaseDuration api.DurationConfig `json:"reconnect_base_duration"`

	// ConnectTimeoutDuration specifies the timeout for establishing a connection and initializing the agent
	ConnectTimeoutDuration api.DurationConfig `json:"connect_timeout_duration"`
	CredentialPolicy       string             `json:"credential_policy"`
	// GracefulCloseMaxWaitDuration specifies the maximum waiting time to close conn gracefully
	GracefulCloseMaxWaitDuration api.DurationConfig `json:"graceful_close_max_wait_duration"`

	TLSContext *v2.TLSConfig `json:"tls_context"`
}

func init() {
	v2.RegisterParseExtendConfig("tunnel_agent", func(config json.RawMessage) error {
		var conf AgentBootstrapConfig
		err := json.Unmarshal(config, &conf)
		if err != nil {
			log.DefaultLogger.Errorf("[tunnel agent] failed to parse agent bootstrap config: %v", err.Error())
			return err
		}
		if conf.Enable {
			utils.GoWithRecover(func() {
				bootstrap(&conf)
			}, nil)

			stagemanager.OnGracefulStop(func() error {
				stopAllPeers()
				return nil
			})
		}
		return nil
	})
}

var stopChan = make(chan struct{})

func stopAllPeers() {
	close(stopChan)
	peerMap.Range(func(key, value interface{}) bool {
		log.DefaultLogger.Infof("[agent] try to delete all peer, key: %v", key)
		peerMap.Delete(key)
		value.(*AgentPeer).Stop()
		return true
	})
}

func bootstrap(conf *AgentBootstrapConfig) {

	if conf.DynamicServerListConfig.DynamicServerLister != "" {
		utils.GoWithRecover(func() {
			lister := ext.GetServerLister(conf.DynamicServerListConfig.DynamicServerLister)
			ch := lister.List(conf.Cluster)
			for {
				select {
				case <-stopChan:
					log.DefaultLogger.Infof("[tunnel agent] agent has been stopped, lister exit")
					return
				case servers := <-ch:
					// Compute the diff between new and old server list
					intersection := make(map[string]bool)
					for i := range servers {
						if _, ok := peerMap.Load(servers[i]); ok {
							intersection[servers[i]] = true
						}
					}
					increased := make([]string, 0)
					for _, addr := range servers {
						if _, ok := intersection[addr]; !ok {
							increased = append(increased, addr)
							func(address string) {
								utils.GoWithRecover(func() {
									connectServer(conf, address)
								}, nil)
							}(addr)
						}
					}
					decreased := make([]string, 0)
					peerMap.Range(func(key, value interface{}) bool {
						addr := key.(string)
						_, ok := intersection[addr]
						if !ok {
							decreased = append(decreased, addr)
						}
						return true
					})
					for _, addr := range decreased {
						val, ok := peerMap.Load(addr)
						if !ok {
							continue
						}
						func(a *AgentPeer) {
							utils.GoWithRecover(func() {
								a.Stop()
							}, nil)
						}(val.(*AgentPeer))

						peerMap.Delete(addr)
					}
					log.DefaultLogger.Infof("[tunnel agent] tunnel server list changed, update success, increased: %+v, decreased: %+v", increased, decreased)

				}
			}
		}, nil)
	}

	for _, serverAddress := range conf.StaticServerList {
		host, port, err := net.SplitHostPort(serverAddress)
		if err != nil {
			log.DefaultLogger.Fatalf("server address invalid format, address: %v", serverAddress)
		}
		addrs, err := net.LookupHost(host)
		if err != nil {
			log.DefaultLogger.Fatalf("[tunnel agent] failed to lookup host by domain: %v", host)
		}
		for _, addr := range addrs {
			func(address string) {
				utils.GoWithRecover(func() {
					connectServer(conf, net.JoinHostPort(address, port))
				}, nil)
			}(addr)
		}
	}
}

var peerMap = &sync.Map{}

var (
	defaultReconnectBaseDuration     = time.Second * 3
	defaultConnectTimeoutDuration    = time.Second * 15
	defaultGracefulCloseWaitDuration = time.Second * 15
	defaultConnectMaxRetryTimes      = -1
)

func connectServer(conf *AgentBootstrapConfig, address string) {
	listener := server.GetServer().Handler().FindListenerByName(conf.HostingListener)
	if listener == nil {
		return
	}
	config := &ConnectionConfig{
		Address:                      address,
		ClusterName:                  conf.Cluster,
		Weight:                       10,
		ReconnectBaseDuration:        conf.ReconnectBaseDuration.Duration,
		ConnectTimeoutDuration:       conf.ConnectTimeoutDuration.Duration,
		ConnectRetryTimes:            conf.ConnectRetryTimes,
		CredentialPolicy:             conf.CredentialPolicy,
		ConnectionNumPerAddress:      conf.ConnectionNum,
		GracefulCloseMaxWaitDuration: conf.GracefulCloseMaxWaitDuration.Duration,
		TLSContext:                   conf.TLSContext,
	}
	if config.Network == "" {
		config.Network = "tcp"
	}
	if config.ReconnectBaseDuration == 0 {
		config.ReconnectBaseDuration = defaultReconnectBaseDuration
	}
	if config.ConnectRetryTimes == 0 {
		config.ConnectRetryTimes = defaultConnectMaxRetryTimes
	}
	if config.ConnectTimeoutDuration == 0 {
		config.ConnectTimeoutDuration = defaultConnectTimeoutDuration
	}
	if config.GracefulCloseMaxWaitDuration == 0 {
		config.GracefulCloseMaxWaitDuration = defaultGracefulCloseWaitDuration
	}
	peer := &AgentPeer{
		conf:     config,
		listener: listener,
	}
	peer.Start()
	peerMap.Store(address, peer)
}

type ConnectionConfig struct {
	Address           string `json:"address"`
	ClusterName       string `json:"cluster_name"`
	Weight            int64  `json:"weight"`
	ConnectRetryTimes int    `json:"connect_retry_times"`
	// ConnectTimeoutDuration specifies the timeout for establishing a connection and initializing the agent
	ConnectTimeoutDuration       time.Duration `json:"connect_timeout_duration"`
	Network                      string        `json:"network"`
	ReconnectBaseDuration        time.Duration `json:"reconnect_base_duration"`
	CredentialPolicy             string        `json:"credential_policy"`
	ConnectionNumPerAddress      int           `json:"connection_num_per_address"`
	GracefulCloseMaxWaitDuration time.Duration `json:"graceful_close_max_wait_duration"`

	TLSContext *v2.TLSConfig `json:"tls_context"`
}

type AgentPeer struct {
	conf        *ConnectionConfig
	connections []*AgentClientConnection
	// asideConn only used to send some control data to server
	asideConn *AgentAsideConnection
	listener  types.Listener
}

func (a *AgentPeer) Start() {
	connList := make([]*AgentClientConnection, 0, a.conf.ConnectionNumPerAddress)
	for i := 0; i < a.conf.ConnectionNumPerAddress; i++ {
		conn := NewAgentCoreConnection(*a.conf, a.listener)
		err := conn.initConnection()
		if err == nil {
			connList = append(connList, conn)
		}
	}
	a.connections = connList
	a.initAside()
}

func (a *AgentPeer) initAside() {
	asideConn := NewAgentAsideConnection(*a.conf, a.listener)
	err := asideConn.initConnection()
	if err != nil {
		return
	}
	a.asideConn = asideConn
}

func (a *AgentPeer) Stop() {
	closeWait := time.NewTimer(a.conf.GracefulCloseMaxWaitDuration)
	utils.GoWithRecover(func() {
		if a.asideConn == nil {
			a.initAside()
		}
		// Init aside connection still fails
		if a.asideConn == nil {
			return
		}
		addresses := make([]string, 0, len(a.connections))
		for _, conn := range a.connections {
			localAddr := conn.rawc.LocalAddr().String()
			conn.PrepareClose()
			addresses = append(addresses, localAddr)
		}
		// Send oneway request to notify server server to offline conn gracefully
		err := a.asideConn.Write(&GracefulCloseOnewayRequest{
			Addresses:   addresses,
			ClusterName: a.conf.ClusterName,
		})
		if err != nil {
			log.DefaultLogger.Errorf("[tunnel agent] write graceful close request error, err: %+v", err)
		}
		_ = a.asideConn.Close()
	}, nil)

	select {
	case <-closeWait.C:
		log.DefaultLogger.Warnf("[tunnel agent] waiting for graceful closing timeout, try to close directly")
		for _, conn := range a.connections {
			err := conn.Close()
			if err != nil {
				log.DefaultLogger.Errorf("[tunnel agent] failed to stop connection, err: %+v", err)
			}
		}
		_ = a.asideConn.Close()
		closeWait.Stop()
	}
}
