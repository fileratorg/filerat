package db

import (
	"github.com/fileratorg/filerat/db/models"
	"github.com/fileratorg/filerat"
		"github.com/fileratorg/filerat/db/neo4j_driver"
		"github.com/satori/go.uuid"
)

func (conn *DbConnector) SaveFolder(obj *models.RatFolder) (*models.RatFolder){
	conn.Open(filerat.BoltPath, filerat.Port)
	defer conn.Close()
	obj.InitParent()
	conn.Save(obj)
	return obj
}


func (conn *DbConnector) GetFolder(uniqueId uuid.UUID) (models.RatFolder){
	conn.Open(filerat.BoltPath, filerat.Port)
	defer conn.Close()
	obj := models.RatFolder{}
	obj.Model = new(neo4j_driver.Model)
	conn.Get(&obj, uniqueId)

	return obj
}


func (conn *DbConnector) DeleteFolder(obj *models.RatFolder, soft bool) bool {
	conn.Open(filerat.BoltPath, filerat.Port)
	defer conn.Close()

	conn.Delete(&obj, obj.UniqueId, soft)
	return true
}

func (conn *DbConnector) SaveFile(obj *models.RatFile) (*models.RatFile){
	conn.Open(filerat.BoltPath, filerat.Port)
	defer conn.Close()
	obj.InitParent()
	conn.Save(obj)
	return obj
}


func (conn *DbConnector) GetFile(uniqueId uuid.UUID) (models.RatFile){
	conn.Open(filerat.BoltPath, filerat.Port)
	defer conn.Close()
	obj := models.RatFile{}
	obj.Model = new(neo4j_driver.Model)
	conn.Get(&obj, uniqueId)

	return obj
}


func (conn *DbConnector) DeleteFile(obj *models.RatFile, soft bool) bool {
	conn.Open(filerat.BoltPath, filerat.Port)
	defer conn.Close()

	conn.Delete(&obj, obj.UniqueId, soft)
	return true
}