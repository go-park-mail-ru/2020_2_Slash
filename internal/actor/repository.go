package actor

import "github.com/go-park-mail-ru/2020_2_Slash/internal/models"

type ActorRepository interface {
	Insert(actor *models.Actor) error
	Update(actor *models.Actor) error
	DeleteById(id uint64) error
	SelectById(id uint64) (*models.Actor, error)
	SelectWhereNameLike(name string, limit, offset uint64) ([]*models.Actor, error)
	SelectAll(pgnt *models.Pagination) ([]*models.Actor, error)
}
