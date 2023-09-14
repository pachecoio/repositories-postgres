package main

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/pachecoio/repositories-postgres/base"
	"github.com/stretchr/testify/assert"
)

func TestPostgresRepository(t *testing.T) {

	sqlxDB, err := sqlx.Connect("sqlite3", ":memory:")
	assert.Nil(t, err)
	db, err := NewPostgresDB(sqlxDB)
	assert.Nil(t, err)
	defer db.Disconnect()

	db.DB.MustExec(sampleSchema)

	repo := NewRepository[Sample](db, "sample")

	t.Run("Should create a sample record", func(t *testing.T) {

		// insert values
		entity := Sample{Name: "John Doe"}
		_, err = repo.Create(&entity)
		assert.Nil(t, err)

		// query values
		sample, err := repo.Get(1)
		assert.Nil(t, err)
		assert.Equal(t, "John Doe", sample.Name)
	})

	t.Run("Should create and update a sample record", func(t *testing.T) {

		// insert values
		entity := Sample{Name: "Arya Stark"}
		id, err := repo.Create(&entity)
		assert.Nil(t, err)

		// query values
		sample, err := repo.Get(id)
		assert.Nil(t, err)
		assert.Equal(t, "Arya Stark", sample.Name)

		// update values
		assert.Nil(t, err)
		updateData := UpdateSample{Name: "Jon Snow"}
		err = repo.Update(id, &updateData)
		assert.Nil(t, err)

		// query values
		sample, err = repo.Get(id)
		assert.Nil(t, err)
		assert.Equal(t, "Jon Snow", sample.Name)
	})

	t.Run("Should create and delete a sample record", func(t *testing.T) {

		// insert values
		entity := Sample{Name: "John Doe"}
		_, err := repo.Create(&entity)
		assert.Nil(t, err)

		// query values
		sample, err := repo.Get(1)
		assert.Nil(t, err)
		assert.Equal(t, "John Doe", sample.Name)

		// delete values
		err = repo.Delete(1)
		assert.Nil(t, err)

		// query values
		sample, err = repo.Get(1)
		assert.NotNil(t, err)
		assert.Equal(t, "", sample.Name)
	})

	t.Run("should create and filter a sample record", func(t *testing.T) {
		// insert values
		entity := Sample{Name: "Jorah Mormont"}
		_, err := repo.Create(&entity)
		assert.Nil(t, err)

		samples, err := repo.Filter(&SampleFilters{Name: "Jorah Mormont"})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(samples))
		assert.Equal(t, "Jorah Mormont", samples[0].Name)
	})

	t.Run("should create records and filter with options", func(t *testing.T) {
		// insert values
		entity := Sample{Name: "John Doe"}
		_, err := repo.Create(&entity)
		assert.Nil(t, err)

		entity = Sample{Name: "Jane Doe"}
		_, err = repo.Create(&entity)
		assert.Nil(t, err)

		filters := SampleFilters{Name: "John Doe"}
		samples, err := repo.Filter(&filters, base.FilterOptions{Limit: 1})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(samples))
		assert.Equal(t, "John Doe", samples[0].Name)
	})

	t.Run("should create and count records", func(t *testing.T) {
		// insert values
		entity := Sample{Name: "Samwell Tarly"}
		_, err := repo.Create(&entity)
		assert.Nil(t, err)

		entity = Sample{Name: "Tyrion Lannister"}
		_, err = repo.Create(&entity)
		assert.Nil(t, err)

		filters := SampleFilters{Name: "Tyrion Lannister"}
		count, err := repo.Count(&filters)
		assert.Nil(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("should create and find single record", func(t *testing.T) {
		// insert values
		entity := Sample{Name: "Sansa Stark"}
		_, err := repo.Create(&entity)
		assert.Nil(t, err)

		entity = Sample{Name: "Arya Stark"}
		_, err = repo.Create(&entity)
		assert.Nil(t, err)

		filters := SampleFilters{Name: "Arya Stark"}
		sample, err := repo.FindOne(&filters)
		assert.Nil(t, err)
		assert.Equal(t, "Arya Stark", sample.Name)
	})
}

type UpdateSample struct {
	Name string `db:"name"`
	ID   any    `db:"id"`
}

func (u *UpdateSample) ToUpdate() any {
	return map[string]any{
		"name": u.Name,
	}
}

var sampleSchema = `
CREATE TABLE IF NOT EXISTS sample (
  id INTEGER PRIMARY KEY,
  name TEXT
);
`

type SampleFilters struct {
	Name string `db:"name"`
}

func (f *SampleFilters) ToQuery() any {
	return " WHERE name = '" + f.Name + "'"
}
