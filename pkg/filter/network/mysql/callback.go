package mysql

import "mosn.io/mosn/pkg/protocol/mysql"

type Callback struct {
	//state *State
}

func (c Callback) OnProtocolError() {
	//TODO implement me
}

func (c Callback) OnNewMessage(state mysql.State) {
	if state == mysql.ChallengeReq {
		c.state.login_attempts.Inc(1)
	}
}

func (c Callback) OnClientLogin(clientLogin *mysql.ClientLoginResponse) {
}

func (c Callback) OnClientLoginResponse() {
}

func (c Callback) OnClientSwitchResponse(cc *mysql.Command) {
	//TODO implement me
}

func (c Callback) OnMoreClientLoginResponse(cr *mysql.CommandResponse) {
	//TODO implement me
}

func (c Callback) OnCommand(*mysql.Command) {
	//TODO implement me
}

func (c Callback) OnCommandResponse(*mysql.CommandResponse) {
	//TODO implement me
}

func (c Callback) OnServerGreeting() {
	c.state.UpgradedToSsl.Inc(1)
}
