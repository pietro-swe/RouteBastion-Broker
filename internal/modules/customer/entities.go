package customer

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Customer struct {
	ID                 uuid.UUID
	BusinessIdentifier string
	Name               string
	APIKey             string
	CreatedAt          time.Time
	UpdatedAt          *time.Time
	DeletedAt          *time.Time
}

func NewCustomer(
	id uuid.UUID,
	name,
	businessIdentifier,
	apiKey string,
	createdAt time.Time,
	updatedAt,
	deletedAt *time.Time,
) *Customer {
	return &Customer{
		ID:                 id,
		BusinessIdentifier: businessIdentifier,
		Name:               name,
		APIKey:             apiKey,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
		DeletedAt:          deletedAt,
	}
}
