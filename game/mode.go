package game

import rl "github.com/gen2brain/raylib-go/raylib"

type GameConstants struct {
	MoveSpeed float32
}

func NewGameConstants() GameConstants {
	return GameConstants{
		MoveSpeed: 12,
	}
}

type Mode interface {
	UpdateMode(*App)
}

type ModeGame struct {
	Camera
	World         *World
	GameConstants GameConstants
}

func NewModeGame() Mode {
	gameModePtr := &ModeGame{
		Camera:        NewCamera(),
		World:         NewWorld("World"),
		GameConstants: NewGameConstants(),
	}
	return Mode(gameModePtr)
}

func (mode *ModeGame) UpdateMode(app *App) {
	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)

	mode.Draw(app)
	mode.HandleInput()

	rl.EndDrawing()
}

func (game *ModeGame) HandleInput() {

	if rl.IsKeyDown(rl.KeyW) {
		game.Camera.position.Y -= game.GameConstants.MoveSpeed
	}
	if rl.IsKeyDown(rl.KeyS) {
		game.Camera.position.Y += game.GameConstants.MoveSpeed
	}
	if rl.IsKeyDown(rl.KeyA) {
		game.Camera.position.X -= game.GameConstants.MoveSpeed
	}
	if rl.IsKeyDown(rl.KeyD) {
		game.Camera.position.X += game.GameConstants.MoveSpeed
	}

}

type ModeMenu struct{}

func NewModeMenu() Mode {
	return Mode(&ModeMenu{})
}

func (*ModeMenu) UpdateMode(app *App) {
	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)

	rl.DrawText("Menu", 10, 10, 20, rl.Black)
	rl.DrawText("Press Enter to start", 10, 40, 20, rl.Black)

	if rl.IsKeyPressed(rl.KeyEnter) {
		app.Mode = NewModeGame()
	}

	rl.EndDrawing()
}
