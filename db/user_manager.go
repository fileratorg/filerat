package db

import (
	"github.com/fileratorg/filerat/db/models"
	"github.com/fileratorg/filerat"
	"github.com/fileratorg/filerat/utils"
	"github.com/fileratorg/filerat/db/neo4j_driver"
	"time"
)

func (conn *DbConnector) Create(user *models.AuthUser) (*models.AuthUser){
	session, _ := conn.Open(filerat.BoltPath, filerat.Port)
	defer conn.Close()

	user.Password = utils.GetPasswordHash(user.Password)
	now := time.Now()
	user.Model = new(neo4j_driver.Model)
	user.Model.CreatedAt = now
	user.Model.UpdatedAt = now
	neo4j_driver.Save(session, "AuthUser", user)
	return user
}
