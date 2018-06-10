package neo4j_driver

import (
	"time"
	"github.com/satori/go.uuid"
)

type Model struct {
	UniqueId	uuid.UUID
	CreatedAt	time.Time
	UpdatedAt	time.Time
	DeletedAt	time.Time
}