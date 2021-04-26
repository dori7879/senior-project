package api

import (
	"context"
	"mime/multipart"
	"time"
)

// All constants.
const (
	MaxHashLen      = 40
	MaxExtensionLen = 10
)

// Attachment is a set of fields that are mapped to the database table attachments
type Attachment struct {
	Hash      string
	Extension string
	CreatedAt time.Time
	Counter   int
}

// Validate returns an error if the attachment contains invalid fields.
// This only performs basic validation.
func (i *Attachment) Validate() error {
	if i.Hash == "" {
		return Errorf(EINVALID, "Hash required.")
	} else if len(i.Hash) > MaxHashLen {
		return Errorf(EINVALID, "Hash too long.")
	} else if len(i.Extension) > MaxExtensionLen {
		return Errorf(EINVALID, "Extension too long.")
	} else if i.Counter < 0 {
		return Errorf(EINVALID, "Counter could not be less than 0.")
	}

	return nil
}

// AttachmentService represents a service for managing attachments.
type AttachmentService interface {
	// Retrieves a attachment by ID.
	// Returns ENOTFOUND if attachment does not exist.
	FindAttachmentByHash(ctx context.Context, hash string) (*Attachment, error)

	// Retrieves a list of attachments by filter. Also returns total count of matching
	// attachments which may differ from returned results if filter.Limit is specified.
	FindAttachments(ctx context.Context, filter AttachmentFilter) ([]*Attachment, int, error)

	// Creates a new attachment.
	CreateAttachment(ctx context.Context, attachment *Attachment) error

	// Updates a attachment object. Returns EUNAUTHORIZED if current attachment is not
	// the attachment that is being updated. Returns ENOTFOUND if attachment does not exist.
	UpdateAttachment(ctx context.Context, hash string, upd AttachmentUpdate) (*Attachment, error)

	// Permanently deletes a attachment. Returns EUNAUTHORIZED
	// if current attachment is not the attachment being deleted. Returns ENOTFOUND if
	// attachment does not exist.
	DeleteAttachment(ctx context.Context, hash string) error
}

// AttachmentFilter represents a filter passed to FindAttachments().
type AttachmentFilter struct {
	// Filtering fields.
	Hash *string `json:"hash"`

	// Restrict to subset of results.
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// AttachmentUpdate represents a set of fields to be updated via UpdateAttachment().
type AttachmentUpdate struct {
	Counter *int
}

// FileService represents a service for managing files, usually related to AttachmentService.
type FileService interface {
	GetPathByHash(ctx context.Context, hash, extension string) string

	CreateFile(ctx context.Context, header *multipart.FileHeader) (string, string, error)

	GetFileHash(ctx context.Context, header *multipart.FileHeader) (string, error)

	DeleteFile(ctx context.Context, hash, extension string) error
}

// AttachmentUploader represents a combination of AttachmentService and FileService that manages
// file creation and attachments database table
type AttachmentUploader struct {
	AttachmentService AttachmentService
	FileService       FileService
}

// UploadAttachment creates the file and adds related row in attachments table
func (u *AttachmentUploader) UploadAttachment(ctx context.Context, header *multipart.FileHeader) (*Attachment, error) {
	hash, extension, err := u.FileService.CreateFile(ctx, header)
	if err != nil {
		return nil, err
	}

	img := Attachment{Hash: hash, Extension: extension, Counter: 0}

	err = u.AttachmentService.CreateAttachment(ctx, &img)
	if err != nil {
		return nil, err
	}
	return &img, nil
}

// DecrementAttachmentCounter updates attachment and if the counter is 0 then deletes the file
func (u *AttachmentUploader) DecrementAttachmentCounter(ctx context.Context, hash string) error {
	oldImg, err := u.AttachmentService.FindAttachmentByHash(ctx, hash)
	if err != nil {
		return err
	}
	newCounter := oldImg.Counter - 1

	// If counter is less than 1, then UpdateAttachment will remove corresponding row
	img, err := u.AttachmentService.UpdateAttachment(ctx, hash, AttachmentUpdate{Counter: &newCounter})
	if err != nil {
		return err
	}

	if img.Counter < 1 {
		err = u.FileService.DeleteFile(ctx, img.Hash, img.Extension)
		if err != nil {
			return err
		}
	}

	return nil
}
