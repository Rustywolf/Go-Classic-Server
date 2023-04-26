package packets

import (
	. "classicserver/classic/constants"
	"net"
)

//
// Server Identification
//

type DownstreamServerIdentification struct {
	DownstreamPacket
	Version    uint8
	ServerName string
	ServerMOTD string
	PlayerMode PlayerMode
}

func NewDownstreamServerIdentification(serverName string, serverMOTD string, playerMode PlayerMode) DownstreamServerIdentification {
	packet := DownstreamServerIdentification{
		DownstreamPacket: DownstreamPacket{
			Id: DOWNSTREAM_SERVER_IDENTIFICATION,
		},
		Version:    0x07,
		ServerName: serverName,
		ServerMOTD: serverMOTD,
		PlayerMode: playerMode,
	}
	return packet
}

func (packet DownstreamServerIdentification) Write(conn net.Conn) error {
	buffer := []byte{byte(packet.Id)}
	buffer = append(buffer, byte(packet.Version))
	buffer = append(buffer, writeString(packet.ServerName)...)
	buffer = append(buffer, writeString(packet.ServerMOTD)...)
	buffer = append(buffer, byte(packet.PlayerMode))
	_, err := conn.Write(buffer)
	return err
}

//
// Client Ping
//

type DownstreamPing struct {
	DownstreamPacket
}

func NewDownstreamPing() DownstreamPing {
	return DownstreamPing{
		DownstreamPacket: DownstreamPacket{
			Id: DOWNSTREAM_PING,
		},
	}
}

func (packet *DownstreamPing) Write(conn net.Conn) error {
	buffer := []byte{byte(packet.Id)}
	_, err := conn.Write(buffer)
	return err
}

//
// Level Data Transfer Initialization
//

type DownstreamLevelInit struct {
	DownstreamPacket
}

func NewDownstreamLevelInit() DownstreamLevelInit {
	return DownstreamLevelInit{
		DownstreamPacket: DownstreamPacket{
			Id: DOWNSTREAM_LEVEL_INIT,
		},
	}
}

func (packet DownstreamLevelInit) Write(conn net.Conn) error {
	buffer := []byte{byte(packet.Id)}
	_, err := conn.Write(buffer)
	return err
}

//
// Level Data Transfer Chunk
//

type DownstreamLevelChunk struct {
	DownstreamPacket
	Length  int16
	Data    []uint8
	Percent uint8
}

func NewDownstreamLevelChunk(data []uint8, percent uint8) DownstreamLevelChunk {
	dataSize := len(data)
	if dataSize < 1024 {
		data = append(data, make([]uint8, 1024-dataSize)...)
	}

	return DownstreamLevelChunk{
		DownstreamPacket: DownstreamPacket{
			Id: DOWNSTREAM_LEVEL_CHUNK,
		},
		Length:  int16(dataSize),
		Data:    data,
		Percent: percent,
	}
}

func (packet DownstreamLevelChunk) Write(conn net.Conn) error {
	buffer := []byte{byte(packet.Id)}
	buffer = append(buffer, writeShort(int16(packet.Length))...)
	buffer = append(buffer, packet.Data...)
	buffer = append(buffer, packet.Percent)
	_, err := conn.Write(buffer)
	return err
}

//
// Level Data Transfer Finalize
//

type DownstreamLevelFinalize struct {
	DownstreamPacket
	X int16
	Y int16
	Z int16
}

func NewDownstreamLevelFinalize(x int16, y int16, z int16) DownstreamLevelFinalize {
	return DownstreamLevelFinalize{
		DownstreamPacket: DownstreamPacket{
			Id: DOWNSTREAM_LEVEL_FINALIZE,
		},
		X: x,
		Y: y,
		Z: z,
	}
}

func (packet DownstreamLevelFinalize) Write(conn net.Conn) error {
	buffer := []byte{byte(packet.Id)}
	buffer = append(buffer, writeShort(packet.X)...)
	buffer = append(buffer, writeShort(packet.Y)...)
	buffer = append(buffer, writeShort(packet.Z)...)
	_, err := conn.Write(buffer)
	return err
}

//
// Set Block
//

type DownstreamSetBlock struct {
	DownstreamPacket
	X     int16
	Y     int16
	Z     int16
	Block uint8
}

func NewDownstreamSetBlock(x int16, y int16, z int16, block Block) DownstreamSetBlock {
	return DownstreamSetBlock{
		DownstreamPacket: DownstreamPacket{
			Id: DOWNSTREAM_SET_BLOCK,
		},
		X:     x,
		Y:     y,
		Z:     z,
		Block: block,
	}
}

func (packet DownstreamSetBlock) Write(conn net.Conn) error {
	buffer := []byte{byte(packet.Id)}
	buffer = append(buffer, writeShort(packet.X)...)
	buffer = append(buffer, writeShort(packet.Y)...)
	buffer = append(buffer, writeShort(packet.Z)...)
	buffer = append(buffer, packet.Block)
	_, err := conn.Write(buffer)
	return err
}

//
// Spawn Player
//

type DownstreamSpawnPlayer struct {
	DownstreamPacket
	PlayerId   int8
	PlayerName string
	X          FPShort
	Y          FPShort
	Z          FPShort
	Yaw        uint8
	Pitch      uint8
}

func NewDownstreamSpawnPlayer(playerId int8, playerName string, x FPShort, y FPShort, z FPShort, yaw uint8, pitch uint8) DownstreamSpawnPlayer {
	return DownstreamSpawnPlayer{
		DownstreamPacket: DownstreamPacket{
			Id: DOWNSTREAM_SPAWN_PLAYER,
		},
		PlayerId:   playerId,
		PlayerName: playerName,
		X:          x,
		Y:          y,
		Z:          z,
		Yaw:        yaw,
		Pitch:      pitch,
	}
}

func (packet DownstreamSpawnPlayer) Write(conn net.Conn) error {
	buffer := []byte{byte(packet.Id)}
	buffer = append(buffer, byte(packet.PlayerId))
	buffer = append(buffer, writeString(packet.PlayerName)...)
	buffer = append(buffer, writeFPShort(packet.X)...)
	buffer = append(buffer, writeFPShort(packet.Y)...)
	buffer = append(buffer, writeFPShort(packet.Z)...)
	buffer = append(buffer, packet.Yaw)
	buffer = append(buffer, packet.Pitch)
	_, err := conn.Write(buffer)
	return err
}

//
// Set Position
//

type DownstreamSetPosition struct {
	DownstreamPacket
	PlayerId int8
	X        FPShort
	Y        FPShort
	Z        FPShort
	Yaw      uint8
	Pitch    uint8
}

func NewDownstreamSetPosition(playerId int8, x FPShort, y FPShort, z FPShort, yaw uint8, pitch uint8) DownstreamSetPosition {
	return DownstreamSetPosition{
		DownstreamPacket: DownstreamPacket{
			Id: DOWNSTREAM_SET_POSITION,
		},
		PlayerId: playerId,
		X:        x,
		Y:        y,
		Z:        z,
		Yaw:      yaw,
		Pitch:    pitch,
	}
}

func (packet DownstreamSetPosition) Write(conn net.Conn) error {
	buffer := []byte{byte(packet.Id)}
	buffer = append(buffer, byte(packet.PlayerId))
	buffer = append(buffer, writeFPShort(packet.X)...)
	buffer = append(buffer, writeFPShort(packet.Y)...)
	buffer = append(buffer, writeFPShort(packet.Z)...)
	buffer = append(buffer, packet.Yaw)
	buffer = append(buffer, packet.Pitch)
	_, err := conn.Write(buffer)
	return err
}

//
// Despawn Player
//

type DownstreamDespawnPlayer struct {
	DownstreamPacket
	PlayerId int8
}

func NewDownstreamDespawnPlayer(playerId int8) DownstreamDespawnPlayer {
	return DownstreamDespawnPlayer{
		DownstreamPacket: DownstreamPacket{
			Id: DOWNSTREAM_DESPAWN_PLAYER,
		},
		PlayerId: playerId,
	}
}

func (packet DownstreamDespawnPlayer) Write(conn net.Conn) error {
	buffer := []byte{byte(packet.Id)}
	buffer = append(buffer, byte(packet.PlayerId))
	_, err := conn.Write(buffer)
	return err
}

//
// Send Message
//

type DownstreamMessage struct {
	DownstreamPacket
	PlayerId int8
	Message  string
}

func NewDownstreamMessage(senderId int8, message string) DownstreamMessage {
	if len(message) > 64 {
		message = message[:63]
	}

	runes := []rune(message)
	if runes[len(runes)-1] == '&' {
		message = string(runes[:len(runes)-1])
	}

	return DownstreamMessage{
		DownstreamPacket: DownstreamPacket{
			Id: DOWNSTREAM_MESSAGE,
		},
		PlayerId: senderId,
		Message:  message,
	}
}

func (packet DownstreamMessage) Write(conn net.Conn) error {
	buffer := []byte{byte(packet.Id)}
	buffer = append(buffer, byte(packet.PlayerId))
	buffer = append(buffer, writeString(packet.Message)...)
	_, err := conn.Write(buffer)
	return err
}

//
// Disconnect Client
//

type DownstreamDisconnectPlayer struct {
	DownstreamPacket
	Reason string
}

func NewDownstreamDisconnectPlayer(reason string) DownstreamDisconnectPlayer {
	return DownstreamDisconnectPlayer{
		DownstreamPacket: DownstreamPacket{
			Id: DOWNSTREAM_DISCONNECT_PLAYER,
		},
		Reason: reason,
	}
}

func (packet DownstreamDisconnectPlayer) Write(conn net.Conn) error {
	buffer := []byte{byte(packet.Id)}
	buffer = append(buffer, writeString(packet.Reason)...)
	_, err := conn.Write(buffer)
	return err
}

//
// Update Player Mode
//

type DownstreamUpdatePlayerMode struct {
	DownstreamPacket
	PlayerMode PlayerMode
}

func NewDownstreamUpdatePlayerMode(mode PlayerMode) DownstreamUpdatePlayerMode {
	return DownstreamUpdatePlayerMode{
		DownstreamPacket: DownstreamPacket{
			Id: DOWNSTREAM_UPDATE_PLAYER_MODE,
		},
		PlayerMode: mode,
	}
}

func (packet DownstreamUpdatePlayerMode) Write(conn net.Conn) error {
	buffer := []byte{byte(packet.Id)}
	buffer = append(buffer, byte(packet.PlayerMode))
	_, err := conn.Write(buffer)
	return err
}
