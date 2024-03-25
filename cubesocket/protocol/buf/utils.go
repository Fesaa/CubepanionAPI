package buf

import (
	"github.com/go-netty/go-netty/utils"
)

func GetVarIntSize(i int) int {
	for size := 1; size < 5; size++ {
		if i&(-1<<(size*7)) == 0 {
			return size
		}
	}

	return 5
}

type ByteWriter interface {
	WriteByte(byte) error
}

type ByteReader interface {
	ReadByte() (byte, error)
}

func WriteVarInt(b ByteWriter, i int) {
	for i&-128 != 0 {
		utils.Assert(b.WriteByte(byte((i & 127) | 128)))
		i = i >> 7
	}
	utils.Assert(b.WriteByte(byte(i)))
}

func ReadVarInt(r ByteReader) int {
	var numRead, result int = 0, 0

	var read byte
	var err error
	for {
		read, err = r.ReadByte()
		utils.Assert(err)

		result = result | ((int(read) & 127) << (numRead * 7))
		numRead++

		if numRead > 5 {
			panic("VarInt is too big")
		}

		if read&128 != 128 {
			break
		}
	}

	return result
}
