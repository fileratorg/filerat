package db

import (
	"testing"
	"github.com/fileratorg/filerat/db/models"
	"github.com/satori/go.uuid"
	"time"
	"github.com/fileratorg/filerat/db/neo4j_driver"
	"github.com/fileratorg/filerat/utils"
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
		now := time.Now()
		user.Model = new(neo4j_driver.Model)
		user.Model.CreatedAt = now
		user.Model.UpdatedAt = now
		user.Model.UniqueId = uniqueId
		user.Password = utils.GetPasswordHash(user.Password)
		conn.SaveUser(&user)
	})

}

func TestUserGet(t *testing.T) {

	t.Run("Verify neo4j connector path", func(t *testing.T) {
		conn := new(DbConnector)
		db := conn.Open(boltPath, port)
		if db.Error != nil {
			t.Errorf("%s.", db.Error)
		}

		conn.GetUser(uniqueId)
	})

}

func TestUserUpdate(t *testing.T) {

	t.Run("Verify neo4j connector path", func(t *testing.T) {
		conn := new(DbConnector)
		db := conn.Open(boltPath, port)
		if db.Error != nil {
			t.Errorf("%s.", db.Error)
		}

		user := conn.GetUser(uniqueId)
		user.Username = "changed"
		conn.SaveUser(&user)
		if user.Username != "changed" {
			t.Errorf("username not changed")
		}
	})

}

func TestUserSoftDelete(t *testing.T) {

	t.Run("Verify neo4j connector path", func(t *testing.T) {
		conn := new(DbConnector)
		db := conn.Open(boltPath, port)
		if db.Error != nil {
			t.Errorf("%s.", db.Error)
		}

		user := conn.GetUser(uniqueId)
		conn.DeleteUser(&user, true)
	})

}

func TestUserHardDelete(t *testing.T) {

	t.Run("Verify neo4j connector path", func(t *testing.T) {
		conn := new(DbConnector)
		db := conn.Open(boltPath, port)
		if db.Error != nil {
			t.Errorf("%s.", db.Error)
		}

		user := conn.GetUser(uniqueId)
		conn.DeleteUser(&user, false)
	})

}