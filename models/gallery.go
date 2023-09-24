package models

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

type Image struct {
	GalleryID int
	Path      string
	FileName  string
}

type Gallery struct {
	ID     int
	UserID int
	Title  string
}

type GalleryService struct {
	DB *sql.DB

	// ImagesDir is used to tell the GalleryServices where to store and locate
	// images. If not set, the GalleryServices will default to using the "images"
	// directory
	ImagesDir string
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
			// import sql as a dependency to check for the errors
			// we only have to import ErrNotFound variable to check
			// for the error
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query gallery by id: %w", err)
	}

	return &gallery, nil
}

func (gs *GalleryService) BYUserID(userID int) ([]Gallery, error) {
	rows, err := gs.DB.Query(`
    SELECT id, title
    FROM galleries
    WHERE user_id = $1;
  `, userID)
	if err != nil {
		return nil, fmt.Errorf("query galleries by user: %w", err)
	}

	var galleries []Gallery

	for rows.Next() {
		gallery := Gallery{
			UserID: userID,
		}

		err := rows.Scan(&gallery.ID, &gallery.Title)
		if err != nil {
			return nil, fmt.Errorf("query galleries by user: %w", err)
		}

		galleries = append(galleries, gallery)
	}

	// if error in for loop
	if rows.Err() != nil {
		return nil, fmt.Errorf("query galleries by user: %w", err)
	}

	return galleries, nil
}

func (gs *GalleryService) Update(gallery *Gallery) error {
	// don't want to return any values
	_, err := gs.DB.Exec(`
		UPDATE galleries
		SET title = $2
		WHERE id = $1;
	`, gallery.ID, gallery.Title)
	if err != nil {
		return fmt.Errorf("update gallery: %w", err)
	}
	return nil
}

func (gs *GalleryService) Delete(id int) error {
	_, err := gs.DB.Exec(`
		DELETE FROM galleries
		WHERE id = $1;
	`, id)
	if err != nil {
		return fmt.Errorf("delete gallery: %w", err)
	}
	return nil
}

func (gs *GalleryService) Images(galleryID int) ([]Image, error) {
	// looking for this pattern
	// "images/gallery-%d/*"
	globPattern := filepath.Join(gs.galleryDir(galleryID), "*")

	allFiles, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("retrieving gallery images: %w", err)
	}

	var images []Image
	for _, file := range allFiles {
		if hasExtension(file, gs.extensions()) {
			images = append(images, Image{Path: file})
		}
	}
	return images, nil
}

func (gs *GalleryService) extensions() []string {
	return []string{".png", ".jpg", ".jpeg", ".gif"}
}

func (gs *GalleryService) galleryDir(galleryID int) string {
	imagesDir := gs.ImagesDir
	if imagesDir == "" {
		imagesDir = "images"
	}

	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%d", galleryID))
}

// utility function to check if the file path has an extension
func hasExtension(file string, extension []string) bool {
	for _, ext := range extension {
		file = strings.ToLower(file)
		ext = strings.ToLower(ext)
		if filepath.Ext(file) == ext {
			return true
		}
	}
	return false
}
