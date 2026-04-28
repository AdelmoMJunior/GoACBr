package storage

import (
	"context"
	"io"
)

// Provider interface defines the contract for file storage operations.
type Provider interface {
	// Upload stores a file and returns its unique key/path.
	Upload(ctx context.Context, key string, data io.Reader, contentType string) (string, error)
	
	// Download retrieves a file by its key.
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	
	// Delete removes a file by its key.
	Delete(ctx context.Context, key string) error
	
	// GetURL returns a presigned URL or public URL for the file (if applicable).
	GetURL(ctx context.Context, key string) (string, error)
}
