package cmddb

import (
	"fmt"
	"time"
)

// SlackClient is the slack_clients db model definition
type SlackClient struct {
	ID        int       `db:"id" json:"id"`
	Token     string    `db:"token" json:"token"`
	Name      string    `db:"name" json:"name"`
	Data      JSONRaw   `db:"data" json:"data"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// GetSlackClients retrieves all the slack_clients from db
func (db *DB) GetSlackClients() ([]SlackClient, error) {

	var clients []SlackClient

	err := db.Sesh.Collection("slack_clients").Find().All(&clients)

	if err != nil {
		return clients, err
	}

	return clients, nil

}

// CreateSlackClientsTable syncs the slack_client table
func (db *DB) CreateSlackClientsTable() {

	fmt.Println("Attempting sync with slack_clients")
	query := `
		CREATE TABLE IF NOT EXISTS "slack_clients" (
     "id"         SERIAL PRIMARY KEY,
     "token"      TEXT UNIQUE NOT NULL,
     "name"       TEXT,
     "data"       JSONB NOT NULL,
     "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
     "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`

	db.Sesh.Query(query)

}
