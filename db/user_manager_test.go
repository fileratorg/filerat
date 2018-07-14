package db

import (
	"testing"
	"github.com/fileratorg/filerat/db/models"
	"github.com/satori/go.uuid"
	"time"
	"github.com/fileratorg/filerat/db/neo4j_driver"
	"github.com/fileratorg/filerat/utils"
	)

var uniqueId, _ = uuid.NewV4()

func TestUserCreate(t *testing.T) {

	t.Run("Verify neo4j connector path", func(t *testing.T) {
		conn := new(DbConnector)

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

		conn.GetUser(uniqueId)
	})

}

func TestUserUpdate(t *testing.T) {

	t.Run("Verify neo4j connector path", func(t *testing.T) {
		conn := new(DbConnector)

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

		user := conn.GetUser(uniqueId)
		conn.DeleteUser(&user, true)
	})

}

func TestUserHardDelete(t *testing.T) {

	t.Run("Verify neo4j connector path", func(t *testing.T) {
		conn := new(DbConnector)

		user := conn.GetUser(uniqueId)
		conn.DeleteUser(&user, false)
	})

}


func TestUserCreateWithOrg(t *testing.T) {

	t.Run("Verify neo4j connector path", func(t *testing.T) {
		conn := new(DbConnector)
		now := time.Now()

		var uniqueIdOrg, _ = uuid.NewV4()
		var uniqueIdWithOrg, _ = uuid.NewV4()

		org := models.RatOrg{Name:"Test"}
		org.Model = new(neo4j_driver.Model)
		org.Model.CreatedAt = now
		org.Model.UpdatedAt = now
		org.Model.UniqueId = uniqueIdOrg
		conn.SaveRatOrg(&org)

		user := models.AuthUser{Username:"test", Email:"Test@test.com", Password:"password"}
		user.Model = new(neo4j_driver.Model)
		user.Model.CreatedAt = now
		user.Model.UpdatedAt = now
		user.Model.UniqueId = uniqueIdWithOrg
		user.RatOrg = org
		user.Password = utils.GetPasswordHash(user.Password)
		conn.SaveUser(&user)


		conn.DeleteRatOrg(&org, false)
		conn.DeleteUser(&user, false)
	})

}