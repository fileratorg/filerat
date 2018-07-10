package db

import (
	"github.com/fileratorg/filerat/db/models"
	"github.com/fileratorg/filerat"
		"github.com/fileratorg/filerat/db/neo4j_driver"
		"github.com/satori/go.uuid"
)

func (conn *DbConnector) SaveUser(obj *models.AuthUser) (*models.AuthUser){
	conn.Open(filerat.BoltPath, filerat.Port)
	defer conn.Close()
	obj.InitParent()
	//if user.Model.UniqueId == uuid.Nil{
	//	var uniqueId, _ = uuid.NewV4()
	//	user.Model.UniqueId = uniqueId
	//}
	conn.Save(obj)
	return obj
}


func (conn *DbConnector) GetUser(uniqueId uuid.UUID) (models.AuthUser){
	conn.Open(filerat.BoltPath, filerat.Port)
	defer conn.Close()
	obj := models.AuthUser{}
	obj.Model = new(neo4j_driver.Model)
	conn.Get(&obj, uniqueId)

	return obj
}


func (conn *DbConnector) DeleteUser(obj *models.AuthUser, soft bool) bool {
	conn.Open(filerat.BoltPath, filerat.Port)
	defer conn.Close()

	conn.Delete(&obj, obj.UniqueId, soft)
	return true
}