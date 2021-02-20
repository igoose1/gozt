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

type set map[int64]struct{}

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

func drawDot(graph []set, connections int, maxConnections int) string {
	result := make([]byte, 0, 1<<5+2*connections*(5+4))
	result = append(result, "graph zt {\n"...)
	for from, linkedTo := range graph {
		penwidth := 15 * float64(len(linkedTo)) / float64(maxConnections)
		if penwidth < 1 {
			penwidth = 1
		}
		attrs := fmt.Sprintf("\t%04d [penwidth = %f]\n", from, penwidth)
		result = append(result, attrs...)
		if len(linkedTo) == 0 {
			continue
		}
		result = append(result, fmt.Sprintf("\t%04d", from)...)
		for to := range linkedTo {
			result = append(result, fmt.Sprintf(" -- %04d", to)...)
		}
		result = append(result, ";\n"...)
	}
	result = append(result, "}"...)
	return string(result)
}

func g() (string, error) {
	allNotes, err := a()
	if err != nil {
		return "", err
	}

	graph := make([]set, len(allNotes))
	for i := 0; i < len(allNotes); i++ {
		graph[i] = make(set)
	}

	connections, maxConnections := 0, 0
	var exists = struct{}{}
	for _, noteName := range allNotes {
		data, err := ioutil.ReadFile(noteName)
		if err != nil {
			return "", err
		}
		text := string(data)
		for _, linkedTo := range validNoteLink.FindAllString(text, -1) {
			from, _ := strconv.ParseInt(noteName, 10, 32)
			to, _ := strconv.ParseInt(linkedTo[1:], 10, 32)

			connections++
			graph[from][to] = exists
			graph[to][from] = exists
			if len(graph[from]) > maxConnections {
				maxConnections = len(graph[from])
			}
			if len(graph[to]) > maxConnections {
				maxConnections = len(graph[to])
			}
		}
	}

	result := drawDot(graph, connections, maxConnections)

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
