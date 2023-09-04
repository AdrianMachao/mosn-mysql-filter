package mysql

import (
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/types"
)

type ServerGreeting struct {
	MySQLCodec
	protocol        uint8
	version         string
	threadId        uint32
	authPluginData1 []uint8
	authPluginData2 []uint8
	serverCap       uint
	serverCharset   uint8
	serverStatus    uint16
	authPluginName  string
}

func (sg *ServerGreeting) getProtocol() uint8 {
	return sg.protocol
}

func (sg *ServerGreeting) getVersion() string {
	return sg.version
}

func (sg *ServerGreeting) getThreadId() uint32 {
	return sg.threadId
}

func (sg *ServerGreeting) getAuthPluginData1() []uint8 {
	return sg.authPluginData1
}

func (sg *ServerGreeting) getAuthPluginData2() []uint8 {
	return sg.authPluginData2
}

func (sg *ServerGreeting) getAuthPluginData() []uint8 {
	if (sg.serverCap&CLIENT_PLUGIN_AUTH) > 0 || (sg.serverCap&CLIENT_SECURE_CONNECTION) > 0 {
		res := make([]uint8, 0, len(sg.authPluginData1)+len(sg.authPluginData2))
		res = append(res, sg.authPluginData1...)
		res = append(res, sg.authPluginData2...)
		return res
	}

	return sg.authPluginData1
}

func (sg *ServerGreeting) setProtocol(protocol uint8) {
	sg.protocol = protocol
}

func (sg *ServerGreeting) setVersion(version string) {
	sg.version = version
}

func (sg *ServerGreeting) setThreadId(threadId uint32) {
	sg.threadId = threadId
}

func (sg *ServerGreeting) setServerCap(serverCap uint) {
	sg.serverCap = serverCap
}

func (sg *ServerGreeting) setBaseServerCap(baseServerCap uint16) {
	sg.serverCap = sg.serverCap & 0xffffffff00000000
	sg.serverCap = sg.serverCap | uint(baseServerCap)
}

func (sg *ServerGreeting) setExtServerCap(extServerCap uint16) {
	ext := extServerCap
	sg.serverCap = sg.serverCap & 0xffffffff00000000
	sg.serverCap = sg.serverCap | (uint(ext) << 16)
}

func (sg *ServerGreeting) setAuthPluginName(name string) {
	sg.authPluginName = name
}

func (sg *ServerGreeting) setAuthPluginData(salt []uint8) {
	if len(salt) <= 8 {
		sg.authPluginData1 = salt
		return
	}

	copy(sg.authPluginData1, salt[:8])
	copy(sg.authPluginData1, salt)
}

func (sg *ServerGreeting) setAuthPluginData1(salt []uint8) {
	sg.authPluginData1 = salt
}

func (sg *ServerGreeting) setAuthPluginData2(salt []uint8) {
	sg.authPluginData2 = salt
}

func (sg *ServerGreeting) setServerCharset(serverLanguage uint8) {
	sg.serverCharset = serverLanguage
}

func (sg *ServerGreeting) setServerStatus(serverStatus uint16) {
	sg.serverStatus = serverStatus
}

func (sg *ServerGreeting) check() DecodeStatus {
	return 0
}

func (sg *ServerGreeting) decode(data types.IoBuffer, seq uint8, length int) DecodeStatus {
	sg.seq = seq
	return sg.parseMessage(data, length)
}

func (sg *ServerGreeting) parseMessage(buffer types.IoBuffer, length int) DecodeStatus {
	protocol, status := readUint8(buffer)
	if status != Success {
		log.DefaultLogger.Debugf("error when parsing protocol in mysql greeting msg")
		return Failure
	}
	sg.protocol = protocol

	version, status := readString(buffer)
	if status != Success {
		log.DefaultLogger.Debugf("error when parsing version in mysql greeting msg")
		return Failure
	}
	sg.version = version

	threadId, status := readUint32(buffer)
	if status != Success {
		log.DefaultLogger.Debugf("error when parsing threadId in mysql greeting msg")
		return Failure
	}
	sg.threadId = threadId

	authPluginData1, status := readBytesBySize(buffer, 8)
	if status != Success {
		log.DefaultLogger.Debugf("error when parsing authPluginData1 in mysql greeting msg")
		return Failure
	}
	sg.authPluginData1 = authPluginData1

	if skipBytes(buffer, 1) != Success {
		log.DefaultLogger.Debugf("error skipping bytes in mysql greeting msg")
		return Failure
	}

	if sg.protocol == MYSQL_PROTOCOL_9 {
		return Success
	}

	var baseServerCap uint16 = 0
	baseServerCap, status = readUint16(buffer)
	if status != Success {
		log.DefaultLogger.Debugf("error when parsing cap flag[lower 2 bytes]  in mysql greeting msg")
		return Failure
	}
	sg.setBaseServerCap(baseServerCap)

	serverCharset, status := readUint8(buffer)
	if status != Success {
		log.DefaultLogger.Debugf("error when parsing server charset  in mysql greeting msg")
		return Failure
	}
	sg.serverCharset = serverCharset

	serverStatus, status := readUint16(buffer)
	if status != Success {
		log.DefaultLogger.Debugf("error when parsing server status in mysql greeting msg")
		return Failure
	}
	sg.serverStatus = serverStatus

	extServerCap, status := readUint16(buffer)
	if status != Success {
		log.DefaultLogger.Debugf("error when parsing server cap in mysql greeting msg")
		return Failure
	}
	sg.setExtServerCap(extServerCap)

	authPluginDataLen, status := readUint8(buffer)
	if status != Success {
		log.DefaultLogger.Debugf("error when parsing length of auth plugin data of mysql greeting msg")
		return Failure
	}

	if status := skipBytes(buffer, 10); status != Success {
		log.DefaultLogger.Debugf("error when parsing reserved bytes of mysql greeting msg")
		return Failure
	}

	if (sg.serverCap & CLIENT_PLUGIN_AUTH) > 0 {
		var authPluginDataLen2 uint8 = 0
		if authPluginDataLen > 8 {
			authPluginDataLen2 = authPluginDataLen - 8
		}

		data, status := readBytesBySize(buffer, int64(authPluginDataLen2))
		if status != Success {
			log.DefaultLogger.Debugf("error when reading bytes of mysql greeting msg")
			return Failure
		}
		copy(sg.authPluginData2, data)

		var skipedBytes uint8 = 0
		if 13 > authPluginDataLen2 {
			skipedBytes = uint8(13) - authPluginDataLen2
		}

		if status := skipBytes(buffer, int64(skipedBytes)); status != Success {
			log.DefaultLogger.Debugf("error when parsing reserved bytes of mysql greeting msg")
			return Failure
		}

		authPluginName, status := readString(buffer)
		if status != Success {
			log.DefaultLogger.Debugf("error when parsing auth plugin name of  mysql greeting msg")
			return Failure
		}
		sg.authPluginName = authPluginName
	} else if (sg.serverCap & CLIENT_SECURE_CONNECTION) > 0 {
		data, status := readBytesBySize(buffer, 12)
		if status != Success {
			log.DefaultLogger.Debugf("error when reading bytes of mysql greeting msg")
			return Failure
		}
		copy(sg.authPluginData2, data)

		if status := skipBytes(buffer, 1); status != Success {
			log.DefaultLogger.Debugf("error when parsing reserved bytes of mysql greeting msg")
			return Failure
		}
	}

	authPluginLen := len(sg.authPluginData1) + len(sg.authPluginData2)
	if (sg.serverCap & CLIENT_PLUGIN_AUTH) > 0 {
		if authPluginLen != int(authPluginDataLen) {
			log.DefaultLogger.Debugf("error when final check failure of mysql greeting msg")
			return Failure
		} else if (sg.serverCap & CLIENT_SECURE_CONNECTION) > 0 {
			if authPluginLen != 20 && authPluginDataLen != 0 {
				log.DefaultLogger.Debugf("error when final check failure of mysql greeting msg")
				return Failure
			}
		}
	}

	return Success
}
