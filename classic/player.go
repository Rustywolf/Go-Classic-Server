package classic

import (
	. "classicserver/classic/constants"
	"classicserver/classic/packets"
	"fmt"
	"log"
	"net"
)

type Player struct {
	Server   *ClassicServer
	Id       int8
	Conn     net.Conn
	Username string
	Mode     PlayerMode
	X        FPShort
	Y        FPShort
	Z        FPShort
	Yaw      uint8
	Pitch    uint8
}

func NewPlayer(server *ClassicServer, id int8, conn net.Conn, username string) *Player {
	return &Player{
		Server:   server,
		Id:       id,
		Conn:     conn,
		Username: username,
	}
}

func (player *Player) Write(packet packets.DownstreamPacketInterface) error {
	err := packet.Write(player.Conn)
	if err != nil {
		log.Println(err)
		player.Disconnect()
	}

	return err
}

func (player *Player) SendMessage(message string, args ...any) {
	player.Write(packets.NewDownstreamMessage(-1, fmt.Sprintf(message, args...)))
}

func (player *Player) SetMode(mode PlayerMode) {
	player.Mode = mode
	player.Write(packets.NewDownstreamUpdatePlayerMode(mode))
}

func (player *Player) Teleport(x FPShort, y FPShort, z FPShort, yaw uint8, pitch uint8) {
	player.X = x
	player.Y = y
	player.Z = z
	player.Yaw = yaw
	player.Pitch = pitch
	player.Server.SetPosition(player, x, y, z, yaw, pitch)
	player.Write(packets.NewDownstreamSetPosition(-1, x, y, z, yaw, pitch))
}

func (player *Player) Kick(reason string) {
	player.Write(packets.NewDownstreamDisconnectPlayer(reason))
	player.Disconnect()
}

func (player *Player) Disconnect() {
	player.Server.DisconnectPlayer(player, false)
}
