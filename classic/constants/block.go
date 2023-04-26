package constants

type Block = uint8

const (
	BUILD_DESTROY = 0
	BUILD_PLACE   = 1
)

var names map[Block]string

const (
	BLOCK_AIR Block = iota
	BLOCK_STONE
	BLOCK_GRASS_BLOCK
	BLOCK_DIRT
	BLOCK_COBBLESTONE
	BLOCK_PLANKS
	BLOCK_SAPLING
	BLOCK_BEDROCK
	BLOCK_WATER_FLOWING
	BLOCK_WATER_STATIONARY
	BLOCK_LAVA_FLOWING
	BLOCK_LAVA_STATIONARY
	BLOCK_SAND
	BLOCK_GRAVEL
	BLOCK_GOLD_ORE
	BLOCK_IRON_ORE
	BLOCK_COAL_ORE
	BLOCK_WOOD
	BLOCK_LEAVES
	BLOCK_SPONGE
	BLOCK_GLASS
	BLOCK_CLOTH_RED
	BLOCK_CLOTH_ORANGE
	BLOCK_CLOTH_YELLOW
	BLOCK_CLOTH_CHARTREUSE
	BLOCK_CLOTH_GREEN
	BLOCK_CLOTH_SPRING_GREEN
	BLOCK_CLOTH_CYAN
	BLOCK_CLOTH_CAPRI
	BLOCK_CLOTH_ULTRAMARINE
	BLOCK_CLOTH_VIOLET
	BLOCK_CLOTH_PURPLE
	BLOCK_CLOTH_MAGENTA
	BLOCK_CLOTH_ROSE
	BLOCK_CLOTH_DARK_GRAY
	BLOCK_CLOTH_LIGHT_GRAY
	BLOCK_CLOTH_WHITE
	BLOCK_FLOWER
	BLOCK_ROSE
	BLOCK_BROWN_MUSHROOM
	BLOCK_RED_MUSHROOM
	BLOCK_BLOCK_OF_GOLD
	BLOCK_BLOCK_OF_IRON
	BLOCK_SLAB_DOUBLE
	BLOCK_SLAB
	BLOCK_BRICKS
	BLOCK_TNT
	BLOCK_BOOKSHELF
	BLOCK_MOSSY_COBBLESTONE
	BLOCK_OBSIDIAN
)

func init() {
	names = make(map[Block]string, BLOCK_OBSIDIAN+1)
	names[BLOCK_AIR] = "Air"
	names[BLOCK_STONE] = "Stone"
	names[BLOCK_GRASS_BLOCK] = "Grass Block"
	names[BLOCK_DIRT] = "Dirt"
	names[BLOCK_COBBLESTONE] = "Cobblestone"
	names[BLOCK_PLANKS] = "Planks"
	names[BLOCK_SAPLING] = "Sapling"
	names[BLOCK_BEDROCK] = "Bedrock"
	names[BLOCK_WATER_FLOWING] = "Flowing Water"
	names[BLOCK_WATER_STATIONARY] = "Stationary Water"
	names[BLOCK_LAVA_FLOWING] = "Flowing Lava"
	names[BLOCK_LAVA_STATIONARY] = "Stationary Lava"
	names[BLOCK_SAND] = "Sand"
	names[BLOCK_GRAVEL] = "Gravel"
	names[BLOCK_GOLD_ORE] = "Gold Ore"
	names[BLOCK_IRON_ORE] = "Iron Ore"
	names[BLOCK_COAL_ORE] = "Coal Ore"
	names[BLOCK_WOOD] = "Wood"
	names[BLOCK_LEAVES] = "Leaves"
	names[BLOCK_SPONGE] = "Sponge"
	names[BLOCK_GLASS] = "Glass"
	names[BLOCK_CLOTH_RED] = "Cloth (Red)"
	names[BLOCK_CLOTH_ORANGE] = "Cloth (Orange)"
	names[BLOCK_CLOTH_YELLOW] = "Cloth (Yellow)"
	names[BLOCK_CLOTH_CHARTREUSE] = "Cloth (Chartreuse)"
	names[BLOCK_CLOTH_GREEN] = "Cloth (Green)"
	names[BLOCK_CLOTH_SPRING_GREEN] = "Cloth (Spring Green)"
	names[BLOCK_CLOTH_CYAN] = "Cloth (Cyan)"
	names[BLOCK_CLOTH_CAPRI] = "Cloth (Capri)"
	names[BLOCK_CLOTH_ULTRAMARINE] = "Cloth (Ultramarine)"
	names[BLOCK_CLOTH_VIOLET] = "Cloth (Violet)"
	names[BLOCK_CLOTH_PURPLE] = "Cloth (Purple)"
	names[BLOCK_CLOTH_MAGENTA] = "Cloth (Magenta)"
	names[BLOCK_CLOTH_ROSE] = "Cloth (Rose)"
	names[BLOCK_CLOTH_DARK_GRAY] = "Cloth (Dark Gray)"
	names[BLOCK_CLOTH_LIGHT_GRAY] = "Cloth (Light Gray)"
	names[BLOCK_CLOTH_WHITE] = "Cloth (White)"
	names[BLOCK_FLOWER] = "Flower"
	names[BLOCK_ROSE] = "Rose"
	names[BLOCK_BROWN_MUSHROOM] = "Brown Mushroom"
	names[BLOCK_RED_MUSHROOM] = "Red Mushroom"
	names[BLOCK_BLOCK_OF_GOLD] = "Block of Gold"
	names[BLOCK_BLOCK_OF_IRON] = "Block of Iron"
	names[BLOCK_SLAB_DOUBLE] = "Slab Double"
	names[BLOCK_SLAB] = "Slab"
	names[BLOCK_BRICKS] = "Bricks"
	names[BLOCK_TNT] = "TNT"
	names[BLOCK_BOOKSHELF] = "Bookshelf"
	names[BLOCK_MOSSY_COBBLESTONE] = "Mossy Cobblestone"
	names[BLOCK_OBSIDIAN] = "Obsidian"
}
