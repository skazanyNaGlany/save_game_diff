package main

// E:\projects\save_game_diff>go run save_game_diff.go --files "C:\Program Files (x86)\Steam\userdata\33882143\239160\remote\Thief7C.sav" "C:\Program Files (x86)\Steam\userdata\33882143\239160\remote\Thief8C.sav" "C:\Program Files (x86)\Steam\userdata\33882143\239160\remote\Thief9C.sav" --values 2 1 0

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/thoas/go-funk"
)

const AppName = "SAVE_GAME_DIFF"
const AppVersion = "0.1"

var exeDir = filepath.Dir(os.Args[0])

func duplicateLog() {
	logFilename := filepath.Base(os.Args[0]) + ".txt"
	logFile, err := os.OpenFile(logFilename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)

	if err != nil {
		panic(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)

	log.SetOutput(mw)
}

func getFullAppName() string {
	return fmt.Sprintf("%v v%v", AppName, AppVersion)
}

func printAppName() {
	log.Println(
		getFullAppName())
	log.Println()
	log.Println("Compare each file to find specific differences.")
	log.Println()
}

func printAppInfo() {
}

func printUsages() {
	log.Printf("Usage: %v --files <file1> <file2> [...] --values <value1> <value2> [...]", os.Args[0])

	log.Println()
	log.Println("Options:")

	log.Println("\t--files <file1> <file2> [...]")
	log.Println("\t\t\t files to compare")
	log.Println()
	log.Println("\t--values <value1> <value2> [...]")
	log.Println("\t\t\t values to search in each file")
	log.Println()
	log.Println("\t--help")
	log.Println("\t\t\t this help")
	log.Println()
}

func shouldPrintUsages() bool {
	lenArgs := len(os.Args)

	if lenArgs == 1 {
		return true
	}

	if lenArgs < 6 {
		return true
	}

	hasFiles := false
	hasValues := false

	for i, arg := range os.Args {
		if i < 0 {
			continue
		}

		if arg == "--files" {
			hasFiles = true
		} else if arg == "--values" {
			hasValues = true
		} else if arg == "--help" {
			return true
		}

		if hasFiles && hasValues {
			return false
		}
	}

	return true
}

func changeCurrentWorkingDir() {
	os.Chdir(exeDir)
}

func getFiles() []string {
	files := make([]string, 0)

	index := funk.IndexOfString(os.Args, "--files") + 1

	for ; index < len(os.Args); index++ {
		value := os.Args[index]

		if strings.HasPrefix(value, "--") {
			break
		}

		files = append(files, value)
	}

	return files
}

func getValues() []string {
	values := make([]string, 0)

	index := funk.IndexOfString(os.Args, "--values") + 1

	for ; index < len(os.Args); index++ {
		value := os.Args[index]

		if strings.HasPrefix(value, "--") {
			break
		}

		values = append(values, value)
	}

	return values
}

func _main() {
	files := getFiles()
	values := getValues()

	if len(files) != len(values) {
		printAppInfo()
		printUsages()

		return
	}

	log.Println("Files to compare:", files)
	log.Println("Values to search:", values)
}

func main() {
	changeCurrentWorkingDir()
	duplicateLog()
	printAppName()

	if shouldPrintUsages() {
		printAppInfo()
		printUsages()

		os.Exit(1)
	}

	_main()
}
