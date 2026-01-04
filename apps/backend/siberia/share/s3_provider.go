package share

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// StorageProvider defines the interface for file storage backends
type StorageProvider interface {
	Upload(key string, data []byte, contentType string) (string, error)
}

// S3Provider implements StorageProvider for S3-compatible stores (MinIO, AWS)
type S3Provider struct {
	client     *minio.Client
	bucketName string
	publicURL  string
}

// NewS3Provider creates a new S3 provider
func NewS3Provider(endpoint, accessKeyID, secretAccessKey, bucketName string, useSSL bool) (*S3Provider, error) {
	// Initialize MinIO client object
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	return &S3Provider{
		client:     minioClient,
		bucketName: bucketName,
		// For local MinIO inside Docker, the browser accesses it via localhost:9000
		// But in prod this might be different.
		// We'll construct the link based on the endpoint for now.
		publicURL: fmt.Sprintf("http://%s/%s", endpoint, bucketName),
	}, nil
}

// Upload uploads data to the configured bucket and returns a public link
func (p *S3Provider) Upload(key string, data []byte, contentType string) (string, error) {
	ctx := context.Background()

	// Ensure bucket exists (MVP convenience)
	exists, err := p.client.BucketExists(ctx, p.bucketName)
	if err != nil {
		return "", fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		// Try to create it
		err = p.client.MakeBucket(ctx, p.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return "", fmt.Errorf("failed to create bucket: %w", err)
		}
		// Set public policy (read-only for everyone) - Basic version
		policy := fmt.Sprintf(`{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::%s/*"]}]}`, p.bucketName)
		if err := p.client.SetBucketPolicy(ctx, p.bucketName, policy); err != nil {
			log.Printf("Warning: failed to set bucket policy: %v", err)
		}
	}

	// Upload object
	reader := bytes.NewReader(data)
	objectSize := int64(len(data))

	_, err = p.client.PutObject(ctx, p.bucketName, key, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object: %w", err)
	}

	// Construct Link
	link := fmt.Sprintf("%s/%s", p.publicURL, key)
	return link, nil
}
