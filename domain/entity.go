package domain

import (
	"time"
)

// Entity is the base for all domain entities with identity and lifecycle.
type Entity interface {
	ID() string
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

// BaseEntity provides common entity fields.
type BaseEntity struct {
	id        string
	createdAt time.Time
	updatedAt time.Time
}

func NewBaseEntity(id string) *BaseEntity {
	now := time.Now().UTC()
	return &BaseEntity{
		id:        id,
		createdAt: now,
		updatedAt: now,
	}
}

func (e *BaseEntity) ID() string              { return e.id }
func (e *BaseEntity) CreatedAt() time.Time    { return e.createdAt }
func (e *BaseEntity) UpdatedAt() time.Time    { return e.updatedAt }
func (e *BaseEntity) Touch()                  { e.updatedAt = time.Now().UTC() }
