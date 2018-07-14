package models

import (
	"github.com/fileratorg/filerat/db/neo4j_driver"
		)

type RatFile struct {
	*neo4j_driver.Model
	Name		string
	Path		string
	Owner		*AuthUser `cypher:"relation_name:BELONGS_TO"`
	// RatFolderId uuid.UUID //`sql:"unique_index:uix_multipleindexes_user_name,uix_multipleindexes_user_email;index:idx_multipleindexes_user_other"`
	RatFolder	*RatFolder `cypher:"relation_name:INSIDE_FOLDER"`
}

type RatFolder struct {
	*neo4j_driver.Model
	Name		string
	Path		string
	Owner		*AuthUser `cypher:"relation_name:BELONGS_TO"`
	RatACL		RatACL `cypher:"relation_name:ACL"`
	RatFolder	*RatFolder `cypher:"relation_name:INSIDE_FOLDER"`
}

type RatACL struct {
	*neo4j_driver.Model
	IsPrivate			bool
	ProtectedShareType	int64
	PublicShareType		int64
}

type RatShare struct {
	IsRecursive bool
	ShareType	int64 `field:"options:read,write"`
	ShareUrl	int64

}