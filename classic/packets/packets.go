package packets

import (
	"bufio"
	. "classicserver/classic/constants"
	"encoding/binary"
	"math"
	"net"
	"strings"
)

// Upstream Packets
// Client -> Server

type UpstreamPacketID uint8

type UpstreamPacket struct {
	Id UpstreamPacketID
}

const (
	UPSTREAM_PLAYER_IDENTIFICATION UpstreamPacketID = 0x00
	UPSTREAM_SET_BLOCK             UpstreamPacketID = 0x05
	UPSTREAM_SET_POSITION          UpstreamPacketID = 0x08
	UPSTREAM_MESSAGE               UpstreamPacketID = 0x0D
)

// Downstream Packets
// Server -> Client

type DownstreamPacketID uint8

type DownstreamPacketInterface interface {
	Write(conn net.Conn) error
}

type DownstreamPacket struct {
	DownstreamPacketInterface
	Id DownstreamPacketID
}

const (
	DOWNSTREAM_SERVER_IDENTIFICATION DownstreamPacketID = 0x00
	DOWNSTREAM_PING                  DownstreamPacketID = 0x01
	DOWNSTREAM_LEVEL_INIT            DownstreamPacketID = 0x02
	DOWNSTREAM_LEVEL_CHUNK           DownstreamPacketID = 0x03
	DOWNSTREAM_LEVEL_FINALIZE        DownstreamPacketID = 0x04
	DOWNSTREAM_SET_BLOCK             DownstreamPacketID = 0x06
	DOWNSTREAM_SPAWN_PLAYER          DownstreamPacketID = 0x07
	DOWNSTREAM_SET_POSITION          DownstreamPacketID = 0x08
	DOWNSTREAM_DESPAWN_PLAYER        DownstreamPacketID = 0x0C
	DOWNSTREAM_MESSAGE               DownstreamPacketID = 0x0D
	DOWNSTREAM_DISCONNECT_PLAYER     DownstreamPacketID = 0x0E
	DOWNSTREAM_UPDATE_PLAYER_MODE    DownstreamPacketID = 0x0F
)

// Util functions

func readString(reader *bufio.Reader) (string, error) {
	buffer := [64]uint8{}
	for i := 0; i < 64; i++ {
		read, err := reader.ReadByte()
		if err != nil {
			return "", err
		}
		buffer[i] = read
	}

	return strings.TrimRight(string(buffer[:]), " "), nil
}

func writeString(value string) []uint8 {
	buffer := [64]uint8{}
	for i := 0; i < 64; i++ {
		if i >= len(value) {
			buffer[i] = ' '
		} else {
			buffer[i] = value[i]
		}
	}

	return buffer[:]
}

func readShort(reader *bufio.Reader) (int16, error) {
	high, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	low, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	return int16(binary.BigEndian.Uint16([]uint8{high, low})), nil
}

func writeShort(value int16) []uint8 {
	bytes := make([]uint8, 2)
	binary.BigEndian.PutUint16(bytes, uint16(value))
	return bytes
}

func readFPShort(reader *bufio.Reader) (FPShort, error) {
	high, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	low, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	whole := int(high)<<3 + int((low&0b11100000)>>5)
	decimal := (low & 0b00011111)
	f := float64(whole)
	for i := 0; i < 5; i++ {
		set := (decimal & (0b1 << (4 - i))) >> (4 - i)
		f += float64(set) * (float64(1) / math.Pow(2, float64(i+1)))
	}

	return FPShort(f), nil
}

func writeFPShort(value FPShort) []uint8 {
	whole := math.Floor(float64(value))
	decimal := float64(value) - whole
	decimalBytes := byte(0)

	for i := 0; i < 5; i++ {
		decimalBytes <<= 1
		decimalUnit := float64(1) / math.Pow(2, float64(i+1))
		if decimal >= decimalUnit {
			decimal -= decimalUnit
			decimalBytes += 1
		}
	}

	short := int16(whole)
	short <<= 5
	short += int16(decimalBytes)
	return writeShort(short)
}
