package googlebucket

import (
	"context"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/storage"
)

// StoreObject stores object in a Google bucket
func StoreObject(filePath, objectName string) bool {
	googleCredentials := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	if googleCredentials == "" {
		fmt.Println("GOOGLE_APPLICATION_CREDENTIALS not found")

		return false
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	defer client.Close()

	if err != nil {
		fmt.Println(err)

		return false
	}

	if err := write(client, "fellrace-finder", filePath, objectName); err != nil {
		fmt.Println("Cannot write to bucket: "+filePath, err)

		return false
	}

	return true
}

func write(client *storage.Client, bucket, filePath, object string) error {
	ctx := context.Background()

	// [START upload_file]
	f, err := os.Open(filePath)

	if err != nil {
		fmt.Println(err)

		return err
	}
	defer f.Close()

	wc := client.Bucket(bucket).Object("maps/" + object).NewWriter(ctx)

	if _, err = io.Copy(wc, f); err != nil {
		fmt.Println(err)

		return err
	}

	if err := wc.Close(); err != nil {
		fmt.Println(err)

		return err
	}

	// fmt.Println("Written " + object)

	// [END upload_file]
	return nil
}
