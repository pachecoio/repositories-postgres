package postgres

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pachecoio/repositories-postgres/adapters/postgres"
	"github.com/stretchr/testify/assert"
)

func TestPostgresDB_unit(t *testing.T) {

	sqlxConn, err := sqlx.Connect("sqlite3", ":memory:")

	db, err := postgres.NewPostgresDB(sqlxConn)
	defer db.Disconnect()

	assert.Nil(t, err)
	assert.NotNil(t, db)

	t.Run("Should create a sample table", func(t *testing.T) {

		// create table
		r, err := sqlxConn.MustExec(schema).RowsAffected()
		assert.Nil(t, err)
		assert.Equal(t, int64(0), r)

		tx := sqlxConn.MustBegin()

		// insert values
		_, err = tx.NamedExec("INSERT INTO sample (name) VALUES (:name)", &Sample{Name: "John Doe"})
		assert.Nil(t, err)

		// commit transaction
		err = tx.Commit()

		// query values
		var sample Sample
		err = sqlxConn.Get(&sample, "SELECT * FROM sample WHERE id = 1")
		assert.Nil(t, err)
		assert.Equal(t, "John Doe", sample.Name)

	})
}

var schema = `
CREATE TABLE IF NOT EXISTS sample (
	id INTEGER PRIMARY KEY,
	name TEXT
);
`

type Sample struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
