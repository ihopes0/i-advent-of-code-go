package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// https://adventofcode.com/2015/day/2
// --- Day 2: I Was Told There Would Be No Math ---

// -- Part 1 --

// The elves are running low on wrapping paper, and so they need to submit an order for more. They have a list of the dimensions (length l, width w, and height h) of each present, and only want to order exactly as much as they need.

// Fortunately, every present is a box (a perfect right rectangular prism), which makes calculating the required wrapping paper for each gift a little easier: find the surface area of the box, which is 2*l*w + 2*w*h + 2*h*l. The elves also need a little extra paper for each present: the area of the smallest side.

// For example:

//     A present with dimensions 2x3x4 requires 2*6 + 2*12 + 2*8 = 52 square feet of wrapping paper plus 6 square feet of slack, for a total of 58 square feet.
//     A present with dimensions 1x1x10 requires 2*1 + 2*10 + 2*10 = 42 square feet of wrapping paper plus 1 square foot of slack, for a total of 43 square feet.

// All numbers in the elves' list are in feet. How many total square feet of wrapping paper should they order?

func main() {
	start := time.Now()
	inputOne, err := getInput()
	if err != nil {
		log.Fatalf("Cannot parse input.md. Err: %s", err)
	}
	ansOne, err := resolveProblemPartOne(inputOne)
	if err != nil {
		log.Fatalf("Err during part one: %s", err)
	}
	fmt.Printf("Answer1      : %d | took: %v\n", ansOne, time.Since(start))

	start = time.Now()
	ansOneOnTheFly, err := resolveProblemPartOneOnTheFly()
	if err != nil {
		log.Fatalf("Err during part one otf: %s", err)
	}
	fmt.Printf("Answer1 [OTF]: %d | took: %v\n", ansOneOnTheFly, time.Since(start))

	start = time.Now()
	ansConc, err := resolveProblemPartOneConcurrent()
	if err != nil {
		log.Fatalf("Err during part one concurrent: %s", err)
	}
	fmt.Printf("Answer1 [CON]: %d | took: %v\n", ansConc, time.Since(start))
}

func getInput() (input [][3]int64, err error) {
	file, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	row := 0
	for scanner.Scan() {
		dims, err := parseDimensions(scanner.Text())
		if err != nil {
			return input, err
		}
		input = append(input, dims)
		row += 1
	}

	return input, nil
}

func resolveProblemPartOne(input [][3]int64) (ans int64, err error) {
	for _, row := range input {
		res, err := getTotalSpace(row[0], row[1], row[2])
		if err != nil {
			return 0, err
		}
		ans += res
	}

	return ans, nil
}

func resolveProblemPartOneOnTheFly() (int64, error) {
	file, err := os.Open("input.txt")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var ans int64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dims, err := parseDimensions(scanner.Text())
		if err != nil {
			return 0, err
		}

		res, err := getTotalSpace(dims[0], dims[1], dims[2])
		if err != nil {
			return 0, err
		}

		ans += res
	}

	return ans, nil
}

func resolveProblemPartOneConcurrent() (int64, error) {
	file, err := os.Open("input.txt")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	const workers = 4

	jobs := make(chan string)
	results := make(chan int64)
	errChan := make(chan error, 1)

	var wg sync.WaitGroup

	// Workers
	for range workers {
		wg.Go(func() {

			for row := range jobs {
				dims, err := parseDimensions(row)
				if err != nil {
					errChan <- err
					return
				}

				res, err := getTotalSpace(dims[0], dims[1], dims[2])
				if err != nil {
					errChan <- err
					return
				}

				results <- res
			}
		})
	}

	// Reader
	go func() {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			jobs <- scanner.Text()
		}
		close(jobs)
	}()

	// Close results ONLY after workers finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect
	var ans int64
	for res := range results {
		ans += res
	}

	select {
	case err := <-errChan:
		return 0, err
	default:
	}

	return ans, nil
}

func parseDimensions(row string) (dims [3]int64, err error) {
	for i, num := range strings.Split(row, "x") {
		// fmt.Println(i, " - ", num)
		dims[i], err = strconv.ParseInt(num, 0, 64)
		if err != nil {
			return dims, err
		}
	}
	return dims, nil
}

func getTotalSpace(l, w, h int64) (int64, error) {
	ar1, ar2, ar3 := l*w, w*h, l*h
	minAr := math.Min(math.Min(float64(ar1), float64(ar2)), float64(ar3))

	return applyFormula(ar1, ar2, ar3, int64(minAr))
}

func applyFormula(ar1, ar2, ar3, minAr int64) (int64, error) {
	return 2*(ar1+ar2+ar3) + minAr, nil
}
