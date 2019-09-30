package ym

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"regexp"
)

const checkString = "LeOnArD!"

// SupportedYmVersions contains the liset of YM file version this
// library supports
var SupportedYmVersions = map[string]bool{"YM5!": true, "YM6!": true}

var ymMarkerRegExp = regexp.MustCompile("YM[23456]!|YM3b|MIX1|YMT1|YMT2")

// DigiDrumSample stores a single DigiDrum sample
type DigiDrumSample struct {
	Size uint32
	Data []byte
}

// Header represents a YM file header
type Header struct {
	Version               string
	FrameCount            uint32
	Attributes            uint32
	DigiDrumSamplesNumber uint16
	MasterClock           uint32
	PlayerFrame           uint16
	LoopFrame             uint32
	DigiDrums             []DigiDrumSample
	SongName              string
	Author                string
	Comment               string
	skipBytes             uint16
	skippedData           []byte
}

// UnsupportedVersionError occurs when a version of the YM file is not supported
// by this library
type UnsupportedVersionError struct {
	UnsupportedVersion string
}

func (e UnsupportedVersionError) Error() string {
	return fmt.Sprintf("unsupported YM version: %v", e.UnsupportedVersion)
}

// ReadHeader reads an YM header from a Reader
func ReadHeader(reader *bufio.Reader) (*Header, error) {
	// Read magic marker
	buf := make([]byte, 4)
	if _, err := reader.Read(buf); err != nil {
		return nil, err
	}

	version := ymMarkerRegExp.FindString(string(buf))
	if version == "" {
		return nil, fmt.Errorf("Unable to extract YM marker or YM version. Marker string: '%v'", string(buf))
	}

	if _, ok := SupportedYmVersions[version]; !ok {
		return nil, UnsupportedVersionError{UnsupportedVersion: version}
	}

	ymHeader := Header{
		Version: version,
	}

	// Read check string
	buf = make([]byte, 8)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return nil, err
	}
	if checkString != string(buf) {
		return nil, fmt.Errorf("Invalid check string: '%v'", string(buf))
	}

	binary.Read(reader, binary.BigEndian, &ymHeader.FrameCount)
	binary.Read(reader, binary.BigEndian, &ymHeader.Attributes)
	binary.Read(reader, binary.BigEndian, &ymHeader.DigiDrumSamplesNumber)
	binary.Read(reader, binary.BigEndian, &ymHeader.MasterClock)
	binary.Read(reader, binary.BigEndian, &ymHeader.PlayerFrame)
	binary.Read(reader, binary.BigEndian, &ymHeader.LoopFrame)

	bufferedReader := bufio.NewReader(reader)
	binary.Read(reader, binary.BigEndian, &ymHeader.skipBytes)
	ymHeader.skippedData = make([]byte, ymHeader.skipBytes)
	io.ReadFull(reader, ymHeader.skippedData)

	for i := 0; i < int(ymHeader.DigiDrumSamplesNumber); i++ {
		dd := DigiDrumSample{}
		binary.Read(bufferedReader, binary.BigEndian, &dd.Size)
		dd.Data = make([]byte, dd.Size)
		io.ReadFull(bufferedReader, dd.Data)
	}

	// Read string information, stripping the trailing \0
	str, _ := bufferedReader.ReadString(0)
	ymHeader.SongName = str[:len(str)-1]
	str, _ = bufferedReader.ReadString(0)
	ymHeader.Author = str[:len(str)-1]
	str, _ = bufferedReader.ReadString(0)
	ymHeader.Comment = str[:len(str)-1]

	return &ymHeader, nil
}

func (header Header) String() string {
	out := fmt.Sprintf("Song name       : %v\n", header.SongName)
	out += fmt.Sprintf("Author          : %v\n", header.Author)
	out += fmt.Sprintf("Comment         : %v\n", header.Comment)
	out += fmt.Sprintf("YM version      : %v\n", header.Version)
	out += fmt.Sprintf("Number of frames: %v\n", header.FrameCount)
	out += fmt.Sprintf("Loop frame      : %v\n", header.LoopFrame)
	out += fmt.Sprintf("Attributes      : %032b\n", header.Attributes)
	out += fmt.Sprintf("Digidrum samples: %v\n", header.DigiDrumSamplesNumber)
	out += fmt.Sprintf("Master clock    : %vHz\n", header.MasterClock)
	out += fmt.Sprintf("Player frame    : %vHz\n", header.PlayerFrame)

	return out
}

// MarshalBinary marshals a YM file header into binary
func (ymHeader Header) MarshalBinary() (data []byte, err error) {
	buf := bytes.Buffer{}

	buf.WriteString(ymHeader.Version)
	buf.WriteString(checkString)

	binary.Write(&buf, binary.BigEndian, ymHeader.FrameCount)
	binary.Write(&buf, binary.BigEndian, ymHeader.Attributes)
	binary.Write(&buf, binary.BigEndian, ymHeader.DigiDrumSamplesNumber)
	binary.Write(&buf, binary.BigEndian, ymHeader.MasterClock)
	binary.Write(&buf, binary.BigEndian, ymHeader.PlayerFrame)
	binary.Write(&buf, binary.BigEndian, ymHeader.LoopFrame)

	binary.Write(&buf, binary.BigEndian, ymHeader.skipBytes)
	buf.Write(ymHeader.skippedData)

	for i := 0; i < int(ymHeader.DigiDrumSamplesNumber); i++ {
		dd := ymHeader.DigiDrums[i]
		binary.Write(&buf, binary.BigEndian, &dd.Size)
		buf.Write(dd.Data)
	}

	buf.WriteString(ymHeader.SongName + "\x00")
	buf.WriteString(ymHeader.Author + "\x00")
	buf.WriteString(ymHeader.Comment + "\x00")

	return buf.Bytes(), nil

}
