package db

import (
	"testing"
	"github.com/fileratorg/filerat/db/models"
		"github.com/fileratorg/filerat/db/neo4j_driver"
	"github.com/satori/go.uuid"
	"time"
	"github.com/fileratorg/filerat/utils"
)

var uniqueIdUser, _ = uuid.NewV4()
var uniqueIdFile, _ = uuid.NewV4()
var uniqueIdFile2, _ = uuid.NewV4()
var uniqueIdFolder, _ = uuid.NewV4()

func TestRatFileCreate(t *testing.T) {

	t.Run("Test file creation", func(t *testing.T) {
		conn := new(DbConnector)
		now := time.Now()

		user := models.AuthUser{Username:"test", Email:"Test@test.com", Password:"password"}
		user.Model = new(neo4j_driver.Model)
		user.Model.CreatedAt = now
		user.Model.UpdatedAt = now
		user.Model.UniqueId = uniqueIdUser
		user.Password = utils.GetPasswordHash(user.Password)
		conn.SaveUser(&user)

		// file without a folder
		file := models.RatFile{}
		file.Model = new(neo4j_driver.Model)
		file.Model.CreatedAt = now
		file.Model.UpdatedAt = now
		file.Model.UniqueId = uniqueIdFile
		file.Owner = &user
		conn.SaveFile(&file)
	})
}

func TestRatFileCreateWithFolder(t *testing.T) {

	t.Run("Test file creation with folder", func(t *testing.T) {
		conn := new(DbConnector)
		now := time.Now()

		user := conn.GetUser(uniqueIdUser)

		folder := models.RatFolder{}
		folder.Model = new(neo4j_driver.Model)
		folder.Model.CreatedAt = now
		folder.Model.UpdatedAt = now
		folder.Model.UniqueId = uniqueIdFolder
		folder.Owner = &user
		conn.SaveFolder(&folder)

		// file without a folder
		file := models.RatFile{}
		file.Model = new(neo4j_driver.Model)
		file.Model.CreatedAt = now
		file.Model.UpdatedAt = now
		file.Model.UniqueId = uniqueIdFile2
		file.Owner = &user
		file.RatFolder = &folder
		conn.SaveFile(&file)
	})
}

func TestRatFileGet(t *testing.T) {

	t.Run("Test file get", func(t *testing.T) {
		conn := new(DbConnector)

		conn.GetFile(uniqueIdFile)
	})

}

func TestRatFileUpdate(t *testing.T) {

	t.Run("Test file update", func(t *testing.T) {
		conn := new(DbConnector)

		obj := conn.GetFile(uniqueIdFile)
		obj.Name = "changed"
		conn.SaveFile(&obj)
		if obj.Name != "changed" {
			t.Errorf("name not changed")
		}
	})
}

func TestRatFileSoftDelete(t *testing.T) {

	t.Run("Test file soft delete", func(t *testing.T) {
		conn := new(DbConnector)

		obj := conn.GetFile(uniqueIdFile)
		conn.DeleteFile(&obj, true)
	})
}

func TestDeleteObjects(t *testing.T) {

	t.Run("Deleting objects for files", func(t *testing.T) {
		conn := new(DbConnector)

		objUser := conn.GetUser(uniqueIdUser)
		conn.DeleteUser(&objUser, false)

		obj := conn.GetFile(uniqueIdFile)
		conn.DeleteFile(&obj, false)

		obj2 := conn.GetFile(uniqueIdFile2)
		conn.DeleteFile(&obj2, false)

		objFolder := conn.GetFolder(uniqueIdFolder)
		conn.DeleteFolder(&objFolder, false)
	})
}
