package mysql

import (
	"encoding/binary"

	"github.com/mattn/go-gnulib/endian"
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
	session   Session
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
	// check frame size
	var length uint32
	var seq uint8

	val := binary.LittleEndian.Uint32(data.Bytes())

	if data.Len() < binary.Size(data) {
		return false
	}
	seq = uint8(endian.Htobe32(val) & MYSQL_HDR_SEQ_MASK)
	length = val & MYSQL_HDR_PKT_SIZE_MASK

	//endian.Le32toh()
	// uint32 -> size = 4
	if 4+int(length) > data.Len() {
		return false
	}

	data.Drain(4)
	d.Callbacks.OnNewMessage(d.session.getState())

	if seq != d.session.getExpectedSeq() {
		if d.session.getState() == ReqResp && seq == MYSQL_REQUEST_PKT_NUM {
			d.session.setExpectedSeq(MYSQL_REQUEST_PKT_NUM)
			d.session.setState(Req)
		} else {
			log.DefaultLogger.Debugf("mysql_proxy: ignoring out-of-sync packet")
			d.Callbacks.OnProtocolError()
			data.Drain(int(length))
			return true
		}
	}

	d.session.setState(State(seq + 1))
	dataLen := data.Len()
	d.parseMessage(data, seq, dataLen)

	consumedLen := dataLen - data.Len()
	data.Drain(int(length) - consumedLen)
	log.DefaultLogger.Debugf("mysql_proxy: %d bytes remaining in buffer", data.Len())
	return false
}

func (d *DecoderImpl) parseMessage(data types.IoBuffer, seq uint8, len int) {
	if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("")
	}
	switch d.session.getState() {
	case Init:
	case ChallengeReq:
	case ChallengeResp41:
	case ChallengeResp320:
	case SslPt:
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

//func (dfi *DecodeFactoryImpl) create(callbacks *DecoderCallbacks) *Decoder {
//	return nil
//}
