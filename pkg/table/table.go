package table

import (
	"fmt"
	"log"
	"strings"
)

const (
	TL string = "┌" // top-left
	TM string = "┬" // top-mid
	TR string = "┐" // top-right
	ML string = "├" // mid-left
	MM string = "┼" // mid-mid
	MR string = "┤" // mid-right
	BL string = "└" // bottom-left
	BM string = "┴" // bottom-mid
	BR string = "┘" // bottom-right
	VE string = "─" // vertical
	HO string = "│" // horizontal
)

func adjustLeft(s string, l int) string {
	sLen := len(s)
	p := ""
	p = fmt.Sprint(s, strings.Repeat(" ", l-sLen))
	return p
}

func adjustRight(s string, l int) string {
	sLen := len(s)
	p := ""
	p = fmt.Sprint(strings.Repeat(" ", l-sLen), s)
	return p
}

type Table struct {
	Header           []string
	Rows             [][]string
	Size             int
	MaxColumnSizes   []int
	Padding          int
	ColumnAlignments []int
}

func NewTable(s int, alignments []int) *Table {
	t := &Table{Header: []string{}, Rows: [][]string{}, Size: s}
	t.MaxColumnSizes = make([]int, t.Size)
	t.Padding = 1
	t.ColumnAlignments = alignments
	return t
}

func (t *Table) SetHeader(columns []string) {
	if len(columns) != t.Size {
		log.Fatal("Number of columns does not match with size ", t.Size)
	}
	t.Header = columns
	t.UpdateMaxColumnSizes(columns)
}

func (t *Table) AppendRow(row []string) {
	if len(row) != t.Size {
		log.Fatal("Number of columns does not match with size ", t.Size)
	}
	t.Rows = append(t.Rows, row)
	t.UpdateMaxColumnSizes(row)
}

func (t *Table) UpdateMaxColumnSizes(records []string) {
	for i := 0; i < t.Size; i++ {
		if len(records[i]) > t.MaxColumnSizes[i] {
			t.MaxColumnSizes[i] = len(records[i])
		}
	}
}

func (t *Table) BuildBorder(left string, middle string, right string) string {
	cols := make([]string, t.Size)
	for i := 0; i < t.Size; i++ {
		cols[i] = fmt.Sprint(strings.Repeat(VE, t.MaxColumnSizes[i]+t.Padding*2))
	}
	return fmt.Sprint(left, strings.Join(cols, middle), right)
}

func (t *Table) BuildRowLine(row []string, columnAlignments []int) string {
	cols := make([]string, t.Size)
	for i := 0; i < t.Size; i++ {
		alignment := columnAlignments[i]
		switch alignment {
		case 0:
			cols[i] = adjustLeft(row[i], t.MaxColumnSizes[i])
		case 1:
			cols[i] = adjustRight(row[i], t.MaxColumnSizes[i])
		}
		cols[i] = fmt.Sprint(strings.Repeat(" ", t.Padding), cols[i], strings.Repeat(" ", t.Padding))
	}
	return fmt.Sprint(HO, strings.Join(cols, HO), HO)
}

func (t *Table) Display() {
	// Top border
	topBorder := t.BuildBorder(TL, TM, TR)
	fmt.Println(topBorder)
	// Header
	header := t.BuildRowLine(t.Header, []int{0, 0})
	fmt.Println(header)
	// Middle border
	middleBorder := t.BuildBorder(ML, MM, MR)
	fmt.Println(middleBorder)
	// Rows
	for _, row := range t.Rows {
		fmt.Println(t.BuildRowLine(row, t.ColumnAlignments))
	}
	// Bottom border
	bottomBorder := t.BuildBorder(BL, BM, BR)
	fmt.Println(bottomBorder)
}
