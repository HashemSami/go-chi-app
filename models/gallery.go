package models

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
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

	// Delete all the images inside this gallery as well
	dir := gs.galleryDir(id)
	fmt.Print(dir)
	err = os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("deleting images: %w", err)
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
			images = append(images, Image{
				GalleryID: galleryID,
				Path:      file,
				FileName:  filepath.Base(file),
			})
		}
	}
	return images, nil
}

func (gs *GalleryService) extensions() []string {
	return []string{".png", ".jpg", ".jpeg", ".gif"}
}

func (gs *GalleryService) imagesContentTypes() []string {
	return []string{"image/png", "image/jpeg", "image/gif"}
}

func (gs *GalleryService) Image(galleryID int, filename string) (Image, error) {
	imagePath := filepath.Join(gs.galleryDir(galleryID), filename)

	// get the status of the file
	_, err := os.Stat(imagePath)
	if err != nil {
		// return a custom error if the image does not exists
		if errors.Is(err, fs.ErrNotExist) {
			return Image{}, ErrNotFound
		}
		return Image{}, fmt.Errorf("querying for image: %w", err)
	}

	return Image{
		FileName:  filename,
		GalleryID: galleryID,
		Path:      imagePath,
	}, nil
}

func (gs *GalleryService) CreateImage(galleryID int, filename string,
	contents io.ReadSeeker,
) error {
	// checking for the extension and the file type
	// before creating
	err := checkContentType(contents, gs.imagesContentTypes())
	if err != nil {
		return fmt.Errorf("creating image %v: %w", filename, err)
	}

	err = checkExtension(filename, gs.extensions())
	if err != nil {
		return fmt.Errorf("creating image %v: %w", filename, err)
	}

	galleryDir := gs.galleryDir(galleryID)
	// make the directory if not exists, along with the other parents
	err = os.MkdirAll(galleryDir, 0o755)
	if err != nil {
		return fmt.Errorf("creating gallery-%d images directory: %w", galleryID, err)
	}

	imagePath := filepath.Join(galleryDir, filename)
	dst, err := os.Create(imagePath)
	if err != nil {
		return fmt.Errorf("creating image file: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, contents)

	if err != nil {
		return fmt.Errorf("copying contents to images: %w", err)
	}
	return nil
}

func (gs *GalleryService) DeleteImage(galleryID int, filename string) error {
	image, err := gs.Image(galleryID, filename)
	if err != nil {
		return fmt.Errorf("deleting image: %w", err)
	}

	err = os.Remove(image.Path)
	if err != nil {
		return fmt.Errorf("deleting images: %w", err)
	}

	return nil
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
