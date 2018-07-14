package db

import (
	"github.com/fileratorg/filerat/db/models"
	"github.com/fileratorg/filerat"
		"github.com/fileratorg/filerat/db/neo4j_driver"
		"github.com/satori/go.uuid"
)

func (conn *DbConnector) SaveRatOrg(obj *models.RatOrg) (*models.RatOrg){
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


func (conn *DbConnector) GetRatOrg(uniqueId uuid.UUID) (models.RatOrg){
	conn.Open(filerat.BoltPath, filerat.Port)
	defer conn.Close()
	obj := models.RatOrg{}
	obj.Model = new(neo4j_driver.Model)
	conn.Get(&obj, uniqueId)

	return obj
}


func (conn *DbConnector) DeleteRatOrg(obj *models.RatOrg, soft bool) bool {
	conn.Open(filerat.BoltPath, filerat.Port)
	defer conn.Close()

	conn.Delete(&obj, obj.UniqueId, soft)
	return true
}