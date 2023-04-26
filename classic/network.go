package classic

import (
	"bufio"
	. "classicserver/classic/constants"
	"classicserver/classic/packets"
	"log"
	"net"
)

type PlayerChannels struct {
	getPlayerId chan int8
	connect     chan PlayerIdentificationChannel
	message     chan MessageChannel
	setPosition chan SetPositionChannel
	setBlock    chan SetBlockChannel
	disconnect  chan DisconnectChannel
}

type PlayerIdentificationChannel struct {
	playerId int8
	conn     net.Conn
	packet   packets.UpstreamPlayerIdentification
}

type MessageChannel struct {
	playerId int8
	packet   packets.UpstreamMessage
}

type SetPositionChannel struct {
	playerId int8
	packet   packets.UpstreamSetPosition
}

type SetBlockChannel struct {
	playerId int8
	packet   packets.UpstreamSetBlock
}

type DisconnectChannel struct {
	playerId int8
}

func (server *ClassicServer) Run() {
	log.Println("Server ready")

	defer server.Listener.Close()

	heartbeatTicker := NewHeartbeatTicker()
	defer heartbeatTicker.Stop()

	worldSaveTicker := NewWorldSaveTicker()
	defer worldSaveTicker.Stop()

	worldLavaTicker := NewWorldLavaTicker()
	defer worldLavaTicker.Stop()

	worldWaterTicker := NewWorldWaterTicker()
	defer worldWaterTicker.Stop()

	channels := PlayerChannels{
		getPlayerId: server.PlayerIdChannel,
		connect:     make(chan PlayerIdentificationChannel),
		message:     make(chan MessageChannel),
		setPosition: make(chan SetPositionChannel),
		setBlock:    make(chan SetBlockChannel),
		disconnect:  make(chan DisconnectChannel),
	}

	heartbeat := func() {
		if server.Settings.Online {
			go Heartbeat(
				server.Settings.Port,
				server.Settings.PlayerCount,
				server.Settings.Name,
				server.Settings.Public,
				server.Salt,
				len(server.Players),
			)
		}
	}

	saveWorld := func() {
		world := *server.World
		go SaveWorld(&world)
	}

	go server.listen(&channels)

	heartbeat()

	for {
		select {

		case <-heartbeatTicker.C:
			heartbeat()

		case <-worldSaveTicker.C:
			saveWorld()

		case <-worldLavaTicker.C:
			server.World.UpdateLava(server)

		case <-worldWaterTicker.C:
			server.World.UpdateWater(server)

		case inc := <-channels.connect:
			player := NewPlayer(server, inc.playerId, inc.conn, inc.packet.Username)
			server.ConnectPlayer(player, inc.packet.Verification)

		case inc := <-channels.disconnect:
			player := server.GetPlayer(inc.playerId)
			if player != nil {
				player.Disconnect()
			}

		case inc := <-channels.message:
			player := server.GetPlayer(inc.playerId)
			server.HandleChat(player, inc.packet.Message)

		case inc := <-channels.setPosition:
			player := server.GetPlayer(inc.playerId)
			player.X = inc.packet.X
			player.Y = inc.packet.Y
			player.Z = inc.packet.Z
			player.Yaw = inc.packet.Yaw
			player.Pitch = inc.packet.Pitch
			server.SetPosition(player, inc.packet.X, inc.packet.Y, inc.packet.Z, inc.packet.Yaw, inc.packet.Pitch)

		case inc := <-channels.setBlock:
			if inc.packet.Mode == BUILD_PLACE {
				server.SetBlock(inc.packet.X, inc.packet.Y, inc.packet.Z, inc.packet.Block)
			} else if inc.packet.Mode == BUILD_DESTROY {
				server.SetBlock(inc.packet.X, inc.packet.Y, inc.packet.Z, BLOCK_AIR)
			}

		}
	}
}

func (server *ClassicServer) listen(channels *PlayerChannels) {
	for {
		conn, err := server.Listener.Accept()
		if err != nil {
			log.Println(err)
		} else {
			go server.HandleConnection(conn, channels)
		}
	}
}

func (server *ClassicServer) HandleConnection(conn net.Conn, channels *PlayerChannels) {
	reader := bufio.NewReader(conn)
	packetId, err := packets.ReadPacketID(reader)
	if err != nil || packetId != packets.UPSTREAM_PLAYER_IDENTIFICATION {
		conn.Close()
		return
	}

	playerIdentificationPacket, err := packets.ReadUpstreamPlayerIdentification(reader)
	if err != nil {
		conn.Close()
		return
	}

	if len(channels.getPlayerId) == 0 {
		packets.NewDownstreamDisconnectPlayer("Server is full").Write(conn)
		conn.Close()
		return
	}

	playerId := <-channels.getPlayerId

	channels.connect <- PlayerIdentificationChannel{
		playerId: playerId,
		conn:     conn,
		packet:   *playerIdentificationPacket,
	}

	disconnect := func() {
		channels.disconnect <- DisconnectChannel{
			playerId,
		}
	}

	defer disconnect()

	for {
		packetId, err := packets.ReadPacketID(reader)
		if err != nil {
			return
		}

		switch packetId {

		case packets.UPSTREAM_MESSAGE:
			packet, err := packets.ReadUpstreamMessage(reader)
			if err != nil {
				return
			}
			channels.message <- MessageChannel{
				playerId: playerId,
				packet:   *packet,
			}

		case packets.UPSTREAM_SET_BLOCK:
			packet, err := packets.ReadUpstreamSetBlock(reader)
			if err != nil {
				return
			}
			channels.setBlock <- SetBlockChannel{
				playerId: playerId,
				packet:   *packet,
			}
			break

		case packets.UPSTREAM_SET_POSITION:
			packet, err := packets.ReadUpstreamSetPosition(reader)
			if err != nil {
				return
			}
			channels.setPosition <- SetPositionChannel{
				playerId: playerId,
				packet:   *packet,
			}

		case packets.UPSTREAM_PLAYER_IDENTIFICATION:
			disconnect()

		}
	}
}
