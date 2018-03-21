package uefi

import (
	"bytes"
	"fmt"
)

// FlashSignature is the sequence of bytes that a Flash image is expected to
// start with.
var (
	FlashSignature = []byte{0x5a, 0xa5, 0xf0, 0x0f}
)

// FlashImage is the main structure that represents an Intel Flash image. It
// implements the Firmware interface.
type FlashImage struct {
	// Holds the raw buffer
	buf                []byte
	DescriptorMapStart uint
	RegionStart        uint
	MasterStart        uint
	DescriptorMap      FlashDescriptorMap
	Region             FlashRegionSection
	Master             FlashMasterSection
}

// IsPCH returns whether the flash image has the more recent PCH format, or not.
// PCH images have the first 16 bytes reserved, and the 4-bytes signature starts
// immediately after. Older images (ICH8/9/10) have the signature at the
// beginning.
func (f FlashImage) IsPCH() bool {
	return bytes.Equal(f.buf[16:16+len(FlashSignature)], FlashSignature)
}

// FindSignature looks for the Intel flash signature, and returns its offset
// from the start of the image. The PCH images are located at offset 16, while
// in ICH8/9/10 they start at 0. If no signature is found, it returns -1.
func (f FlashImage) FindSignature() int {
	if bytes.Equal(f.buf[16:16+len(FlashSignature)], FlashSignature) {
		// 16 + 4 since the descriptor starts after the signature
		return 20
	}
	if bytes.Equal(f.buf[:len(FlashSignature)], FlashSignature) {
		// + 4 since the descriptor starts after the signature
		return 4
	}
	return -1
}

// Validate runs a set of checks on the flash image and returns a list of
// errors specifying what is wrong.
func (f FlashImage) Validate() []error {
	errors := make([]error, 0)
	start := f.FindSignature()
	if start < 0 {
		errors = append(errors, ErrFlashSignatureNotFound)
	}
	errors = append(errors, f.DescriptorMap.Validate()...)
	// TODO also validate regions, masters, etc
	return errors
}

func (f FlashImage) String() string {
	return fmt.Sprintf("FlashImage{Size=%v, Descriptor=%v, Region=%v, Master=%v}",
		len(f.buf),
		f.DescriptorMap.String(),
		f.Region.String(),
		f.Master.String(),
	)
}

// Summary prints a multi-line description of the flash image
func (f FlashImage) Summary() string {
	return fmt.Sprintf("FlashImage{\n"+
		"    Size=%v\n"+
		"    DescriptorMapStart=%v\n"+
		"    RegionStart=%v\n"+
		"    MasterStart=%v\n"+
		"    Descriptor=%v\n"+
		"    Region=%v\n"+
		"    Master=%v\n"+
		"}",
		len(f.buf),
		f.DescriptorMapStart,
		f.RegionStart,
		f.MasterStart,
		Indent(f.DescriptorMap.Summary(), 4),
		Indent(f.Region.Summary(), 4),
		Indent(f.Master.Summary(), 4),
	)
}

// NewFlashImage tries to create a FlashImage structure, and returns a FlashImage
// and an error if any. This only works with images that operate in Descriptor
// mode.
func NewFlashImage(buf []byte) (*FlashImage, error) {
	if len(buf) < FlashDescriptorMapSize {
		return nil, ErrImageTooSmall
	}
	flash := FlashImage{buf: buf}
	descriptorMapStart := flash.FindSignature()
	if descriptorMapStart < 0 {
		return nil, ErrFlashSignatureNotFound
	}
	flash.DescriptorMapStart = uint(descriptorMapStart)

	// Descriptor Map
	desc, err := NewFlashDescriptorMap(buf[flash.DescriptorMapStart : flash.DescriptorMapStart+FlashDescriptorMapSize])
	if err != nil {
		return nil, err
	}
	flash.DescriptorMap = *desc

	// Region
	flash.RegionStart = uint(flash.DescriptorMap.RegionBase) * 0x10
	region, err := NewFlashRegionSection(buf[flash.RegionStart : flash.RegionStart+uint(FlashRegionSectionSize)])
	if err != nil {
		return nil, err
	}
	flash.Region = *region

	// Master
	flash.MasterStart = uint(flash.DescriptorMap.MasterBase) * 0x10
	master, err := NewFlashMasterSection(buf[flash.MasterStart : flash.MasterStart+uint(FlashMasterSectionSize)])
	if err != nil {
		return nil, err
	}
	flash.Master = *master

	return &flash, nil
}
