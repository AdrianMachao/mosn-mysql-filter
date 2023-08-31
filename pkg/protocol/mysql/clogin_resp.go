package mysql

type ClientLoginResponse struct {
	MySQLCodec
	respCode uint8
}

type AuthMoreMessage struct {
	*ClientLoginResponse
	morePluginData []uint8
}

type AuthSwitchMessage struct {
	*ClientLoginResponse
	isOldAuthSwitch bool
	authPluginData  []uint8
	authPluginName  string
}

type OkMessage struct {
	*ClientLoginResponse
	affectedRows        uint64
	lastInsertId        uint64
	status              uint16
	warnings            uint16
	info                string
	sessionStateChanges string
}

type ErrMessage struct {
	*ClientLoginResponse
	marker       uint8
	errorCode    uint16
	sqlState     string
	errorMessage string
}

func (clr *ClientLoginResponse) setRespCode(code uint8) {
}

func (clr *ClientLoginResponse) GetRespCode() uint8 {
	return 0
}

func (amm *AuthMoreMessage) setAuthMoreData(data []uint8) {

}

func (amm *AuthMoreMessage) getAuthMoreData() []uint8 {
	return nil
}

func (asm *AuthSwitchMessage) isIldAuthSwitch() {
}

func (asm *AuthSwitchMessage) getAuthPluginData() []uint8 {
	return nil
}

func (asm *AuthSwitchMessage) getAuthPluginName() string {
	return ""
}

func (asm *AuthSwitchMessage) setIsOldAuthSwitch(old bool) {
}

func (asm *AuthSwitchMessage) setAuthPluginData(data []uint8) {
}

func (asm *AuthSwitchMessage) setAuthPluginName(name string) {
}

func (em *ErrMessage) setErrorCode(errCode uint16) {
}

func (em *ErrMessage) setSqlStateMarker(marker uint8) {
}

func (em *ErrMessage) setSqlState(state string) {
}

func (em *ErrMessage) setErrorMessage(msg string) {
}

func (em *ErrMessage) getErrorCode() uint16 {
	return 0
}

func (em *ErrMessage) getSqlStateMarker() uint8 {
	return 0
}

func (em *ErrMessage) getSqlState() string {
	return ""
}

func (em *ErrMessage) getErrorMessage() string {
	return ""
}
