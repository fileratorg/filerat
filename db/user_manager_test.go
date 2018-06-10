package db

import (
	"testing"
	"github.com/fileratorg/filerat/db/models"
	"github.com/satori/go.uuid"
)

var boltPath = "bolt://neo4j:admin@0.0.0.0"
var port = 7687
var uniqueId, _ = uuid.NewV4()

func TestUserCreate(t *testing.T) {

	t.Run("Verify neo4j connector path", func(t *testing.T) {
		conn := new(DbConnector)
		db := conn.Open(boltPath, port)
		if db.Error != nil {
			t.Errorf("%s.", db.Error)
		}

		user := models.AuthUser{Username:"test", Email:"Test@test.com", Password:"password"}
		user.Model.UniqueId = uniqueId
		conn.Save(&user)
	})

}

func TestUserGet(t *testing.T) {

	t.Run("Verify neo4j connector path", func(t *testing.T) {
		conn := new(DbConnector)
		db := conn.Open(boltPath, port)
		if db.Error != nil {
			t.Errorf("%s.", err)
		}

		//user := conn.Get(uniqueId)
	})

}

func TestUserUpdate(t *testing.T) {

	t.Run("Verify neo4j connector path", func(t *testing.T) {
		conn := new(DbConnector)
		db := conn.Open(boltPath, port)
		if db.Error != nil {
			t.Errorf("%s.", db.Error)
		}

		//user := conn.Get(uniqueId)
		user := models.AuthUser{Username:"test2", Email:"Test2@test.com", Password:"password2"}
		conn.Save(&user)
	})

}

func TestUserSoftDelete(t *testing.T) {

	t.Run("Verify neo4j connector path", func(t *testing.T) {
		conn := new(DbConnector)
		db := conn.Open(boltPath, port)
		if db.Error != nil {
			t.Errorf("%s.", db.Error)
		}

		user := models.AuthUser{Username:"test2", Email:"Test2@test.com", Password:"password2"}
		conn.Delete(&user, true)
	})

}

func TestUserHardDelete(t *testing.T) {

	t.Run("Verify neo4j connector path", func(t *testing.T) {
		conn := new(DbConnector)
		db := conn.Open(boltPath, port)
		if db.Error != nil {
			t.Errorf("%s.", db.Error)
		}

		user := models.AuthUser{Username:"test2", Email:"Test2@test.com", Password:"password2"}
		conn.Delete(&user, false)
	})

}