package mysql

import (
	"encoding/binary"
	"mosn.io/mosn/pkg/types"
)

const (
	UNSET_BYTES = 23
)

func addUint8(buf types.IoBuffer, val uint8) {
	buf.WriteByte(val)
}

func addUint16(buf types.IoBuffer, val uint16) {
	buf.WriteUint16(val)
}
func addUint24(buf types.IoBuffer, val uint32) {

}

func addUint32(buf types.IoBuffer, val uint32) {
	buf.WriteUint32(val)
}

func addLengthEncodedInteger(buf types.IoBuffer, val uint64) {
}

func addBytes(buf types.IoBuffer, data []byte) {
	buf.Write(data)
}

func addString(buf types.IoBuffer, str string) {
	buf.WriteString(str)
}

func addVector(buf types.IoBuffer, data []uint8) {
	buf.Write(data)
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

func readUint16(buf types.IoBuffer) (uint16, DecodeStatus) {
	if buf.Len() < 2 {
		return 0, Failure
	}

	data := buf.Peek(2)
	buf.Drain(1)

	return uint16(data[0]) | uint16(data[1])<<8, Success
}

func readUint24(buf types.IoBuffer) (uint32, DecodeStatus) {
	val := binary.LittleEndian.Uint32(buf.Peek(3))
	buf.Drain(3)
	return val, Success
}

func readUint32(buf types.IoBuffer) (uint32, DecodeStatus) {
	val := binary.LittleEndian.Uint32(buf.Peek(4))
	buf.Drain(4)
	return val, Success
}

func readBytesBySize(buf types.IoBuffer, length int64) ([]byte, DecodeStatus) {
	if buf.Len() < int(length) {
		return nil, Failure
	}

	data := buf.Peek(int(length))
	buf.Drain(int(length))

	return data, Success
}

func readLengthEncodedInteger(buf types.IoBuffer) (uint64, DecodeStatus) {
	var val uint64
	byteVal, status := readUint8(buf)
	if status == Failure {
		return 0, Failure
	}
	switch byteVal {
	case LENENCODINT_2BYTES:
		val = uint64(binary.LittleEndian.Uint16(buf.Peek(2)))
		buf.Drain(2)
	case LENENCODINT_3BYTES:
		val = uint64(binary.LittleEndian.Uint32(buf.Peek(3)))
		buf.Drain(3)
	case LENENCODINT_8BYTES:
		val = binary.LittleEndian.Uint64(buf.Peek(8))
		buf.Drain(8)
	default:
		return 0, Failure
	}
	return val, Success
}

func skipBytes(buf types.IoBuffer, skipBytes int64) DecodeStatus {
	if buf.Len() < int(skipBytes) {
		return Failure
	}

	buf.Drain(int(skipBytes))
	return Success

}

func readString(buf types.IoBuffer) (string, DecodeStatus) {
	index := -1
	for i, v := range buf.Bytes() {
		if v == MYSQL_STR_END {
			index = i
			break
		}
	}
	if index == -1 {
		return "", Failure
	}

	str := string(buf.Peek(index + 1))
	buf.Drain(index + 1)

	return str, Success
}
func readVectorBySize(buf types.IoBuffer, length int) ([]byte, DecodeStatus) {
	if buf.Len() < length {
		return []byte{}, Failure
	}
	data := buf.Peek(length)
	buf.Drain(length)
	return data, Success
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
