package repository

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/subscription"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
)

type SubscriptionPgRepository struct {
	db *sql.DB
}

func NewSubscriptionPgRepository(db *sql.DB) subscription.SubscriptionRepository {
	return &SubscriptionPgRepository{db: db}
}

func (rep *SubscriptionPgRepository) Insert(subscription *models.Subscription) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	err = tx.QueryRow(`
		INSERT INTO subscriptions(owner, expires, is_paid, is_canceled)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		subscription.UserID, subscription.Expires,
		subscription.IsPaid, subscription.IsCanceled).
		Scan(&subscription.ID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logger.Error(rollbackErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (rep *SubscriptionPgRepository) Update(subscription *models.Subscription) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE subscriptions
		SET expires = $2,
		    is_paid = $3,
		    is_canceled = $4
		WHERE id = $1`,
		subscription.ID, subscription.Expires,
		subscription.IsPaid, subscription.IsCanceled)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logger.Error(rollbackErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (rep *SubscriptionPgRepository) SelectByUserID(userID uint64) (*models.Subscription, error) {
	subscription := &models.Subscription{}
	err := rep.db.QueryRow(`
		SELECT id, owner, expires, is_paid, is_canceled
		FROM subscriptions
		WHERE owner=$1`, userID).
		Scan(&subscription.ID, &subscription.UserID,
			&subscription.Expires, &subscription.IsPaid,
			&subscription.IsCanceled)
	if err != nil {
		return nil, err
	}
	return subscription, nil
}

func (rep *SubscriptionPgRepository) DeleteByID(id uint64) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM subscriptions
		WHERE id = $1`, id)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logger.Error(rollbackErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
