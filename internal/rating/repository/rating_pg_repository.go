package repository

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/rating"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
)

type RatingPgRepository struct {
	dbConn *sql.DB
}

func NewRatingPgRepository(db *sql.DB) rating.RatingRepository {
	return &RatingPgRepository{dbConn: db}
}

func (rep *RatingPgRepository) Insert(rating *models.Rating) error {
	tx, err := rep.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO rates(user_id, content_id, likes)
		VALUES ($1, $2, $3)`,
		rating.UserID, rating.ContentID, rating.Likes)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logger.Error(rollbackErr.Error())
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (rep *RatingPgRepository) SelectByUserIDContentID(userID uint64, contentID uint64) (*models.Rating, error) {
	rates := &models.Rating{}

	err := rep.dbConn.QueryRow(`
		SELECT user_id, content_id, likes
		FROM rates
		WHERE user_id=$1 AND content_id=$2`, userID, contentID).
		Scan(&rates.UserID, &rates.ContentID, &rates.Likes)
	if err != nil {
		return nil, err
	}

	return rates, nil
}

func (rep *RatingPgRepository) Update(rating *models.Rating) error {
	tx, err := rep.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE rates
		SET likes=$1
		WHERE user_id=$2 AND content_id=$3`,
		rating.Likes, rating.UserID, rating.ContentID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logger.Error(rollbackErr.Error())
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (rep *RatingPgRepository) SelectRatesCount(contentID uint64) (int, error) {
	var ratesCount int
	err := rep.dbConn.QueryRow(`
		SELECT COUNT(*)
		FROM rates
		WHERE content_id=$1
		GROUP BY content_id`, contentID).Scan(&ratesCount)
	if err != nil {
		return 0, err
	}
	return ratesCount, nil
}

func (rep *RatingPgRepository) SelectLikesCount(contentID uint64) (int, error) {
	var likesCount int
	err := rep.dbConn.QueryRow(`
		SELECT COUNT(*)
		FROM rates
		WHERE likes=true AND content_id=$1
		GROUP BY content_id`, contentID).
		Scan(&likesCount)
	if err != nil {
		return 0, err
	}
	return likesCount, nil
}

func (rep *RatingPgRepository) Delete(rating *models.Rating) error {
	tx, err := rep.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM rates
		WHERE user_id=$1 AND content_id=$2`,
		rating.UserID, rating.ContentID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logger.Error(rollbackErr.Error())
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
