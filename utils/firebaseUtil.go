package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/storage"
	"github.com/Crade47/medium-blog-scraper/models"
	"google.golang.org/api/option"
)

type FirestoreClient struct {
	fsClient *firestore.Client
}

type StorageClient struct {
	storageClient *storage.Client
}

var firestoreClient *FirestoreClient
var storageClient *StorageClient

func Initialize(ctx context.Context, firebaseConfigPath string) error {
	config := &firebase.Config{
		StorageBucket: "go-scrape-test.appspot.com",
	}

	opt := option.WithCredentialsFile(firebaseConfigPath)
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return err
	}

	//------------------------------CLIENT OBJECTS INIT------------------------------
	firestoreClient = &FirestoreClient{}
	storageClient = &StorageClient{}

	//------------------------------ FIRESTORE CLIENT------------------------------
	firestoreClient.fsClient, err = app.Firestore(ctx)
	if err != nil {
		return err
	}

	//------------------------------ STORAGE CLIENT AND BUCKET------------------------------
	storageClient.storageClient, err = app.Storage(ctx)
	if err != nil {
		return err
	}

	return nil
}

// ------------------------------Returning the client objects------------------------------
func GetFirestoreClient() *FirestoreClient {
	return firestoreClient
}

func GetFirebaseStorage() *StorageClient {
	return storageClient
}

// ------------------------------ADD DOCUMENT METHOD IN FIRESTORE CLIENT------------------------------
func (fc *FirestoreClient) AddDocument(ctx context.Context, id string, data *models.Blog) {
	_, err := fc.fsClient.Collection("blogs").Doc(id).Set(ctx, data)
	if err != nil {
		log.Fatalf("Failed document addition")
	}
}

func (fc *FirestoreClient) GetDocument(ctx context.Context, id string) (map[string]interface{}, error) {
	dsnap, err := fc.fsClient.Collection("blogs").Doc(id).Get(ctx)
	if err != nil {
		fmt.Println("Failed to get document from firestore")
	}
	data := dsnap.Data()

	return data, nil
}

//  ------------------------------UPLOAD METHOD IN STORAGE------------------------------

func (s *StorageClient) UploadFile(ctx context.Context, localFilePath string, destFileName string) (string, error) {
	bucket, err := s.storageClient.DefaultBucket()
	if err != nil {
		return "", err
	}

	file, err := os.Open(localFilePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	obj := bucket.Object(destFileName)

	w := obj.NewWriter(ctx)
	w.ContentType = "text/html"
	w.ObjectAttrs.ContentDisposition = "attachment; filename=\"" + destFileName + "\""

	if _, err := io.Copy(w, file); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}
	// attrs, err := obj.Attrs(ctx)
	// if err != nil {
	// 	return "", err
	// }
	downloadUrl := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media", "go-scrape-test.appspot.com", destFileName)
	return downloadUrl, nil
}
