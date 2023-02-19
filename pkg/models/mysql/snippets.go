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

	query := "SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?"

	row := spp.DB.QueryRow(query, id)

	snippet := models.Snippet{}

	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return &snippet, nil
}

// Return the 10 latest snippet from the database
func (spp *SnippetModel) Lastest() ([]*models.Snippet, error) {
	query := "SELECT id,title, content,created,expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10"

	rows, err := spp.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []*models.Snippet

	for rows.Next() {
		snippet := &models.Snippet{}

		err := rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)

		if err != nil {
			return nil, err
		}

		snippets = append(snippets, snippet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
