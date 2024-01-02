package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Clause []int

var clauses []Clause

func intSeq() func() int {
	c := 0
	return func() int {
		c++
		return c
	}
}

var seq = intSeq()


func spiteColumnRow(ts [][]int) ([][]int, [][]int) {
	/* file format ex
	2 3 
	1 0
	1 0
	1 0
	0 0
	0 1
	*/

	cn := ts[0][1]
	rn := ts[0][0]

	cs := make ([][]int, cn)
	ii := 0
	for i:= 1; i < 1+cn; i++ {
		cs[ii] = ts[i]
		ii++
	}

	ii = 0
	rs := make([][]int, rn)
	for i:= cn+1; i < cn+rn+1; i++ {
		rs[ii] = ts[i]
		ii++
	}
	return cs, rs

}
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: puzzle <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]

	dboard := board(filename)

	cRules,rRules := spiteColumnRow(dboard)
	columnNum := len(cRules)
	rowNum := len(rRules)
	// fmt.Println(cRules)
	// fmt.Println()
	// fmt.Println(rRules)


	bVars := make([][]int, rowNum)

	for i := range bVars {
		bVars[i] = make([]int, columnNum)
		for j := range bVars[i] {
			bVars[i][j] = seq()
		}
	}


	// fmt.Println()
	// fmt.Println(bVars)


	/*
		dboard = [][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		}
	*/


	//  fmt.Println(dboard)
	// ==== for later check//
	file, err := os.Create("tvars")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, row := range bVars {
		line := ""
		for _, v := range row {
			line += strconv.Itoa(v) + " "
		}
		_, err := file.WriteString(line + "\n")
		if err != nil {
			panic(err)
		}
	}
	// ==== //


	for i,v:= range rRules {
		// rules := []int{3, 1, 2}
		// width := 10
		pattern := make([]bool, columnNum/*width*/) 
		ys := generatePatterns(pattern, 0, v, 0, bVars[i], []int{})
		clauses = append(clauses, ys)
		// fmt.Println(ys)
	}
	for i,v:= range cRules {
		// rules := []int{3, 1, 2}
		// width := 10
		pattern := make([]bool, rowNum/*height*/) 
		tVars := make([]int, rowNum)
		for ii := 0; ii<rowNum; ii++ {
			tVars[ii] = bVars[ii][i]
		}
		ys := generatePatterns(pattern, 0, v, 0, tVars, []int{})
		fmt.Println(ys)
		// clauses = append(clauses, ys)
		
	}



	/*
	for i, v := range expandedDBoard {
		for j, vv := range v {
			if vv == -1 {
				continue
			}

			// vars around target cell
			tVars := []int{}
			for ii := -1; ii <= 1; ii++ {
				for jj := -1; jj <= 1; jj++ {
					tVars = append(tVars, expandedDBoardVars[i+ii][j+jj])
				}
			}
			// Determine
			if vv == 0 {
				for _, v := range tVars {
					clauses = append(clauses, Clause{-v})
				}
				continue
			}
			// Determine
			if vv == 9 {
				for _, v := range tVars {
					clauses = append(clauses, Clause{v})
				}
				continue
			}

			if vv == 8 {
				c := []int{}
				for _, v := range tVars {
					c = append(c, -v)
				}
				clauses = append(clauses, c)
			}

			// true isNot Positive
			for k := vv; k < 8; k++ {
				comb(k+1, 9-(k+1), false, tVars)
			}

			if vv == 1 {
				c := []int{}
				for _, v := range tVars {
					c = append(c, v)
				}
				clauses = append(clauses, c)
			}

			// false is Postive
			for k := 9 - vv; k < 8; k++ {
				comb(k+1, 9-(k+1), true, tVars)
			}
		}
	}
	*/

	cnf := clausesToString(clauses, seq()-1)

	// Save the CNF to a text file
	fOut := "r_cnf.txt"
	if err := os.WriteFile(fOut, []byte(cnf), 0644); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	fmt.Println("CNF file generated successfully:", fOut)
}

func board(fName string) [][]int {
	var board [][]int
	file, err := os.Open(fName)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {continue}
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



func clausesToString(clauses []Clause, varCount int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("p cnf %d %d\n", varCount, len(clauses)))
	for _, clause := range clauses {
		for _, lit := range clause {
			sb.WriteString(fmt.Sprintf("%d ", lit))
		}
		sb.WriteString("0\n")
	}
	return sb.String()
}



/*
func main() {
	rules := []int{3, 1, 2}
	width := 10
	pattern := make([]bool, width) // trueは塗りつぶされたマスを表す
	generatePatterns(pattern, 0, rules, 0)
}
*/

// generatePatterns は、与えられたルールに基づいてパターンを生成する再帰関数です。
func generatePatterns(pattern []bool, position int, rules []int, ruleIndex int, tVars, ys []int) ([]int){
	if ruleIndex == len(rules) {
		// すべてのルールを配置した後、パターンを出力
		ys = printPattern(pattern, tVars,ys)
		return ys
	}

	// 現在のルールを配置するために必要なスペース
	spaceNeeded := rules[ruleIndex] + 1 // 1は次のブロックのための空間
	if ruleIndex == len(rules)-1 {
		spaceNeeded-- // 最後のルールの場合は追加の空間は不要
	}

	for i := position; i <= len(pattern)-spaceNeeded; i++ {
		// 現在のルールを配置
		newPattern := make([]bool, len(pattern))
		copy(newPattern, pattern)
		for j := 0; j < rules[ruleIndex]; j++ {
			newPattern[i+j] = true
		}

		// 次のルールの配置のために再帰呼び出し
		nextPosition := i + rules[ruleIndex] + 1
		ys = generatePatterns(newPattern, nextPosition, rules, ruleIndex+1, tVars, ys)
	}

	return ys
}

// printPattern はパターンを出力するヘルパー関数です。
func printPattern(pattern []bool, tVars, ys []int) ([]int){
	var sb strings.Builder
	c := Clause{}
	t := seq()
	for i, filled := range pattern {
		if filled {
			//c = append(c, tVars[i])
			clauses = append(clauses, Clause{-t,tVars[i]})
			sb.WriteString("■ ")
		} else {
			c = append(c, -tVars[i])
	//		clauses = append(clauses, Clause{t})
			clauses = append(clauses, Clause{-t,-tVars[i]})
			sb.WriteString("□ ")
		}
	}
	ys = append(ys, t)
	// fmt.Println(sb.String())
	return ys

}
