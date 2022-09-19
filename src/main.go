package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

const Version = "1.2"

const (
	CmdAdd   = "add"
	CmdInit  = "init"
	CmdMan   = "man"
	CmdShelf = "shelf"
)

var GlobalOpts GlobalOptions

func main() {
	var opts Options
	p := flags.NewParser(&GlobalOpts, flags.Default)

	p.ShortDescription = "Terminal Library"
	p.LongDescription = "Light-weight CLI-based library manager utilizing the OpenLibrary API (https://openlibrary.org/)"
	p.SubcommandsOptional = true // --version doesn't need one

	if _, err := p.AddCommand(CmdInit, "Initialize a new library", "", &opts.InitOpts); err != nil {
		panic(err)
	}
	if _, err := p.AddCommand(CmdAdd, "Add books to library", "", &opts.AddOpts); err != nil {
		panic(err)
	}
	if _, err := p.AddCommand(CmdShelf, "Print books and other library-wide operations", "", &opts.ShelfOpts); err != nil {
		panic(err)
	}

	man, err := p.AddCommand(CmdMan, "Write man page", "", &EmptyOptions{})
	if err != nil {
		panic(err)
	}
	man.Hidden = true

	// p.ParseArgs(strings.Split("add 9780743496735", " "))
	// GlobalOpts.FilePath = "/home/josh/src/tlib/tlib.json"
	p.Parse()

	if GlobalOpts.Version {
		Print(Version)
		return
	}

	if p.Active == nil {
		p.WriteHelp(os.Stdout)
		return
	}

	switch p.Active.Name {
	case CmdAdd:
		RunAdd(&opts.AddOpts)
	case CmdInit:
		RunInit(&opts.InitOpts)
	case CmdMan:
		p.WriteManPage(os.Stdout)
	case CmdShelf:
		RunShelf(&opts.ShelfOpts)
	default:
		p.WriteHelp(os.Stdout)
	}
}

func Debug(msg string, a ...any) {
	if GlobalOpts.Debug {
		fmt.Fprintf(os.Stderr, "DEBUG: "+msg+"\n", a...)
	}
}
func Error(msg string, a ...any) {
	fmt.Fprintf(os.Stderr, "ERROR: "+msg+"\n", a...)
}
func Warn(msg string, a ...any) {
	if !GlobalOpts.Quiet {
		fmt.Fprintf(os.Stderr, "WARN:  "+msg+"\n", a...)
	}
}
func Print(msg string, a ...any) {
	fmt.Printf(msg+"\n", a...)
}

type GlobalOptions struct {
	FilePath string `long:"lib" description:"Library file path" default:"tlib.json"`
	Debug    bool   `long:"debug" description:"Print debug info"`
	Quiet    bool   `short:"q" long:"quiet" description:"Only print error messages"`
	Version  bool   `short:"v" long:"version" description:"Print tlib version"`

	Pretty          bool `short:"p" long:"pretty" description:"Write library JSON in human-readable format"`
	AlwaysPretty    bool `long:"always-pretty" description:"Update JSON library to always save pretty"`
	NotAlwaysPretty bool `long:"not-always-pretty" description:"Undo --always-pretty"`
}

type Options struct {
	AddOpts   AddOptions
	InitOpts  InitOptions
	ShelfOpts ShelfOptions
}

type EmptyOptions struct{}
