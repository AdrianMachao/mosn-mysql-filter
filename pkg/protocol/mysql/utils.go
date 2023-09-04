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

func readUint16(buf types.IoBuffer) (uint16, DecodeStatus) {
	if buf.Len() < 2 {
		return 0, Failure
	}

	data := buf.Peek(2)
	buf.Drain(1)

	return uint16(data[0]) | uint16(data[1])<<8, Success
}

func readUint24(buf types.IoBuffer, val uint32) DecodeStatus {
	return 0
}

func readUint32(buf types.IoBuffer) (uint32, DecodeStatus) {
	data := buf.Peek(4)
	val := uint32(uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3]<<24))
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

func readLengthEncodedInteger(buf types.IoBuffer) (uint8, DecodeStatus) {
	byteVal, status := readUint8(buf)
	if status == Failure {
		return 0, Failure
	}

	if byteVal < LENENCODINT_1BYTE {
		return byteVal, Success
	}

	if byteVal == LENENCODINT_2BYTES {
		
	}
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
		if v == 0 {
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
