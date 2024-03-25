package buf

import (
	"fmt"
	"runtime"
)

const (
	READABLE_BYTES_FORMAT = "Not enough bytes to read %d. Have %d left"
	WRITABLE_BYTES_FORMAT = "Not enough capacity to write %d bytes. Have %d left"
)

type packetBufferImpl struct {
	_bytes []byte
	read   int
	write  int
}

const defaultSize = 1024

func NewPacketBuffer(bytes []byte) PacketBuffer {
	var write int
	if bytes == nil {
		bytes = make([]byte, defaultSize)
		write = 0
	} else {
		write = len(bytes)
	}
	return &packetBufferImpl{
		_bytes: bytes,
		read:   0,
		write:  write,
	}
}

func (b *packetBufferImpl) checkReadableBytes(min int) {
	if b.write-b.read < min {
		fmt.Println(runtime.Caller(4))
		fmt.Println(runtime.Caller(3))
		fmt.Println(runtime.Caller(2))
		fmt.Println(runtime.Caller(1))
		panic(fmt.Sprintf(READABLE_BYTES_FORMAT, min, b.write-b.read))
	}
}

func (b *packetBufferImpl) checkWritableBytes(min int) {
	if len(b._bytes)-b.write < min {
		fmt.Println(runtime.Caller(4))
		fmt.Println(runtime.Caller(3))
		fmt.Println(runtime.Caller(2))
		fmt.Println(runtime.Caller(1))
		panic(fmt.Sprintf(WRITABLE_BYTES_FORMAT, min, len(b._bytes)-b.write))
	}
}

func (b *packetBufferImpl) ReadBytes(p []byte) {
	b.checkReadableBytes(len(p))
	b.read += len(p)
	copy(p, b._bytes[b.read-len(p):b.read])
}

func (b *packetBufferImpl) WriteBytes(bytes []byte) {
	b.checkWritableBytes(len(bytes))
	for _, byte := range bytes {
		b.WriteByte(byte)
	}
}

func (b *packetBufferImpl) ReadByte() byte {
	b.checkReadableBytes(1)
	byte := b._bytes[b.read]
	b.read++
	return byte
}

func (b *packetBufferImpl) WriteByte(byte byte) {
	b.checkWritableBytes(1)
	b._bytes[b.write] = byte
	b.write++
}

func (b *packetBufferImpl) ReadBool() bool {
	return b.ReadByte() != 0
}

func (b *packetBufferImpl) WriteBool(bool bool) {
	if bool {
		b.WriteByte(1)
	} else {
		b.WriteByte(0)
	}
}

func (b *packetBufferImpl) ReadInt() int {
	b.checkReadableBytes(4)
	var result int = 0
	for i := 0; i < 4; i++ {
		read := b.ReadByte()
		result += int(read) << ((3 - i) * 8)
	}
	return result
}

func (b *packetBufferImpl) WriteInt(i int) {
	b.checkWritableBytes(4)
	bytes := make([]byte, 4)
	for j := 0; j < 4; j++ {
		bytes[j] = byte(i >> ((3 - j) * 8))
	}
	b.WriteBytes(bytes)
}

func (b *packetBufferImpl) ReadVarInt() int {
	var numRead, result int = 0, 0

	var read byte
	for {
		read = b.ReadByte()
		result = result | ((int(read) & 127) << (numRead * 7))
		numRead++

		if numRead > 5 {
			panic("VarInt is too big")
		}

		if read&128 == 0 {
			break
		}
	}

	return result
}

func (b *packetBufferImpl) WriteVarInt(i int) {
	for i&-128 != 0 {
		b.WriteByte(byte((i & 127) | 128))
		i = i >> 7
	}
	b.WriteByte(byte(i))
}

func (b *packetBufferImpl) ReadLong() int64 {
	b.checkReadableBytes(8)
	var result int64 = 0
	for i := 0; i < 8; i++ {
		read := b.ReadByte()
		result += int64(read) << ((7 - i) * 8)
	}
	return result
}

func (b *packetBufferImpl) WriteLong(i int64) {
	b.checkWritableBytes(8)
	bytes := make([]byte, 8)
	for j := 0; j < 8; j++ {
		bytes[j] = byte(i >> ((7 - j) * 8))
	}
	b.WriteBytes(bytes)
}

func (b *packetBufferImpl) ReadString() string {
	length := b.ReadInt()
	b.checkReadableBytes(length)
	bytes := make([]byte, length)

	for i := 0; i < length; i++ {
		bytes[i] = b.ReadByte()
	}
	return string(bytes)
}

func (b *packetBufferImpl) WriteString(str string) {
	bytes := []byte(str)
	b.WriteInt(len(bytes))
	b.WriteBytes(bytes)
}

func (b *packetBufferImpl) ReadUUID() string {
	return b.ReadString()
}

func (b *packetBufferImpl) WriteUUID(uuid string) {
	b.WriteString(uuid)
}

func (b *packetBufferImpl) ReadableBytes() int {
	return b.write - b.read
}

func (b *packetBufferImpl) GetBytes() []byte {
	return b._bytes[b.read:b.write]
}

func (b *packetBufferImpl) Grow(n int) {
	if len(b._bytes) < n {
		return
	}

	newBytes := make([]byte, len(b._bytes)+n)
	copy(newBytes, b._bytes)
	b._bytes = newBytes
}
