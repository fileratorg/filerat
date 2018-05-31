package storage

//import (
//	"github.com/minio/minio-go"
//	"log"
//	"fmt"
//	"net/url"
//	"time"
//)
//
//func storage() {
//	endpoint := "play.minio.io:9000"
//	accessKeyID := "Q3AM3UQ867SPQQA43P2F"
//	secretAccessKey := "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
//	useSSL := true
//
//	// Initialize minio client object.
//	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	// Make a new bucket called mymusic.
//	bucketName := "filerat"
//	location := "us-east-1"
//
//	err = minioClient.MakeBucket(bucketName, location)
//	if err != nil {
//		// Check to see if we already own this bucket (which happens if you run this twice)
//		exists, err := minioClient.BucketExists(bucketName)
//		if err == nil && exists {
//			log.Printf("We already own %s\n", bucketName)
//		} else {
//			log.Fatalln(err)
//		}
//	}
//	log.Printf("Successfully created %s\n", bucketName)
//
//	// Upload the zip file
//	objectName := "test.zip"
//	filePath := "./tmp/test.zip"
//	contentType := "application/zip"
//
//	// Upload the zip file with FPutObject
//	n, err := minioClient.FPutObject(
//		bucketName,
//		objectName,
//		filePath,
//		minio.PutObjectOptions{
//			ContentType:contentType,
//		})
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)
//
//
//	// Create a done channel to control 'ListObjects' go routine.
//	doneCh := make(chan struct{})
//
//	// Indicate to our routine to exit cleanly upon return.
//	defer close(doneCh)
//
//	// List all objects from a bucket-name with a matching prefix.
//	for object := range minioClient.ListObjects(bucketName, "", true, doneCh) {
//		if object.Err != nil {
//			fmt.Println(object.Err)
//			return
//		}
//		fmt.Println(object)
//	}
//	TestPresignGet(bucketName, objectName)
//
//}
//
//func TestPresign(){
//	endpoint := "play.minio.io:9000"
//	accessKeyID := "Q3AM3UQ867SPQQA43P2F"
//	secretAccessKey := "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
//	useSSL := true
//
//	// Initialize minio client object.
//	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	// Set request parameters for content-disposition.
//	reqParams := make(url.Values)
//	reqParams.Set("response-content-disposition", "attachment; filename=\"your-filename.txt\"")
//
//	// Generates a presigned url which expires in a day.
//	presignedURL, err := minioClient.PresignedGetObject("mybucket", "myobject", time.Second * 24 * 60 * 60, reqParams)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println("Successfully generated presigned URL", presignedURL)
//}
//
//func TestPresignGet(bucketName, objectName string){
//
//	endpoint := "play.minio.io:9000"
//	accessKeyID := "Q3AM3UQ867SPQQA43P2F"
//	secretAccessKey := "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
//	useSSL := true
//
//	// Initialize minio client object.
//	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
//
//	// Set request parameters for content-disposition.
//	reqParams := make(url.Values)
//	reqParams.Set("response-content-disposition", `attachment; filename=\"{objectName}\"`)
//
//	// Generates a presigned url which expires in a day.
//	presignedURL, err := minioClient.PresignedGetObject(bucketName, objectName, time.Second * 24 * 60 * 60, reqParams)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println("Successfully generated presigned URL", presignedURL)
//}