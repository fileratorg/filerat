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

func (obj *Model) InitParent() *Model {
	now := time.Now()
	if obj == nil{
		obj = new(Model)
		obj.CreatedAt = now
		obj.UpdatedAt = now
		var uniqueId, _ = uuid.NewV4()
		obj.UniqueId = uniqueId
	}
	return obj
}
