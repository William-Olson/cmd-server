package cmddb

import (
	"fmt"
	"time"
)

// SlackClient is the slack_clients db model definition
type SlackClient struct {
	ID          int       `db:"id" json:"id"`
	Token       string    `db:"token" json:"-"`
	Host        string    `db:"host" json:"host"`
	VersionPath string    `db:"version_path" json:"version_path"`
	Name        string    `db:"name" json:"name"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	slugs       []SlackSlug
}

// SlackSlug is the model definition for slack_slugs table entries
type SlackSlug struct {
	ID            int       `db:"id" json:"id"`
	SlackClientID int       `db:"slack_client_id" json:"slack_client_id"`
	Name          string    `db:"name" json:"name"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

// GetSlackClientsOrErr retrieves all the slack_clients from db
func (db *DB) GetSlackClientsOrErr() ([]SlackClient, error) {

	var clients []SlackClient

	err := db.Sesh.Collection("slack_clients").Find().All(&clients)

	if err != nil {
		return clients, err
	}

	return clients, nil

}

// CreateSlackClientsTable syncs the slack_client table
func (db *DB) CreateSlackClientsTable() {

	fmt.Println("Attempting sync with slack_clients & slack_slugs")
	clientQuery := `
		CREATE TABLE IF NOT EXISTS "slack_clients" (
			"id"           SERIAL PRIMARY KEY,
			"token"        TEXT UNIQUE NOT NULL,
			"host"         TEXT NOT NULL,
			"version_path" TEXT,
			"name"         TEXT,
			"created_at"   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			"updated_at"   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`

	slugQuery := `
		CREATE TABLE IF NOT EXISTS "slack_slugs" (
			"id"               SERIAL PRIMARY KEY,
			"slack_client_id"  SERIAL NOT NULL,
			"name"             TEXT NOT NULL,
			"created_at"       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			"updated_at"       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`

	db.Sesh.Query(clientQuery)
	db.Sesh.Query(slugQuery)

}

// GetSlackClientByTokenOrErr attempts to fetch a slack_client by its token
func (db *DB) GetSlackClientByTokenOrErr(token string) (SlackClient, error) {

	var slackClient SlackClient

	res := db.Sesh.Collection("slack_clients").Find("token", token)
	err := res.One(&slackClient)

	if err != nil {
		return slackClient, err
	}

	// get client slugs
	var slackSlugs []SlackSlug
	slugRes := db.Sesh.Collection("slack_slugs").Find("slack_client_id", slackClient.ID)
	err = slugRes.All(&slackSlugs)
	slackClient.slugs = slackSlugs

	return slackClient, err

}

// GetSlugs returns the slice of slugs for the slack_client
func (s *SlackClient) GetSlugs() []SlackSlug {

	return s.slugs

}

// CreateSlackSlugOrErr inserts a new slack_slug into the db
func (db *DB) CreateSlackSlugOrErr(newEntry SlackSlug) error {

	createTime := time.Now()
	_, err := db.Sesh.InsertInto("slack_slugs").
		Columns(
			"slack_client_id",
			"name",
			"created_at",
			"updated_at",
		).
		Values(
			newEntry.SlackClientID,
			newEntry.Name,
			createTime,
			createTime,
		).
		Exec()

	return err

}
