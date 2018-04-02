package main

import (
	"HelloWorld/pipeline"
	"fmt"
	"os"
	"bufio"
)

// generate a source file and print the 100 rows of file
func main() {
	const filename  = "small.in"
	const n = 64
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	p := pipeline.RandomSource(n)
	// input or output with buffer
	writer := bufio.NewWriter(file)
	pipeline.WriterSink(writer, p)
	writer.Flush()

	file, err = os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	p = pipeline.ReaderSource(bufio.NewReader(file), -1)

	// only output 100 rows
	count := 0
	for v := range p {
		fmt.Println(v)
		count ++
		if count > 100 {
			break
		}
	}
}

func mergeDemo() {
	p := pipeline.Merge(pipeline.InMemSort(pipeline.ArraySource(1,4,5,6,2143,11)),
		pipeline.InMemSort(pipeline.ArraySource(1,4,5,6,2143,11)))
	for v := range p {
		fmt.Println(v)
	}
}
