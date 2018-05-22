package uefi

import (
	"fmt"
	"strings"
)

type BiosRegion struct {
	FirmwareVolumes []FirmwareVolume
}

func (br BiosRegion) Summary() string {
	var fvols []string
	for _, fv := range br.FirmwareVolumes {
		fvols = append(fvols, fv.Summary())
	}
	return fmt.Sprintf("BiosRegion{\n"+
		"    FirmwareVolumes=[\n"+
		"        %v\n"+
		"    ]\n"+
		"}", Indent(strings.Join(fvols, "\n"), 8))
}

func NewBiosRegion(data []byte) (*BiosRegion, error) {
	var br BiosRegion
	for {
		offset := FindFirmwareVolumeOffset(data)
		if offset == -1 {
			// no firmware volume found, stop searching
			break
		}
		fv, err := NewFirmwareVolume(data[offset:])
		if err != nil {
			return nil, err
		}
		data = data[uint64(offset)+fv.Length:]
		br.FirmwareVolumes = append(br.FirmwareVolumes, *fv)
		// FIXME remove the `break` and move the offset to the next location to
		// search for FVs (i.e. offset + fv.size)
	}
	return &br, nil
}
