package ym_test

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/nguillaumin/ymtool/ym"
)

func TestReadSimpleHeader(t *testing.T) {

	header := []byte{
		0x59, 0x4d, 0x35, 0x21, 0x4c, 0x65, 0x4f, 0x6e, 0x41, 0x72, 0x44, 0x21, 0x00, 0x00, 0x1c, 0xb0,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x1e, 0x84, 0x80, 0x00, 0x32, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x41, 0x6d, 0x62, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x20, 0x6d, 0x65, 0x6e, 0x75, 0x20,
		0x32, 0x00, 0x4c, 0x6f, 0x72, 0x64, 0x00, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64,
		0x20, 0x77, 0x69, 0x74, 0x68, 0x20, 0x53, 0x61, 0x69, 0x6e, 0x54, 0x00}

	buf := bytes.NewBuffer(header)

	ymHeader, err := ym.ReadHeader(bufio.NewReader(buf))

	if err != nil {
		t.Errorf("Unexpected error while reading header: %v", err)
	}

	if ymHeader.Version != "YM5!" {
		t.Errorf("Expected version YM5!, but got %v", ymHeader.Version)
	}
	if ymHeader.FrameCount != 7344 {
		t.Errorf("Expected 7344 frames, but got %v", ymHeader.FrameCount)
	}
	if ymHeader.Attributes != 1 {
		t.Errorf("Expected attributes = 1, but got %v", ymHeader.Attributes)
	}
	if ymHeader.DigiDrumSamplesNumber != 0 {
		t.Errorf("Expected 0 DigiDrum samples, but got: %v", ymHeader.DigiDrumSamplesNumber)
	}
	if ymHeader.MasterClock != 2000000 {
		t.Errorf("Expected master clock = 2000000, but got: %v", ymHeader.MasterClock)
	}
	if ymHeader.PlayerFrame != 50 {
		t.Errorf("Expected player frame = 50, but got %v", ymHeader.PlayerFrame)
	}
	if ymHeader.LoopFrame != 0 {
		t.Errorf("Expected loop frame = 0, but got %v", ymHeader.LoopFrame)
	}
	if len(ymHeader.DigiDrums) != 0 {
		t.Errorf("Expected no DigiDrum samples, but got %v", len(ymHeader.DigiDrums))
	}
	if ymHeader.SongName != "Ambition menu 2" {
		t.Errorf("Unexpected song name '%v'", ymHeader.SongName)
	}
	if ymHeader.Author != "Lord" {
		t.Errorf("Unexpected author '%v'", ymHeader.Author)
	}
	if ymHeader.Comment != "Generated with SainT" {
		t.Errorf("Unexpected comment '%v'", ymHeader.Comment)
	}
}

func TestNoStrings(t *testing.T) {

	header := []byte{
		0x59, 0x4d, 0x35, 0x21, 0x4c, 0x65, 0x4f, 0x6e, 0x41, 0x72, 0x44, 0x21, 0x00, 0x00, 0x1c, 0xb0,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x1e, 0x84, 0x80, 0x00, 0x32, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00}

	buf := bytes.NewBuffer(header)

	ymHeader, err := ym.ReadHeader(bufio.NewReader(buf))

	if err != nil {
		t.Errorf("Unexpected error while reading header: %v", err)
	}
	if ymHeader.SongName != "" {
		t.Errorf("Unexpected song name '%v'", ymHeader.SongName)
	}
	if ymHeader.Author != "" {
		t.Errorf("Unexpected author '%v'", ymHeader.Author)
	}
	if ymHeader.Comment != "" {
		t.Errorf("Unexpected comment '%v'", ymHeader.Comment)
	}

}

func TestReadNoYMMarker(t *testing.T) {
	header := []byte("Not an YM file")

	buf := bytes.NewBuffer(header)

	ymHeader, err := ym.ReadHeader(bufio.NewReader(buf))
	if err == nil {
		t.Errorf("Expected an error, but got a valid YM header: %v", ymHeader)
	}
}

func TestMarshalSimpleHeader(t *testing.T) {

	header := []byte{
		0x59, 0x4d, 0x35, 0x21, 0x4c, 0x65, 0x4f, 0x6e, 0x41, 0x72, 0x44, 0x21, 0x00, 0x00, 0x1c, 0xb0,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x1e, 0x84, 0x80, 0x00, 0x32, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x41, 0x6d, 0x62, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x20, 0x6d, 0x65, 0x6e, 0x75, 0x20,
		0x32, 0x00, 0x4c, 0x6f, 0x72, 0x64, 0x00, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64,
		0x20, 0x77, 0x69, 0x74, 0x68, 0x20, 0x53, 0x61, 0x69, 0x6e, 0x54, 0x00}

	buf := bytes.NewBuffer(header)
	ymHeader, err := ym.ReadHeader(bufio.NewReader(buf))
	if err != nil {
		t.Errorf("Unexpected error while reading header: %v", err)
	}

	written, err := ymHeader.MarshalBinary()
	if err != nil {
		t.Errorf("Unexpected error while marshalling header: %v", err)
	}
	if !bytes.Equal(header, written) {
		t.Errorf("Written header differs from read one.\nRead: %v\nWritten: %v", header, written)
	}

}
