package packets

import (
	"bufio"
	. "classicserver/classic/constants"
	"errors"
)

//
// Player Identification
//

type UpstreamPlayerIdentification struct {
	UpstreamPacket
	Version      uint8
	Username     string
	Verification string
}

func ReadPacketID(reader *bufio.Reader) (UpstreamPacketID, error) {
	id, err := reader.ReadByte()
	if err != nil {
		return 0xFF, err
	}

	if id == byte(UPSTREAM_PLAYER_IDENTIFICATION) || id == byte(UPSTREAM_MESSAGE) || id == byte(UPSTREAM_SET_BLOCK) || id == byte(UPSTREAM_SET_POSITION) {
		return UpstreamPacketID(id), nil
	} else {
		return 0xFF, errors.New("Unknown Packet ID")
	}
}

func ReadUpstreamPlayerIdentification(reader *bufio.Reader) (*UpstreamPlayerIdentification, error) {
	version, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	username, err := readString(reader)
	if err != nil {
		return nil, err
	}

	verification, err := readString(reader)
	if err != nil {
		return nil, err
	}

	// Read empty byte
	_, err = reader.ReadByte()
	if err != nil {
		return nil, err
	}

	packet := UpstreamPlayerIdentification{
		UpstreamPacket: UpstreamPacket{
			Id: UPSTREAM_PLAYER_IDENTIFICATION,
		},
		Version:      version,
		Username:     username,
		Verification: verification,
	}

	return &packet, nil
}

//
// Player Set Block
//

type UpstreamSetBlock struct {
	UpstreamPacket
	X     int16
	Y     int16
	Z     int16
	Mode  uint8
	Block uint8
}

func ReadUpstreamSetBlock(reader *bufio.Reader) (*UpstreamSetBlock, error) {
	x, err := readShort(reader)
	if err != nil {
		return nil, err
	}

	y, err := readShort(reader)
	if err != nil {
		return nil, err
	}

	z, err := readShort(reader)
	if err != nil {
		return nil, err
	}

	mode, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	block, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	packet := UpstreamSetBlock{
		UpstreamPacket: UpstreamPacket{
			Id: UPSTREAM_SET_BLOCK,
		},
		X:     x,
		Y:     y,
		Z:     z,
		Block: block,
		Mode:  mode,
	}

	return &packet, nil
}

//
// Set Position
//

type UpstreamSetPosition struct {
	UpstreamPacket
	PlayerId uint8
	X        FPShort
	Y        FPShort
	Z        FPShort
	Yaw      uint8
	Pitch    uint8
}

func ReadUpstreamSetPosition(reader *bufio.Reader) (*UpstreamSetPosition, error) {
	playerId, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	x, err := readFPShort(reader)
	if err != nil {
		return nil, err
	}

	y, err := readFPShort(reader)
	if err != nil {
		return nil, err
	}

	z, err := readFPShort(reader)
	if err != nil {
		return nil, err
	}

	yaw, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	pitch, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	packet := UpstreamSetPosition{
		UpstreamPacket: UpstreamPacket{
			Id: UPSTREAM_SET_POSITION,
		},
		PlayerId: playerId,
		X:        x,
		Y:        y,
		Z:        z,
		Yaw:      yaw,
		Pitch:    pitch,
	}

	return &packet, nil
}

//
// Message
//

type UpstreamMessage struct {
	UpstreamPacket
	PlayerId int8
	Message  string
}

func ReadUpstreamMessage(reader *bufio.Reader) (*UpstreamMessage, error) {
	playerId, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	message, err := readString(reader)
	if err != nil {
		return nil, err
	}

	packet := UpstreamMessage{
		UpstreamPacket: UpstreamPacket{
			Id: UPSTREAM_MESSAGE,
		},
		PlayerId: int8(playerId),
		Message:  message,
	}

	return &packet, nil
}
