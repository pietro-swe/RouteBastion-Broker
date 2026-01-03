package provider

import (
	"time"

	uuid "github.com/satori/go.uuid"
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
	return &Provider{
		ID:        uuid.NewV4(),
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
