package cmddb

import (
	"github.com/william-olson/cmd-server/cmddeps"
	"github.com/william-olson/cmd-server/cmdutils"
	"gopkg.in/matryer/try.v1"
	"math"
	"time"
	upper "upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/postgresql"
)

const (
	maxRetries  = 10
	retryFactor = 1.7
)

// DB is the database helper and holds the session to db
type DB struct {
	Deps   *cmddeps.Deps
	Sesh   upper.Database
	logger cmdutils.Logger
}

// Connect establishes a connection to the database
func (db *DB) Connect() {

	db.logger = cmdutils.NewLogger("db")
	err := try.Do(db.connect())

	if err != nil {
		db.logger.Error("could not connect to db")
		panic(err)
	}

	// sync the slack_clients tables
	db.CreateSlackClientsTable()

}

/*

	Closure function to retry connection to db

*/
func (db *DB) connect() try.Func {

	config := db.Deps.Get("config").(*cmdutils.Config)
	dbUrl := getDbCreds(config)

	return func(attempt int) (bool, error) {

		var err error
		shouldRetry := attempt <= maxRetries
		timeout := time.Second * time.Duration(math.Pow(retryFactor, float64(attempt)))

		// connect
		db.logger.KV("attempt", attempt).Log("connecting to db")
		sesh, err := postgresql.Open(dbUrl)

		// connect err
		if err != nil {
			time.Sleep(timeout)
			return shouldRetry, err
		}

		// check ping
		err = sesh.Ping()
		if err != nil {
			time.Sleep(timeout)
		}

		db.Sesh = sesh
		return shouldRetry, err

	}

}

/*

	Retrieve a connection url for the postgres db based on app config

*/
func getDbCreds(config *cmdutils.Config) postgresql.ConnectionURL {

	return postgresql.ConnectionURL{
		Host:     config.Get("DB_HOST"),
		Database: config.Get("DB_DBNAME"),
		User:     config.Get("DB_USER"),
		Password: config.Get("DB_PW"),
	}

}
