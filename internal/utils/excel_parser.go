package utils

import (
	"fmt"
	"strconv"
	"strings"
)

type ExcelSoalRow struct {
	RowIndex   int
	NoSoal     int
	Soal       string
	OpsiA      string
	OpsiB      string
	OpsiC      string
	OpsiD      string
	OpsiE      string
	Kunci      string
	GambarSoal string
	GambarA    string
	GambarB    string
	GambarC    string
	GambarD    string
	GambarE    string
}

// ParseExcelRow mengextract data dari row excel dan convert ke struct
func ParseExcelRow(values []string, rowIndex int) (*ExcelSoalRow, error) {
	// Helper function untuk safely get value dari slice
	getValueAt := func(idx int) string {
		if idx < len(values) {
			return toStringFromValue(values[idx])
		}
		return ""
	}

	// Parse no_soal (Column A, index 0)
	noSoalStr := getValueAt(0)
	var noSoal int
	if noSoalStr != "" {
		parsed, err := strconv.Atoi(noSoalStr)
		if err != nil {
			return nil, fmt.Errorf("no_soal harus berupa angka, got: %v", noSoalStr)
		}
		noSoal = parsed
	}

	return &ExcelSoalRow{
		RowIndex:   rowIndex,
		NoSoal:     noSoal,
		Soal:       getValueAt(1),
		OpsiA:      getValueAt(2),
		OpsiB:      getValueAt(3),
		OpsiC:      getValueAt(4),
		OpsiD:      getValueAt(5),
		OpsiE:      getValueAt(6),
		Kunci:      toUpperStringFromValue(getValueAt(7)),
		// Column I (index 8) skipped
		GambarSoal: getValueAt(9),
		// Column K (index 10) skipped
		GambarA: getValueAt(11),
		GambarB: getValueAt(12),
		GambarC: getValueAt(13),
		GambarD: getValueAt(14),
		GambarE: getValueAt(15),
	}, nil
}

func toStringFromValue(v string) string {
	return strings.TrimSpace(v)
}

func toUpperStringFromValue(v string) string {
	return strings.ToUpper(strings.TrimSpace(v))
}
