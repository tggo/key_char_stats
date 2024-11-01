package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <dir>")
		return
	}

	// get dir from command line first argument
	dir := os.Args[1]

	var str string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			str += strings.ToLower(string(content))
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error reading files:", err)
		return
	}

	charCount := make(map[rune]int)
	nextCharCount := make(map[rune]map[rune]int)

	for i, char := range str {
		charCount[char]++
		if i < len(str)-1 {
			nextChar := rune(str[i+1])
			if nextCharCount[char] == nil {
				nextCharCount[char] = make(map[rune]int)
			}
			nextCharCount[char][nextChar]++
		}
	}

	type charInfo struct {
		char  rune
		count int
	}

	var charList []charInfo
	for char, count := range charCount {
		charList = append(charList, charInfo{char, count})
	}

	sort.Slice(charList, func(i, j int) bool {
		return charList[i].count > charList[j].count
	})

	fmt.Println("Character counts (sorted by count):")
	for _, ci := range charList {
		buf := ""

		buf += fmt.Sprintf("%s%s%s: %d \t ", Yellow, convertCharToString(ci.char), Reset, ci.count)
		counted := 0

		if nextChars, exists := nextCharCount[ci.char]; exists {
			var nextCharList []charInfo
			for nextChar, count := range nextChars {
				nextCharList = append(nextCharList, charInfo{nextChar, count})
			}

			sort.Slice(nextCharList, func(i, j int) bool {
				return nextCharList[i].count > nextCharList[j].count
			})

			idx := 0
			for _, nc := range nextCharList {
				if nc.count < 3 {
					continue
				}

				idx++

				// return only top 5
				if idx > 5 {
					continue
				}

				buf += fmt.Sprintf("%s%s%s:%s%d%s  ", Yellow, convertCharToString(nc.char), Reset, Green, nc.count, Reset)
				counted += nc.count
			}
		}

		if counted <= 2 {
			continue
		}
		fmt.Printf("%s\n", buf)
	}

}

func convertCharToString(char rune) string {
	if char == '\t' {
		return "↹"
	}
	if char == '\n' {
		return "↲"
	}
	if char == '\r' {
		return "↲"
	}
	if char == ' ' {
		return "␣"
	}
	return string(char)
}
