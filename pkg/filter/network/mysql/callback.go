package mysql

import (
	"mosn.io/mosn/pkg/protocol/mysql"
)

type Callback struct {
	state *State
}

func (c Callback) OnProtocolError() {
	c.state.ProtocolErrors.Inc(1)
}

func (c Callback) OnNewMessage(state mysql.State) {
	if state == mysql.ChallengeReq {
		c.state.LoginAttempts.Inc(1)
	}
}

func (c Callback) OnClientLogin(cl *mysql.ClientLogin) {
	if cl.IsSSLRequest() {
		c.state.UpgradedToSsl.Inc(1)
	}
}

func (c Callback) OnClientLoginResponse(clr *mysql.ClientLoginResponse) {
	if clr.GetRespCode() == mysql.MYSQL_RESP_AUTH_SWITCH {
		c.state.AuthSwitchRequest.Inc(1)
	} else if clr.GetRespCode() == mysql.MYSQL_RESP_ERR {
		c.state.LoginFailures.Inc(1)
	}
}

func (c Callback) OnServerGreeting(sg *mysql.ServerGreeting) {
	c.state.UpgradedToSsl.Inc(1)
}

func (c Callback) OnClientSwitchResponse(cc *mysql.Command) {
}

func (c Callback) OnMoreClientLoginResponse(clr *mysql.ClientLoginResponse) {
	if clr.GetRespCode() == mysql.MYSQL_RESP_ERR {
		c.state.LoginFailures.Inc(1)
	}
}

func (c Callback) OnCommand(*mysql.Command) {
	// TODO implement me
}

func (c Callback) OnCommandResponse(*mysql.CommandResponse) {
}
