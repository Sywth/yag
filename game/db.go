package game

import (
	"database/sql"

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

func (backend *Backend) GetChunk(x, y int32) (*Chunk, bool) {
	rows, _ := backend.Database.Query(`
		SELECT chunkData FROM chunk WHERE x = ? AND y = ?
	`, x, y)
	defer rows.Close()

	if !rows.Next() {
		return nil, false
	}

	var serializedChunk []byte
	rows.Scan(&serializedChunk)
	chunk := NewChunkFromSerialized(serializedChunk)
	return chunk, true
}
