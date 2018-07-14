package neo4j_driver

import (
	"testing"
	"fmt"
)

func TestConnector(t *testing.T) {
	boltPath := "bolt://neo4j:admin@0.0.0.0"
	port := 7687

	assertCorrectFieldAssignment := func(t *testing.T, expected, outcome string) {
		if expected != outcome {
			t.Errorf("Incorrect assignment: got %s and expected %s.", expected, outcome)
		}
	}

	t.Run("Test neo4j connector path", func(t *testing.T) {
		conn := new(Connector)
		defer conn.Close()
		conn.Open(boltPath, port)
		testUri := fmt.Sprintf("%s:%d", boltPath, port)
		assertCorrectFieldAssignment(t, testUri, conn.getFullPath())
	})
}
