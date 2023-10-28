package main

import (
	"fmt"
	"strings"

	"github.com/pachecoio/repositories-postgres/base"
)

type Repository[T base.Model] struct {
	table string
	db    *PostgresDB
}

func NewRepository[T base.Model](db *PostgresDB, table string) *Repository[T] {
	return &Repository[T]{db: db, table: table}
}

func (r *Repository[T]) Get(id any) (T, error) {
	var entity T
	err := r.db.DB.Get(&entity, "SELECT * FROM "+r.table+" WHERE id = $1", id)
	return entity, err
}

func (r *Repository[T]) Create(model *T) (any, error) {
	tx := r.db.DB.MustBegin()
	res, err := tx.NamedExec("INSERT INTO "+r.table+" (name) VALUES (:name)", model)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return "", err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *Repository[T]) Update(id any, data base.PartialUpdate[T]) error {
	tx := r.db.DB.MustBegin()

	var fields []string
	dataMap := make(map[string]any)
	for k, v := range data.ToUpdate().(map[string]interface{}) {
		fields = append(fields, k+"= :"+k)
		dataMap[k] = v
	}
	// join fields to stirng

	fieldsString := strings.Join(fields, ", ")

	dataMap["id"] = id

	statement := "UPDATE " + r.table + " SET " + fieldsString + " WHERE id = :id"
	_, err := tx.NamedExec(statement, dataMap)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository[T]) Delete(id any) error {
	tx := r.db.DB.MustBegin()
	statement := "DELETE FROM " + r.table + " WHERE id = $1"
	_, err := tx.Exec(statement, id)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

type DefaultFilters struct{}

func (f *DefaultFilters) ToQuery() any {
	return ""
}

func (r *Repository[T]) Filter(filters base.Filters, options ...base.FilterOptions) ([]T, error) {
	if filters == nil {
		filters = &DefaultFilters{}
	}
	var entities []T
	statement := filters.ToQuery().(string)

	if options != nil {
		for _, option := range options {
			if option.Limit != 0 {
				statement = statement + fmt.Sprintf(" LIMIT %d", option.Limit)
			}
			if option.Offset != 0 {
				statement = statement + fmt.Sprintf(" OFFSET %d", option.Offset)
			}
			if option.Sort != nil {
				statement = statement + fmt.Sprintf(" ORDER BY %s", option.Sort)
			}
		}
	}

	err := r.db.DB.Select(&entities, fmt.Sprintf("SELECT * FROM %s %s", r.table, statement))

	return entities, err
}

func (r *Repository[T]) FindOne(filters base.Filters) (*T, error) {
	if filters == nil {
		filters = &DefaultFilters{}
	}
	var entity T
	statement := filters.ToQuery().(string)
	err := r.db.DB.Get(&entity, fmt.Sprintf("SELECT * FROM %s %s", r.table, statement))
	return &entity, err
}

func (r *Repository[T]) Count(filters base.Filters) (int, error) {
	if filters == nil {
		filters = &DefaultFilters{}
	}
	var count int
	statement := filters.ToQuery().(string)
	err := r.db.DB.Get(&count, fmt.Sprintf("SELECT COUNT(*) FROM %s %s", r.table, statement))
	return count, err
}
