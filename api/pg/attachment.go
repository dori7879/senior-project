package pg

import (
	"context"
	"strings"

	"github.com/dori7879/senior-project/api"
)

// Ensure service implements interface.
var _ api.AttachmentService = (*AttachmentService)(nil)

// AttachmentService represents a service for managing attachments.
type AttachmentService struct {
	db *DB
}

// NewAttachmentService returns a new instance of AttachmentService.
func NewAttachmentService(db *DB) *AttachmentService {
	return &AttachmentService{db: db}
}

// FindAttachmentByHash retrieves a attachment by hash.
// Returns ENOTFOUND if attachment does not exist.
func (s *AttachmentService) FindAttachmentByHash(ctx context.Context, hash string) (*api.Attachment, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Fetch attachment.
	attachment, err := findAttachmentByHash(ctx, tx, hash)
	if err != nil {
		return nil, err
	}
	return attachment, nil
}

// FindAttachments retrieves a list of attachments by filter. Also returns total count of
// matching attachments which may differ from returned results if filter.Limit is specified.
func (s *AttachmentService) FindAttachments(ctx context.Context, filter api.AttachmentFilter) ([]*api.Attachment, int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()
	return findAttachments(ctx, tx, filter)
}

// CreateAttachment creates a new attachment.
func (s *AttachmentService) CreateAttachment(ctx context.Context, attachment *api.Attachment) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create a new attachment object and attach associated OAuth objects.
	if err := createAttachment(ctx, tx, attachment); err != nil {
		return err
	}
	return tx.Commit()
}

// UpdateAttachment updates a attachment object. Returns EUNAUTHORIZED if current attachment is
// not the attachment that is being updated. Returns ENOTFOUND if attachment does not exist.
func (s *AttachmentService) UpdateAttachment(ctx context.Context, hash string, upd api.AttachmentUpdate) (*api.Attachment, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Update attachment.
	attachment, err := updateAttachment(ctx, tx, hash, upd)
	if err != nil {
		return attachment, err
	} else if err := tx.Commit(); err != nil {
		return attachment, err
	}
	return attachment, nil
}

// DeleteAttachment permanently deletes a attachment.
// Returns EUNAUTHORIZED if current attachment is not the attachment being deleted.
// Returns ENOTFOUND if attachment does not exist.
func (s *AttachmentService) DeleteAttachment(ctx context.Context, hash string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Verify object exists.
	if _, err := findAttachmentByHash(ctx, tx, hash); err != nil {
		return err
	}

	if err := deleteAttachment(ctx, tx, hash); err != nil {
		return err
	}
	return tx.Commit()
}

// findAttachmentByHash is a helper function to fetch a attachment by hash.
// Returns ENOTFOUND if attachment does not exist.
func findAttachmentByHash(ctx context.Context, tx *Tx, hash string) (*api.Attachment, error) {
	a, _, err := findAttachments(ctx, tx, api.AttachmentFilter{Hash: &hash})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &api.Error{Code: api.ENOTFOUND, Message: "Выбранная картинка не найдена."}
	}
	return a[0], nil
}

// findAttachments returns a list of attachments matching a filter. Also returns a count of
// total matching attachments which may differ if filter.Limit is set.
func findAttachments(ctx context.Context, tx *Tx, filter api.AttachmentFilter) (_ []*api.Attachment, n int, err error) {
	// Build WHERE clause.
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.Hash; v != nil {
		where, args = append(where, "hash = ?"), append(args, *v)
	}

	// Execute query to fetch attachment rows.
	rows, err := tx.QueryContext(ctx, `
		SELECT 
		    hash,
		    extension,
		    created_at,
		    counter,
		    COUNT(*) OVER()
		FROM attachments
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY created_at DESC
		`+FormatLimitOffset(filter.Limit, filter.Offset),
		args...,
	)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	// Deserialize rows into Attachment objects.
	attachments := make([]*api.Attachment, 0)
	for rows.Next() {
		var attachment api.Attachment
		if err := rows.Scan(
			&attachment.Hash,
			&attachment.Extension,
			&attachment.CreatedAt,
			&attachment.Counter,
			&n,
		); err != nil {
			return nil, 0, err
		}

		attachments = append(attachments, &attachment)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return attachments, n, nil
}

// createAttachment creates a new attachment. Sets the new database Hash to attachment.Hash and sets
// the timestamps to the current time.
func createAttachment(ctx context.Context, tx *Tx, attachment *api.Attachment) error {
	// Set timestamps to the current time.
	attachment.CreatedAt = tx.now

	// Perform basic field validation.
	if err := attachment.Validate(); err != nil {
		return err
	}

	// Execute insertion query.
	_, err := tx.ExecContext(ctx, `
		INSERT INTO attachments (
			hash,
			extension,
			created_at,
			counter
		)
		VALUES (?, ?, ?, ?)
	`,
		attachment.Hash,
		attachment.Extension,
		attachment.CreatedAt,
		attachment.Counter,
	)
	if err != nil {
		return FormatError(err)
	}

	return nil
}

// updateAttachment updates fields on a attachment object. Returns EUNAUTHORIZED if current
// attachment is not the attachment being updated.
func updateAttachment(ctx context.Context, tx *Tx, hash string, upd api.AttachmentUpdate) (*api.Attachment, error) {
	// Fetch current object state.
	attachment, err := findAttachmentByHash(ctx, tx, hash)
	if err != nil {
		return attachment, err
	}

	// Update fields.
	if v := upd.Counter; v != nil {
		attachment.Counter = *v
	}

	// Perform basic field validation.
	if err := attachment.Validate(); err != nil {
		return attachment, err
	}

	// Execute update query.
	if _, err := tx.ExecContext(ctx, `
		UPDATE attachments
		SET counter = ?
		WHERE hash = ?
	`,
		attachment.Counter,
		attachment.Hash,
	); err != nil {
		return attachment, FormatError(err)
	}

	if attachment.Counter < 1 {
		err = deleteAttachment(ctx, tx, hash)
		if err != nil {
			return attachment, err
		}
	}

	return attachment, nil
}

// deleteAttachment permanently removes a attachment by Hash. Returns EUNAUTHORIZED if current
// attachment is not the one being deleted.
func deleteAttachment(ctx context.Context, tx *Tx, hash string) error {
	// Remove row from database.
	if _, err := tx.ExecContext(ctx, `DELETE FROM attachments WHERE hash = ?`, hash); err != nil {
		return FormatError(err)
	}
	return nil
}
