package db

import (
	"testing"
	"github.com/fileratorg/filerat/db/models"
		"github.com/fileratorg/filerat/db/neo4j_driver"
	"github.com/satori/go.uuid"
	"time"
)

var uniqueIdOrg, _ = uuid.NewV4()

func TestRatOrgCreate(t *testing.T) {

	t.Run("Test rat org create", func(t *testing.T) {
		conn := new(DbConnector)
		now := time.Now()
		org := models.RatOrg{Name:"Test"}
		org.Model = new(neo4j_driver.Model)
		org.Model.CreatedAt = now
		org.Model.UpdatedAt = now
		org.Model.UniqueId = uniqueIdOrg
		conn.SaveRatOrg(&org)
	})
}

func TestRatOrgGet(t *testing.T) {

	t.Run("Test rat org get", func(t *testing.T) {
		conn := new(DbConnector)

		conn.GetRatOrg(uniqueIdOrg)
	})

}

func TestRatOrgUpdate(t *testing.T) {

	t.Run("Test rat org update", func(t *testing.T) {
		conn := new(DbConnector)

		org := conn.GetRatOrg(uniqueIdOrg)
		org.Name = "changed"
		conn.SaveRatOrg(&org)
		if org.Name != "changed" {
			t.Errorf("username not changed")
		}
	})
}

func TestRatOrgSoftDelete(t *testing.T) {

	t.Run("Test rat org soft delete", func(t *testing.T) {
		conn := new(DbConnector)

		org := conn.GetRatOrg(uniqueIdOrg)
		conn.DeleteRatOrg(&org, true)
	})
}

func TestRatOrgHardDelete(t *testing.T) {

	t.Run("Test rat org hard delete", func(t *testing.T) {
		conn := new(DbConnector)

		org := conn.GetRatOrg(uniqueIdOrg)
		conn.DeleteRatOrg(&org, false)
	})
}
