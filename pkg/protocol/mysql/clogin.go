package mysql

import (
	"mosn.io/mosn/pkg/types"
	"mosn.io/pkg/log"
	"unsafe"
)

type ClientLogin struct {
	MySQLCodec
	clientCap      uint
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
	cl.clientCap &= 0xffffffff00000000
	cl.clientCap = cl.clientCap | uint(baseCap)
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
	if uint(baseCap)&CLIENT_SSL == 1 {
		return cl.parseMessageSsl(buffer)
	}

	if uint(baseCap)&CLIENT_PROTOCOL_41 == 1 {
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
	var (
		extCap    uint16
		maxPacket uint32
		charset   uint8
		status    DecodeStatus
	)
	if extCap, status = readUint16(buffer); status != Success {
		log.DefaultLogger.Debugf("error when parsing cap flag of client ssl message")
		return Failure
	}
	cl.setExtendedClientCap(extCap)

	if maxPacket, status = readUint32(buffer); status != Success {
		log.DefaultLogger.Debugf("error when parsing max packet length of client ssl message")
		return Failure
	}
	cl.setMaxPacket(maxPacket)

	if charset, status = readUint8(buffer); status != Success {
		log.DefaultLogger.Debugf("error when parsing character of client ssl message")
		return Failure
	}
	cl.setCharset(charset)

	if skipBytes(buffer, UNSET_BYTES) != Success {
		log.DefaultLogger.Debugf("error when parsing reserved bytes of client ssl message")
		return Failure
	}
	return Success
}

func (cl *ClientLogin) parseMessage41(buffer types.IoBuffer) DecodeStatus {
	var (
		total     int
		extCap    uint16
		maxPacket uint32
		charset   uint8
		userName  string
		authLen   uint64
		authResp  []uint8
		status    DecodeStatus
	)

	total = buffer.Len()
	if extCap, status = readUint16(buffer); status != Success {
		log.DefaultLogger.Debugf("error when parsing client cap flag of client login message")
		return Failure
	}
	cl.setExtendedClientCap(extCap)

	if maxPacket, status = readUint32(buffer); status != Success {
		log.DefaultLogger.Debugf("error when parsing max packet length of client login message")
		return Failure
	}
	cl.setMaxPacket(maxPacket)

	if charset, status = readUint8(buffer); status != Success {
		log.DefaultLogger.Debugf("error when parsing charset of client login message")
		return Failure
	}
	cl.setCharset(charset)

	if status = skipBytes(buffer, UNSET_BYTES); status != Success {
		log.DefaultLogger.Debugf("error when skipping bytes of client login message")
		return Failure
	}

	if userName, status = readString(buffer); status != Success {
		log.DefaultLogger.Debugf("error when parsing username of client login message")
		return Failure
	}
	cl.setUsername(userName)

	if cl.clientCap&CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA == 1 {
		if authLen, status = readLengthEncodedInteger(buffer); status != Success {
			log.DefaultLogger.Debugf("error when parsing length of auth response of client login message")
			return Failure
		}

		if authResp, status = readVectorBySize(buffer, int(authLen)); status != Success {
			log.DefaultLogger.Debugf("error when parsing auth response of client login message")
			return Failure
		}
		cl.setAuthResp(authResp)

	} else if cl.clientCap&CLIENT_SECURE_CONNECTION == 1 {
		var al uint8
		if al, status = readUint8(buffer); status != Success {
			log.DefaultLogger.Debugf("error when parsing length of auth response of client login message")
			return Failure
		}
		if authResp, status = readVectorBySize(buffer, int(al)); status != Success {
			log.DefaultLogger.Debugf("error when parsing auth response of client login message")
			return Failure
		}
		cl.setAuthResp(authResp)

	} else {
		if _, status = readVectorBySize(buffer, len(authResp)); status != Success {
			log.DefaultLogger.Debugf("error when parsing auth response of client login message")
			return Failure
		}
	}

	if cl.clientCap&CLIENT_CONNECT_WITH_DB == 1 {
		if _, status = readString(buffer); status != Success {
			log.DefaultLogger.Debugf("error when parsing db name of client login message")
			return Failure
		}
	}

	if cl.clientCap&CLIENT_PLUGIN_AUTH == 1 {
		var authPluginName string
		if authPluginName, status = readString(buffer); status != Success {
			log.DefaultLogger.Debugf("error when parsing auth plugin name of client login message")
			return Failure
		}
		cl.setAuthPluginName(authPluginName)
	}

	var kvsLen uint64
	if cl.clientCap&CLIENT_CONNECT_ATTRS == 1 {
		if kvsLen, status = readLengthEncodedInteger(buffer); status != Success {
			log.DefaultLogger.Debugf("error when parsing length of all key-values in connection attributes of client login message")
			return Failure
		}

		for kvsLen > 0 {
			var strLen, prevLen uint64
			prevLen = uint64(buffer.Len())

			if strLen, status = readLengthEncodedInteger(buffer); status != Success {
				log.DefaultLogger.Debugf("error when parsing total length of connection attribute key in connection attributes of client login message")
				return Failure
			}
			var key string
			if key, status = readStringBySize(buffer, int64(strLen)); status != Success {
				log.DefaultLogger.Debugf("error when parsing connection attribute key in connection attributes of client login message")
				return Failure
			}

			if strLen, status = readLengthEncodedInteger(buffer); status != Success {
				log.DefaultLogger.Debugf("error when parsing length of connection attribute value in connection attributes of client login message")
				return Failure
			}
			var val string
			if val, status = readStringBySize(buffer, int64(strLen)); status != Success {
				log.DefaultLogger.Debugf("error when parsing connection attribute val in connection attributes of client login message")
				return Failure
			}
			cl.connAttr = append(cl.connAttr, []string{key, val})
			kvsLen -= prevLen - uint64(buffer.Len())
		}
	}
	log.DefaultLogger.Debugf("parsed client login protocol 41, consumed len %d, remain len %d", total-buffer.Len(), buffer.Len())
	return Success
}

func (cl *ClientLogin) parseMessage320(buffer types.IoBuffer, remainLen uint32) DecodeStatus {
	var (
		status    DecodeStatus
		maxPacket uint32
		userName  string
		authResp  []uint8
		db        string
	)

	originLen := buffer.Len()
	if maxPacket, status = readUint24(buffer); status != Success {
		log.DefaultLogger.Debugf("error when parsing max packet length of client login message")
		return Failure
	}
	cl.setMaxPacket(maxPacket)

	if userName, status = readString(buffer); status != Success {
		log.DefaultLogger.Debugf("error when parsing username of client login message")
		return Failure
	}
	cl.setUsername(userName)

	if cl.clientCap&CLIENT_CONNECT_WITH_DB == 1 {
		if authResp, status = readVectorBySize(buffer, len(cl.authResp)); status != Success {
			log.DefaultLogger.Debugf("error when parsing auth response of client login message")
			return Failure
		}
		cl.setAuthResp(authResp)

		if db, status = readString(buffer); status != Success {
			log.DefaultLogger.Debugf("error when parsing db name of client login message")
			return Failure
		}
		cl.setDb(db)

	} else {
		var comsumedLen int = originLen - buffer.Len()
		if authResp, status = readVectorBySize(buffer, int(remainLen)-comsumedLen); status != Success {
			log.DefaultLogger.Debugf("error when parsing auth response of client login message")
			return Failure
		}
	}
	return Success
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
