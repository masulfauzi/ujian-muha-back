package utils

import (
	"strings"
)

// ValidateSoalRow memvalidasi satu row dari excel
func ValidateSoalRow(row *ExcelSoalRow) []string {
	var errors []string

	kunci := strings.TrimSpace(row.Kunci)
	if kunci == "" {
		errors = append(errors, "kunci tidak boleh kosong")
	} else if !IsValidKunci(kunci) {
		errors = append(errors, "kunci harus berupa A, B, C, D, atau E")
	}

	return errors
}

// IsValidKunci checks if kunci is valid (A-E)
func IsValidKunci(k string) bool {
	k = strings.ToUpper(strings.TrimSpace(k))
	return k == "A" || k == "B" || k == "C" || k == "D" || k == "E"
}
