package mysql

import "mosn.io/mosn/pkg/protocol/mysql"

type Callback struct {
	//state *State
}

func (c Callback) OnProtocolError() {
	//TODO implement me
}

func (c Callback) OnNewMessage(state mysql.State) {
	//TODO implement me
}

func (c Callback) OnServerGreeting(sg *mysql.ServerGreeting) {
	//TODO implement me
}

func (c Callback) OnClientLogin(cl *mysql.ClientLogin) {
	//TODO implement me
}

func (c Callback) OnClientLoginResponse(clr *mysql.ClientLoginResponse) {
	//TODO implement me
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
