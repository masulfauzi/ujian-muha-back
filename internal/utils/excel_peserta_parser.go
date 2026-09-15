package utils

type ExcelPesertaRow struct {
	RowIndex  int
	Nama      string
	Username  string
	Password  string
	NamaKelas string
}

// ParseExcelPesertaRow mengextract data peserta dari satu row excel.
// Urutan kolom: A=nama, B=username, C=password, D=kelas.
func ParseExcelPesertaRow(values []string, rowIndex int) *ExcelPesertaRow {
	getValueAt := func(idx int) string {
		if idx < len(values) {
			return toStringFromValue(values[idx])
		}
		return ""
	}

	return &ExcelPesertaRow{
		RowIndex:  rowIndex,
		Nama:      getValueAt(0),
		Username:  getValueAt(1),
		Password:  getValueAt(2),
		NamaKelas: getValueAt(3),
	}
}
