package game

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Backend struct {
	Database *sql.DB
}

func NewBackend() *Backend {
	database, _ := sql.Open("sqlite3", "./data/game.db")
	createChunkTable, _ := database.Prepare(`
		CREATE TABLE IF NOT EXISTS chunk (
			x INTEGER,
			y INTEGER,
			chunkData BLOB,
			PRIMARY KEY (x, y)
		)
	`)
	createChunkTable.Exec()
	return &Backend{
		Database: database,
	}
}

func (backend *Backend) SaveChunk(x, y int32, chunk *Chunk) {
	serializedChunk := chunk.Serialize()
	backend.Database.Exec(`
		INSERT OR REPLACE INTO chunk (x, y, chunkData) VALUES (?, ?, ?)
	`, x, y, serializedChunk)
}

func (backend *Backend) GetChunk(x, y int32) (bool, *Chunk) {
	row := backend.Database.QueryRow("SELECT chunkData FROM chunk WHERE x = ? AND y = ?", x, y)
	var serializedChunk []byte
	err := row.Scan(&serializedChunk)
	if err != nil {
		fmt.Println(fmt.Errorf("error getting chunk: %v", err))
		return false, nil
	}
	return true, NewChunkFromSerialized(serializedChunk)
}
