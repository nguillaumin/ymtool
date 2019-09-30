package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"os"

	"github.com/nguillaumin/ymtool/ym"
)

const exitCodeCmdParsing = 1
const exitCodeIOError = 2
const exitCodeUnsupportedVersion = 3

func main() {

	flag.Parse()

	if len(flag.Args()) < 1 {
		Usage()
		os.Exit(exitCodeCmdParsing)
	}

	command := flag.Arg(0)
	switch command {
	case "info":
		InfoCmd()
	case "md5":
		Md5Cmd()
	case "update":
		UpdateCmd()
	default:
		Usage()
		os.Exit(exitCodeCmdParsing)
	}

}

// Usage prints the command line utility usage
func Usage() {
	fmt.Println("YM Tool")
	fmt.Println("")
	fmt.Println("Usage: ymtool <command> [arg...]")
	fmt.Println("")

	fmt.Println("Available commands:")
	fmt.Println("  - info <file.ym>             : Get metadata of a YM file")
	fmt.Println("  - md5 <file.ym>              : Compute the MD5 sum of a YM song data")
	fmt.Println("  - update <options> <file.ym> : Update YM file metadata")
	fmt.Println("")
	fmt.Println("Run a command without argument to get more information.")
	fmt.Println("")

	fmt.Print("Supported YM versions: ")
	for version := range ym.SupportedYmVersions {
		fmt.Printf("%v ", version)
	}
	fmt.Println("")
	fmt.Println("")

	fmt.Println("Exit codes:")
	fmt.Printf("  - %v: Error parsing command line.\n", exitCodeCmdParsing)
	fmt.Printf("  - %v: I/O error reading or writing files.\n", exitCodeIOError)
	fmt.Printf("  - %v: Unsupported YM version.\n", exitCodeUnsupportedVersion)
	fmt.Println("")

	fmt.Println("WARNING: This tool doesn't unpack LHA/LZH compressed YM files.")
	fmt.Println("The YM file is expected to be already unpacked.")

	flag.PrintDefaults()
}

// InfoCmd shows information about a YM file
func InfoCmd() {
	if len(flag.Args()) < 2 {

		fmt.Println("YM Tool - Show information about a song")
		fmt.Println("")
		fmt.Println("Usage: ymtool info <file.ym>")

		os.Exit(exitCodeCmdParsing)
	}

	filePath := flag.Arg(1)
	ymFile, err := ym.NewFile(filePath, false)
	switch err := err.(type) {
	case nil:
	case ym.UnsupportedVersionError:
		fmt.Printf("Error opening file %v: %v\n", filePath, err)
		os.Exit(exitCodeUnsupportedVersion)
	default:
		fmt.Printf("Error opening file %v: %v\n", filePath, err)
		os.Exit(exitCodeIOError)
	}

	fmt.Printf("Information for %v:\n\n", filePath)
	fmt.Println(ymFile.Header)
}

// Md5Cmd computes the MD5 hash of a YM file
func Md5Cmd() {
	if len(flag.Args()) < 2 {
		fmt.Println("YM Tool - Compute MD5 hash of a song data")
		fmt.Println("")
		fmt.Println("  This will compute the MD5 hash of a song data (i.e. the YM registers")
		fmt.Println("  ignoring the song metadata. Useful to find duplicate songs that have")
		fmt.Println("  different metadata")
		fmt.Println("")
		fmt.Println("Usage: ymtool md5 <file.ym>")

		os.Exit(exitCodeCmdParsing)
	}

	filePath := flag.Arg(1)
	ymFile, err := ym.NewFile(filePath, false)
	switch err := err.(type) {
	case nil:
	case ym.UnsupportedVersionError:
		fmt.Printf("Error opening file %v: %v\n", filePath, err)
		os.Exit(exitCodeUnsupportedVersion)
	default:
		fmt.Printf("Error opening file %v: %v\n", filePath, err)
		os.Exit(exitCodeIOError)
	}

	fmt.Printf("%x\t%v\n", md5.Sum(ymFile.Frames), filePath)
}

// UpdateCmd updates the metadata of a YM file
func UpdateCmd() {

	fs := flag.NewFlagSet("Update options", flag.ExitOnError)
	songName := fs.String("song-name", "", "Name of the song")
	author := fs.String("author", "", "Song author")
	comment := fs.String("comment", "", "Song comment")

	err := fs.Parse(os.Args[2:])
	if err != nil {
		fmt.Printf("Error parsing command line flags: %v", err)
		os.Exit(exitCodeCmdParsing)
	}

	if len(fs.Args()) < 1 || (*songName == "" && *author == "" && *comment == "") {
		fmt.Println("YM Tool - Update metadata of a song")
		fmt.Println("")
		fmt.Println("Usage: ymtool update [options...] <file.ym>")
		fmt.Println("")
		fmt.Println("Options:")
		fs.PrintDefaults()
		fmt.Println("")
		fmt.Println("Example:")
		fmt.Println("  ymtool update -author=Joe -song-name=\"Joe's song\" file.ym")

		os.Exit(exitCodeCmdParsing)
	}

	filePath := fs.Arg(0)
	ymFile, err := ym.NewFile(filePath, false)
	if err != nil {
		fmt.Printf("Error opening file %v: %v\n", filePath, err)
		os.Exit(exitCodeIOError)
	}

	updated := false
	if *songName != "" {
		ymFile.Header.SongName = *songName
		updated = true
	}
	if *author != "" {
		ymFile.Header.Author = *author
		updated = true
	}
	if *comment != "" {
		ymFile.Header.Comment = *comment
		updated = true
	}

	if updated {
		data, err := ymFile.MarshalBinary()
		if err != nil {
			fmt.Printf("Error marshalling updated file to binary: %v", err)
			os.Exit(exitCodeIOError)
		}

		f, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("Error opening output file %v for writing: %v", filePath, err)
			os.Exit(exitCodeIOError)
		}
		defer f.Close()

		_, err = f.Write(data)
		if err != nil {
			fmt.Printf("Error writing output file: %v", err)
			os.Exit(exitCodeIOError)
		}

	}
}
