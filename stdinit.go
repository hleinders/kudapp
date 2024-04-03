package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/fatih/color"
	"golang.org/x/term"
)

// color funcs
// var mkNormal = color.New(color.Reset).SprintFunc()

var mkBold = color.New(color.Bold).SprintFunc()

// var mkFaint = color.New(color.Faint).SprintFunc()
// var mkItalic = color.New(color.Italic).SprintFunc()
var mkUnderline = color.New(color.Underline).SprintFunc()

// var mkStrike = color.New(color.CrossedOut).SprintFunc()
// var mkBlink = color.New(color.BlinkSlow).SprintFunc()

var mkRed = color.New(color.FgRed).SprintFunc()
var mkGreen = color.New(color.FgGreen).SprintFunc()
var mkYellow = color.New(color.FgYellow).SprintFunc()

// var mkBlue = color.New(color.FgBlue).SprintFunc()
// var mkMagenta = color.New(color.FgMagenta).SprintFunc()
// var mkCyan = color.New(color.FgCyan).SprintFunc()
// var mkWhite = color.New(color.FgWhite).SprintFunc()

var bulletChar = "•"

// Application Details
const (
	appName    = "KuDAPP"
	appVersion = "1.2 (2024-04-03)"
	appAuthor  = "Harald Leinders"
	appEMail   = "harald@leinders.de"
	maxLineLen = 80
	maxDecLen  = 50
)

// globla vars
var programName string

// IsTTY checks for interactive terminal
func IsTTY() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// HasColor checks color capabilities
func HasColor() bool {
	return os.Getenv("TERM") != "dumb" && IsTTY()
}

// GetSize retruen the current terminal size
func GetSize() (int, int, error) {
	return term.GetSize(int(os.Stdout.Fd()))
}

func IsColorTerm() bool {
	noColor := false

	if HasColor() {
		noColor = false
	} else {
		noColor = true
	}

	if !IsTTY() {
		noColor = true
	}

	return noColor
}
func init() {
	// Use available Cores for Goroutines
	runtime.GOMAXPROCS(runtime.NumCPU())

	programName = filepath.Base(os.Args[0])

	// Detect color mode
	color.NoColor = IsColorTerm()

	// Detect locale for printing:
	if runtime.GOOS == "windows" {
		bulletChar = "*"
	}

	defer handleExit()
}

type Exit struct{ Code int }

// exit code handler honoring deferred calls
func handleExit() {
	if e := recover(); e != nil {
		if exit, ok := e.(Exit); ok {
			os.Exit(exit.Code)
		}
		panic(e) // not an Exit, bubble up
	}
}

func check(e error, rcode int) {
	if e != nil {
		fmt.Fprintf(os.Stderr, mkRed("*** Error: %+v\n"), e)
		panic(Exit{Code: rcode})
	}
}
