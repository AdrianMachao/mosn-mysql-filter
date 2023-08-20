package mysql

import "mosn.io/mosn/pkg/types"

type MySQLCodec struct {
	seq uint8
}

type DecodeStatus uint8

const (
	Success DecodeStatus = iota
	Failure
)

func (m *MySQLCodec) parseMessage(data types.IoBuffer) DecodeStatus {
	return 0
}
