package uefi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

// FlashRegionSectionSize is the size of the Region descriptor. It is made up by 18 fields, each 16-bits large.
const FlashRegionSectionSize = 36

type FlashRegionSection struct {
	Reserved            uint16
	FlashBlockEraseSize uint16
	BiosBase, BiosLimit uint16
	MeBase, MeLimit     uint16
	GbeBase, GbeLimit   uint16
	PdrBase, PdrLimit   uint16
}

// AvailableRegions returns a list of names of the regions with non-zero size.
func (f FlashRegionSection) AvailableRegions() []string {
	var regions []string
	if f.BiosLimit > 0 {
		regions = append(regions, "BIOS")
	}
	if f.MeLimit > 0 {
		regions = append(regions, "ME")
	}
	if f.GbeLimit > 0 {
		regions = append(regions, "GbE")
	}
	if f.PdrLimit > 0 {
		regions = append(regions, "PDR")
	}
	return regions
}

func (f FlashRegionSection) String() string {
	return fmt.Sprintf("FlashRegionSection{Regions=%v}",
		strings.Join(f.AvailableRegions(), ","),
	)
}

func (f FlashRegionSection) Summary() string {
	return fmt.Sprintf("FlashRegionSection{\n"+
		"    Regions=%v\n"+
		"    BiosBase=%v (size: %v)\n"+
		"    MeBase=%v (size: %v)\n"+
		"    GbeBase=%v (size: %v)\n"+
		"    PdrBase=%v (size: %v)\n"+
		"}",
		strings.Join(f.AvailableRegions(), ","),
		f.BiosBase, f.BiosLimit,
		f.MeBase, f.MeLimit,
		f.GbeBase, f.GbeLimit,
		f.PdrBase, f.PdrLimit,
	)
}

// NewFlashRegionSection initializes a FlashRegionSection from a slice of bytes
func NewFlashRegionSection(data []byte) (*FlashRegionSection, error) {
	if len(data) < FlashRegionSectionSize {
		return nil, ErrImageTooSmall
	}
	var region FlashRegionSection
	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.LittleEndian, &region); err != nil {
		return nil, err
	}
	return &region, nil
}
