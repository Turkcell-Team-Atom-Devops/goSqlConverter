package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

const (
	newLine                      string = "\r\n"
	viewHeaderCreateOrReplace    string = "CREATE VIEW"
	viewHeaderCreateOrReplaceNew string = "CREATE OR REPLACE VIEW"
	viewFunctionIsNull           string = "ISNULL"
	viewFunctionCoalesce         string = "COALESCE"
)

func main() {
	fptr := flag.String("fpath", "./files/dummy.txt", "file path to read from")
	flag.Parse()
	f, err := os.Open(*fptr)
	check(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var sb strings.Builder
	for scanner.Scan() {
		sb.Write([]byte(scanner.Text() + " " + newLine))
	}

	rawSQL := sb.String()
	_ = view(rawSQL)
}

func view(rawSQL string) string {

	lines := strings.Split(rawSQL, newLine)
	var sb strings.Builder
	for _, v := range lines {
		temp := v
		// Header
		if strings.ContainsAny(temp, viewHeaderCreateOrReplace) {
			temp = strings.Replace(temp, viewHeaderCreateOrReplace, viewHeaderCreateOrReplaceNew, 1)
		}

		// ][ change
		if strings.ContainsAny(temp, "[") || strings.ContainsAny(temp, "]") {
			temp = strings.Replace(temp, "[", "\"", -1)
			temp = strings.Replace(temp, "]", "\"", -1)
		}

		if strings.ContainsAny(temp, ".dbo") {
			temp = strings.ReplaceAll(temp, ".dbo", "")
		}

		if strings.ContainsAny(temp, ".") {
			reg := regexp.MustCompile(`\.(\w+)`)
			words := reg.FindAllString(temp, -1)
			for _, v := range words {
				word := strings.Replace(v, ".", ".\"", 1)
				word = word + "\""
				temp = strings.Replace(temp, v, word, 1)
			}
		}

		if strings.ContainsAny(temp, ",") {
			reg := regexp.MustCompile(`(\w+),`)
			word := reg.FindString(temp)
			word = strings.Replace(word, ",", "\",", 1)
			word = "\"" + word
			temp = reg.ReplaceAllString(temp, word)
		}

		if strings.ContainsAny(strings.ToUpper(temp), viewFunctionIsNull) {
			temp = strings.Replace(temp, viewFunctionIsNull, viewFunctionCoalesce, 1)
		}

		if strings.ContainsAny(temp, "WITH (NOLOCK)") {
			temp = strings.ReplaceAll(temp, "WITH (NOLOCK)", "")
		}

		if strings.ContainsAny(temp, "COLLATE") {
			temp = strings.ReplaceAll(temp, "COLLATE", "")
		}

		if strings.ContainsAny(temp, "<>") {
			temp = strings.ReplaceAll(temp, "<>", "!=")
		}

		if strings.ContainsAny(temp, "+") {
			temp = strings.ReplaceAll(temp, "+", "||")
		}

		sb.Write([]byte(temp + " " + newLine))
	}

	fmt.Println(sb.String())
	return sb.String()
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
		panic(e)
	}
}
