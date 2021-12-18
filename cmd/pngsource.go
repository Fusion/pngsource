package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fusion/pngsource/lib"
)

func setConfig(l *log.Logger) (*lib.Config, error) {
	config := &lib.Config{Display: true}
	fs := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)
	fs.BoolVar(&config.Verbose, "verbose", false, "Verbose display")
	fs.BoolVar(&config.Verbose, "v", false, "Verbose display")
	fs.BoolVar(&config.ReWriteCRC, "crc", false, "Write newly computed CRC values")
	fs.BoolVar(&config.ReWriteCRC, "c", false, "Write newly computed CRC values")
	fs.BoolVar(&config.Lenient, "lenient", false, "Do not reject incorrect CRC values")
	fs.BoolVar(&config.Lenient, "l", false, "Do not reject incorrect CRC values")
	fs.BoolVar(&config.Overwrite, "overwrite", false, "Allow overwriting previous embed if present")
	fs.BoolVar(&config.Overwrite, "o", false, "Allow overwriting previous embed if present")
	fs.BoolVar(&config.InPlace, "inplace", false, "When embedding, embed in current image (no 'to' file)")
	fs.BoolVar(&config.InPlace, "i", false, "When embedding, embed in current image (no 'to' file)")
	fs.StringVar(&config.ActionRead, "read", "", "Read image's embedded code (use '-' for stdin)")
	fs.StringVar(&config.ActionRead, "r", "", "Read image's embedded code (use '-' for stdin)")
	fs.StringVar(&config.ActionWrite, "embed", "", "Embed code in image")
	fs.StringVar(&config.ActionWrite, "e", "", "Embed code in image")
	fs.StringVar(&config.SourceFile, "source", "", "File containing code to embed")
	fs.StringVar(&config.SourceFile, "s", "", "File containing code to embed")
	fs.StringVar(&config.DestFile, "to", "", "New file to be created with embedded code (use '-' for stdout)")
	fs.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage for %s:\n\n", os.Args[0])
		fs.PrintDefaults()
		fmt.Fprintln(w, "\nExamples:\n")
		fmt.Fprintf(w, "> %s --embed <your image.png> --source <source code.txt> --to <new image.png>\n", os.Args[0])
		fmt.Fprintf(w, "> cat <your image.png> | %s --embed - --source <source code.txt> --to <new image.png>\n", os.Args[0])
		fmt.Fprintf(w, "> %s --read <image.png> --verbose\n", os.Args[0])
	}

	err := fs.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}
	if config.ActionRead != "" && config.ActionWrite != "" {
		l.Fatal("You cannot provide both and read and a write action.")
	}
	if config.ActionRead == "" && config.ActionWrite == "" {
		l.Fatal("Please specify a '--read' or '--write' action.")
	}
	if config.ActionWrite != "" && config.SourceFile == "" {
		l.Fatal("You must specify a '--source' file to perform a '--write' action.")
	}
	if config.ActionWrite != "" && config.DestFile == "" && !config.InPlace {
		l.Fatal("You must specify a '--to' file or '--inplace' when embedding.")
	}
	if config.DestFile != "" && config.InPlace {
		l.Fatal("You cannot specify '--inplace' and a '--to' file together.")
	}
	return config, nil
}

func main() {
	l := log.New(os.Stderr, "", 0)
	config, err := setConfig(l)
	if err != nil {
		println(err)
		return
	}
	if config.ActionRead != "" {
		lib.Read_content(l, config)
	}
	if config.ActionWrite != "" {
		lib.Write_content(l, config)
	}
}
