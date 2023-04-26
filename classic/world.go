package classic

import (
	"bytes"
	. "classicserver/classic/constants"
	"classicserver/classic/packets"
	"compress/gzip"
	"encoding/binary"
	"io"
	"log"
	"math"
	"os"
	"time"
)

const WORLD_FILENAME string = "world.gw"

type World struct {
	SizeX      int16
	SizeY      int16
	SizeZ      int16
	SpawnX     float64
	SpawnY     float64
	SpawnZ     float64
	SpawnYaw   uint8
	SpawnPitch uint8
	Blocks     [][][]Block // [Y][Z][X]
}

func NewWorld(sizeX int16, sizeY int16, sizeZ int16) *World {

	blocks := make([][][]Block, sizeY)
	for y := int16(0); y < sizeY; y++ {
		blocks[y] = make([][]Block, sizeZ)
		for z := int16(0); z < sizeZ; z++ {
			blocks[y][z] = make([]Block, sizeX)
			for x := int16(0); x < sizeX; x++ {
				blocks[y][z][x] = BLOCK_AIR
			}
		}
	}

	world := &World{
		SizeX:      sizeX,
		SizeY:      sizeY,
		SizeZ:      sizeZ,
		SpawnX:     0.0,
		SpawnY:     0.0,
		SpawnZ:     0.0,
		SpawnYaw:   0,
		SpawnPitch: 0,
		Blocks:     blocks,
	}

	GenerateWorld(world)

	return world
}

func GenerateWorld(world *World) {
	for y := int16(0); y <= (world.SizeY/2)+1; y++ {
		top := y == (world.SizeY/2)+1
		for z := int16(0); z < world.SizeZ; z++ {
			for x := int16(0); x < world.SizeX; x++ {
				if top {
					world.Blocks[y][z][x] = BLOCK_GRASS_BLOCK
				} else {
					world.Blocks[y][z][x] = BLOCK_DIRT
				}
			}
		}
	}

	world.SpawnX = math.Floor(float64(world.SizeX)/2.0) + 0.5
	world.SpawnZ = math.Floor(float64(world.SizeZ)/2.0) + 0.5
	for y := world.SizeY - 1; y >= 0; y-- {
		if world.GetBlock(int16(world.SpawnX), y, int16(world.SpawnZ)) != BLOCK_AIR {
			world.SpawnY = float64(y) + 2
			return
		}
	}
}

func LoadWorld(sizeX int16, sizeY int16, sizeZ int16) (*World, error) {
	world := NewWorld(sizeX, sizeY, sizeZ)
	if _, err := os.Stat(WORLD_FILENAME); err != nil {
		if os.IsNotExist(err) {
			log.Println("Generating new world...")
			GenerateWorld(world)

			if err := SaveWorld(world); err != nil {
				return nil, err
			}

			return world, nil
		} else {
			return nil, err
		}
	}
	log.Println("Loading world...")

	contents, err := os.ReadFile(WORLD_FILENAME)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	b.Write(contents)
	gz, err := gzip.NewReader(&b)
	defer gz.Close()

	if err != nil {
		return nil, err
	}

	readShort := func() int16 {
		bytes := make([]uint8, 2)
		gz.Read(bytes)
		return int16(binary.BigEndian.Uint16(bytes))
	}

	world.SizeX = readShort()
	world.SizeY = readShort()
	world.SizeZ = readShort()

	readFloat := func() float64 {
		bytes := make([]uint8, 8)
		gz.Read(bytes)
		i64 := binary.BigEndian.Uint64(bytes)
		return math.Float64frombits(i64)
	}

	world.SpawnX = readFloat()
	world.SpawnY = readFloat()
	world.SpawnZ = readFloat()

	readByte := func() uint8 {
		bytes := make([]uint8, 1)
		gz.Read(bytes)
		return bytes[0]
	}

	world.SpawnYaw = readByte()
	world.SpawnPitch = readByte()

	for y := int16(0); y < world.SizeY; y++ {
		for z := int16(0); z < world.SizeZ; z++ {
			for x := int16(0); x < world.SizeX; x++ {
				b := make([]uint8, 1)
				if _, err := gz.Read(b); err != nil {
					if x != world.SizeX-1 || y != world.SizeY-1 || z != world.SizeZ-1 || err != io.EOF {
						return nil, err
					}
				}
				world.SetBlock(x, y, z, b[0])
			}
		}
	}

	return world, nil
}

func SaveWorld(world *World) error {
	log.Println("-- Saving world...")
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	writeShort := func(val int16) {
		bytes := make([]uint8, 2)
		binary.BigEndian.PutUint16(bytes, uint16(val))
		gz.Write(bytes)
	}

	writeShort(world.SizeX)
	writeShort(world.SizeY)
	writeShort(world.SizeZ)

	writeFloat := func(val float64) {
		i64 := math.Float64bits(val)
		bytes := make([]uint8, 8)
		binary.BigEndian.PutUint64(bytes, i64)
		gz.Write(bytes)
	}

	writeFloat(world.SpawnX)
	writeFloat(world.SpawnY)
	writeFloat(world.SpawnZ)

	gz.Write([]byte{world.SpawnYaw})
	gz.Write([]byte{world.SpawnPitch})

	if err := gzipWorld(world, gz); err != nil {
		return err
	}

	if err := gz.Close(); err != nil {
		return err
	}

	data := b.Bytes()

	if err := os.WriteFile(WORLD_FILENAME, data, 0644); err != nil {
		return err
	}

	log.Println("-- World Saved Successfully")

	return nil
}

func gzipWorld(world *World, gz *gzip.Writer) error {
	for y := int16(0); y < world.SizeY; y++ {
		for z := int16(0); z < world.SizeZ; z++ {
			for x := int16(0); x < world.SizeX; x++ {
				if _, err := gz.Write([]uint8{world.Blocks[y][z][x]}); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (world *World) SendWorld(player *Player) error {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	levelSize := make([]uint8, 4)
	sizeCubed := uint32(world.SizeX) * uint32(world.SizeY) * uint32(world.SizeZ)
	binary.BigEndian.PutUint32(levelSize, sizeCubed)
	if _, err := gz.Write(levelSize); err != nil {
		return err
	}

	err := gzipWorld(world, gz)
	if err != nil {
		log.Println(err)
		return err
	}

	if err := gz.Close(); err != nil {
		log.Println(err)
		return err
	}

	levelInitPacket := packets.NewDownstreamLevelInit()
	if err := player.Write(levelInitPacket); err != nil {
		return err
	}

	data := b.Bytes()
	chunks := [][]uint8{}
	for i := 0; i < len(data); i += 1024 {
		end := i + 1024
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, data[i:end])
	}

	for i, chunk := range chunks {
		levelChunkPacket := packets.NewDownstreamLevelChunk(chunk, uint8((255*i)/len(chunks)))
		if err := player.Write(levelChunkPacket); err != nil {
			return err
		}
	}

	levelFinalizePacket := packets.NewDownstreamLevelFinalize(world.SizeX, world.SizeY, world.SizeZ)
	if err := player.Write(levelFinalizePacket); err != nil {
		return err
	}

	return nil
}

func (world *World) SetBlock(x int16, y int16, z int16, block Block) {
	world.Blocks[y][z][x] = block
}

func (world *World) GetBlock(x int16, y int16, z int16) Block {
	return world.Blocks[y][z][x]
}

func (world *World) ValidBlock(x int16, y int16, z int16) bool {
	return x >= 0 && y >= 0 && z >= 0 && x < world.SizeX && y < world.SizeY && z < world.SizeZ
}

func tickLava(server *ClassicServer, world *World, x int16, y int16, z int16) {
	if world.ValidBlock(x, y-1, z) {
		block := world.GetBlock(x, y-1, z)
		if block == BLOCK_AIR {
			server.SetBlock(x, y-1, z, BLOCK_LAVA_FLOWING)
		} else if block == BLOCK_WATER_FLOWING || block == BLOCK_WATER_STATIONARY {
			server.SetBlock(x, y-1, z, BLOCK_STONE)
		}
	}

	if world.ValidBlock(x+1, y, z) {
		block := world.GetBlock(x+1, y, z)
		if block == BLOCK_AIR {
			server.SetBlock(x+1, y, z, BLOCK_LAVA_FLOWING)
		} else if block == BLOCK_WATER_FLOWING || block == BLOCK_WATER_STATIONARY {
			server.SetBlock(x+1, y, z, BLOCK_STONE)
		}
	}

	if world.ValidBlock(x-1, y, z) {
		block := world.GetBlock(x-1, y, z)
		if block == BLOCK_AIR {
			server.SetBlock(x-1, y, z, BLOCK_LAVA_FLOWING)
		} else if block == BLOCK_WATER_FLOWING || block == BLOCK_WATER_STATIONARY {
			server.SetBlock(x-1, y, z, BLOCK_STONE)
		}
	}

	if world.ValidBlock(x, y, z+1) {
		block := world.GetBlock(x, y, z+1)
		if block == BLOCK_AIR {
			server.SetBlock(x, y, z+1, BLOCK_LAVA_FLOWING)
		} else if block == BLOCK_WATER_FLOWING || block == BLOCK_WATER_STATIONARY {
			server.SetBlock(x, y, z+1, BLOCK_STONE)
		}
	}

	if world.ValidBlock(x, y, z-1) {
		block := world.GetBlock(x, y, z-1)
		if block == BLOCK_AIR {
			server.SetBlock(x, y, z-1, BLOCK_LAVA_FLOWING)
		} else if block == BLOCK_WATER_FLOWING || block == BLOCK_WATER_STATIONARY {
			server.SetBlock(x, y, z-1, BLOCK_STONE)
		}
	}

}

func (world *World) UpdateLava(server *ClassicServer) {
	for y := int16(0); y < world.SizeY; y++ {
		for z := int16(0); z < world.SizeZ; z++ {
			for x := int16(0); x < world.SizeX; x++ {
				block := world.Blocks[y][z][x]
				if block == BLOCK_LAVA_FLOWING {
					defer tickLava(server, world, x, y, z)
				}
			}
		}
	}
}

func tickWater(server *ClassicServer, world *World, x int16, y int16, z int16) {
	if world.ValidBlock(x, y-1, z) {
		block := world.GetBlock(x, y-1, z)
		if block == BLOCK_AIR {
			server.SetBlock(x, y-1, z, BLOCK_WATER_FLOWING)
		} else if block == BLOCK_LAVA_FLOWING || block == BLOCK_LAVA_STATIONARY {
			server.SetBlock(x, y-1, z, BLOCK_STONE)
		}
	}

	if world.ValidBlock(x+1, y, z) {
		block := world.GetBlock(x+1, y, z)
		if block == BLOCK_AIR {
			server.SetBlock(x+1, y, z, BLOCK_WATER_FLOWING)
		} else if block == BLOCK_LAVA_FLOWING || block == BLOCK_LAVA_STATIONARY {
			server.SetBlock(x+1, y, z, BLOCK_STONE)
		}
	}

	if world.ValidBlock(x-1, y, z) {
		block := world.GetBlock(x-1, y, z)
		if block == BLOCK_AIR {
			server.SetBlock(x-1, y, z, BLOCK_WATER_FLOWING)
		} else if block == BLOCK_LAVA_FLOWING || block == BLOCK_LAVA_STATIONARY {
			server.SetBlock(x-1, y, z, BLOCK_STONE)
		}
	}

	if world.ValidBlock(x, y, z+1) {
		block := world.GetBlock(x, y, z+1)
		if block == BLOCK_AIR {
			server.SetBlock(x, y, z+1, BLOCK_WATER_FLOWING)
		} else if block == BLOCK_LAVA_FLOWING || block == BLOCK_LAVA_STATIONARY {
			server.SetBlock(x, y, z+1, BLOCK_STONE)
		}
	}

	if world.ValidBlock(x, y, z-1) {
		block := world.GetBlock(x, y, z-1)
		if block == BLOCK_AIR {
			server.SetBlock(x, y, z-1, BLOCK_WATER_FLOWING)
		} else if block == BLOCK_LAVA_FLOWING || block == BLOCK_LAVA_STATIONARY {
			server.SetBlock(x, y, z-1, BLOCK_STONE)
		}
	}

}

func (world *World) UpdateWater(server *ClassicServer) {
	defer world.UpdateSponge(server)

	for y := int16(0); y < world.SizeY; y++ {
		for z := int16(0); z < world.SizeZ; z++ {
			for x := int16(0); x < world.SizeX; x++ {
				block := world.Blocks[y][z][x]
				if block == BLOCK_WATER_FLOWING {
					defer tickWater(server, world, x, y, z)
				}
			}
		}
	}
}

func tickSponge(server *ClassicServer, world *World, x int16, y int16, z int16) {
	for dy := int16(-2); dy <= 2; dy++ {
		for dz := int16(-2); dz <= 2; dz++ {
			for dx := int16(-2); dx <= 2; dx++ {
				ny := y + dy
				nz := z + dz
				nx := x + dx
				if world.ValidBlock(nx, ny, nz) {
					block := world.GetBlock(nx, ny, nz)
					if block == BLOCK_WATER_FLOWING || block == BLOCK_WATER_STATIONARY {
						server.SetBlock(nx, ny, nz, BLOCK_AIR)
					}
				}
			}
		}
	}
}

func (world *World) UpdateSponge(server *ClassicServer) {
	for y := int16(0); y < world.SizeY; y++ {
		for z := int16(0); z < world.SizeZ; z++ {
			for x := int16(0); x < world.SizeX; x++ {
				block := world.Blocks[y][z][x]
				if block == BLOCK_SPONGE {
					tickSponge(server, world, x, y, z)
				}
			}
		}
	}
}

func canSandPass(block Block) bool {
	return block == BLOCK_AIR ||
		block == BLOCK_WATER_FLOWING ||
		block == BLOCK_WATER_STATIONARY ||
		block == BLOCK_LAVA_FLOWING ||
		block == BLOCK_LAVA_STATIONARY
}

func (world *World) FallSand(server *ClassicServer, x int16, y int16, z int16) {
	blockType := world.GetBlock(x, y, z)
	if world.ValidBlock(x, y-1, z) && canSandPass(world.GetBlock(x, y-1, z)) {
		server.SetBlock(x, y, z, BLOCK_AIR)

		for ny := y - 1; ny >= 0; ny-- {
			under := world.GetBlock(x, ny, z)
			if !canSandPass(under) && world.ValidBlock(x, ny+1, z) {
				server.SetBlock(x, ny+1, z, blockType)
				return
			}
		}

		server.SetBlock(x, 0, z, blockType)
	}
}

func (world *World) UpdateBlock(server *ClassicServer, x int16, y int16, z int16) {
	block := world.GetBlock(x, y, z)
	if block == BLOCK_SAND || block == BLOCK_GRAVEL {
		world.FallSand(server, x, y, z)
	} else if block == BLOCK_SPONGE {
		tickSponge(server, world, x, y, z)
	} else if block == BLOCK_WATER_FLOWING || block == BLOCK_WATER_STATIONARY {
		world.UpdateSponge(server)
	}

	if world.ValidBlock(x, y+1, z) {
		block := world.GetBlock(x, y+1, z)
		if block == BLOCK_SAND || block == BLOCK_GRAVEL {
			world.FallSand(server, x, y+1, z)
		}
	}
}

func NewWorldSaveTicker() *time.Ticker {
	return time.NewTicker(time.Second * 60)
}

func NewWorldWaterTicker() *time.Ticker {
	return time.NewTicker(200 * time.Millisecond)
}

func NewWorldLavaTicker() *time.Ticker {
	return time.NewTicker(1500 * time.Millisecond)
}
