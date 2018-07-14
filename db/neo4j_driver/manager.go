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
	Error		error
}

type Relation struct {
	FromNode	interface{}
	ToNode		interface{}
	RelName		string
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
	params := map[string]interface{}{
		"UniqueId": uniqueId.String(),
	}
	if soft {
		query = fmt.Sprintf("MATCH (n:%s { UniqueId: {UniqueId} }) SET n.DeletedAt={DeletedAt}", modelTypeName)
		params["DeletedAt"] = time.Now().Unix()
	}else {
		query = fmt.Sprintf("MATCH (n:%s { UniqueId: {UniqueId} }) DETACH DELETE n", modelTypeName)
	}
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

	_, relations := db.SaveNode(model)
	for _, relation := range relations {
		db.SaveRelation(relation)
	}
	return db
}

func (db *Connector)SaveNode(model interface{}) (*Connector, []Relation){
	modelTypeName := getModelName(model)

	query := fmt.Sprintf("MERGE (n:%s {", modelTypeName)
	params, uniqueId, relations := getModelFields(model)

	query += fmt.Sprintf("%s: \"%s\"}) SET ", "UniqueId", uniqueId.String())

	for field_name, _ := range params{
		query += fmt.Sprintf(" n.%s = {%s},", field_name, field_name)
	}

	query = query[:len(query) - 1]
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

	return db, relations
}

//func (db *Connector)SaveRelation(model_from interface{}, model_to interface{}, relation_name string) *Connector{
func (db *Connector)SaveRelation(relation Relation) *Connector{
	if relation.FromNode == nil || relation.ToNode == nil {
		return db
	}
	modelSourceTypeName := getModelName(relation.FromNode)
	modelDestinyTypeName := getModelName(relation.ToNode)

	uniqueIdSource := getUniqueId(relation.FromNode)
	uniqueIdDestiny := getUniqueId(relation.ToNode)
	if uniqueIdSource == uuid.Nil || uniqueIdDestiny == uuid.Nil{
		return db
	}
	params := map[string]interface{}{
		"UniqueIdSource": uniqueIdSource.String(),
		"UniqueIdDestiny": uniqueIdDestiny.String(),
	}

	query := fmt.Sprintf(
		"MATCH (s:%s {UniqueId: {UniqueIdSource}}), (d:%s {UniqueId: {UniqueIdDestiny}})",
		modelSourceTypeName, modelDestinyTypeName)
	query += fmt.Sprintf(" MERGE (s)-[r:%s]->(d)", relation.RelName)

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
	log.Printf("CREATED RELATIONS: %d\n", numResult) // CREATED ROWS: 1
	db.Error = err
	return db
}


//TODO: clean this up
func getModelFields(model interface{}) (map[string]interface {}, uuid.UUID, []Relation) {
	var uniqueId uuid.UUID
	s := reflect.ValueOf(model).Elem()
	typeOfT := s.Type()
	fields := make(map[string]interface {})
	var relations []Relation
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fieldName := typeOfT.Field(i).Name
		if fieldName == "Model"{

			sEmbedded := reflect.ValueOf(f.Interface()).Elem()
			typeOfTEmbedded := sEmbedded.Type()
			for z := 0; z < sEmbedded.NumField(); z++ {
				fEmbedded := sEmbedded.Field(z)
				field_type := typeOfTEmbedded.Field(z).Tag.Get("type")
				fieldName := typeOfTEmbedded.Field(z).Name

				if field_type == "unique_id" {
					reflect.ValueOf(&uniqueId).Elem().Set(reflect.ValueOf(fEmbedded.Interface()))
					continue
				}
				if field_type == "id" {
					continue
				}
				//field_type := fEmbedded.Type()
				switch fEmbedded.Interface().(type) {
				case time.Time:
					var tmpVal time.Time
					reflect.ValueOf(&tmpVal).Elem().Set(reflect.ValueOf(fEmbedded.Interface()))
					fields[fieldName] = tmpVal.Unix()
				case uuid.UUID:
					fields[fieldName] = fmt.Sprintf("%s", fEmbedded.Interface())
				default:
					fields[fieldName] = fEmbedded.Interface()
				}
			}

		}else{
			//field_type := f.Type()

			fieldMetaType := typeOfT.Field(i).Tag.Get("cypher")
			fieldMetaCypher := typeOfT.Field(i).Tag.Get("cypher")
			fieldName := typeOfT.Field(i).Name
			fieldValue := reflect.ValueOf(f.Interface())

			if fieldMetaType == "unique_id" {
				reflect.ValueOf(&uniqueId).Elem().Set(fieldValue)
				continue
			}
			if fieldMetaType == "id" {
				continue
			}
			if len(fieldMetaCypher) > 0 {
				// model_from interface{}, model_to interface{}, relation_name string
				var fromNode = model
				var toNode = f.Interface()
				if toNode == nil {
					continue
				}
				relName := fieldMetaCypher[strings.Index(fieldMetaCypher, ":") + 1:]
				relation := Relation{
					FromNode:fromNode,
					ToNode:toNode,
					RelName:relName,
				}
				relations = append(relations, relation)
				continue
			}

			//field_type := fEmbedded.Type()
			switch f.Interface().(type) {
			case time.Time:
				var tmpVal time.Time
				reflect.ValueOf(&tmpVal).Elem().Set(reflect.ValueOf(f.Interface()))
				fields[fieldName] = tmpVal.Unix()
			case uuid.UUID:
				fields[fieldName] = fmt.Sprintf("%s", f.Interface())
			default:
				fields[fieldName] = f.Interface()
			}


			//field_value := f.Interface()
			//fields[field_name] = field_value
		}
	}
	return fields, uniqueId, relations
}

func getModelRelations(model interface{}) (map[string]interface {}, uuid.UUID) {
	var uniqueId uuid.UUID
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
				field_type := typeOfTEmbedded.Field(z).Tag.Get("type")
				field_name := typeOfTEmbedded.Field(z).Name

				if field_type == "unique_id" {
					reflect.ValueOf(&uniqueId).Elem().Set(reflect.ValueOf(fEmbedded.Interface()))
					continue
				}
				if field_type == "id" {
					continue
				}
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
			field_type := typeOfT.Field(i).Tag.Get("cypher")
			log.Println(field_type)
			//field_type := f.Type()
			field_value := f.Interface()
			fields[field_name] = field_value
		}
	}
	return fields, uniqueId
}

func getUniqueId(model interface{}) uuid.UUID {
	var uniqueId uuid.UUID
	s := indirect(reflect.ValueOf(model))
	if !s.CanAddr() {
		return uuid.Nil
	}
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		field_name := typeOfT.Field(i).Name
		if field_name == "Model"{
			sEmbedded := reflect.ValueOf(f.Interface()).Elem()
			typeOfTEmbedded := sEmbedded.Type()
			for z := 0; z < sEmbedded.NumField(); z++ {
				fEmbedded := sEmbedded.Field(z)
				field_type := typeOfTEmbedded.Field(z).Tag.Get("type")

				if field_type == "unique_id" {
					reflect.ValueOf(&uniqueId).Elem().Set(reflect.ValueOf(fEmbedded.Interface()))
					return uniqueId
				}
			}
		}
	}
	return uniqueId
}

func loadModel(model interface{}, node graph.Node) interface{} {
	clone := reflect.ValueOf(model).Elem()
	clone.FieldByName("Id").Set(reflect.ValueOf(node.NodeIdentity))
	for field_name, field_value := range node.Properties {
		if field_name == "Id" {
			//TODO: maybe use struct field tags for this
			continue
		}
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
