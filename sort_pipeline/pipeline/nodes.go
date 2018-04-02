package pipeline

import (
	"sort"
	"io"
	"encoding/binary"
	"math/rand"
)

// Get data from array , then put them into channel and return channel
func ArraySource(a ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, v := range a {
			out <- v
		}
		close(out)
	}()
	return out
}

// Sort data from channel, then put sorted data to the channel ,and return it
func InMemSort(in <-chan int) <-chan int {
	out := make(chan  int)
	go func() {
		// Read into memory
		a := []int{}
		for v := range in {
			a = append(a, v)
		}
		// Sort
		sort.Ints(a)
		//Output
		for _,v :=  range a {
			out <- v
		}
		close(out)
	}()
	return out
}

// Merge 2 sorted channel
func Merge(in1, in2 <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		v1, ok1 := <-in1
		v2, ok2 := <-in2
		for ok1 || ok2 {
			if !ok2 || (ok1 && v1 <= v2) {
				out <- v1
				v1, ok1 = <-in1
			}else {
				out <- v2
				v2, ok2 = <-in2
			}
		}
		defer close(out)
	}()
	return out
}

// Read data from file source, then put the data into channel
func ReaderSource(reader io.Reader, chunkSize int) <-chan int {
	out := make(chan int)
	go func() {
		buffer := make([]byte, 8)

		byteRead := 0
		for {
			n, err := reader.Read(buffer)
			byteRead += n
			if n > 0 {
				// Deserialize buffer data to int type, then send into channel
				v := int(binary.BigEndian.Uint64(buffer))
				out <- v
			}
			if err != nil || (chunkSize != -1 && byteRead >= chunkSize) {
				break
			}
		}
		defer close(out)
	}()
	return out
}

func WriterSink(writer io.Writer, in <-chan int) {
	for v := range in {
		buffer := make([]byte, 8)
		// Serialize channel data to buffer for writing
		binary.BigEndian.PutUint64(buffer, uint64(v))
		writer.Write(buffer)
	}
}

func RandomSource(count int) <-chan int {
	out := make(chan  int)
	go func() {
		for i := 0; i < count; i++ {
			out <- rand.Int()
		}
		close(out)
	}()
	return out
}

func MergeN(inputs ...<-chan int) <-chan int {
	if len(inputs) == 1 {
		return inputs[0]
	}
	m := len(inputs) / 2
	return Merge(MergeN(inputs[:m] ...),MergeN(inputs[m:] ...))
}