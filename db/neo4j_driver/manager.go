package neo4j_driver

import (
	"fmt"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"reflect"
	"time"
	"log"
)

type DatabaseManager interface {
	Save() Model
}

type Connector struct {
	BoltPath 	string
	Port 		int
	Session bolt.Conn
}

func (conn *Connector) Open(path string, port int) (bolt.Conn, error) {
	conn.BoltPath = path
	conn.Port = port
	session, ex := start(conn.getFullPath())
	conn.Session = session
	return conn.Session, ex
}
func (conn *Connector) Close() {
	conn.Session.Close()
}

func (conn *Connector) getFullPath() string {
	uri := fmt.Sprintf("%s:%d", conn.BoltPath,  conn.Port)
	return uri
}

func start(uri string) (bolt.Conn, error) {
	driver := bolt.NewDriver()

	conn, err := driver.OpenNeo(uri)
	if err != nil {
		panic(err)
	}
	return conn, err
}

func Save(conn bolt.Conn, labelName string, model interface{}) interface{}{

	query := fmt.Sprintf("CREATE (n:%s {", labelName)
	params := getModelFields(model)

	for field_name, _ := range params{
		query += fmt.Sprintf("%s: {%s},", field_name, field_name)
	}

	query = query[:len(query) - 1] + "})"
	stmt, err := conn.PrepareNeo(query)
	if err != nil {
		panic(err)
	}

	// Executing a statement just returns summary information
	result, err := stmt.ExecNeo(params)
	if err != nil {
		panic(err)
	}
	numResult, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	log.Printf("CREATED ROWS: %d\n", numResult) // CREATED ROWS: 1

	// Closing the statment will also close the rows
	stmt.Close()

	return model
}

//TODO: clean this up
func getModelFields(model interface{}) map[string]interface {} {

	s := reflect.ValueOf(model).Elem()
	typeOfT := s.Type()
	fields := make(map[string]interface {})
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		field_name := typeOfT.Field(i).Name
		if field_name == "Model"{

			sEmbeded := reflect.ValueOf(f.Interface()).Elem()
			typeOfTEmbeded := sEmbeded.Type()
			for z := 0; z < sEmbeded.NumField(); z++ {
				fEmbeded := sEmbeded.Field(z)
				field_name := typeOfTEmbeded.Field(z).Name
				field_type := fEmbeded.Type()
				log.Println(field_type)
				switch ftype := fEmbeded.Interface().(type) {
				case time.Time:
					log.Println(ftype)
					var test time.Time
					reflect.ValueOf(&test).Elem().Set(reflect.ValueOf(fEmbeded.Interface()))
					fields[field_name] = test.Unix()
					log.Println(test)
				default:
					fields[field_name] = fEmbeded.Interface()
				}
			}

		}else{
			//field_type := f.Type()
			field_value := f.Interface()
			fields[field_name] = field_value
		}
	}
	return fields
}