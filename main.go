package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
)

var (
	Cwd           string
	validNoteName = regexp.MustCompile(`\d{4}`)
	validNoteLink = regexp.MustCompile(`\$\d{4}`)
)

const notesLessThan = 1 << 14 // log2(10 ^ 4)

func isSynopsis(note string) bool {
	return note[0] == '9'
}

func a(withSynopsis bool) ([]string, error) {
	dir, err := os.Open(Cwd)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	notes := make([]string, 0, notesLessThan)
	for {
		files, err := dir.Readdir(notesLessThan)
		if err != nil {
			if err == io.EOF {
				break
			}
			sort.Strings(notes)
			return notes, err
		}
		for _, file := range files {
			if !file.IsDir() && validNoteName.MatchString(file.Name()) {
				if !withSynopsis && isSynopsis(file.Name()) {
					continue
				}
				notes = append(notes, file.Name())
			}
		}
	}
	sort.Strings(notes)
	return notes, dir.Close()
}

type edge struct {
	from, to int64
}

func drawDot(graph map[edge]struct{}, linkNumber []int, maxPrivateConnections int) string {
	result := make([]byte, 0, 1<<5*(2+len(linkNumber)+len(graph)))

	result = append(
		result,
		"digraph zt {\n"+
			"\tgraph [truecolor=true bgcolor=pink model=subset]\n"+
			"\tnode [colorscheme=ylgn9 style=filled shape=circle];\n"...,
	)
	for i, ln := range linkNumber {
		if ln == -1 {
			continue
		}
		color := int(8 * float64(ln) / float64(maxPrivateConnections))
		if color < 1 {
			color = 1
		}
		result = append(
			result,
			fmt.Sprintf("\t%04d [color=%d];\n", i, color)...,
		)
	}

	for e := range graph {
		result = append(
			result,
			fmt.Sprintf("\t%04d -> %04d;\n", e.from, e.to)...,
		)
	}
	result = append(result, "}"...)
	return string(result)
}

func contains(where []string, what string) bool {
	for _, element := range where {
		if element == what {
			return true
		}
	}
	return false
}

// Do not play with that a lot, child! It's O(N^2).
func unique(array []string) []string {
	result := make([]string, 0, len(array))
	for _, element := range array {
		if !contains(result, element) {
			result = append(result, element)
		}
	}
	return result
}

func maxStringSlice(array []string) string {
	ans_index := 0
	for i, value := range array {
		if value > array[ans_index] {
			ans_index = i
		}
	}
	return array[ans_index]
}

func maxIntSlice(array []int) int {
	ans_index := 0
	for i, value := range array {
		if value > array[ans_index] {
			ans_index = i
		}
	}
	return array[ans_index]
}

func g() (string, error) {
	allNotes, err := a(true)
	if err != nil {
		return "", err
	}

	graph := make(map[edge]struct{})
	maxNode, _ := strconv.ParseInt(maxStringSlice(allNotes), 10, 64)
	linkNumber := make([]int, maxNode+1)
	for i := range linkNumber {
		linkNumber[i] = -1
	}
	for _, note := range allNotes {
		noteNumber, _ := strconv.ParseInt(note, 10, 64)
		linkNumber[noteNumber] = 0
	}

	for _, noteName := range allNotes {
		data, err := ioutil.ReadFile(noteName)
		if err != nil {
			return "", err
		}
		text := string(data)
		for _, linkedTo := range unique(validNoteLink.FindAllString(text, -1)) {
			from, _ := strconv.ParseInt(linkedTo[1:], 10, 64)
			to, _ := strconv.ParseInt(noteName, 10, 64)

			graph[edge{from, to}] = struct{}{}
			linkNumber[from]++
		}
	}

	result := drawDot(graph, linkNumber, maxIntSlice(linkNumber))

	return result, err
}

func nextToLast(delta int) (string, error) {
	notes, err := a(false)
	if err != nil {
		return "", err
	}
	number, err := strconv.ParseInt(notes[len(notes)-1], 10, 32)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%04d", int(number)+delta), nil
}

func main() {
	var err error

	Cwd, err = os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	cmdName := path.Base(os.Args[0])
	switch cmdName {
	case "zt-a":
		notes, err := a(true)
		if err != nil {
			log.Fatalln(err)
		}
		for _, note := range notes {
			fmt.Println(note)
		}
	case "zt-l":
		note, err := nextToLast(0)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(note)
	case "zt-ll":
		note, err := nextToLast(-1)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(note)
	case "zt-n":
		note, err := nextToLast(1)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(note)
	case "zt-g":
		graph, err := g()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(graph)
	default:
		log.Fatalln("unknown command", cmdName)
	}
}
