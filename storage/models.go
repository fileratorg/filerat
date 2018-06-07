package storage

import "github.com/minio/minio-go"

type Storage struct {
	Client minio.Client
	Location string
}

type StorageManager interface {
	New()
	PreSignPut()
	PreSignGet()
}
