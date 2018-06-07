package storage

import (
	"github.com/minio/minio-go"
	"log"
	"net/url"
	"time"
	"fmt"
	"path/filepath"
)

func New() (Storage, error) {

	endpoint := "play.minio.io:9000"
	accessKeyID := "Q3AM3UQ867SPQQA43P2F"
	secretAccessKey := "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
	location := "us-west-2"
	useSSL := true

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}
	storageClient := Storage{*minioClient, location }

	return storageClient, err
}

func (client Storage) MakeBucket(bucketName string) bool {
	err := client.Client.MakeBucket(bucketName, client.Location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, err := client.Client.BucketExists(bucketName)
		if err == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	}
	log.Printf("Successfully created %s\n", bucketName)
	return true
}

func (client Storage) Put(bucketName, filePath, contentType string) (string, error) {

	// Upload the zip file
	objectName := filepath.Base(filePath)

	// Upload the zip file with FPutObject
	_, err := client.Client.FPutObject(
		bucketName,
		objectName,
		filePath,
		minio.PutObjectOptions{
			ContentType:contentType,
		})
	if err != nil {
		log.Fatalln(err)
	}

	return objectName, err
}

func (client Storage) PreSignGet(bucketName, objectName string) *url.URL {

	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", `attachment; filename=\"{objectName}\"`)

	// Generates a presigned url which expires in a day.
	presignedURL, err := client.Client.PresignedGetObject(bucketName, objectName, time.Second * 24 * 60 * 60, reqParams)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println("Successfully generated presigned URL", presignedURL)

	return presignedURL
}


func (client Storage) PreSignPut(bucketName, objectName string) *url.URL {

	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", `attachment; filename=\"{objectName}\"`)

	// Generates a presigned url which expires in a day.
	presignedURL, err := client.Client.PresignedPutObject(bucketName, objectName, time.Second * 24 * 60 * 60)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println("Successfully generated presigned URL", presignedURL)

	return presignedURL
}