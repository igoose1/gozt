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

func a() ([]string, error) {
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
	result := make([]byte, 0, 1<<5*(2 + len(linkNumber) + len(graph)))

	result = append(result, "graph zt {\n"...)
	for i, ln := range linkNumber {
		penwidth := 15 * float64(ln) / float64(maxPrivateConnections)
		if penwidth < 1 {
			penwidth = 1
		}
		attrs := fmt.Sprintf("\t%04d [penwidth = %f]\n", i, penwidth)
		result = append(result, attrs...)
	}

	for e := range graph {
		result = append(
			result,
			fmt.Sprintf("\t%04d -- %04d;\n", e.from, e.to)...,
		)
	}
	result = append(result, "}"...)
	return string(result)
}

func g() (string, error) {
	allNotes, err := a()
	if err != nil {
		return "", err
	}

	maxPrivateConnections := 0
	graph := make(map[edge]struct{})
	linkNumber := make([]int, len(allNotes))

	for _, noteName := range allNotes {
		data, err := ioutil.ReadFile(noteName)
		if err != nil {
			return "", err
		}
		text := string(data)
		for _, linkedTo := range validNoteLink.FindAllString(text, -1) {
			from, _ := strconv.ParseInt(noteName, 10, 64)
			to, _ := strconv.ParseInt(linkedTo[1:], 10, 64)

			if from > to {
				from, to = to, from
			}
			log.Print(from, to)

			graph[edge{from, to}] = struct{}{}
			linkNumber[from]++
			linkNumber[to]++

			if linkNumber[from] > maxPrivateConnections {
				maxPrivateConnections = linkNumber[from]
			}
			if linkNumber[to] > maxPrivateConnections {
				maxPrivateConnections = linkNumber[to]
			}
		}
	}

	result := drawDot(graph, linkNumber, maxPrivateConnections)

	return result, err
}

func nextToLast(delta int) (string, error) {
	notes, err := a()
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
		notes, err := a()
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
