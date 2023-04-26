package classic

import (
	. "classicserver/classic/constants"
	"classicserver/classic/packets"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"
)

type ClassicServer struct {
	Listener        net.Listener
	Players         map[int8]*Player
	World           *World
	PlayerIdChannel chan int8
	Settings        *Settings
	Salt            string
	OPs             []string
	Bans            []string
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func generateSalt() string {
	var sb strings.Builder

	for i := 0; i < 16; i++ {
		sb.WriteByte(byte(letters[rand.Intn(len(letters))]))
	}

	return sb.String()
}

func NewClassicServer() (*ClassicServer, error) {
	settings, err := LoadSettings()

	log.Println("Server Started on Port", settings.Port)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", settings.IP, settings.Port))
	if err != nil {
		return nil, err
	}

	log.Println("Name:", settings.Name)
	log.Println("MOTD:", settings.MOTD)
	log.Println("Online:", settings.Online)

	players := make(map[int8]*Player)
	salt := generateSalt()
	log.Println("Salt:", salt)

	world, err := LoadWorld(settings.WorldX, settings.WorldY, settings.WorldZ)
	if err != nil {
		return nil, err
	}

	ops, err := LoadOPs()
	if err != nil {
		return nil, err
	}

	bans, err := LoadBans()
	if err != nil {
		return nil, err
	}

	log.Printf("Allocating for %d players\n", settings.PlayerCount)
	playerIdChannel := make(chan int8, settings.PlayerCount)
	for i := uint8(0); i < settings.PlayerCount && i < 128; i++ {
		playerIdChannel <- int8(i)
	}

	server := ClassicServer{
		Listener:        listener,
		Players:         players,
		World:           world,
		PlayerIdChannel: playerIdChannel,
		Settings:        settings,
		Salt:            salt,
		OPs:             ops,
		Bans:            bans,
	}

	return &server, nil
}

func (server *ClassicServer) GetPlayer(playerId int8) *Player {
	return server.Players[playerId]
}

func (server *ClassicServer) GetPlayerFromName(name string) *Player {
	for _, player := range server.Players {
		if player.Username == name {
			return player
		}
	}

	return nil
}

func (server *ClassicServer) ConnectPlayer(player *Player, verification string) error {
	server.Players[player.Id] = player

	deny := func(reason string) error {
		player.Write(packets.NewDownstreamDisconnectPlayer(reason))
		server.DisconnectPlayer(player, true)
		return errors.New(reason)
	}

	if server.Settings.Online {
		bytes := md5.Sum([]byte(server.Salt + player.Username))
		hash := hex.EncodeToString(bytes[:])
		if strings.ToLower(verification) != strings.ToLower(hash) {
			return deny("Invalid verification provided")
		}
	} else {
		if server.Settings.Password != "" && verification != server.Settings.Password {
			return deny("Incorrect password")
		}
	}

	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_]{3,16}$`)
	if !validUsername.MatchString(player.Username) {
		return deny("Invalid username provided (Letters, numbers and _)")
	}

	for _, other := range server.Players {
		if other.Username == player.Username && other.Id != player.Id {
			return deny("Username in use")
		}
	}

	if slices.Contains(server.Bans, player.Username) {
		return deny("You have been banned")
	}

	mode := MODE_NORMAL
	if slices.Contains(server.OPs, player.Username) {
		mode = MODE_OP
	}
	player.Mode = mode

	serverIdentificationPacket := packets.NewDownstreamServerIdentification(
		server.Settings.Name,
		server.Settings.MOTD,
		mode,
	)
	if err := player.Write(serverIdentificationPacket); err != nil {
		return err
	}

	if err := server.World.SendWorld(player); err != nil {
		return err
	}

	joinMsg := fmt.Sprintf("%s has joined", player.Username)
	log.Println(joinMsg)
	server.BroadcastMessage(-1, joinMsg)

	spawnPlayerPacket := packets.NewDownstreamSpawnPlayer(
		player.Id,
		player.Username,
		FPShort(server.World.SpawnX),
		FPShort(server.World.SpawnY),
		FPShort(server.World.SpawnZ),
		server.World.SpawnYaw,
		server.World.SpawnPitch,
	)
	for _, other := range server.Players {
		if other.Id == player.Id {
			other.Write(packets.NewDownstreamSpawnPlayer(
				-1,
				player.Username,
				FPShort(server.World.SpawnX),
				FPShort(server.World.SpawnY),
				FPShort(server.World.SpawnZ),
				server.World.SpawnYaw,
				server.World.SpawnPitch,
			))
		} else {
			player.Write(
				packets.NewDownstreamSpawnPlayer(
					other.Id,
					other.Username,
					other.X,
					other.Y,
					other.Z,
					other.Yaw,
					other.Pitch,
				),
			)
			other.Write(spawnPlayerPacket)
		}
	}

	return nil
}

func (server *ClassicServer) SetPosition(player *Player, x FPShort, y FPShort, z FPShort, yaw uint8, pitch uint8) {
	setPositionPacket := packets.NewDownstreamSetPosition(player.Id, x, y, z, yaw, pitch)
	for _, other := range server.Players {
		if other.Id != player.Id {
			other.Write(setPositionPacket)
		}
	}
}

func (server *ClassicServer) SetBlock(x int16, y int16, z int16, block Block) {
	server.World.SetBlock(x, y, z, block)

	setBlockPacket := packets.NewDownstreamSetBlock(x, y, z, block)
	for _, player := range server.Players {
		player.Write(setBlockPacket)
	}

	server.World.UpdateBlock(server, x, y, z)
}

func (server *ClassicServer) Ban(username string, reason string) {
	if !slices.Contains(server.Bans, username) {
		server.Bans = append(server.Bans, username)
		SaveBans(server.Bans)
	}

	player := server.GetPlayerFromName(username)
	if player != nil {
		player.Kick(reason)
	}
}

func (server *ClassicServer) Unban(username string) {
	if slices.Contains(server.Bans, username) {
		n := 0
		for _, ban := range server.Bans {
			if ban != username {
				server.Bans[n] = ban
				n++
			}
		}
		server.Bans = server.Bans[:n]
		SaveBans(server.Bans)
	}
}

func (server *ClassicServer) AddOP(username string) {
	if !slices.Contains(server.OPs, username) {
		server.OPs = append(server.OPs, username)
		SaveOPs(server.OPs)
	}

	player := server.GetPlayerFromName(username)
	if player != nil {
		player.SetMode(MODE_OP)
	}
}

func (server *ClassicServer) RemoveOP(username string) {
	if slices.Contains(server.OPs, username) {
		n := 0
		for _, op := range server.OPs {
			if op != username {
				server.OPs[n] = op
				n++
			}
		}
		server.OPs = server.OPs[:n]
		SaveOPs(server.OPs)
	}

	player := server.GetPlayerFromName(username)
	if player != nil {
		player.SetMode(MODE_NORMAL)
	}
}

func (server *ClassicServer) BroadcastMessage(senderId int8, message string) {
	despawnPacket := packets.NewDownstreamMessage(senderId, message)
	for _, other := range server.Players {
		if other.Id == senderId {
			other.Write(packets.NewDownstreamMessage(-1, message))
		} else {
			other.Write(despawnPacket)
		}
	}
}

func (server *ClassicServer) DisconnectPlayer(player *Player, silent bool) {
	if player.Conn == nil {
		return
	}

	server.PlayerIdChannel <- player.Id

	err := player.Conn.Close()
	if err != nil {
		log.Println(err)
	}

	player.Conn = nil

	if !silent {
		leaveMsg := fmt.Sprintf("%s has disconnected", player.Username)
		log.Println(leaveMsg)

		despawnPacket := packets.NewDownstreamDespawnPlayer(player.Id)
		for _, other := range player.Server.Players {
			if other.Id != player.Id {
				other.Write(despawnPacket)
				other.SendMessage(leaveMsg)
			}
		}
	}

	delete(server.Players, player.Id)
}
