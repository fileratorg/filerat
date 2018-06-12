package neo4j_driver

import (
	"time"
	"github.com/satori/go.uuid"
)

type Model struct {
	Id			int64 `type:"id"`
	UniqueId	uuid.UUID `type:"unique_id"`
	CreatedAt	time.Time
	UpdatedAt	time.Time
	DeletedAt	time.Time
}