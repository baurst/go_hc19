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
	"time"
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

func readDatasetFile(datasetFile string) []photo {

	var photos []photo
	fmt.Println("Reading: ", datasetFile)
	file, err := os.Open(datasetFile)
	if err != nil {
		return photos
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
			var pic photo
			chunks := strings.Split(line, " ")
			for i := range chunks {
				chunks[i] = strings.TrimSpace(chunks[i])
			}

			if len(chunks) > 1 {

				orientation := strings.TrimSpace(chunks[0])

				pic.isHorizontal = orientation == "H"
				// idx, _ := strconv.Atoi(strings.TrimSpace(chunks[1]))

				pic.idx = lineIdx - 1
				tags := chunks[2:]

				var finalTags []string

				for _, tag := range tags {
					finalTags = append(finalTags, strings.TrimSpace(tag))
				}

				pic.tags = finalTags

				photos = append(photos, pic)

				if err != nil {
					break
				}
			} else {
				break
			}
		}
		lineIdx++
	}
	if err != io.EOF {
		fmt.Printf(" > Failed with error: %v\n", err)
		return photos
	}
	fmt.Println("Num Photos: ", numPhotos)

	return photos
}

func scoreAllSlides(slides []slide) int {
	score := 0
	if len(slides) > 1 {
		prevSlideTags := slides[0].tags
		for i := 1; i < len(slides); i++ {
			curSlideTags := slides[i].tags
			score += evaluateTags(prevSlideTags, curSlideTags)
			prevSlideTags = curSlideTags
		}
	}
	return score
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func evaluateTags(prevSlide, curSlide []string) int {
	mapIntersect := map[string]int{} // present in both tag-arrays
	scoreIntersect := 0
	mapTagsPrev := map[string]int{} // present only in tag-array i
	scoreTagsPrev := 0
	mapTagsCur := map[string]int{} // present only in tag-array i+1
	scoreTagsCur := 0

	for _, prevVal := range prevSlide {
		mapIntersect[prevVal] = 1
		mapTagsPrev[prevVal] = 1
		mapTagsCur[prevVal] = 0
	}
	for _, curVal := range curSlide {
		mapIntersect[curVal] = mapIntersect[curVal] + 1
		mapTagsPrev[curVal] = 0
		if _, present := mapTagsCur[curVal]; !present {
			mapTagsCur[curVal] = 1
		}
	}

	for _, mapIntersectVal := range mapIntersect {
		if mapIntersectVal == 2 {
			scoreIntersect++
		}
	}

	for _, mapTagsPrevVal := range mapTagsPrev {
		if mapTagsPrevVal == 1 {
			scoreTagsPrev++
		}
	}

	for _, mapTagsCurVal := range mapTagsCur {
		if mapTagsCurVal == 1 {
			scoreTagsCur++
		}
	}

	return min(scoreIntersect, min(scoreTagsPrev, scoreTagsCur))
}

func writeSolution(result []slide, outDir string, score int, origDatasetPath string) {
	currentTime := time.Now().Format("20060102150405")
	var outFile = filepath.Join(outDir, currentTime+"_"+strings.TrimSuffix(filepath.Base(origDatasetPath), filepath.Ext(origDatasetPath))+"_"+fmt.Sprint(score)+".txt")
	fmt.Println("Writing result to: ", outFile)
	f, err := os.Create(outFile)

	if err == nil {
		fmt.Fprintln(f, len(result))
		for _, el := range result {
			var indices []string
			for _, img := range el.photos {
				indices = append(indices, fmt.Sprint(img.idx))
			}

			fmt.Fprintln(f, strings.Join(indices, " "))
		}
	}
	fmt.Println("Wrote result to: ", outFile)

}

func createDumbSolution(dataset []photo) []slide {
	solution := make([]slide, len(dataset))
	for i, p := range dataset {
		solution[i].photos = []photo{p}
		solution[i].tags = p.tags
	}
	return solution
}

func main() {
	flag.Parse()
	datasets := getDatasetFiles(*dataDir, *datasets)
	fmt.Println("Processing: ", datasets)
	for _, ds := range datasets {
		pics := readDatasetFile(ds)
		solution := createDumbSolution(pics)
		score := scoreAllSlides(solution)
		writeSolution(solution, "out", score, ds)

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
