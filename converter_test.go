package main

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

func TestViewConverter(t *testing.T) {
	f, err := os.Open("./files/dummy.txt")
	check(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var sb strings.Builder
	for scanner.Scan() {
		sb.Write([]byte(scanner.Text() + " " + newLine))
	}

	rawSQL := sb.String()
	sql := view(rawSQL)
	if sql == "" {
		t.Log("error")
	}
}
