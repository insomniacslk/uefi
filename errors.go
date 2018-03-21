package uefi

import "fmt"

// Errors used by FlashImage
var (
	// ErrInvalidImageSize is used when an image's size is not exactly what
	// expected (i.e .smaller or larger)
	ErrInvalidImageSize = fmt.Errorf("Invalid image size")
	// ErrImageTooSmall is used when a firmware image's size does not
	// have the minimum required size for a firmware type
	ErrImageTooSmall = fmt.Errorf("Image size too small")
	// ErrFlashSignatureNotFound is used when a firmware image does not
	// have a valid Flash signature where it is expected to be
	ErrFlashSignatureNotFound = fmt.Errorf("Flash signature not found")
	// ErrUnknownFirmwareType is used when a firmware image does not match any
	// known firmware type
	ErrUnknownFirmwareType = fmt.Errorf("Unknown firmware type")
)

// Errors used by FlashDescriptor

// ErrInvalidBaseAddr is used when a region in a FlashDescriptor has a base
// address larger than FlashDescriptorMaxBase
type ErrInvalidBaseAddr struct {
	message string
}

// NewErrInvalidBaseAddr returns an ErrInvalidBaseAddr with a custom
// message
func NewErrInvalidBaseAddr(message string) *ErrInvalidBaseAddr {
	return &ErrInvalidBaseAddr{message: message}
}

func (e ErrInvalidBaseAddr) Error() string {
	return e.message
}
