package domain

import "errors"

var (
	ErrInvalidPlatformID   = errors.New("invalid platform id")
	ErrInvalidPlatformName = errors.New("invalid platform name")

	ErrInvalidBridgeMode   = errors.New("invalid bridge mode")
	ErrSameSourceTarget    = errors.New("source and target cannot be the same")

	ErrRequestExpired      = errors.New("connection request expired")
	ErrRequestNotPending   = errors.New("connection request is not pending")
	ErrTargetNotAssigned   = errors.New("target chat not assigned")
)