package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type Node struct {
	data  []string
	left  *Node
	right *Node
}
type Tree struct {
	root *Node
}

var (
	readfrom       *os.File
	sorted         [][]string
	inputFileName  = flag.String("i", "", "Use a file with the name file-name as an input")
	outputFileName = flag.String("o", "", "Use a file with the name file-name as an output")
	headerFlag     = flag.Bool("h", false, "Remove headers from sorting")
	fieldFlag      = flag.Int("f", 0, "Sort input lines by value number N.")
	reverseFlag    = flag.Bool("r", false, "Sort input lines in reverse order")
	algorithmFlag  = flag.Int("a", 1, "Sorting algorithm: 1 - built in, 2 - Tree Sort")
)

func main() {
	flag.Parse()

	readfrom = inputAsFile()
	header := *headerFlag
	field := *fieldFlag
	reverse := *reverseFlag
	sortAlgorithm := *algorithmFlag
	content := readContent(readfrom)

	sortContent(content, header, field, reverse, sortAlgorithm)
	output(sorted)
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func inputAsFile() *os.File {
	var answer *os.File
	if isFlagPassed("i") {
		f, err := os.Open(*inputFileName)
		if err != nil {
			log.Fatal(err)
		}
		answer = f
		defer f.Close()
	} else {
		answer = os.Stdin
	}
	return answer
}

func output(text [][]string) {
	if isFlagPassed("o") {
		f, err := os.Create(*outputFileName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintln(f, text)
		fmt.Printf("Output is written to file %s\n", *outputFileName)
		defer f.Close()
	} else {
		fmt.Printf("Result: %v\n", text)
	}
}

func readContent(readfrom *os.File) [][]string {
	n := 0
	content := [][]string{}
	s := bufio.NewScanner(readfrom)

	if s.Err() != nil {
		log.Fatal(s.Err())
	}

	for s.Scan() {
		line := s.Text()
		row := strings.Split(line, ",")
		if line == "" {
			break
		}
		if n == 0 {
			n = len(row)
		}
		if n != len(row) {
			log.Fatal("ERROR: The number of columns is not equal to the number of rows")
		}
		content = append(content, row)
	}
	return content
}

func sortContent(content [][]string, header bool, field int, reverse bool, sortAlgorithm int) {
	h := 0
	if header {
		h = 1
	}
	switch sortAlgorithm {
	case 1:
		sort.Slice(content[h:], func(i, j int) bool {
			if reverse {
				return content[i+h][field] > content[j+h][field]
			}
			return content[i+h][field] < content[j+h][field]
		})
		sorted = content
	case 2:
		// tree sort
		t := &Tree{}
		for i := h; i < len(content); i++ {
			t.insert(content[i], field)
		}
		t.root.rewriteTree()
	}
}

func (t *Tree) insert(data []string, field int) *Tree {
	if t.root == nil {
		t.root = &Node{data: data, left: nil, right: nil}
	} else {
		t.root.insert(data, field)
	}
	return t
}

func (n *Node) insert(data []string, field int) {
	if n == nil {
		return
	} else if data[field] <= n.data[field] {
		if n.left == nil {
			n.left = &Node{data: data, left: nil, right: nil}
		} else {
			n.left.insert(data, field)
		}
	} else {
		if n.right == nil {
			n.right = &Node{data: data, left: nil, right: nil}
		} else {
			n.right.insert(data, field)
		}
	}
}

func (node *Node) rewriteTree() {
	if node.left != nil {
		node.left.rewriteTree()
	}
	sorted = append(sorted, node.data)
	if node.right != nil {
		node.right.rewriteTree()
	}
}
