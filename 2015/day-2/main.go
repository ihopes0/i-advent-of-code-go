package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
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
	inputOne, err := getInput()
	if err != nil {
		log.Fatalf("Cannot parse input.md. Err: %s", err)
	}

	// fmt.Println("Starting solving problem part 1 with input: ", inputOne)

	ansOne, err := resolveProblemPartOne(inputOne)
	if err != nil {
		log.Fatalf("Err during part one: %s", err)
	}

	fmt.Printf("Answer1: %d\n", ansOne)
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
		fmt.Println(scanner.Text())
		var r [3]int64
		input = append(input, r)
		for i, num := range strings.Split(scanner.Text(), "x") {
			fmt.Println(i, " - ", num)
			input[row][i], err = strconv.ParseInt(num, 0, 64)
			if err != nil {
				return nil, err
			}
		}
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

func getTotalSpace(l, w, h int64) (int64, error) {
	ar1, ar2, ar3 := l*w, w*h, l*h
	minAr := math.Min(math.Min(float64(ar1), float64(ar2)), float64(ar3))

	return applyFormula(ar1, ar2, ar3, int64(minAr))
}

func applyFormula(ar1, ar2, ar3, minAr int64) (int64, error) {
	return 2*(ar1+ar2+ar3) + minAr, nil
}
