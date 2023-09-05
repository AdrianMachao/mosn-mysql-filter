package mysql

import (
	"mosn.io/mosn/pkg/types"
	"unsafe"
)

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

func (cl *ClientLogin) parseMessage(buffer types.IoBuffer, length uint32) DecodeStatus {
	var baseCap uint16
	var status DecodeStatus
	if baseCap, status = readUint16(buffer); status != Success {
		return Failure
	}

	cl.setBaseClientCap(baseCap)
	if uint32(baseCap)&CLIENT_SSL == 1 {
		return cl.parseMessageSsl(buffer)
	}

	if uint32(baseCap)&CLIENT_PROTOCOL_41 == 1 {
		return cl.parseMessage41(buffer)
	}

	return cl.parseMessage320(buffer, length-uint32(unsafe.Sizeof(baseCap)))
	// TODO resp code
	//respCode, status := readUint8(buffer)
	//if status != Success {
	//	return Failure
	//}
	//
	//if respCode != MYSQL_RESP_OK {
	//	return Failure
	//}
	//var val uint64
	//val, status = readLengthEncodedInteger(buffer)
	//if status == Failure {
	//	return Failure
	//}
	//
	//val, status = readLengthEncodedInteger(buffer)
	//if status == Failure {
	//	return Failure
	//}

	return 0
}

func (m *ClientLogin) decode(data types.IoBuffer, seq uint8, length uint32) DecodeStatus {
	m.seq = seq
	return m.parseMessage(data, length)
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

func (cl *ClientLogin) addConnectionAttribute() {
}

func (cl *ClientLogin) isResponse41() bool {
	return true
}

func (cl *ClientLogin) isResponse320() bool {
	return true
}
func (cl *ClientLogin) IsSSLRequest() bool {
	return true
}
func (cl *ClientLogin) isConnectWithDb() bool {
	return true
}
func (cl *ClientLogin) isClientAuthLenClData() bool {
	return true
}
func (cl *ClientLogin) isClientSecureConnection() bool {
	return true
}
