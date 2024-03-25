package buf

type PacketBuffer interface {
	ReadVarInt() int
	WriteVarInt(int)

	ReadBytes([]byte)
	WriteBytes([]byte)

	ReadByte() byte
	WriteByte(byte)

	ReadUUID() string
	WriteUUID(string)

	ReadBool() bool
	WriteBool(bool)

	ReadInt() int
	WriteInt(int)

	ReadLong() int64
	WriteLong(int64)

	ReadString() string
	WriteString(string)

	ReadableBytes() int
	GetBytes() []byte
	Grow(int)
}
