package provider

import (
	"time"

	"github.com/google/uuid"
)

type CommunicationMethod string

const (
	CommunicationMethodREST            CommunicationMethod = "rest"
	CommunicationMethodProtocolBuffers CommunicationMethod = "protocol_buffers"
)

type Provider struct {
	ID         uuid.UUID
	Name       string
	CreatedAt  time.Time
	ModifiedAt *time.Time
	DeletedAt  *time.Time
}

func NewProvider(name string) *Provider {
	id, err := uuid.NewV7()
	if err != nil {
		panic("failed to generate uuid v7 for provider")
	}

	return &Provider{
		ID:        id,
		Name:      name,
		CreatedAt: time.Now(),
	}
}

func FromDatabase(id uuid.UUID, name string, createdAt time.Time, modifiedAt, deletedAt *time.Time) *Provider {
	return &Provider{
		ID:         id,
		Name:       name,
		CreatedAt:  createdAt,
		ModifiedAt: modifiedAt,
		DeletedAt:  deletedAt,
	}
}
