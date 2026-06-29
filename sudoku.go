// sudoku.go
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	reset  = "\033[0m"
	red    = "\033[91m"
	green  = "\033[92m"
	yellow = "\033[93m"
	blue   = "\033[94m"
	cyan   = "\033[96m"
	gray   = "\033[90m"
	bold   = "\033[1m"
)

func colorize(text, color string) string {
	return color + text + reset
}

type Sudoku struct {
	board [9][9]int
	steps int
}

func NewSudoku() *Sudoku {
	return &Sudoku{}
}

func (s *Sudoku) String() string {
	return s.Render(nil)
}

func (s *Sudoku) Render(highlight [][2]int) string {
	var lines []string
	for i := 0; i < 9; i++ {
		if i%3 == 0 && i > 0 {
			lines = append(lines, colorize("┃━━━┃━━━┃━━━┃", gray))
		}
		var row []string
		for j := 0; j < 9; j++ {
			if j%3 == 0 && j > 0 {
				row = append(row, colorize("┃", gray))
			}
			val := s.board[i][j]
			if val == 0 {
				row = append(row, " ")
			} else {
				col := green
				if highlight != nil {
					for _, h := range highlight {
						if h[0] == i && h[1] == j {
							col = cyan
							break
						}
					}
				}
				row = append(row, colorize(strconv.Itoa(val), col))
			}
		}
		lines = append(lines, strings.Join(row, " "))
	}
	return strings.Join(lines, "\n")
}

func (s *Sudoku) ParseString(str string) error {
	if len(str) != 81 {
		return fmt.Errorf("строка должна содержать ровно 81 символ")
	}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			ch := str[i*9+j]
			if ch >= '1' && ch <= '9' {
				s.board[i][j] = int(ch - '0')
			} else if ch == '0' {
				s.board[i][j] = 0
			} else {
				return fmt.Errorf("недопустимый символ: %c", ch)
			}
		}
	}
	return nil
}

func (s *Sudoku) ExportString() string {
	var b strings.Builder
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if s.board[i][j] == 0 {
				b.WriteByte('0')
			} else {
				b.WriteByte(byte(s.board[i][j] + '0'))
			}
		}
	}
	return b.String()
}

func (s *Sudoku) IsValid() bool {
	// Проверка строк
	for i := 0; i < 9; i++ {
		seen := make(map[int]bool)
		for j := 0; j < 9; j++ {
			val := s.board[i][j]
			if val != 0 {
				if val < 1 || val > 9 || seen[val] {
					return false
				}
				seen[val] = true
			}
		}
	}
	// Проверка столбцов
	for j := 0; j < 9; j++ {
		seen := make(map[int]bool)
		for i := 0; i < 9; i++ {
			val := s.board[i][j]
			if val != 0 {
				if seen[val] {
					return false
				}
				seen[val] = true
			}
		}
	}
	// Проверка блоков
	for br := 0; br < 3; br++ {
		for bc := 0; bc < 3; bc++ {
			seen := make(map[int]bool)
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					val := s.board[br*3+i][bc*3+j]
					if val != 0 {
						if seen[val] {
							return false
						}
						seen[val] = true
					}
				}
			}
		}
	}
	return true
}

func (s *Sudoku) findBestEmpty() (int, int, []int) {
	bestR, bestC := -1, -1
	minCandidates := 10
	var bestCands []int
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if s.board[i][j] == 0 {
				cands := s.getCandidates(i, j)
				if len(cands) < minCandidates {
					minCandidates = len(cands)
					bestR, bestC = i, j
					bestCands = cands
					if minCandidates == 1 {
						return bestR, bestC, bestCands
					}
				}
			}
		}
	}
	return bestR, bestC, bestCands
}

func (s *Sudoku) getCandidates(row, col int) []int {
	if s.board[row][col] != 0 {
		return []int{}
	}
	used := make(map[int]bool)
	for j := 0; j < 9; j++ {
		used[s.board[row][j]] = true
	}
	for i := 0; i < 9; i++ {
		used[s.board[i][col]] = true
	}
	br, bc := row/3*3, col/3*3
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			used[s.board[br+i][bc+j]] = true
		}
	}
	var cands []int
	for v := 1; v <= 9; v++ {
		if !used[v] {
			cands = append(cands, v)
		}
	}
	return cands
}

func (s *Sudoku) Solve(animate bool, delay float64) (bool, float64) {
	s.steps = 0
	start := time.Now()
	var solve func() bool
	solve = func() bool {
		s.steps++
		r, c, cands := s.findBestEmpty()
		if r == -1 {
			return true
		}
		for _, val := range cands {
			s.board[r][c] = val
			if animate {
				s.printAnimated(r, c)
				time.Sleep(time.Duration(delay * float64(time.Second)))
			}
			if solve() {
				return true
			}
			s.board[r][c] = 0
		}
		return false
	}
	solved := solve()
	elapsed := time.Since(start).Seconds()
	return solved, elapsed
}

func (s *Sudoku) printAnimated(row, col int) {
	fmt.Print("\033[H\033[2J")
	fmt.Println(s.Render([][2]int{{row, col}}))
	fmt.Println(colorize("Шаг: "+strconv.Itoa(s.steps), yellow))
}

func (s *Sudoku) CountSolutions(maxCount int) int {
	count := 0
	var countFn func()
	countFn = func() {
		if count >= maxCount {
			return
		}
		r, c, cands := s.findBestEmpty()
		if r == -1 {
			count++
			return
		}
		for _, val := range cands {
			s.board[r][c] = val
			countFn()
			s.board[r][c] = 0
			if count >= maxCount {
				return
			}
		}
	}
	countFn()
	return count
}

func (s *Sudoku) Generate(level string) {
	// Сначала решаем пустую доску
	s.solveInternal()
	filled := s.board
	// Уровни
	var minEmpty, maxEmpty int
	switch level {
	case "easy":
		minEmpty, maxEmpty = 30, 35
	case "medium":
		minEmpty, maxEmpty = 40, 45
	case "hard":
		minEmpty, maxEmpty = 50, 55
	default:
		minEmpty, maxEmpty = 30, 35
	}
	target := minEmpty + rand.Intn(maxEmpty-minEmpty+1)
	removed := 0
	attempts := 0
	for removed < target && attempts < 10000 {
		attempts++
		i, j := rand.Intn(9), rand.Intn(9)
		if filled[i][j] != 0 {
			backup := filled[i][j]
			filled[i][j] = 0
			test := &Sudoku{board: filled}
			if test.CountSolutions(2) == 1 {
				removed++
			} else {
				filled[i][j] = backup
			}
		}
	}
	s.board = filled
}

func (s *Sudoku) solveInternal() {
	var solve func() bool
	solve = func() bool {
		r, c, cands := s.findBestEmpty()
		if r == -1 {
			return true
		}
		for _, val := range cands {
			s.board[r][c] = val
			if solve() {
				return true
			}
			s.board[r][c] = 0
		}
		return false
	}
	solve()
}

func loadBoard(source string) (*Sudoku, error) {
	var content string
	if _, err := os.Stat(source); err == nil {
		data, err := os.ReadFile(source)
		if err != nil {
			return nil, err
		}
		content = string(data)
	} else {
		content = source
	}
	// Фильтруем только цифры
	var digits []rune
	for _, ch := range content {
		if ch >= '0' && ch <= '9' {
			digits = append(digits, ch)
		}
	}
	if len(digits) != 81 {
		return nil, fmt.Errorf("должно быть ровно 81 цифра")
	}
	s := NewSudoku()
	if err := s.ParseString(string(digits)); err != nil {
		return nil, err
	}
	return s, nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println(colorize("Usage: sudoku <generate|solve|check|count|export> [options]", yellow))
		fmt.Println("  generate -l <easy|medium|hard>")
		fmt.Println("  solve -i <file|string> [-a] [--delay <sec>]")
		fmt.Println("  check -i <file|string>")
		fmt.Println("  count -i <file|string> [--max <N>]")
		fmt.Println("  export -i <file|string>")
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "generate":
		level := "easy"
		for i := 2; i < len(os.Args); i++ {
			if os.Args[i] == "-l" && i+1 < len(os.Args) {
				level = os.Args[i+1]
				i++
			}
		}
		s := NewSudoku()
		s.Generate(level)
		fmt.Println(s.Render(nil))
		empty := 0
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				if s.board[i][j] == 0 {
					empty++
				}
			}
		}
		fmt.Println(colorize("Уровень: "+level+", пустых ячеек: "+strconv.Itoa(empty), yellow))

	case "solve":
		var input string
		animate := false
		delay := 0.1
		for i := 2; i < len(os.Args); i++ {
			if os.Args[i] == "-i" && i+1 < len(os.Args) {
				input = os.Args[i+1]
				i++
			} else if os.Args[i] == "-a" {
				animate = true
			} else if os.Args[i] == "--delay" && i+1 < len(os.Args) {
				delay, _ = strconv.ParseFloat(os.Args[i+1], 64)
				i++
			}
		}
		if input == "" {
			fmt.Println(colorize("Укажите входную доску через -i", red))
			os.Exit(1)
		}
		s, err := loadBoard(input)
		if err != nil {
			fmt.Println(colorize("Ошибка: "+err.Error(), red))
			os.Exit(1)
		}
		if !s.IsValid() {
			fmt.Println(colorize("Доска невалидна!", red))
			os.Exit(1)
		}
		solved, elapsed := s.Solve(animate, delay)
		if solved {
			fmt.Println(s.Render(nil))
			fmt.Println(colorize("Решение найдено за "+strconv.FormatFloat(elapsed, 'f', 3, 64)+" сек, шагов: "+strconv.Itoa(s.steps), green))
		} else {
			fmt.Println(colorize("Решения не существует!", red))
		}

	case "check":
		var input string
		for i := 2; i < len(os.Args); i++ {
			if os.Args[i] == "-i" && i+1 < len(os.Args) {
				input = os.Args[i+1]
				i++
			}
		}
		if input == "" {
			fmt.Println(colorize("Укажите входную доску через -i", red))
			os.Exit(1)
		}
		s, err := loadBoard(input)
		if err != nil {
			fmt.Println(colorize("Ошибка: "+err.Error(), red))
			os.Exit(1)
		}
		if s.IsValid() {
			fmt.Println(colorize("✅ Доска корректна.", green))
		} else {
			fmt.Println(colorize("❌ Доска невалидна.", red))
		}

	case "count":
		var input string
		maxCount := 100
		for i := 2; i < len(os.Args); i++ {
			if os.Args[i] == "-i" && i+1 < len(os.Args) {
				input = os.Args[i+1]
				i++
			} else if os.Args[i] == "--max" && i+1 < len(os.Args) {
				maxCount, _ = strconv.Atoi(os.Args[i+1])
				i++
			}
		}
		if input == "" {
			fmt.Println(colorize("Укажите входную доску через -i", red))
			os.Exit(1)
		}
		s, err := loadBoard(input)
		if err != nil {
			fmt.Println(colorize("Ошибка: "+err.Error(), red))
			os.Exit(1)
		}
		count := s.CountSolutions(maxCount)
		if count >= maxCount {
			fmt.Println(colorize("Количество решений >= "+strconv.Itoa(maxCount)+" (ограничено)", yellow))
		} else {
			fmt.Println(colorize("Количество решений: "+strconv.Itoa(count), green))
		}

	case "export":
		var input string
		for i := 2; i < len(os.Args); i++ {
			if os.Args[i] == "-i" && i+1 < len(os.Args) {
				input = os.Args[i+1]
				i++
			}
		}
		if input == "" {
			fmt.Println(colorize("Укажите входную доску через -i", red))
			os.Exit(1)
		}
		s, err := loadBoard(input)
		if err != nil {
			fmt.Println(colorize("Ошибка: "+err.Error(), red))
			os.Exit(1)
		}
		fmt.Println(s.ExportString())

	default:
		fmt.Println(colorize("Неизвестная команда: "+cmd, red))
	}
}
