package mysql

import (
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/types"
)

//type DecodeFactoryImpl struct {
//	DecoderFactory
//}
//
//type DecoderFactory interface {
//	create(callbacks *DecoderCallbacks) *Decoder
//}

type Decoder struct {
	attributes map[string]string
}

type DecoderImpl struct {
	Decoder
	Callbacks DecoderCallbacks
	session   *Session
}

type DecoderCallbacks interface {
	OnProtocolError()
	OnNewMessage(state State)
	OnServerGreeting(sg *ServerGreeting)
	OnClientLogin(cl *ClientLogin)
	OnClientLoginResponse(clr *ClientLoginResponse)
	OnClientSwitchResponse(c *Command)
	OnMoreClientLoginResponse(cr *CommandResponse)
	OnCommand(c *Command)
	OnCommandResponse(cr *CommandResponse)
}

func (d *DecoderImpl) OnData(data types.IoBuffer) {
	for data.Len() != 0 && d.decode(data) {
	}
}

func (d *DecoderImpl) decode(data types.IoBuffer) bool {
	//var length uint32
	var seq uint8
	//val := binary.LittleEndian.Uint32(data.Bytes())
	//if data.Len() < binary.Size(data) {
	//	return false
	//}
	//seq = int(endian.Htobe32(val) & MYSQL_HDR_SEQ_MASK)
	//seq = int(uint32(dataBytes[1]) & MYSQL_HDR_SEQ_MASK)
	//length = val & MYSQL_HDR_PKT_SIZE_MASK

	//endian.Le32toh()
	// uint32 -> size = 4

	// read packet header
	header := data.Peek(4)
	seq = header[3]
	pktLen := int(uint32(header[0]) | uint32(header[1])<<8 | uint32(header[2])<<16)

	println(pktLen)
	//if 4+int(length) > data.Len() {
	//	return false
	//}
	data.Drain(4)
	d.Callbacks.OnNewMessage(d.session.getState())

	if seq != d.session.getExpectedSeq() {
		if d.session.getState() == ReqResp && seq == MYSQL_REQUEST_PKT_NUM {
			d.session.setExpectedSeq(MYSQL_REQUEST_PKT_NUM)
			d.session.setState(Req)
		} else {
			log.DefaultLogger.Debugf("mysql_proxy: ignoring out-of-sync packet")
			d.Callbacks.OnProtocolError()
			data.Drain(pktLen)
			return true
		}
	}

	d.session.setState(State(seq + 1))
	// read packet body [pktLen bytes]
	d.parseMessage(data, seq, pktLen)
	data.Drain(pktLen)
	log.DefaultLogger.Debugf("mysql_proxy: %d bytes remaining in buffer", pktLen)
	return false
}

func (d *DecoderImpl) parseMessage(data types.IoBuffer, seq uint8, length int) {
	if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("")
	}
	switch d.session.getState() {
	case Init:
		var greeting ServerGreeting
		greeting.decode(data, seq, length)
		d.session.setState(ChallengeReq)
		d.Callbacks.OnServerGreeting(&greeting)
		break
	case ChallengeReq:
		var clientLogin ClientLogin
		clientLogin.decode(data, seq, length)
		if clientLogin.isSSLRequest() {
			d.session.setState(SslPt)
		} else if clientLogin.isResponse41() {
			d.session.setState(ChallengeResp41)
		} else {
			d.session.setState(ChallengeResp320)
		}
		d.Callbacks.OnClientLogin(&clientLogin)
		break
	case ChallengeResp41:
	case ChallengeResp320:
		//var respCode int
		// TODO read buf
		d.session.setState(NotHandled)
		break
	case SslPt:
		data.Drain(data.Len())
		break
	case AuthSwitchReq:

	case AuthSwitchReqOld:
	case AuthSwitchResp:
	case AuthSwitchMore:
	case ReqResp:
	case Req:
	case Resync:
	case NotHandled:
	case Error:
	}
}
