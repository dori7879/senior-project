package fs

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dori7879/senior-project/api"
)

// Ensure service implements interface.
var _ api.FileService = (*FileService)(nil)

// FileService represents a service for managing files.
type FileService struct {
	baseDir string
	HashKey []byte
}

// NewFileService returns a new instance of FileService.
func NewFileService(baseDir string, hashKey []byte) *FileService {
	return &FileService{baseDir: baseDir, HashKey: hashKey}
}

// GetPathByHash retrieves the filepath of an image by hash.
func (f *FileService) GetPathByHash(ctx context.Context, hash, extension string) string {
	return getPathByHash(ctx, f.baseDir, hash, extension)
}

// CreateFile creates a new image file.
func (f *FileService) CreateFile(ctx context.Context, header *multipart.FileHeader) (string, string, error) {
	hash, extension, err := createFile(ctx, f.baseDir, f.HashKey, header)
	if err != nil {
		return "", "", err
	}

	return hash, extension, nil
}

// GetFileHash computes a hash of a file.
func (f *FileService) GetFileHash(ctx context.Context, header *multipart.FileHeader) (string, error) {
	hash, err := hashFile(f.HashKey, header)
	if err != nil {
		return "", api.Errorf(api.EINTERNAL, "Could not hash file")
	}
	return hash, nil
}

// DeleteFile permanently deletes a image file.
// Returns EUNAUTHORIZED if current image is not the image being deleted.
// Returns ENOTFOUND if image does not exist.
func (f *FileService) DeleteFile(ctx context.Context, hash, extension string) error {
	if err := deleteFile(ctx, f.baseDir, hash, extension); err != nil {
		return err
	}
	return nil
}

// getPathByHash is a helper function to fetch a image by hash.
func getPathByHash(ctx context.Context, baseDir, hash, extension string) string {
	baseImageDir := filepath.Join(baseDir, "media", "files")
	filename := hash + "." + extension

	imgURL := filepath.Join(baseImageDir, hash[:2], hash[2:4], filename)

	return imgURL
}

func hashFile(hashKey []byte, header *multipart.FileHeader) (string, error) {
	f, err := header.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := hmac.New(sha1.New, hashKey)
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	hBytes := h.Sum(nil)[:20]

	return hex.EncodeToString(hBytes), nil
}

func createDirs(hash, baseDir string) (string, error) {
	baseImageDir := filepath.Join(baseDir, "media", "files")

	imgDir := filepath.Join(baseImageDir, hash[:2], hash[2:4])

	err := os.MkdirAll(imgDir, 0766)
	if err != nil {
		return "", err
	}

	return imgDir, nil
}

func getFileContentType(f *multipart.File) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	bytesRead, err := (*f).Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	// Slice to remove fill-up zero values which cause a wrong content type detection in the next step
	buffer = buffer[:bytesRead]

	return http.DetectContentType(buffer), nil
}

// createFile creates a new image file.
func createFile(ctx context.Context, baseDir string, hashKey []byte, header *multipart.FileHeader) (string, string, error) {
	hash, err := hashFile(hashKey, header)
	if err != nil {
		return "", "", api.Errorf(api.EINTERNAL, "Could not hash file")
	}

	extension := strings.Split(header.Filename, ".")[1]

	filename := hash + "." + extension
	imgDir, err := createDirs(hash, baseDir)
	if err != nil {
		return "", "", api.Errorf(api.EINTERNAL, "Could not create directory for new image file")
	}
	imgPath := filepath.Join(imgDir, filename)

	f, err := os.OpenFile(imgPath, os.O_WRONLY|os.O_CREATE, 0766)
	if err != nil {
		return "", "", api.Errorf(api.EINTERNAL, "Could not open blank image file")
	}
	defer f.Close()

	file, err := header.Open()
	if err != nil {
		return "", "", api.Errorf(api.EINTERNAL, "Could not open new header image file")
	}
	defer file.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		return "", "", api.Errorf(api.EINTERNAL, "Could not copy header file to blank file")
	}

	return hash, extension, nil
}

// deleteFile permanently removes a image file by hash. Returns EUNAUTHORIZED if current
// image file is not the one being deleted.
func deleteFile(ctx context.Context, baseDir, hash, extension string) error {
	// Verify object exists.
	filepath := getPathByHash(ctx, baseDir, hash, extension)

	return os.Remove(filepath)
}
