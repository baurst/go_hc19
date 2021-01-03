package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type photo struct {
	tags         []string
	isHorizontal bool
	idx          int
}

type slide struct {
	photos []photo
	tags   []string
}

var dataDir = flag.String("data_dir", "datasets", "The directory where the the input datasets are stored.")
var datasets = flag.String("datasets", "a_example.txt", "List of datasets \"a_example.txt b_lovely_landscapes.txt ...\" to be processed")

func getDatasetFiles(baseDir string, dataSets string) []string {
	var ds = strings.Split(strings.TrimSpace(dataSets), " ")
	var datasetFiles []string

	for _, el := range ds {
		if len(strings.TrimSpace(el)) > 0 {
			dsFile := filepath.Join(baseDir, el)
			datasetFiles = append(datasetFiles, dsFile)
		}
	}
	return datasetFiles

}

func readDatasetFile(datasetFile string) (err error) {
	fmt.Println("Reading: ", datasetFile)
	file, err := os.Open(datasetFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Start reading from the file with a reader.
	reader := bufio.NewReader(file)
	var line string
	lineIdx := 0
	numPhotos := -1
	for {
		line, err = reader.ReadString('\n')
		if err != nil && err != io.EOF {
			break
		}
		if lineIdx == 0 {
			numPhotos, _ = strconv.Atoi(strings.TrimSpace(line))
		} else {

			// Process the line here.
			fmt.Printf(" > Read %d characters\n", len(line))

			if err != nil {
				break
			}
		}
		lineIdx++
	}
	if err != io.EOF {
		fmt.Printf(" > Failed with error: %v\n", err)
		return err
	}
	fmt.Println("Num Photos: ", numPhotos)

	return
}

func main() {
	flag.Parse()
	datasets := getDatasetFiles(*dataDir, *datasets)
	fmt.Println("Processing: ", datasets)
	for _, ds := range datasets {
		readDatasetFile(ds)
	}
	fmt.Printf("Done")
	os.Exit(0)
}

// var wg = sync.WaitGroup{}

// func main() {
// 	start := time.Now()
// 	ch := make(chan int, 100)
// 	wg.Add(1)
// 	go sendInt(ch)
// 	for i := 0; i < 100; i++ {
// 		wg.Add(1)
// 		go receiveInt(ch)
// 	}
// 	wg.Wait()
// 	end := time.Now()

// 	fmt.Println(time.Duration(end.Sub(start)))

// 	fmt.Println(1000 * (10 + 100))

// }

// func receiveInt(ch <-chan int) {
// 	for i := range ch {
// 		fmt.Println("Received ", i)
// 		time.Sleep(100 * time.Millisecond)
// 	}
// 	wg.Done()
// }

// func sendInt(ch chan<- int) {
// 	for i := 0; i < 1000; i++ {
// 		fmt.Println("Sending ", i)
// 		ch <- i
// 		time.Sleep(10 * time.Millisecond)
// 	}
// 	close(ch)
// 	wg.Done()
// }
