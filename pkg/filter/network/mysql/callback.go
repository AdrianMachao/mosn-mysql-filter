package mysql

import "mosn.io/mosn/pkg/protocol/mysql"

type Callback struct {
	state *State
}

func (c Callback) onProtocolError() {
	c.state.ProtocolErrors.Inc(1)
}

func (c Callback) OnNewMessage() {
	c.state.login_attempts.Inc(1)
}

func (c Callback) OnServerGreeting() {
	c.state.UpgradedToSsl.Inc(1)
}

func (c Callback) OnClientLogin(clientLogin *mysql.ClientLoginResponse) {
}

func (c Callback) OnClientLoginResponse() {
}

func (c Callback) OnClientSwitchResponse() {
}

func (c Callback) OnMoreClientLoginResponse() {
}

func (c Callback) OnCommand() {
}

func (c Callback) OnCommandResponse() {
}
