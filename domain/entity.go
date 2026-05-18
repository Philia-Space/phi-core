package domain

import (
	"time"
)

// Entity is the base for all domain entities with identity and lifecycle.
type Entity interface {
	EntityID() string
	EntityCreatedAt() time.Time
	EntityUpdatedAt() time.Time
}

// BaseEntity provides common entity fields.
type BaseEntity struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewBaseEntity(id string) *BaseEntity {
	now := time.Now().UTC()
	return &BaseEntity{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (e *BaseEntity) EntityID() string              { return e.ID }
func (e *BaseEntity) EntityCreatedAt() time.Time    { return e.CreatedAt }
func (e *BaseEntity) EntityUpdatedAt() time.Time    { return e.UpdatedAt }
func (e *BaseEntity) Touch()                        { e.UpdatedAt = time.Now().UTC() }
