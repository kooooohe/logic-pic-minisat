package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)



func main() {

	tVars := sliceFromTxt("tvars") /* target vars */
	result := sliceFromTxt("out") /* minisat output after removing "SAT" string on top*/

	resultM := map[int]struct{}{}
	for _,v := range result {
		for _,vv := range v {
			resultM[vv] = struct{}{}
		}
	}

	for _,v := range tVars {
		for _,vv := range v {
			if _, ok := resultM[vv];ok {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}


	/*
	dboard = [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	*/

}

func sliceFromTxt(fName string) [][]int {
	var board [][]int
	file, err := os.Open(fName)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	bufSize := 10 * 1024 * 1024 // 10MB
	buf := make([]byte, bufSize)
	scanner.Buffer(buf, bufSize)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		numbers := strings.Split(line, " ")

		var row []int
		for _, numStr := range numbers {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				fmt.Printf("Error converting string to int: %s\n", err)
				os.Exit(1)
			}
			row = append(row, num)
		}
		board = append(board, row)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}

	// for _, row := range board {
	// 	fmt.Println(row)
	// }
	return board
}

