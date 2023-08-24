package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type Gallery struct {
	ID     int
	UserID int
	Title  string
}

type GalleryService struct {
	DB *sql.DB
}

func (gs *GalleryService) Create(title string, userID int) (*Gallery, error) {
	gallery := Gallery{
		Title:  title,
		UserID: userID,
	}

	row := gs.DB.QueryRow(`
    INSERT INTO galleries (title, user_id)
    VALUES ($1, $2)
    RETURNING id;
  `, gallery.Title, gallery.UserID)

	err := row.Scan(&gallery.ID)
	if err != nil {
		return nil, fmt.Errorf("create gallery: %w", err)
	}

	return &gallery, nil
}

func (gs *GalleryService) ByID(id int) (*Gallery, error) {
	// TODO: add validation on the ID passed in
	gallery := Gallery{
		ID: id,
	}

	row := gs.DB.QueryRow(`
    SELECT title, user_id
    FROM galleries
    WHERE id = $1;
  `, gallery.ID)

	err := row.Scan(&gallery.Title, &gallery.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// now the controller package don't have to
			// import as a dependency to check for the errors
			// we only have to import ErrNotFound variable to check
			// for the error
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query gallery by id: %w", err)
	}

	return &gallery, nil
}
