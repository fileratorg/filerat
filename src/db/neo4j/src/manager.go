package neo4j

import (
	"fmt"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

type databaseConnection interface {
	start()
	stop()
	create()
	read()
	update()
	delete()
}

type Connector struct {
	BoltPath 	string
	Port 		int

}

func (conn *Connector) Init(path string, port int) {
	start(conn.getFullPath())
	conn.BoltPath = path
	conn.Port = port
	return
}

func (conn *Connector) getFullPath() string {
	uri := fmt.Sprintf("%s:%d", conn.BoltPath,  conn.Port)
	return uri
}

func start(uri string) {
	driver := golangNeo4jBoltDriver.NewDriver()
	fmt.Println(&driver)

	conn, err := driver.OpenNeo(uri)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

}
