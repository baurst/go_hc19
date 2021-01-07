package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
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

// var datasets = flag.String("datasets", "d_pet_pictures.txt", "List of datasets \"a_example.txt b_lovely_landscapes.txt ...\" to be processed")

var datasets = flag.String("datasets", "a_example.txt b_lovely_landscapes.txt c_memorable_moments.txt d_pet_pictures.txt e_shiny_selfies.txt", "List of datasets \"a_example.txt b_lovely_landscapes.txt ...\" to be processed")

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

	if numPhotos != len(photos) {
		panic("Incorrect number of photos read!")
	}
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

func max(a, b int) int {
	if b > a {
		return b
	}
	return a
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
	f, err := os.Create(outFile)

	imageIndices := map[int]bool{}
	if err == nil {
		fmt.Fprintln(f, len(result))
		for _, el := range result {
			var indices []int
			for _, img := range el.photos {
				indices = append(indices, img.idx)
				if _, ok := imageIndices[img.idx]; ok {
					panic("Solution contains duplicates!")
				}
				imageIndices[img.idx] = true
			}

			sep := ","
			n := len(sep) * (len(indices) - 1)
			for i := 0; i < len(indices); i++ {
				n += len(fmt.Sprint(indices[i]))
			}

			var b strings.Builder
			b.Grow(n)
			b.WriteString(fmt.Sprint(indices[0]))
			for _, s := range indices[1:] {
				b.WriteString(sep)
				b.WriteString(fmt.Sprint(s))
			}

			fmt.Fprintln(f, b.String())
		}
	}
	fmt.Println("Wrote result to: ", outFile)

}

func moveRandomSlideBy1(slides []slide, iterations int) []slide {
	amountSlides := len(slides)
	searchScope := 2
	for i := 0; i < iterations; i++ {
		randomSlideIdx := rand.Intn(amountSlides)
		relevantSlides := make([]slide, 2*searchScope+1)
		copy(relevantSlides, slides[max(0, randomSlideIdx-searchScope):min(len(slides), randomSlideIdx+searchScope+1)])
		initialScore := scoreAllSlides(relevantSlides)
		//fmt.Printf("slides before %d iter around %d with score %d: %s\n", i, randomSlideIdx, initialScore, converSlidesToPhotoIdx(relevantSlides))
		swapRightScore := scoreAllSlides(swapSlide(relevantSlides, min(searchScope, randomSlideIdx), -1))
		swapLeftScore := scoreAllSlides(swapSlide(swapSlide(relevantSlides, min(searchScope, randomSlideIdx), -1), min(searchScope, randomSlideIdx), 1))
		if swapRightScore > initialScore {
			//fmt.Println("SwappedRight")
			swapSlide(slides, randomSlideIdx, -1)
		} else if swapLeftScore > initialScore {
			//fmt.Println("SwappedLeft")
			swapSlide(slides, randomSlideIdx, 1)
		}
		//newScore := scoreAllSlides(slides[max(0, randomSlideIdx-searchScope):min(len(slides), randomSlideIdx+searchScope+1)])
		//fmt.Printf("slides after %d iter around %d with score %d: %s\n", i, randomSlideIdx, newScore, converSlidesToPhotoIdx(slides[max(0, randomSlideIdx-searchScope):min(len(slides), randomSlideIdx+searchScope+1)]))
	}
	return slides
}

func converSlidesToPhotoIdx(slides []slide) string {
	photoString := ""
	for _, valS := range slides {
		idxString := ""
		for _, valP := range valS.photos {
			idxString += fmt.Sprint(valP.idx) + "+"
		}
		photoString += idxString + " "
	}
	return photoString
}

func swapSlide(slides []slide, position int, movement int) []slide {
	if position+movement > 0 && position+movement < len(slides) {
		slideAtPosition := slides[position]
		slides[position] = slides[position+movement]
		slides[position+movement] = slideAtPosition
	}
	return slides
}

func createDumbSolution(dataset []photo) []slide {
	solution := make([]slide, len(dataset))
	for i, p := range dataset {
		solution[i].photos = []photo{p}
		solution[i].tags = p.tags
	}
	return solution
}

func stringUnion(a, b []string) []string {
	union := map[string]bool{}
	for _, val := range a {
		union[val] = true
	}
	for _, val := range b {
		union[val] = true
	}
	keys := make([]string, 0, len(union))
	for k := range union {
		keys = append(keys, k)
	}
	return keys
}

func createInitialSlideshowByNumTags(dataset []photo) []slide {
	var horizontalPhotos []photo
	var verticalPhotos []photo

	for _, pic := range dataset {
		if pic.isHorizontal {
			horizontalPhotos = append(horizontalPhotos, pic)
		} else {
			verticalPhotos = append(verticalPhotos, pic)
		}
	}

	// Vertical Slides
	sort.Slice(verticalPhotos, func(i, j int) bool {
		return len(verticalPhotos[i].tags) < len(verticalPhotos[i].tags)
	})
	verticalSlides := arrangeVerticalPhotosByTotalTagsOfN(verticalPhotos, 20)

	// Horizontal Slides
	var horizontalSlides []slide
	for _, pic := range horizontalPhotos {
		tmpSlide := slide{[]photo{pic}, pic.tags}
		horizontalSlides = append(horizontalSlides, tmpSlide)
	}
	var slides []slide
	slides = append(slides, horizontalSlides...)
	slides = append(slides, verticalSlides...)

	return slides

}

func narrowIndicesTo(vs []bool, num int) []int {
	vsf := make([]int, 0, num)
	storedValues := 0
	for i := len(vs) - 1; i >= 0; i-- {
		if !vs[i] {
			vsf = append(vsf, i)
			storedValues++
		}
		if storedValues >= num {
			return vsf
		}
	}
	return vsf
}

func arrangeVerticalPhotosByTotalTagsOfN(verticalPhotos []photo, num int) []slide {
	var verticalSlides []slide
	verticalPhotosUsed := make([]bool, len(verticalPhotos))
	for i := 0; i < len(verticalPhotos); i++ {
		if verticalPhotosUsed[i] == false {
			baseImg := verticalPhotos[i]
			verticalPhotosUsed[i] = true
			maxTagCount := len(baseImg.tags)
			maxTagCountIdx := -1

			for _, j := range narrowIndicesTo(verticalPhotosUsed, num) {
				if verticalPhotosUsed[j] == false {
					nextImg := verticalPhotos[j]
					totalTagsCount := len(stringUnion(baseImg.tags, nextImg.tags))
					if totalTagsCount > maxTagCount {
						maxTagCount = totalTagsCount
						maxTagCountIdx = j
					}
				}
			}
			if maxTagCountIdx >= 0 {
				usedImg := verticalPhotos[maxTagCountIdx]
				verticalPhotosUsed[maxTagCountIdx] = true
				totalTags := stringUnion(baseImg.tags, usedImg.tags)
				verticalSlides = append(verticalSlides, slide{[]photo{baseImg, usedImg}, totalTags})
			} else {
				verticalSlides = append(verticalSlides, slide{[]photo{baseImg}, baseImg.tags})
			}
		}
	}
	return verticalSlides
}

func arrangeVerticalPhotosByTagLength(verticalPhotos []photo) []slide {
	numVertSlides := int(len(verticalPhotos) / 2.0)
	isEvenNumSlides := len(verticalPhotos)%2 == 0

	var verticalSlides []slide
	for i := 0; i < numVertSlides; i++ {
		firstImg := verticalPhotos[i]
		lastImg := verticalPhotos[len(verticalPhotos)-1-i]
		totalTags := stringUnion(firstImg.tags, lastImg.tags)
		verticalSlides = append(verticalSlides, slide{[]photo{firstImg, lastImg}, totalTags})
	}

	if !isEvenNumSlides {
		verticalSlides = append(verticalSlides, slide{[]photo{verticalPhotos[numVertSlides+1]}, verticalPhotos[numVertSlides+1].tags})
	}

	return verticalSlides
}

func sum(array []int) int {
	result := 0
	for _, v := range array {
		result += v
	}
	return result
}

func shuffleSolution(sol []slide) []slide {
	dest := make([]slide, len(sol))
	perm := rand.Perm(len(sol))
	for i, v := range perm {
		dest[v] = sol[i]
	}
	return dest
}

func getMaxAvailableIdx(a []int, slidesTaken []bool) int {
	maxIdx := -1
	max := -1
	for i, e := range a {
		if e > max && !slidesTaken[i] {
			max = e
			maxIdx = i
		}
	}
	return maxIdx
}

func optimizeRandomSubsets(solution []slide, maxIter int) []slide {
	const subsetSize int = 50
	if len(solution) > subsetSize {
		maxIdxUpper := len(solution) - subsetSize
		for i := 0; i < maxIter; i++ {
			startIdx := rand.Intn(maxIdxUpper)
			// println("startIdx:", startIdx)
			sliceToOptimize := solution[startIdx : startIdx+subsetSize]
			// println("Photo indices before", converSlidesToPhotoIdx(sliceToOptimize))
			origSlice := append([]slide(nil), sliceToOptimize...)
			prevScore := scoreAllSlides(origSlice)
			// println("prevScore:", prevScore)

			var transitionScoreMatrix [subsetSize][subsetSize]int

			for firstSlideIdx := 0; firstSlideIdx < subsetSize; firstSlideIdx++ {
				for secondSlideIdx := firstSlideIdx; secondSlideIdx < subsetSize; secondSlideIdx++ {
					if firstSlideIdx == secondSlideIdx {
						transitionScoreMatrix[firstSlideIdx][secondSlideIdx] = -1000
					} else {
						localScore := evaluateTags(sliceToOptimize[firstSlideIdx].tags, sliceToOptimize[secondSlideIdx].tags)
						transitionScoreMatrix[firstSlideIdx][secondSlideIdx] = localScore
						transitionScoreMatrix[secondSlideIdx][firstSlideIdx] = localScore
					}
				}
			}

			slidesTaken := make([]bool, subsetSize)
			slidesTaken[len(slidesTaken)-1] = true
			greedyImprovedSlidesIdx := make([]int, 0, subsetSize)
			startSlice := -1
			bestFitSecondSlideIdx := -1
			for firstSlideIdx := 0; firstSlideIdx < subsetSize-1; firstSlideIdx++ {
				if startSlice < 0 {
					startSlice = firstSlideIdx
					greedyImprovedSlidesIdx = append(greedyImprovedSlidesIdx, startSlice)
					slidesTaken[startSlice] = true
				} else {
					startSlice = bestFitSecondSlideIdx
				}

				bestFitSecondSlideIdx = getMaxAvailableIdx(transitionScoreMatrix[startSlice][:], slidesTaken)

				if bestFitSecondSlideIdx < 0 {
					break
				}

				greedyImprovedSlidesIdx = append(greedyImprovedSlidesIdx, bestFitSecondSlideIdx)
				slidesTaken[bestFitSecondSlideIdx] = true
			}
			greedyImprovedSlidesIdx = append(greedyImprovedSlidesIdx, subsetSize-1)
			var improvedSlideSlice []slide
			for _, greedyIndex := range greedyImprovedSlidesIdx {
				improvedSlideSlice = append(improvedSlideSlice, sliceToOptimize[greedyIndex])
			}

			newScore := scoreAllSlides(improvedSlideSlice)
			// fmt.Printf("greedyImprovedSlidesIdx: %v", greedyImprovedSlidesIdx)
			// println("newScore:", newScore)

			if newScore > prevScore {
				for replacementIdx := 0; replacementIdx < subsetSize; replacementIdx++ {
					solution[replacementIdx+startIdx] = improvedSlideSlice[replacementIdx]
				}
			}
			// println("Photo indices after", converSlidesToPhotoIdx(solution[startIdx:startIdx+subsetSize]))

		}

	}

	return solution

}

func main() {
	flag.Parse()
	datasets := getDatasetFiles(*dataDir, *datasets)
	fmt.Println("Processing: ", datasets)
	scores := make([]int, len(datasets))
	// sem := make(chan empty, len(datasets))
	wg := &sync.WaitGroup{}
	wg.Add(len(datasets))
	for idx, ds := range datasets {
		go func(idx int, ds string) {
			defer wg.Done()
			pics := readDatasetFile(ds)
			solution := createInitialSlideshowByNumTags(pics)
			initScore := scoreAllSlides(solution)
			// solution = shuffleSolution(solution)
			solution = optimizeRandomSubsets(solution, 10000)
			afterOptimizationScore := scoreAllSlides(solution)
			solution = moveRandomSlideBy1(solution, 10000)
			score := scoreAllSlides(solution)
			fmt.Println(ds, ": Initial score: ", initScore, " - After optimization: ", afterOptimizationScore, " - Final: ", score)
			writeSolution(solution, "out", score, ds)
			scores[idx] = score
		}(idx, ds)
	}
	wg.Wait()

	fmt.Println("Done. Total score: ", sum(scores))
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
