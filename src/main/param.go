package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	help   bool
	source string
	target string
	split  string
)

func init() {
	flag.BoolVar(&help, "h", false, "this help")

	flag.StringVar(&source, "s", "", "parse xlsx source path")
	flag.StringVar(&target, "t", "", "parsed csv path")
	flag.StringVar(&split, "p", ",", "split sign, default ,")

	// 改变默认的 Usage
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `xlsx2csv version: 0.0.1
Usage: xlsx2csv [-help] [-s sourcepath] [-t targetpath] [-p splitsign]

Options:\
`)
	flag.PrintDefaults()
}

func ParseArgs() {

	flag.Parse()

	if help {
		flag.Usage()
		return
	} else if source != "" {
		sourcePath = source
	}
	if target != "" {
		targetPath = target
	}
	if split != "" {
		splitSign = split
	}

}
