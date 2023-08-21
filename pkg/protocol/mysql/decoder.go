package mysql

import (
	"encoding/binary"
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
	callBacks *DecoderCallbacks
	session   Session
}

type DecoderCallbacks interface {
	OnNewMessage()
	OnServerGreeting()
	OnClientLogin()
	OnClientLoginResponse()
	OnClientSwitchResponse()
	OnMoreClientLoginResponse()
	OnCommand()
	OnCommandResponse()
}

func (d *DecoderImpl) OnData(data types.IoBuffer) {
	for data.Len() != 0 && d.decode(data) {
	}
}

func (d *DecoderImpl) decode(data types.IoBuffer) bool {
	// check frame size
	payLoadLen := binary.LittleEndian.Uint32(data.Bytes())
	if data.Len() >= int(payLoadLen) {
		return true
	}

	return false
}

func (d *DecoderImpl) parseMessage(data types.IoBuffer, seq uint8, len uint32) {

}

//func (dfi *DecodeFactoryImpl) create(callbacks *DecoderCallbacks) *Decoder {
//	return nil
//}
