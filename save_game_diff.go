package main

// E:\projects\save_game_diff>go run save_game_diff.go --files "C:\Program Files (x86)\Steam\userdata\33882143\239160\remote\Thief7C.sav" "C:\Program Files (x86)\Steam\userdata\33882143\239160\remote\Thief8C.sav" "C:\Program Files (x86)\Steam\userdata\33882143\239160\remote\Thief9C.sav" --values 2 1 0
// E:\projects\save_game_diff>go run save_game_diff.go --files "C:\Program Files (x86)\Steam\userdata\33882143\239160\remote\Thief7I.sav" "C:\Program Files (x86)\Steam\userdata\33882143\239160\remote\Thief8I.sav" "C:\Program Files (x86)\Steam\userdata\33882143\239160\remote\Thief9I.sav" --values 2 1 0
// E:\projects\save_game_diff>go run save_game_diff.go --files "C:\Program Files (x86)\Steam\userdata\33882143\239160\remote\Thief7I.sav" "C:\Program Files (x86)\Steam\userdata\33882143\239160\remote\Thief8I.sav" "C:\Program Files (x86)\Steam\userdata\33882143\239160\remote\Thief9I.sav" --decreased_value_by 10
// E:\projects\save_game_diff>go run save_game_diff.go --files "C:\Program Files (x86)\Steam\userdata\33882143\239160\remote\Thief7I.sav" "C:\Program Files (x86)\Steam\userdata\33882143\239160\remote\Thief8I.sav" "C:\Program Files (x86)\Steam\userdata\33882143\239160\remote\Thief9I.sav" --increased_value_by 10

// TODO need help for each mode (with examples)
// mode0 --values
// mode 1 --decreased_value
// mode 1 --decreased_value_by
// mode 2 --increased_value
// mode 2 --increased_value_by

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
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
	log.Println("\t\t\t files to compare (count of the files must be the same as count of the values)")
	log.Println()
	log.Println("\t--values <value1> <value2> [...]")
	log.Println("\t\t\t values to search in each file")
	log.Println()
	log.Println("\t--help")
	log.Println("\t\t\t this help")
	log.Println()
	log.Println("Example (each save file has less \"health\" in game):")
	log.Println("go run save_game_diff.go --files \"save1.sav\" \"save2.sav\" \"save3.sav\" --values 100 70 30")
	log.Println()
}

func hasArgument(name string) bool {
	for i, argument := range os.Args {
		if i < 0 {
			continue
		}

		if argument == name {
			return true
		}
	}

	return false
}

func shouldPrintUsages() bool {
	lenArgs := len(os.Args)

	if lenArgs == 1 {
		return true
	}

	if hasArgument("--help") {
		return true
	}

	hasFiles := hasArgument("--files")
	hasValues := hasArgument("--values")
	hasDecreasedValue := hasArgument("--decreased_value")
	hasDecreasedValueBy := hasArgument("--decreased_value_by")
	hasIncreasedValue := hasArgument("--increased_value")
	hasIncreasedValueBy := hasArgument("--increased_value_by")

	if !hasFiles {
		return true
	}

	return !hasValues && !hasDecreasedValue && !hasDecreasedValueBy && !hasIncreasedValue && !hasIncreasedValueBy
}

func changeCurrentWorkingDir() {
	os.Chdir(exeDir)
}

func getArgumentValues(name string, allowNoArgument bool) []string {
	values := make([]string, 0)

	index := funk.IndexOfString(os.Args, name)

	if index < 0 {
		if allowNoArgument {
			return values
		}

		log.Panicln("No argument", name)
	}

	index++

	for ; index < len(os.Args); index++ {
		value := os.Args[index]

		if strings.HasPrefix(value, "--") {
			break
		}

		values = append(values, value)
	}

	return values
}

func getArgumentValuesAsBytes(name string, allowNoArgument bool) []byte {
	values := getArgumentValues(name, allowNoArgument)
	byteValues := make([]byte, 0)

	for _, sValue := range values {
		b, err := strconv.ParseInt(sValue, 10, 64)

		if err != nil {
			log.Panicln(sValue, "is not a byte")
		}

		if b < 0 || b > 255 {
			log.Panicln(sValue, "is not a byte")
		}

		byteValues = append(byteValues, byte(b))
	}

	return byteValues
}

func arrayIsDiff(data []byte) bool {
	lenData := len(data)

	for i := 0; i < lenData; i++ {
		for _, b := range data {
			if data[i] != b {
				return true
			}
		}
	}

	return false
}

func arrayContainsArray(valuesToFind []byte, inArray []byte) bool {
	count := 0

	for _, b := range valuesToFind {
		for _, b2 := range inArray {
			if b2 == b {
				count++
				break
			}
		}
	}

	return count == len(inArray)
}

func arraysEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

func getDiffrences(files []string) map[int64][]byte {
	var differences map[int64][]byte
	var handles []*os.File

	differences = make(map[int64][]byte)

	for _, pathname := range files {
		ihandle, err := os.Open(pathname)

		if err != nil {
			log.Println("Cannot open file:", err)
			break
		}

		handles = append(handles, ihandle)
	}

	offset := int64(0)

	for {
		bytesAtOffset := make([]byte, 0)
		offset++

		for _, ihandle := range handles {
			if _, err := ihandle.Seek(offset-1, io.SeekStart); err != nil {
				continue
			}

			b := make([]byte, 1)

			n, err := ihandle.Read(b)

			if err != nil {
				continue
			}

			if n < 1 {
				continue
			}

			bytesAtOffset = append(bytesAtOffset, b[0])
		}

		lenBytesAtOffset := len(bytesAtOffset)

		if lenBytesAtOffset <= 1 {
			break
		}

		if arrayIsDiff(bytesAtOffset) {
			differences[offset-1] = bytesAtOffset
		}
	}

	for _, ihandle := range handles {
		ihandle.Close()
	}

	return differences
}

func mode0_Values() {
	files := getArgumentValues("--files", false)
	values := getArgumentValuesAsBytes("--values", false)

	if len(files) != len(values) {
		printAppInfo()
		printUsages()

		return
	}

	log.Println("Files to compare:", files)
	log.Println("Values to search:", values)
	log.Println("Looking for differences")

	differences := getDiffrences(files)

	log.Println("Count of differenes:", len(differences))

	for offset, data := range differences {
		if arrayContainsArray(values, data) {
			perfectMatch := arraysEqual(values, data)

			if perfectMatch {
				printMatch(offset, data, "[PERFECT MATCH]")
			} else {
				printMatch(offset, data)
			}
		}
	}
}

func printMatch(offset int64, data []byte, additionalText ...string) {
	additional := ""

	if len(additionalText) > 0 {
		additional = fmt.Sprintf("%v", additionalText)
	}

	log.Printf("Match at offset %v (0x%X): %v %v", offset, offset, data, additional)
}

func mode1_DecreasedValue() {
	files := getArgumentValues("--files", false)
	byValue := getArgumentValuesAsBytes("--decreased_value_by", true)

	if len(files) < 2 {
		printAppInfo()
		printUsages()

		return
	}

	lenByValue := len(byValue)

	if lenByValue > 0 {
		if lenByValue != 1 {
			printAppInfo()
			printUsages()

			return
		}
	}

	log.Println("Files to compare:", files)

	if lenByValue > 0 {
		log.Println("Must decrease by:", byValue[0])
	}

	log.Println("Looking for differences")

	differences := getDiffrences(files)

	log.Println("Count of differenes:", len(differences))

	for offset, data := range differences {
		if arrayDecreasedValues(data, byValue...) {
			printMatch(offset, data)
		}
	}
}

// From left to right
func arrayDecreasedValues(data []byte, byValue ...byte) bool {
	lenData := len(data)
	hasByValue := len(byValue) == 1

	for i, value := range data {
		if i+1 == lenData {
			break
		}

		if value <= data[i+1] {
			return false
		}

		if hasByValue {
			if value-data[i+1] != byValue[0] {
				return false
			}
		}
	}

	return true
}

// From right to left
func arrayIncreasedValues(data []byte, byValue ...byte) bool {
	hasByValue := len(byValue) == 1

	for i := len(data) - 1; i >= 0; i-- {
		if i == 0 {
			break
		}

		value := data[i]

		if data[i-1] >= value {
			return false
		}

		if hasByValue {
			if value-data[i-1] != byValue[0] {
				return false
			}
		}
	}

	return true
}

func _main() {
	hasValues := hasArgument("--values")
	hasDecreasedValue := hasArgument("--decreased_value")
	hasDecreasedValueBy := hasArgument("--decreased_value_by")
	hasIncreasedValue := hasArgument("--increased_value")
	hasIncreasedValueBy := hasArgument("--increased_value_by")

	if hasValues {
		mode0_Values()
	} else if hasDecreasedValue || hasDecreasedValueBy {
		mode1_DecreasedValue()
	} else if hasIncreasedValue || hasIncreasedValueBy {
		mode2_IncreasedValue()
	}
}

func mode2_IncreasedValue() {
	files := getArgumentValues("--files", false)
	byValue := getArgumentValuesAsBytes("--increased_value_by", true)

	if len(files) < 2 {
		printAppInfo()
		printUsages()

		return
	}

	lenByValue := len(byValue)

	if lenByValue > 0 {
		if lenByValue != 1 {
			printAppInfo()
			printUsages()

			return
		}
	}

	log.Println("Files to compare:", files)

	if lenByValue > 0 {
		log.Println("Must increase by:", byValue[0])
	}

	log.Println("Looking for differences")

	differences := getDiffrences(files)

	log.Println("Count of differenes:", len(differences))

	for offset, data := range differences {
		if arrayIncreasedValues(data, byValue...) {
			printMatch(offset, data)
		}
	}
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
