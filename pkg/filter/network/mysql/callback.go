package mysql

type Callback struct {
}

func (c Callback) OnNewMessage() {
}

func (c Callback) OnServerGreeting() {
}

func (c Callback) OnClientLogin() {
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
