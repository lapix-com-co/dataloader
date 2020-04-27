package pkg

import (
	"context"
	"time"

	"github.com/lapix-com-co/dataloader/pkg/pagination"
)

// Pet refers to the customer's pet
type Pet struct {
	ID        uint `gorm:"type:int auto_increment;primary_key"`
	Name      string
	OwnerID   string
	CreatedAt time.Time
}

type petDataLoader interface {
	Create(context.Context, *Pet) error
	Update(context.Context, *Pet) error
	Delete(context.Context, *Pet) error
	Find(context.Context, []string) ([]*Pet, error)
	All(context.Context, pagination.Input) (*PetPage, error)

	// Dataloader does not know about this method, it will be skiped.
	AnyOtherMethod(context.Context, []string) error
}
