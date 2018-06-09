package storage

import (
	"testing"
	"path/filepath"
)


//func TestMain(m *testing.M){
	// setup
//	result := m.Run()
	// tear down
//	os.Exit(result)
//}

func createClient(t *testing.T) Storage{
	client, err := New()

	if err != nil{
		t.Error("Failed to create client")
	}
	return client
}
func TestCanCreateClient(t *testing.T){
	createClient(t)
}

func TestCanCreateBucket(t *testing.T){
	client := createClient(t)
	bucketName := "test"
	if !client.MakeBucket(bucketName){
		t.Error("Error creating bucket")
	}
}

func TestCanPreSignPut(t *testing.T){
	client := createClient(t)

	bucketName := "test"
	objectName := "test"
	if client.PreSignPut(bucketName, objectName) == nil {
		t.Error("Error creating signature for put")
	}

	//t.Parallel() // can run in parallel
	//t.Fatal("msg") // stops execution
	//t.Error("msg")
	//if testing.Short(){
	//	// run is -short flag was passed
	//	t.Skip("Integration Test")
	//}

}

func TestCanPut(t *testing.T){
	client := createClient(t)
	bucketName := "test"
	filePath, path_err := filepath.Abs("../tmp/test.zip")
	if path_err != nil{
		t.Fatal("error geting file path")
	}
	contentType := "application/zip"
	_, err := client.Put(bucketName, filePath, contentType)
	if err != nil {
		t.Error("Error uploading file")
	}
}

func TestCanPreSignGet(t *testing.T){
	client := createClient(t)
	bucketName := "test"
	objectName := "test"
	if client.PreSignGet(bucketName, objectName) == nil {
		t.Error("Error creating signaure for get")
	}
}
