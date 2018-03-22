package uefi

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
)

// FirmwareVolume constants
const (
	FirmwareVolumeFixedHeaderSize = 56
	FirmwareVolumeMinSize         = FirmwareVolumeFixedHeaderSize + 8 // +8 for the null block that terminates the block list
)

// Block describes number and size of the firmware volume blocks
type Block struct {
	Count uint32
	Size  uint32
}

// FirmwareVolumeFixedHeader contains the fixed fields of a firmware volume
// header
type FirmwareVolumeFixedHeader struct {
	Zeros          [16]uint8
	FileSystemGUID [16]uint8
	Length         uint64
	Signature      uint32
	AttrMask       uint8
	HeaderLen      uint16
	Checksum       uint16
	Reserved       [3]uint8
	Revision       uint8
	Unused         [3]uint8
}

// FirmwareVolume represents a firmware volume. It combines the fixed header and
// a variable list of blocks
type FirmwareVolume struct {
	FirmwareVolumeFixedHeader
	// there must be at least one that is zeroed and indicates the end of the
	// block list
	Blocks []Block
}

// Summary prints a multi-line representation of a FirmwareVolume object
func (fv FirmwareVolume) Summary() string {
	hexGUID := make([]byte, hex.EncodedLen(len(fv.FileSystemGUID)))
	hex.Encode(hexGUID, []byte(fv.FileSystemGUID[:]))
	return fmt.Sprintf("FirmwareVolume{\n"+
		"    FileSystemGUID=%s\n"+
		"    Length=%v\n"+
		"    Signature=0x%08x\n"+
		"    AttrMask=0x%02x\n"+
		"    HeaderLen=%v\n"+
		"    Checksum=0x%04x\n"+
		"    Revision=%v\n"+
		"    Blocks=%v\n"+
		"}",
		hexGUID,
		fv.Length, fv.Signature, fv.AttrMask,
		fv.HeaderLen, fv.Checksum, fv.Revision,
		fv.Blocks,
	)
}

// FindFirmwareVolumeOffset searches for a firmware volume signature, "_FVH"
// using 8-byte alignment. If found, returns the offset from the start of the
// firmware volume, otherwise returns -1.
func FindFirmwareVolumeOffset(data []byte) int64 {
	if len(data) < 32 {
		return -1
	}
	var (
		offset int64
		fvSig  = []byte("_FVH")
	)
	for offset = 32; offset < int64(len(data)); offset += 8 {
		if bytes.Equal(data[offset:offset+4], fvSig) {
			return offset - 40 // the actual volume starts 40 bytes before the signature
		}
	}
	return -1
}

// NewFirmwareVolume parses a sequence of bytes and returns a FirmwareVolume
// object, if a valid one is passed, or an error
func NewFirmwareVolume(data []byte) (*FirmwareVolume, error) {
	if len(data) < FirmwareVolumeMinSize {
		return nil, ErrImageTooSmall
	}
	var fv FirmwareVolume
	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.LittleEndian, &fv.FirmwareVolumeFixedHeader); err != nil {
		return nil, err
	}
	// read the block map
	blocks := make([]Block, 0)
	for {
		var block Block
		if err := binary.Read(reader, binary.LittleEndian, &block); err != nil {
			return nil, err
		}
		if block.Count == 0 && block.Size == 0 {
			// found the terminating block
			log.Print("Terminating block")
			break
		}
		log.Print("Block")
		blocks = append(blocks, block)
	}
	fv.Blocks = blocks
	return &fv, nil
}
