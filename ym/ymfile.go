package ym

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

const endMarker = "End!"

// File represents an YM file, with header and data
type File struct {
	Header *Header
	Frames []byte
}

// NewFile creates a new ym.File from a file path. If is set to true
// It will attempt to validate the file correctness (for example by
// checking for the presence of the "End!" marker at the end of the file)
func NewFile(path string, strict bool) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file %v: %v", path, err)
	}
	defer f.Close()

	ymFile := File{}

	reader := bufio.NewReader(f)
	ymFile.Header, err = ReadHeader(reader)
	if err != nil {
		return nil, err
	}

	ymFile.Frames = make([]byte, ymFile.Header.FrameCount*16)
	n, err := io.ReadFull(reader, ymFile.Frames)
	if err != nil {
		return nil, fmt.Errorf("error reading frames data: %v. Expected %v bytes, got %v", err, ymFile.Header.FrameCount*16, n)
	}
	if n != int(ymFile.Header.FrameCount)*16 {
		return nil, fmt.Errorf("expected %v bytes to be read, but got %v", ymFile.Header.FrameCount*16, len(ymFile.Frames))
	}

	if strict {
		// Check for end marker
		var buf = make([]byte, 4)
		n, err = reader.Read(buf)
		if err != nil {
			return nil, fmt.Errorf("error reading end marker: %v", err)
		}
		if n != len(endMarker) {
			return nil, fmt.Errorf("expected to read %v bytes for the end marker, but got %v", len(endMarker), n)
		}
		if string(buf) != endMarker {
			return nil, fmt.Errorf("expected to read the end marker, but got: %v", buf)
		}
	}
	return &ymFile, nil
}

// MarshalBinary marshals a YM file into binary
func (f File) MarshalBinary() (data []byte, err error) {
	buf := bytes.Buffer{}

	headerBytes, err := f.Header.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("error marshalling header to bytes: %v", err)
	}

	buf.Write(headerBytes)
	buf.Write(f.Frames)
	buf.WriteString(endMarker)

	return buf.Bytes(), nil
}
