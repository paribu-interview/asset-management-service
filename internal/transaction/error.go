package transaction

import "github.com/pkg/errors"

var (
	ErrTransactionNotFound        = errors.New("transaction not found")
	ErrTransactionCannotBeDeleted = errors.New("transaction cannot be deleted")
	ErrAssetNotFound              = errors.New("asset not found")
	ErrInsufficientBalance        = errors.New("insufficient balance")
)
