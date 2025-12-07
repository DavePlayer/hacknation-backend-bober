package imports

import "encoding/json"

// Table reprezentuje jedną tabelę z pliku (arkusz Excela / wykryta tabela z PDF)
type Table struct {
	FileName string     `json:"file_name"`
	Sheet    string     `json:"sheet"`
	Rows     [][]string `json:"rows"`
}

// FileTables to zbiór tabel z jednego pliku
type FileTables struct {
	FileName string  `json:"file_name"`
	Tables   []Table `json:"tables"`
}

// RowLLMResult – odpowiedź z LLM dla pojedynczego wiersza
type RowLLMResult struct {
	FileName string          `json:"file_name"`
	Sheet    string          `json:"sheet"`
	RowIndex int             `json:"row_index"` // index w tabeli (od 1, bo 0 to nagłówki)
	Output   json.RawMessage `json:"output"`    // JSON zwrócony przez LLM
}

// ImportFilesAIResponse – odpowiedź handlera
type ImportFilesAIResponse struct {
	Status     string         `json:"status"`
	Tables     []FileTables   `json:"tables"`
	PDFFiles   []string       `json:"pdf_files"`
	RowResults []RowLLMResult `json:"row_results"`
	// Jak chcesz, możesz dorzucić jeszcze np. podsumowanie tekstowe:
	LLMResult string `json:"llmResult,omitempty"`
}

// buildRowPrompt – buduje prompt dla pojedynczego wiersza
func buildRowPrompt(schema string, headers []string, row []string) string {
	// z wiersza robimy małego JSON-a: { "Kolumna": "Wartość", ... }
	rowObj := make(map[string]string)
	for i, h := range headers {
		if i < len(row) {
			rowObj[h] = row[i]
		}
	}

	rowJSON, _ := json.Marshal(rowObj)
	var dataStructure = ""

	var prompt = `
	ZADANIE:
		Masz dane wejściowe w formacie JSON oraz docelowy schemat JSON.
		Twoim zadaniem jest PRZETWORZYĆ dane wejściowe do formatu zgodnego z tym schematem.

		WYMAGANIA:
		- Odpowiedź MUSI być **wyłącznie** poprawnym JSON-em.
		- NIE dodawaj żadnego tekstu, komentarzy, markdown, kodu, wyjaśnień.
		- Używaj dokładnie pól z podanego schematu.
		- Jeśli jakiejś wartości nie da się wyznaczyć, użyj null albo pustego stringa.
		- Nie streszczaj danych, tylko je przemapuj.

		schemat wyjściowy: ` + string(dataStructure) + `
		dane: ` + string(rowJSON) + `
		
	`

	return prompt
}
