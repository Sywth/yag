package game

import rl "github.com/gen2brain/raylib-go/raylib"

type TextureOnAtlasInfo struct {
	srcRect rl.Rectangle
}

func GetTileTexture(tileType TileType) TextureOnAtlasInfo {
	textureOffset := rl.NewVector2(0, 0)
	switch tileType {
	case UNDEFINED:
		textureOffset.X = Constants.TEXTURE_SIZE_PX * 0
		textureOffset.Y = Constants.TEXTURE_SIZE_PX * 0
	case WATER:
		textureOffset.X = Constants.TEXTURE_SIZE_PX * 1
		textureOffset.Y = Constants.TEXTURE_SIZE_PX * 0
	case SAND:
		textureOffset.X = Constants.TEXTURE_SIZE_PX * 2
		textureOffset.Y = Constants.TEXTURE_SIZE_PX * 0
	case GRASS:
		textureOffset.X = Constants.TEXTURE_SIZE_PX * 3
		textureOffset.Y = Constants.TEXTURE_SIZE_PX * 0
	case FOREST:
		textureOffset.X = Constants.TEXTURE_SIZE_PX * 4
		textureOffset.Y = Constants.TEXTURE_SIZE_PX * 0
	case MOUNTAIN:
		textureOffset.X = Constants.TEXTURE_SIZE_PX * 5
		textureOffset.Y = Constants.TEXTURE_SIZE_PX * 0
	}

	return TextureOnAtlasInfo{
		srcRect: rl.NewRectangle(
			textureOffset.X,
			textureOffset.Y,
			Constants.TEXTURE_SIZE_PX,
			Constants.TEXTURE_SIZE_PX,
		),
	}
}
