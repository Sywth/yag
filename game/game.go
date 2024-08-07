package game

import (
	"bytes"
	"encoding/gob"
	vl "p3/veclib"

	noise "github.com/KEINOS/go-noise"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var noiseGen, _ = noise.New(noise.Perlin, 1234)
var Constants = struct {
	TILE_SIZE             float32
	CHUNK_SIZE            float32
	MAP_SCALAR            float32
	PATH_TO_TEXTURE_ATLAS string
	TEXTURE_SIZE_PX       float32
	NOISE                 noise.Generator
}{
	TILE_SIZE:             8,
	CHUNK_SIZE:            16,
	MAP_SCALAR:            32,
	PATH_TO_TEXTURE_ATLAS: "assets/texture_atlas.png",
	TEXTURE_SIZE_PX:       32,
	NOISE:                 noiseGen,
}

// AREA : TILE
type TileType int

const (
	UNDEFINED TileType = iota
	WATER     TileType = iota
	SAND      TileType = iota
	GRASS     TileType = iota
	FOREST    TileType = iota
	MOUNTAIN  TileType = iota
)

type Tile struct {
	TileType TileType
}

func (tile *Tile) Draw(sx, sy float32, textureAtlas rl.Texture2D) {
	textureData := GetTileTexture(tile.TileType)
	rl.DrawTexturePro(
		textureAtlas,
		textureData.srcRect,
		rl.NewRectangle(sx, sy, Constants.TILE_SIZE, Constants.TILE_SIZE),
		rl.NewVector2(0, 0),
		0,
		rl.White,
	)
}

// END AREA : TILE

// AREA : CHUNK
type Chunk struct {
	Tiles *vl.Matrix[Tile]
}

func (chunk *Chunk) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(chunk); err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

func NewChunkFromSerialized(bufferBytes []byte) *Chunk {
	buffer := bytes.NewBuffer(bufferBytes)
	decoder := gob.NewDecoder(buffer)
	chunk := &Chunk{}
	if err := decoder.Decode(chunk); err != nil {
		panic(err)
	}
	return chunk
}

type ChunkMap struct {
	chunks map[vl.Vec2Int32]*Chunk
}

func heightToTileType(height float32) TileType {
	if height < 0.2 {
		return WATER
	}
	if height < 0.3 {
		return SAND
	}
	if height < 0.5 {
		return GRASS
	}
	if height < 0.7 {
		return FOREST
	}
	if height <= 1 {
		return MOUNTAIN
	}
	return UNDEFINED
}

func generateChunk(chunkCoord vl.Vec2Int32) *Chunk {
	chunkSize := int32(Constants.CHUNK_SIZE)
	chunk := &Chunk{
		Tiles: vl.NewMatrix[Tile](chunkSize, chunkSize),
	}
	for y := int32(0); y < chunkSize; y++ {
		for x := int32(0); x < chunkSize; x++ {
			height := Constants.NOISE.Eval32(
				float32(chunkCoord.X*chunkSize+x)/Constants.MAP_SCALAR,
				float32(chunkCoord.Y*chunkSize+y)/Constants.MAP_SCALAR,
			)
			chunk.Tiles.Set(y, x, Tile{
				TileType: heightToTileType(height),
			})
		}
	}
	return chunk
}

// Expects tx, ty as tile coordinates (should be integer float32)
func (world *World) get(tx, ty float32) Tile {

	// DEBUG CHUNK BORDERS
	// /*
	if int32(tx)%int32(Constants.CHUNK_SIZE) == 0 || int32(ty)%int32(Constants.CHUNK_SIZE) == 0 {
		return Tile{
			TileType: UNDEFINED,
		}
	}
	//*/
	// END DEBUG CHUNK BORDERS

	key := vl.Vec2Int32{
		X: vl.FloorDiv(tx, Constants.CHUNK_SIZE),
		Y: vl.FloorDiv(ty, Constants.CHUNK_SIZE),
	}

	tileInChunkX := vl.FloatModInt(tx, Constants.CHUNK_SIZE)
	tileInChunkY := vl.FloatModInt(ty, Constants.CHUNK_SIZE)

	// try Main Memory
	chunk, ok := world.chunkMap.chunks[key]
	if ok {
		tile := chunk.Tiles.Get(tileInChunkY, tileInChunkX)
		return tile
	}

	// try Disk
	chunk, ok = world.backend.GetChunk(key.X, key.Y)
	if ok {
		tile := chunk.Tiles.Get(tileInChunkY, tileInChunkX)
		world.chunkMap.chunks[key] = chunk
		return tile
	}

	// All else fails, generate chunk
	chunk = generateChunk(key)
	world.chunkMap.chunks[key] = chunk
	world.backend.SaveChunk(key.X, key.Y, chunk)
	tile := chunk.Tiles.Get(tileInChunkY, tileInChunkX)
	return tile
}

// END AREA : CHUNK

// AREA : ECS
type World struct {
	Name string

	chunkMap ChunkMap
	backend  *Backend
}

func NewWorld(name string) *World {
	return &World{
		Name: name,
		chunkMap: ChunkMap{
			chunks: make(map[vl.Vec2Int32]*Chunk),
		},
		backend: NewBackend(),
	}
}

// END AREA : ECS

// AREA : GAME
type Camera struct {
	position rl.Vector2
}

func NewCamera() Camera {
	return Camera{
		position: rl.NewVector2(0, 0),
	}
}

type AppSettings struct {
	WindowSize rl.Vector2
}

type App struct {
	AppSettings  AppSettings
	WindowTitle  string
	Mode         Mode
	TextureAtlas rl.Texture2D
}

func (game *ModeGame) Draw(app *App) {
	wTopLeft := rl.NewVector2(
		game.Camera.position.X-app.AppSettings.WindowSize.X/2.0,
		game.Camera.position.Y-app.AppSettings.WindowSize.Y/2.0,
	)
	wBottomRight := rl.NewVector2(
		game.Camera.position.X+app.AppSettings.WindowSize.X/2.0,
		game.Camera.position.Y+app.AppSettings.WindowSize.Y/2.0,
	)
	tTopLeft := vl.Vec2Int32{
		X: vl.FloorDiv(wTopLeft.X, Constants.TILE_SIZE),
		Y: vl.FloorDiv(wTopLeft.Y, Constants.TILE_SIZE),
	}
	tBottomRight := vl.Vec2Int32{
		X: vl.FloorDiv(wBottomRight.X, Constants.TILE_SIZE) + 1,
		Y: vl.FloorDiv(wBottomRight.Y, Constants.TILE_SIZE) + 1,
	}
	var sy float32 = -vl.FloatModFloat(wTopLeft.Y, Constants.TILE_SIZE)

	for ty := float32(tTopLeft.Y); ty < float32(tBottomRight.Y); ty += 1 {
		var sx float32 = -vl.FloatModFloat(wTopLeft.X, Constants.TILE_SIZE)
		for tx := float32(tTopLeft.X); tx < float32(tBottomRight.X); tx += 1 {

			tile := game.World.get(tx, ty)
			tile.Draw(sx, sy, app.TextureAtlas)
			// DEBUG ORIGIN
			/*
				if tx == 0 && ty == 0 {
					rl.DrawRectangle(
						int32(sx),
						int32(sy),
						int32(Constants.TILE_SIZE),
						int32(Constants.TILE_SIZE),
						rl.Magenta,
					)
				}
			// */
			// END DEBUG ORIGIN

			sx += Constants.TILE_SIZE
		}
		sy += Constants.TILE_SIZE
	}
}

func NewGame() *App {
	game := &App{
		AppSettings: AppSettings{
			WindowSize: rl.NewVector2(900, 800),
		},
		WindowTitle: "YAG",
		Mode:        NewModeMenu(),
	}
	return game
}

func (app *App) Run() {
	rl.InitWindow(
		int32(app.AppSettings.WindowSize.X),
		int32(app.AppSettings.WindowSize.Y),
		app.WindowTitle,
	)

	// Initialize texture atlas (needs be called after rl.InitWindow)
	app.TextureAtlas = rl.LoadTexture(
		Constants.PATH_TO_TEXTURE_ATLAS,
	)

	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		app.Mode.UpdateMode(app)
	}
}

// END AREA : GAME
