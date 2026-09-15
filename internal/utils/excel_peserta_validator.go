package utils

// ValidatePesertaRow memvalidasi satu row peserta hasil parse excel.
// Pengecekan duplikat username (baik dalam file maupun terhadap database)
// dilakukan di service karena butuh state lintas-row/lintas-tabel.
func ValidatePesertaRow(row *ExcelPesertaRow) []string {
	var errors []string

	if row.Nama == "" {
		errors = append(errors, "nama tidak boleh kosong")
	}

	if row.Username == "" {
		errors = append(errors, "username tidak boleh kosong")
	}

	if row.Password == "" {
		errors = append(errors, "password tidak boleh kosong")
	} else if len(row.Password) < 6 {
		errors = append(errors, "password minimal 6 karakter")
	}

	return errors
}
