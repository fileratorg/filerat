package neo4j

import (
	"testing"
	"fmt"
)

func TestConnector(t *testing.T) {

	boltPath := "localhost"
	port := 7687

	assertCorrectFieldAssignment := func(t *testing.T, expected, outcome string) {
		if expected != outcome {
			t.Errorf("Incorrect assignment: got %s and expected %s.", expected, outcome)
		}
	}

	t.Run("Verify neo4j connector path", func(t *testing.T) {
		conn := new(Connector)
		conn.Init(boltPath, port)
		testUri := fmt.Sprintf("%s:%d", boltPath, port)
		assertCorrectFieldAssignment(t, testUri, conn.getFullPath())
	})

}
