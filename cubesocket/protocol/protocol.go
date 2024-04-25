package protocol

import (
	"github.com/Fesaa/CubepanionAPI/cubesocket/protocol/packets"
)

var protocolRegistry = make(map[int]packets.Packet)

func init() {
	RegisterPacket(0, &packets.PacketPing{})
	RegisterPacket(1, &packets.PacketPong{})
	RegisterPacket(2, &packets.PacketHelloPing{})
	RegisterPacket(3, &packets.PacketHelloPong{})
	RegisterPacket(4, &packets.PacketLocationUpdate{})
	RegisterPacket(5, &packets.PacketPerkUpdate{})
	RegisterPacket(6, &packets.PacketDisconnection{})
	RegisterPacket(7, &packets.PacketLogin{})
	RegisterPacket(8, &packets.PacketLoginComplete{})
	RegisterPacket(9, &packets.PacketSetProtocolVersion{})
}

func RegisterPacket(id int, packet packets.Packet) {
	protocolRegistry[id] = packet
}

func PacketFromId(id int) packets.Packet {
	if packet, ok := protocolRegistry[id]; ok {
		return packet
	}
	return nil
}
