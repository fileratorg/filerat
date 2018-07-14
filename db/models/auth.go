package models

import (
	"github.com/fileratorg/filerat/db/neo4j_driver"
	)

type AuthUser struct {
	*neo4j_driver.Model
	Username string
	Email string
	Password string
	//RatOrgId uuid.UUID
	RatOrg RatOrg `cypher:"relation_name:INSIDE_ORG"`
}

type RatOrg struct {
	*neo4j_driver.Model
	Name string
}