package repository

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"

	"strings"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/actor"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type ActorPgRepository struct {
	db *sql.DB
}

func NewActorPgRepository(conn *sql.DB) actor.ActorRepository {
	return &ActorPgRepository{
		db: conn,
	}
}

func (rep *ActorPgRepository) Insert(actor *models.Actor) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	err = tx.QueryRow(
		`INSERT INTO actors(name)
		VALUES ($1)
		RETURNING id`,
		actor.Name).Scan(&actor.ID)
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

func (rep *ActorPgRepository) Update(actor *models.Actor) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`UPDATE actors
		SET name = $2
		WHERE id = $1`,
		actor.ID, actor.Name)
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

func (rep *ActorPgRepository) DeleteById(id uint64) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`DELETE FROM actors
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

func (rep *ActorPgRepository) SelectById(id uint64) (*models.Actor, error) {
	dbActor := &models.Actor{}

	row := rep.db.QueryRow(
		`SELECT id, name
		FROM actors WHERE id=$1`, id)

	err := row.Scan(&dbActor.ID, &dbActor.Name)
	if err != nil {
		return nil, err
	}
	return dbActor, nil
}

func (rep *ActorPgRepository) SelectWhereNameLike(name string, limit, offset uint64) ([]*models.Actor, error) {
	selectQuery := `
		SELECT id, name
		FROM actors
		WHERE name ILIKE $1
		ORDER BY id`

	var values []interface{}
	searchName := "%" + name + "%"
	values = append(values, searchName)

	var pgntQuery string
	if limit != 0 {
		pgntQuery = "LIMIT $2 OFFSET $3"
		values = append(values, limit, offset)
	}

	resultQuery := strings.Join([]string{
		selectQuery,
		pgntQuery,
	}, " ")

	rows, err := rep.db.Query(resultQuery, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actors []*models.Actor

	for rows.Next() {
		actor := &models.Actor{}

		err := rows.Scan(&actor.ID, &actor.Name)
		if err != nil {
			return nil, err
		}

		actors = append(actors, actor)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return actors, nil
}

func (rep *ActorPgRepository) SelectAll(pgnt *models.Pagination) ([]*models.Actor, error) {
	var values []interface{}

	selectQuery := `
		SELECT id, name
		FROM actors
		ORDER BY name`

	var pgntQuery string
	if pgnt.Count != 0 {
		pgntQuery = "LIMIT $1 OFFSET $2"
		values = append(values, pgnt.Count, pgnt.From)
	}

	resultQuery := strings.Join([]string{
		selectQuery,
		pgntQuery,
	}, " ")

	rows, err := rep.db.Query(resultQuery, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actors []*models.Actor
	for rows.Next() {
		actor := &models.Actor{}
		err := rows.Scan(&actor.ID, &actor.Name)
		if err != nil {
			return nil, err
		}
		actors = append(actors, actor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return actors, nil
}
