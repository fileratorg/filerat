package models

import "github.com/fileratorg/filerat/db/neo4j_driver"

type AuthUser struct {
	*neo4j_driver.Model
	Username string
	Email string
	Password string
}
