package filerat

import "github.com/fileratorg/filerat/utils"

var	BoltPath = utils.GetEnv("DB_URL", "bolt://neo4j:admin@0.0.0.0")
var	Port = 7687
