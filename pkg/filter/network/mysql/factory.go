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

package mysql

import (
	"context"
	"encoding/json"
	"fmt"

	"mosn.io/api"
	v2 "mosn.io/mosn/pkg/config/v2"
)

func init() {
	api.RegisterNetwork(v2.Mysql, CreateMysqlProxyFactory)
}

type mysqlFilterConfigFactory struct {
	Proxy *v2.StreamProxy
}

func (f *mysqlFilterConfigFactory) CreateFilterChain(context context.Context, callbacks api.NetWorkFilterChainFactoryCallbacks) {
	rf := NewProxy(context, f.Proxy)
	callbacks.AddReadFilter(rf)
}

func CreateMysqlProxyFactory(conf map[string]interface{}) (api.NetworkFilterChainFactory, error) {
	p, err := ParseStreamProxy(conf)
	if err != nil {
		return nil, err
	}
	return &mysqlFilterConfigFactory{
		Proxy: p,
	}, nil
}

// ParseStreamProxy
func ParseStreamProxy(cfg map[string]interface{}) (*v2.StreamProxy, error) {
	proxy := &v2.StreamProxy{}
	if data, err := json.Marshal(cfg); err == nil {
		json.Unmarshal(data, proxy)
	} else {
		return nil, fmt.Errorf("[config] config is not a tcp proxy config: %v", err)
	}
	return proxy, nil
}
