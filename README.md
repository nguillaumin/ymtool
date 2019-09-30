# YM tool

YM Tool is a utility to maniplate [YM files](http://leonard.oxg.free.fr/ymformat.html).
YM files store music data from the [YM2149](https://en.wikipedia.org/wiki/General_Instrument_AY-3-8910) 
chip which was used in various computers of the 80's and 90's, notably the
Atari ST.

This tools is intended to make it easy to manipulate and alter YM files on the
command line. It was built while working on the [YM Jukebox](https://nguillaumin.github.io/ym-jukebox/).

## Download

Visit the [releases](../../releases) page.

## Usage

Run `ymtool` without any arguments to see what commands are available:

```
YM Tool

Usage: ymtool <command> [arg...]

Available commands:
  - info <file.ym>             : Get metadata of a YM file
  - md5 <file.ym>              : Compute the MD5 sum of a YM song data
  - update <options> <file.ym> : Update YM file metadata

Run a command without argument to get more information.

Supported YM versions: YM6! YM5! 

Exit codes:
  - 1: Error parsing command line.
  - 2: I/O error reading or writing files.
  - 3: Unsupported YM version.

WARNING: This tool doesn't unpack LHA/LZH compressed YM files.
The YM file is expected to be already unpacked.
```

## Limitations

### LHA/LZH compression

⚠️ YM files are usually compressed with LHA/LZH. This utility will not unpack
such files (as I couldn't find an existing LHA/LZH library in Go), so the file
is expected to be passed in already unpacked.

On Linux, [lhasa](https://github.com/fragglet/lhasa) can unpack LHA files. On
Windows, 7-Zip or WinRAR can be used.

### Supported YM versions

At this time, only YM5 and YM6 formats are supported. Adding support for more
formats should not be difficult but is not implemented yet.