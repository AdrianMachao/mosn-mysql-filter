package mysql

import (
	"fmt"
	"unsafe"

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

func CreateDecoder(callbacks DecoderCallbacks) *DecoderImpl {
	p := &DecoderImpl{
		Callbacks: callbacks,
		session:   &Session{},
	}

	return p
}

type DecoderCallbacks interface {
	OnProtocolError()
	OnNewMessage(state State)
	OnServerGreeting(sg *ServerGreeting)
	OnClientLogin(cl *ClientLogin)
	OnClientLoginResponse(clr *ClientLoginResponse)
	OnClientSwitchResponse(c *Command)
	OnMoreClientLoginResponse(cr *ClientLoginResponse)
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
	pktLen := uint32(header[0]) | uint32(header[1])<<8 | uint32(header[2])<<16

	println(pktLen)
	//if 4+int(length) > data.Len() {
	//	return false
	//}
	data.Drain(4)
	d.Callbacks.OnNewMessage(d.session.getState())
	fmt.Println("------seq:", seq)
	if seq != d.session.getExpectedSeq() {
		if d.session.getState() == ReqResp && seq == MYSQL_REQUEST_PKT_NUM {
			d.session.setExpectedSeq(MYSQL_REQUEST_PKT_NUM)
			d.session.setState(Req)
		} else {
			log.DefaultLogger.Debugf("mysql_proxy: ignoring out-of-sync packet")
			d.Callbacks.OnProtocolError()
			data.Drain(int(unsafe.Sizeof(pktLen)))
			return true
		}
	}

	d.session.setState(State(seq + 1))
	// read packet body [pktLen bytes]
	d.parseMessage(data, seq, pktLen)
	data.Drain(int(unsafe.Sizeof(pktLen)))
	log.DefaultLogger.Debugf("mysql_proxy: %d bytes remaining in buffer", pktLen)
	return false
}

func (d *DecoderImpl) parseMessage(data types.IoBuffer, seq uint8, length uint32) {
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
		if clientLogin.IsSSLRequest() {
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
		header := data.Peek(1)
		if header == nil {
			d.session.setState(NotHandled)
			break
		}
		// resp_code

		switch header[0] {
		case MYSQL_RESP_OK:
			msg := &OkMessage{}
			d.session.setExpectedSeq(MYSQL_REQUEST_PKT_NUM)
			msg.decode(data, seq, length)
			d.session.setState(Req)
			d.Callbacks.OnClientLoginResponse(msg.ClientLoginResponse)
			break
		case MYSQL_RESP_AUTH_SWITCH:
			msg := &AuthMoreMessage{}
			msg.decode(data, seq, length)
			d.session.setState(AuthSwitchResp)
			d.Callbacks.OnClientLoginResponse(msg.ClientLoginResponse)
			break
		case MYSQL_RESP_ERR:
			msg := &ErrMessage{}
			msg.decode(data, seq, length)
			d.session.setState(Error)
			d.Callbacks.OnClientLoginResponse(msg.ClientLoginResponse)
			break
		case MYSQL_RESP_MORE:
			msg := &AuthMoreMessage{}
			msg.decode(data, seq, length)
			d.session.setState(NotHandled)
			d.Callbacks.OnClientLoginResponse(msg.ClientLoginResponse)
			break
		default:
			d.session.setState(NotHandled)
			d.Callbacks.OnClientLoginResponse(&ClientLoginResponse{})
			return
		}

	case SslPt:
		data.Drain(data.Len())
		break
	case AuthSwitchReq:

	case AuthSwitchReqOld:
	case AuthSwitchResp:
		clientSwitchResponse := &ClientSwitchResponse{}
		clientSwitchResponse.decode(data, seq, length)
		d.session.setState(AuthSwitchMore)
		break
	case AuthSwitchMore:
		header := data.Peek(1)
		if header == nil {
			d.session.setState(NotHandled)
			break
		}
		switch header[0] {
		case MYSQL_RESP_OK:
			msg := &OkMessage{}
			d.session.setExpectedSeq(MYSQL_REQUEST_PKT_NUM)
			d.session.setState(Req)
			msg.decode(data, seq, length)
			d.Callbacks.OnMoreClientLoginResponse(msg.ClientLoginResponse)
			break
		case MYSQL_RESP_MORE:
			msg := &AuthMoreMessage{}
			d.session.setState(AuthSwitchResp)
			msg.decode(data, seq, length)
			d.Callbacks.OnMoreClientLoginResponse(msg.ClientLoginResponse)
			break
		case MYSQL_RESP_ERR:
			msg := &ErrMessage{}
			d.session.setExpectedSeq(MYSQL_REQUEST_PKT_NUM)
			d.session.setState(Resync)
			msg.decode(data, seq, length)
			d.Callbacks.OnMoreClientLoginResponse(msg.ClientLoginResponse)
			break
		case MYSQL_RESP_AUTH_SWITCH:
			msg := &AuthSwitchMessage{}
			d.session.setState(NotHandled)
			msg.decode(data, seq, length)
			d.Callbacks.OnMoreClientLoginResponse(msg.ClientLoginResponse)
		default:
			d.session.setState(NotHandled)
			d.Callbacks.OnMoreClientLoginResponse(&ClientLoginResponse{})
			return
		}
	case ReqResp:
		commandResp := &CommandResponse{}
		commandResp.decode(data, seq, length)
		d.Callbacks.OnCommandResponse(commandResp)
	case Req:
		command := &Command{}
		command.decode(data, seq, int(length))
		d.Callbacks.OnCommand(command)
	case Resync:
		// re-sync to MYSQL_REQ state
		// expected seq check succeeded, no need to verify
		d.session.setState(Req)
		fallthrough
	case NotHandled:
	case Error:
	}
}
