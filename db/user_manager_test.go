package db

import (
	"testing"
	"github.com/fileratorg/filerat/db/models"
)

func TestUserCreate(t *testing.T) {
	boltPath := "bolt://neo4j:admin@0.0.0.0"
	port := 7687

	t.Run("Verify neo4j connector path", func(t *testing.T) {
		conn := new(DbConnector)
		_, err := conn.Open(boltPath, port)
		if err != nil {
			t.Errorf("%s.", err)
		}

		user := models.AuthUser{Username:"test2", Email:"Test2@test.com", Password:"password2"}
		conn.Create(&user)
	})

}
