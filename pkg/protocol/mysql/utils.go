package mysql

import (
	"mosn.io/mosn/pkg/types"
)

func addUint8(buf types.IoBuffer, val uint8) {
}

func addUint16(buf types.IoBuffer, val uint16) {
}
func addUint24(buf types.IoBuffer, val uint32) {
}
func addUint32(buf types.IoBuffer, val uint32) {
}
func addLengthEncodedInteger(buf types.IoBuffer, val uint64) {
}

func addBytes(buf types.IoBuffer, data byte, val uint64) {
}
func addString(buf types.IoBuffer, str string) {
}
func addVector(buf types.IoBuffer, data []uint8) {
}
func encodeHdr(buf types.IoBuffer, seq uint8) {
}
func endOfBuffer(buf types.IoBuffer) {
}

func readUint8(buf types.IoBuffer) (uint8, DecodeStatus) {
	if buf.Len() < 1 {
		return 0, Failure
	}

	data := buf.Peek(1)
	buf.Drain(1)

	return data[0], Success
}

func readUint16(buf types.IoBuffer, val uint16) DecodeStatus {
	return 0
}

func readUint24(buf types.IoBuffer, val uint32) DecodeStatus {
	return 0
}

func readUint32(buf types.IoBuffer, val uint32) DecodeStatus {
	return 0
}

func readLengthEncodedInteger(buf types.IoBuffer) (uint8, DecodeStatus) {
	byteVal, status := readUint8(buf)
	if status == Failure {
		return 0, Failure
	}
	switch byteVal {
	case LENENCODINT_2BYTES:
		val := buf.Peek(2)
		buf.Drain(2)
	case LENENCODINT_3BYTES:
		val := buf.Peek(3)
		buf.Drain(3)
	case LENENCODINT_8BYTES:
		val := buf.Peek(8)
		buf.Drain(8)
	default:
		return 0, Failure
	}

}

func skipBytes(buf types.IoBuffer, skipBytes int64) DecodeStatus {
	return 0
}

func readString(buf types.IoBuffer, str string) DecodeStatus {
	return 0
}
func readVector(buf types.IoBuffer, data []uint8) DecodeStatus {
	return 0
}

func readStringBySize(buf types.IoBuffer, length int64) (string, DecodeStatus) {
	if buf.Len() < int(length) {
		return "", Failure
	}

	data := buf.Peek(int(length))
	buf.Drain(int(length))

	return string(data), Success
}

func readAll(buf types.IoBuffer, str string) DecodeStatus {
	return 0
}
func peekUint32(buf types.IoBuffer, val uint32) DecodeStatus {
	//buf.WriteUint32()
	return 0
}
func peekUint8(buf types.IoBuffer, val uint32) DecodeStatus {
	return 0
}
func consumeHdr(buf types.IoBuffer) DecodeStatus {
	return 0
}
func peekHdr(buf types.IoBuffer, length uint32, seq uint8) DecodeStatus {
	//var val uint32
	// TODO 修改成指针
	//buf.WriteUint32()

	return 0
}
