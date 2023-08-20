package mysql

import "mosn.io/mosn/pkg/types"

type ClientLogin struct {
	MySQLCodec
	clientCap      uint32
	maxPacket      uint32
	charset        uint8
	userName       string
	authResp       []uint8
	db             string
	authPluginName string
	connAttr       [][]string
}

func (cl *ClientLogin) setClientCap(clientCap uint32) {
}

func (cl *ClientLogin) setBaseClientCap(baseCap uint16) {
}

func (cl *ClientLogin) setExtendedClientCap(extendedClientCap uint16) {
}

func (cl *ClientLogin) setMaxPacket(maxPacket uint32) { cl.maxPacket = maxPacket }

func (cl *ClientLogin) setCharset(charset uint8) { cl.charset = charset }

func (cl *ClientLogin) setUsername(userName string) {}

func (cl *ClientLogin) setDb(db string) {}

func (cl *ClientLogin) setAuthResp(authResp []uint8) {}

func (cl *ClientLogin) setAuthPluginName(plugin string) {
}

func (cl *ClientLogin) parseMessage(buffer types.IoBuffer, lenth uint32) DecodeStatus {
	return 0
}

func (cl *ClientLogin) parseMessageSsl(buffer types.IoBuffer) DecodeStatus {
	return 0
}

func (cl *ClientLogin) parseMessage41(buffer types.IoBuffer) DecodeStatus {
	return 0
}

func (cl *ClientLogin) parseMessage320(buffer types.IoBuffer, remainLen uint32) DecodeStatus {
	return 0
}

func (cl *ClientLogin) encode(out types.IoBuffer) {
}

func (cl *ClientLogin) encodeResponseSsl(out types.IoBuffer) {
}

func (cl *ClientLogin) encodeResponse41(out types.IoBuffer) {
}

func (cl *ClientLogin) encodeResponse320(out types.IoBuffer) {
}
