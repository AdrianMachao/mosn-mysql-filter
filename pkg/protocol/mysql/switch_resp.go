package mysql

type ClientSwitchResponse struct {
	MySQLCodec
	authPluginResp []uint8
}

func (csp *ClientSwitchResponse) getAuthPluginResp() (authPluginResp []uint8) {
	return nil
}

func (csp *ClientSwitchResponse) setAuthPluginResp(authPluginResp []uint8) {
}
