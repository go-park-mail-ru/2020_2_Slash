package repository

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/favourite"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type FavouritePgRepository struct {
	dbConn *sql.DB
}

func NewFavouritePgRepository(conn *sql.DB) favourite.FavouriteRepository {
	return &FavouritePgRepository{
		dbConn: conn,
	}
}

func (rep *FavouritePgRepository) Insert(favourite *models.Favourite) error {
	tx, err := rep.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO favourites(user_id, content_id, created)
		VALUES ($1, $2, $3)`, favourite.UserID, favourite.ContentID, favourite.Created)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (rep *FavouritePgRepository) Select(favourite *models.Favourite) error {
	dbFavourite := &models.Favourite{}

	row := rep.dbConn.QueryRow(`
		SELECT user_id, content_id, created
		FROM favourites
		WHERE user_id=$1 AND content_id=$2`,
		favourite.UserID, favourite.ContentID)
	err := row.Scan(&dbFavourite.UserID, &dbFavourite.ContentID, &dbFavourite.Created)

	return err
}

func (rep *FavouritePgRepository) SelectFavouriteContent(userID uint64) ([]*models.Content, error) {
	rows, err := rep.dbConn.Query(`
		SELECT c.id, name, original_name, description, short_description, year, images, type
		FROM favourites f
		JOIN content c on f.content_id = c.id
		WHERE user_id=$1
		ORDER BY created DESC`, userID)
	if err != nil {
		return nil, err
	}

	var favouriteContent []*models.Content

	for rows.Next() {
		content := &models.Content{}
		err = rows.Scan(&content.ContentID, &content.Name, &content.OriginalName,
			&content.Description, &content.ShortDescription, &content.Year,
			&content.Images, &content.Type)
		if err != nil {
			return nil, err
		}
		favouriteContent = append(favouriteContent, content)
	}

	return favouriteContent, nil
}

func (rep *FavouritePgRepository) Delete(favourite *models.Favourite) error {
	tx, err := rep.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`DELETE FROM favourites
		WHERE user_id=$1 AND content_id=$2`,
		favourite.UserID, favourite.ContentID)
	if err != nil {
		_ = tx.Rollback()
		return nil
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
