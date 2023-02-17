package mysql

import (
	"database/sql"

	"github.com/caleberi/snippet-bin/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

// Insert inserts a new SnippetModel into the database
func (spp *SnippetModel) Insert(title, content, expires string) (int, error) {
	query := "INSERT INTO snippets (title, content, created, expires) VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))"

	result, err := spp.DB.Exec(query, title, content, expires)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get returns the snippet model instance with the given id
func (spp *SnippetModel) Get(id int) (*models.Snippet, error) {
	return nil, nil
}

// Return the 10 latest snippet from the database
func (spp *SnippetModel) Lastest(id int) ([]*models.Snippet, error) {
	return nil, nil
}
