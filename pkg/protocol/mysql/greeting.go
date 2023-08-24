package mysql

type ServerGreeting struct {
	MySQLCodec
	protocol        uint8
	version         string
	threadId        uint32
	authPluginData1 []uint8
	authPluginData2 []uint8
	serverCap       uint32
	serverCharset   uint8
	serverStatus    uint16
	authPluginName  string
}

func (sg *ServerGreeting) getProtocol() uint8 {
	return 0
}

func (sg *ServerGreeting) getVersion() string {
	return ""
}

func (sg *ServerGreeting) getThreadId() uint32 {
	return 0
}

func (sg *ServerGreeting) getAuthPluginData1() []uint8 {
	return nil
}

func (sg *ServerGreeting) getAuthPluginData2() []uint8 {
	return nil
}

func (sg *ServerGreeting) getAuthPluginData() []uint8 {
	return nil
}

func (sg *ServerGreeting) setProtocol(protocol uint8) {
}

func (sg *ServerGreeting) setVersion(version string) {
}

func (sg *ServerGreeting) setThreadId(threadId uint32) {
}

func (sg *ServerGreeting) setServerCap(serverCap uint32) {
}

func (sg *ServerGreeting) setBaseServerCap(baseServerCap uint16) {
}

func (sg *ServerGreeting) setExtServerCap(extServerCap uint16) {
}

func (sg *ServerGreeting) setAuthPluginName(name string) {
}

func (sg *ServerGreeting) setAuthPluginData(salt []uint8) {
}

func (sg *ServerGreeting) setAuthPluginData1(salt []uint8) {
}

func (sg *ServerGreeting) setAuthPluginData2(salt []uint8) {
}

func (sg *ServerGreeting) setServerCharset(serverLanguage uint8) {
}

func (sg *ServerGreeting) setServerStatus(serverStatus uint16) {
}

func (sg *ServerGreeting) check() DecodeStatus {
	return 0
}
