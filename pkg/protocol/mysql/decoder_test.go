package mysql

import (
	"mosn.io/mosn/pkg/types"
	"mosn.io/pkg/buffer"
	"testing"
)

type Callback struct {
	state *State
}

func (c Callback) OnProtocolError() {
	//TODO implement me
}

func (c Callback) OnNewMessage(state State) {
	//TODO implement me
}

func (c Callback) OnServerGreeting(sg *ServerGreeting) {
	//TODO implement me
}

func (c Callback) OnClientLogin(cl *ClientLogin) {
	//TODO implement me
}

func (c Callback) OnClientLoginResponse(clr *ClientLoginResponse) {
	//TODO implement me
}

func (c Callback) OnClientSwitchResponse(cc *Command) {
	//TODO implement me
}

func (c Callback) OnMoreClientLoginResponse(cr *ClientLoginResponse) {
	//TODO implement me
}

func (c Callback) OnCommand(ccc *Command) {
	//TODO implement me
}

func (c Callback) OnCommandResponse(cr *CommandResponse) {
	//TODO implement me
}

func TestDecoderImpl_OnData(t *testing.T) {
	type fields struct {
		Decoder   Decoder
		Callbacks DecoderCallbacks
		session   *Session
	}
	type args struct {
		data types.IoBuffer
	}
	var pktLen = 8
	data := make([]byte, 512)
	data[0] = byte(pktLen)
	data[1] = byte(pktLen >> 8)
	data[2] = byte(pktLen >> 16)
	data[3] = 0

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name: "bob",
			fields: fields{
				session:   &Session{state: Init},
				Callbacks: Callback{},
			},
			args: args{data: buffer.NewIoBufferBytes(data)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DecoderImpl{
				Decoder:   tt.fields.Decoder,
				Callbacks: tt.fields.Callbacks,
				session:   tt.fields.session,
			}
			d.OnData(tt.args.data)
		})
	}
}

func TestDecoderImpl_decode(t *testing.T) {
	type fields struct {
		Decoder   Decoder
		Callbacks DecoderCallbacks
		session   *Session
	}
	type args struct {
		data types.IoBuffer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DecoderImpl{
				Decoder:   tt.fields.Decoder,
				Callbacks: tt.fields.Callbacks,
				session:   tt.fields.session,
			}
			if got := d.decode(tt.args.data); got != tt.want {
				t.Errorf("decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoderImpl_parseMessage(t *testing.T) {
	type fields struct {
		Decoder   Decoder
		Callbacks DecoderCallbacks
		session   *Session
	}
	type args struct {
		data   types.IoBuffer
		seq    uint8
		length int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DecoderImpl{
				Decoder:   tt.fields.Decoder,
				Callbacks: tt.fields.Callbacks,
				session:   tt.fields.session,
			}
			d.parseMessage(tt.args.data, tt.args.seq, tt.args.length)
		})
	}
}
