package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

const Version = "0.1-alpha"

const (
	CmdAdd   = "add"
	CmdInit  = "init"
	CmdShelf = "shelf"
)

var GlobalOpts GlobalOptions

func main() {
	var opts Options
	p := flags.NewParser(&GlobalOpts, flags.Default)

	if _, err := p.AddCommand(CmdInit, "Initialize a new library", "Initialize a new library", &opts.InitOpts); err != nil {
		panic(err)
	}
	if _, err := p.AddCommand(CmdAdd, "Add books to library", "Add books to library", &opts.AddOpts); err != nil {
		panic(err)
	}
	if _, err := p.AddCommand(CmdShelf, "Print books in shelf order", "Print books in shelf order", &opts.ShelfOpts); err != nil {
		panic(err)
	}

	// p.ParseArgs(strings.Split("shelf --find-all", " "))
	// GlobalOpts.FilePath = "/home/josh/src/tlib"
	p.Parse()

	if p.Active == nil {
		p.WriteHelp(os.Stdout)
		return
	}

	switch p.Active.Name {
	case CmdAdd:
		RunAdd(&opts.AddOpts)
	case CmdInit:
		RunInit(&opts.InitOpts)
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
	fmt.Fprintf(os.Stderr, msg+"\n", a...)
}

type GlobalOptions struct {
	FilePath string `short:"d" long:"dir" description:"Library directory"`
	Pretty   bool   `short:"p" long:"pretty" description:"Write library JSON in human-readable format"`
	Debug    bool   `long:"debug" description:"Print debug info"`
	Quiet    bool   `short:"q" long:"quiet" description:"Only print error messages"`
}

type Options struct {
	AddOpts   AddOptions
	InitOpts  InitOptions
	ShelfOpts ShelfOptions
}
