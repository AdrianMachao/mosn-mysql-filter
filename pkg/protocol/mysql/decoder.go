package mysql

import (
	"encoding/binary"
	"mosn.io/mosn/pkg/types"
)

type Decoder struct {
	attributes_ map[string]string
}

type DecoderImpl struct {
	callBacks DecoderCallbacks
	session   Session
}

type DecoderCallbacks interface {
	onNewMessage()
	onServerGreeting()
	onClientLogin()
	onClientLoginResponse()
	onClientSwitchResponse()
	onMoreClientLoginResponse()
	onCommand()
	onCommandResponse()
}

func (d *DecoderImpl) OnData(data types.IoBuffer) {
	d.decode(data)
}

func (d *DecoderImpl) decode(data types.IoBuffer) DecodeStatus {
	// check frame size
	payLoadLen := binary.LittleEndian.Uint32(data.Bytes())
	if data.Len() >= int(payLoadLen) {
		return 0
	}

	return 0
}

func (d *DecoderImpl) parseMessage(data types.IoBuffer, seq uint8, len uint32) {

}
