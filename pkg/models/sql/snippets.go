package sql

import (
	"database/sql"

	"github.com/caleberi/snippet-bin/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

// Insert inserts a new SnippetModel into the database
func (spp *SnippetModel) Insert(title, content string) (int, error) {
	return 0, nil
}

// Get returns the snippet model instance with the given id
func (spp *SnippetModel) Get(id int) (*models.Snippet, error) {
	return nil, nil
}

// Return the 10 latest snippet from the database
func (spp *SnippetModel) Lastest(id int) ([]*models.Snippet, error) {
	return nil, nil
}
