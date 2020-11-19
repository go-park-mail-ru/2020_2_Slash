package repository

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/actor"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type ActorPgRepository struct {
	dbConn *sql.DB
}

func NewActorPgRepository(conn *sql.DB) actor.ActorRepository {
	return &ActorPgRepository{
		dbConn: conn,
	}
}

func (ar *ActorPgRepository) Insert(actor *models.Actor) error {
	tx, err := ar.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	err = tx.QueryRow(
		`INSERT INTO actors(name)
		VALUES ($1)
		RETURNING id`,
		actor.Name).Scan(&actor.ID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ar *ActorPgRepository) Update(actor *models.Actor) error {
	tx, err := ar.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`UPDATE actors
		SET name = $2
		WHERE id = $1`,
		actor.ID, actor.Name)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ar *ActorPgRepository) DeleteById(id uint64) error {
	tx, err := ar.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`DELETE FROM actors
		WHERE id = $1`, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ar *ActorPgRepository) SelectById(id uint64) (*models.Actor, error) {
	dbActor := &models.Actor{}

	row := ar.dbConn.QueryRow(
		`SELECT id, name
		FROM actors WHERE id=$1`, id)

	err := row.Scan(&dbActor.ID, &dbActor.Name)
	if err != nil {
		return nil, err
	}
	return dbActor, nil
}
