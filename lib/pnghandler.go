package lib

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

var PNGHEADERSTR = "\x89\x50\x4E\x47\x0D\x0A\x1A\x0A"

type Config struct {
	Verbose    bool
	ReWriteCRC bool
	Lenient    bool
	Overwrite  bool
	InPlace    bool
	Display    bool

	ActionRead  string
	ActionWrite string
	SourceFile  string
	DestFile    string
}

type XmlChunk struct {
	Diagram string `xml:"diagram"`
}

type Chunk struct {
	ChunkLen  uint32
	ChunkType string
	ChunkData []byte
	ChunkCRC  uint32
}

func Read_content(l *log.Logger, config *Config) string {
	imgFile, err := open_png_file(l, config.ActionRead)
	if err != nil {
		l.Fatal(err)
	}
	defer imgFile.Close()

	var content string

	for {
		chunk, err := read_chunk(config, imgFile)
		if err != nil {
			break
		}

		if config.Verbose {
			l.Println("Found chunk type: ", chunk.ChunkType)
		}
		if strings.ToLower(chunk.ChunkType) == "text" {
			content, err = decode_chunk(l, config, chunk)
			if err != nil {
				l.Fatal(err)
			}
		}
	}

	if config.Display {
		print(content)
	}
	return content
}

func Write_content_from_data(l *log.Logger, action string, content []byte) string {
	config := new(Config)
	config.Display = false

	workFile, err := ioutil.TempFile("", "*")
	if err != nil {
		l.Fatal("unable to open temporary file")
	}
	_, err = workFile.Write(content)
	config.ActionRead = workFile.Name()
	workFile.Close()

	if action == "read" {
		return Read_content(l, config)
	} else {
		return workFile.Name()
	}
}

func Write_content_dynamic_config(
	l *log.Logger,
	writePath string,
	sourcePath string,
	destPath string) {

	config := new(Config)
	config.Display = false
	config.ActionWrite = writePath
	config.SourceFile = sourcePath
	config.DestFile = destPath
	config.Lenient = true
	config.Overwrite = true
	config.InPlace = false
	config.ReWriteCRC = false
	config.Verbose = false

	Write_content(l, config)
}

func Write_content(l *log.Logger, config *Config) {
	finalizeInPlace := false

	imgFile, err := open_png_file(l, config.ActionWrite)
	if err != nil {
		l.Fatal(err)
	}

	var destFile *os.File

	finalizeDest := func() {
		if finalizeInPlace {
			if err = os.Remove(imgFile.Name()); err != nil {
				l.Fatal(err)
			}
			if err = os.Rename(destFile.Name(), imgFile.Name()); err != nil {
				l.Fatal(err)
			}
		}
	}

	if config.InPlace {
		destFile, err = ioutil.TempFile("", "*")
		if err != nil {
			l.Fatal("unable to open temporary file")
		}
		defer finalizeDest()
	} else {
		if config.DestFile == "-" {
			destFile = os.Stdout
		} else {
			destFile, err = os.Create(config.DestFile)
			if err != nil {
				l.Fatal("unable to open destination file")
			}
		}
		defer destFile.Close()
	}

	// Declare this defer call further down the stack to ensure
	// proper finalization work.
	defer imgFile.Close()

	sourceCode, err := os.ReadFile(config.SourceFile)
	if err != nil {
		l.Fatal(err)
	}

	destFile.Write([]byte(PNGHEADERSTR))

	sourceWritten := false

	for {
		chunk, err := read_chunk(config, imgFile)
		if err != nil {
			break
		}

		normalizedType := strings.ToLower(chunk.ChunkType)
		if config.Verbose {
			l.Println("Chunk type", normalizedType)
		}
		if normalizedType == "text" {
			if !config.Overwrite {
				l.Println("! Not overwriting existing text chunk.")
			} else {
				if config.Verbose {
					l.Println("-> Replacing existing text chunk")
				}
				write_sourcecode(config, destFile, sourceCode)
				sourceWritten = true
				continue
			}
		} else {
			if normalizedType == "iend" {
				if !sourceWritten {
					if config.Verbose {
						l.Println("-> Including new text chunk")
					}
					write_sourcecode(config, destFile, sourceCode)
				}
			}
		}

		write_chunk(config, destFile, chunk)
	}

	if config.InPlace {
		finalizeInPlace = true
	}
}

func write_sourcecode(config *Config, destFile *os.File, sourceCode []byte) {
	var buf bytes.Buffer
	writer, _ := flate.NewWriter(&buf, -1)
	writer.Write(sourceCode)
	writer.Close()
	deflatedSC := buf.Bytes()
	b64SC := base64.StdEncoding.EncodeToString(deflatedSC)
	data := "mxfile" +
		url.QueryEscape(
			`<mxfile host="txt" modified="`+
				time.Now().Format(time.RFC3339)+
				`" agent="1.0 (txt)" etag="1" version="1" type="device"><diagram id="W0tw_o5MNMKq2u_i8K05" name="Page-1">`) +
		string(b64SC) +
		url.QueryEscape("</diagram></mxfile>")

	chunk := new(Chunk)
	chunk.ChunkType = "tExt"
	chunk.ChunkLen = uint32(len(data))
	chunk.ChunkData = []byte(data)
	chunk.ChunkCRC = chunk_CRC([]byte(chunk.ChunkType), chunk.ChunkData)
	write_chunk(config, destFile, chunk)
}

func open_png_file(l *log.Logger, filepath string) (*os.File, error) {
	var err error
	var imgFile *os.File

	if filepath == "-" {
		imgFile = os.Stdin
	} else {
		imgFile, err = os.Open(filepath)
		if err != nil {
			return nil, err
		}
	}

	ihdrStr := make([]byte, 8)
	_, err = io.ReadFull(imgFile, ihdrStr)
	if err != nil {
		return nil, err
	}

	if PNGHEADERSTR != string(ihdrStr) {
		return nil, fmt.Errorf("This is not a PNG file.")
	}

	return imgFile, nil
}

func read_chunk(config *Config, imgFile *os.File) (*Chunk, error) {
	chunk := new(Chunk)
	var err error

	rawChunkLen := make([]byte, 4)
	_, err = io.ReadFull(imgFile, rawChunkLen)
	if err != nil {
		return nil, fmt.Errorf("error reading length: %v", err)
	}
	chunk.ChunkLen = uint32(binary.BigEndian.Uint32(rawChunkLen))

	rawChunkType := make([]byte, 4)
	_, err = io.ReadFull(imgFile, rawChunkType)
	if err != nil {
		return nil, fmt.Errorf("error reading type")
	}
	chunk.ChunkType = string(rawChunkType)

	chunk.ChunkData = make([]byte, chunk.ChunkLen)
	_, err = io.ReadFull(imgFile, chunk.ChunkData)
	if err != nil {
		return nil, fmt.Errorf("error reading data")
	}

	rawChunkCRC := make([]byte, 4)
	_, err = io.ReadFull(imgFile, rawChunkCRC)
	if err != nil {
		return nil, fmt.Errorf("error reading CRC")
	}
	chunk.ChunkCRC = uint32(binary.BigEndian.Uint32(rawChunkCRC))

	if !config.Lenient {
		if chunk.ChunkCRC != chunk_CRC(rawChunkType, chunk.ChunkData) {
			return nil, fmt.Errorf("error incorrect CRC")
		}
	}

	return chunk, nil
}

func write_chunk(config *Config, destFile *os.File, chunk *Chunk) error {
	var rawChunkLen bytes.Buffer
	var rawChunkCRC bytes.Buffer
	if binary.Write(&rawChunkLen, binary.BigEndian, chunk.ChunkLen) != nil {
		return fmt.Errorf("Unable to write chunk length")
	}
	destFile.Write(rawChunkLen.Bytes())

	destFile.Write([]byte(chunk.ChunkType))

	destFile.Write(chunk.ChunkData)

	if config.ReWriteCRC {
		if binary.Write(
			&rawChunkCRC,
			binary.BigEndian,
			chunk_CRC([]byte(chunk.ChunkType), chunk.ChunkData)) != nil {
			return fmt.Errorf("Unable to convert computed CRC size")
		}
	} else {
		if binary.Write(&rawChunkCRC, binary.BigEndian, chunk.ChunkCRC) != nil {
			return fmt.Errorf("Unable to write computed CRC size")
		}
	}
	destFile.Write(rawChunkCRC.Bytes())

	return nil
}

func decode_chunk(l *log.Logger, config *Config, chunk *Chunk) (string, error) {
	if config.Verbose {
		l.Println("Chunk size: ", chunk.ChunkLen)
	}
	// Brutal removal of header and footer
	if chunk.ChunkLen < 22+31 {
		return "", fmt.Errorf("error chunk size too small")
	}
	start := strings.Index(string(chunk.ChunkData), "name%3D%22Page-1%22%3E")
	if start == -1 {
		return "", fmt.Errorf("error no preamble in chunk")
	}
	data := string(chunk.ChunkData[start+22 : chunk.ChunkLen-31])
	debased, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("error not base64 encoded")
	}
	deflated, err := ioutil.ReadAll(flate.NewReader(bytes.NewReader(debased)))
	if err != nil {
		return "", fmt.Errorf("error not deflate mode")
	}
	return string(deflated), nil
}

func chunk_CRC(chunkType []byte, chunkData []byte) uint32 {
	c32 := crc32.NewIEEE()
	c32.Write(chunkType)
	c32.Write(chunkData)
	return c32.Sum32()
}
