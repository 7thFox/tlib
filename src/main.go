package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

const Version = "0.1-alpha"

const (
	CmdAdd  = "add"
	CmdInit = "init"
)

var GlobalOpts GlobalOptions

func Debug(msg string, a ...any) {
	fmt.Fprintf(os.Stderr, "DEBUG: "+msg+"\n", a...)
}

func main() {
	var opts Options
	p := flags.NewParser(&GlobalOpts, flags.Default)

	if _, err := p.AddCommand(CmdInit, "Initialize a new library", "Initialize a new library", &opts.InitOpts); err != nil {
		panic(err)
	}
	if _, err := p.AddCommand(CmdAdd, "Add books to library", "Add books to library", &opts.AddOpts); err != nil {
		panic(err)
	}

	// p.ParseArgs([]string{"add", "9781594037306"})
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
	default:
		p.WriteHelp(os.Stdout)
	}
}

type GlobalOptions struct {
	FilePath string `short:"d" long:"dir" description:"Library directory"`
	Pretty   bool   `short:"p" long:"pretty" description:"Write library JSON in human-readable format"`
}

type Options struct {
	AddOpts  AddOptions
	InitOpts InitOptions
}
