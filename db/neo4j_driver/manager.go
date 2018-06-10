package neo4j_driver

import (
	"fmt"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"reflect"
	"time"
	"log"
	"github.com/satori/go.uuid"
	"strings"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/structures/graph"
)

type DatabaseManager interface {
	Save() Model
}

type Connector struct {
	BoltPath 	string
	Port 		int
	Session 	bolt.Conn
	Error			error
}

func (conn *Connector) Open(path string, port int) *Connector {
	conn.BoltPath = path
	conn.Port = port
	session, err := start(conn.getFullPath())
	conn.Session = session
	conn.Error = err
	return conn
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

func (db *Connector)Delete(model interface{}, uniqueId uuid.UUID, soft bool) *Connector {
	// TODO: find a way to get away without passing uniqueId
	modelTypeName := getModelName(model)
	query := ""
	if soft {
		query = fmt.Sprintf("MATCH (n:%s { Id: {%s} } SET n.DeletedAt=%d", modelTypeName, uniqueId, time.Now().Unix())
	}else {
		query = fmt.Sprintf("MATCH (n:%s { Id: {%s} } DETACH DELETE n", modelTypeName, uniqueId)
	}
	stmt, err := db.Session.PrepareNeo(query)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}

	// Executing a statement just returns summary information
	result, err := stmt.ExecNeo(nil)
	if err != nil {
		panic(err)
	}
	numResult, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	log.Printf("DELETED ROWS: %d\n", numResult) // CREATED ROWS: 1
	db.Error = err

	return db
}

func (db *Connector)Get(model interface{}, uniqueId uuid.UUID) *Connector{
	modelTypeName := getModelName(model)

	query := fmt.Sprintf("MATCH (n:%s { UniqueId: {UniqueId} }) RETURN n", modelTypeName)
	params := map[string]interface{}{"UniqueId": uniqueId.String()}
	stmt, err := db.Session.PrepareNeo(query)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}
	rows, err := stmt.QueryNeo(params)
	if err != nil {
		panic(err)
	}
	data, _, err := rows.NextNeo()
	fmt.Printf("COLUMNS: %#v\n", rows.Metadata()["fields"].([]interface{})) // COLUMNS: n.foo,n.bar
	node := graph.Node{}
	reflect.ValueOf(&node).Elem().Set(reflect.ValueOf(data[0]))
	loadModel(model, node)

	db.Error = err
	return db
}

func (db *Connector)Save(model interface{}) *Connector{
	modelTypeName := getModelName(model)

	query := fmt.Sprintf("MERGE (n:%s {", modelTypeName)
	params := getModelFields(model)

	for field_name, _ := range params{
		query += fmt.Sprintf("%s: {%s},", field_name, field_name)
	}

	query = query[:len(query) - 1] + "})"
	stmt, err := db.Session.PrepareNeo(query)
	defer stmt.Close()
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
	db.Error = err
	return db
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

			sEmbedded := reflect.ValueOf(f.Interface()).Elem()
			typeOfTEmbedded := sEmbedded.Type()
			for z := 0; z < sEmbedded.NumField(); z++ {
				fEmbedded := sEmbedded.Field(z)
				field_name := typeOfTEmbedded.Field(z).Name
				//field_type := fEmbedded.Type()
				switch fEmbedded.Interface().(type) {
				case time.Time:
					var tmpVal time.Time
					reflect.ValueOf(&tmpVal).Elem().Set(reflect.ValueOf(fEmbedded.Interface()))
					fields[field_name] = tmpVal.Unix()
				case uuid.UUID:
					fields[field_name] = fmt.Sprintf("%s", fEmbedded.Interface())
				default:
					fields[field_name] = fEmbedded.Interface()
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

func loadModel(model interface{}, node graph.Node) interface{} {
	clone := reflect.ValueOf(model).Elem()
	for field_name, field_value := range node.Properties {
		v := clone.FieldByName(field_name)
		if v.IsValid() {
			var tmpVal time.Time
			switch v.Interface().(type) {
			case time.Time:
				if v.Interface() != nil{
					var val int64
					reflect.ValueOf(&val).Elem().Set(reflect.ValueOf(field_value))
					tmpVal = time.Unix(val, 0)
					v.Set(reflect.ValueOf(tmpVal).Convert(reflect.TypeOf(tmpVal)))
				}
			case uuid.UUID:
				var tmpVal uuid.UUID
				if v.Interface() != nil{
					var val string
					reflect.ValueOf(&val).Elem().Set(reflect.ValueOf(field_value))
					tmpVal, _ = uuid.FromString(val)
					v.Set(reflect.ValueOf(tmpVal))
				}
			default:
				if v.Interface() != nil {
					v = reflect.ValueOf(field_value)
				}
			}
		}
		clone.FieldByName(field_name).Set(v)
	}
	//reflect.ValueOf(&model).Set(reflect.ValueOf(clone))
	model = &clone
	return model
}

func indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

func getModelName(model interface{}) string {
	modelTypeName := reflect.TypeOf(model).String()
	modelTypeName = modelTypeName[strings.LastIndex(modelTypeName, ".") + 1:]
	return modelTypeName
}