package asset

import "github.com/pkg/errors"

var (
	ErrDuplicateAsset = errors.New("asset already exist")
)
