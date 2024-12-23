package common

import (
	"github.com/pkg/errors"
	"strconv"
	"time"
)

var Now = time.Now

func ParseIntFromString[K int | uint | int32 | int64](s string) (K, error) {
	if s == "" {
		return 0, errors.New("id is required")
	}
	parsedID, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}

	if parsedID <= 0 {
		return 0, errors.New("invalid id")
	}

	return K(parsedID), nil
}
