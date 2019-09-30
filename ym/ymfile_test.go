package ym_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/nguillaumin/ymtool/ym"
)

func TestReadFile(t *testing.T) {
	file, err := ym.NewFile("testdata/song1.ym", true)

	if err != nil {
		t.Errorf("Unexpected error reading file: %v", err)
	}

	if len(file.Frames) != int(file.Header.FrameCount)*16 {
		t.Errorf("Invalid length for the data. Expected %v, but got %v", int(file.Header.FrameCount*16), len(file.Frames))
	}
}

func TestMarshalFiles(t *testing.T) {
	files := []string{"testdata/song1.ym", "testdata/song2.ym"}
	for i := 0; i < len(files); i++ {
		file, err := ym.NewFile(files[i], true)
		if err != nil {
			t.Errorf("Unexpected error reading file %v: %v", files[i], err)
		}

		written, err := file.MarshalBinary()
		if err != nil {
			t.Errorf("Unexpected error marshalling file %v: %v", files[i], err)
		}

		read, err := ioutil.ReadFile(files[i])
		if err != nil {
			t.Errorf("Unexpected error reading test file: %v", err)
		}

		if !bytes.Equal(read, written) {
			t.Errorf("Written file %v differs from read one.\nRead: %v\nWritten: %v", files[i], len(read), len(written))
		}
	}
}
