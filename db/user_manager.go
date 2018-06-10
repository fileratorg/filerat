package db

import (
	"github.com/fileratorg/filerat/db/models"
	"github.com/fileratorg/filerat"
	"github.com/fileratorg/filerat/utils"
	"github.com/fileratorg/filerat/db/neo4j_driver"
	"time"
	"github.com/satori/go.uuid"
)

func (conn *DbConnector) SaveUser(user *models.AuthUser) (*models.AuthUser){
	conn.Open(filerat.BoltPath, filerat.Port)
	defer conn.Close()

	user.Password = utils.GetPasswordHash(user.Password)
	now := time.Now()
	if user.Model == nil{
		user.Model = new(neo4j_driver.Model)
	}
	user.Model.CreatedAt = now
	user.Model.UpdatedAt = now
	if user.Model.UniqueId == uuid.Nil{
		var uniqueId, _ = uuid.NewV4()
		user.Model.UniqueId = uniqueId
	}
	conn.Save(user)
	return user
}


func (conn *DbConnector) GetUser(uniqueId uuid.UUID) (models.AuthUser){
	conn.Open(filerat.BoltPath, filerat.Port)
	defer conn.Close()
	user := models.AuthUser{}
	user.Model = new(neo4j_driver.Model)
	conn.Get(&user, uniqueId)

	return user
}


func (conn *DbConnector) DeleteUser(user *models.AuthUser, soft bool) bool {
	conn.Open(filerat.BoltPath, filerat.Port)
	defer conn.Close()

	conn.Delete(&user, user.UniqueId, soft)
	return true
}