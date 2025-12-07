package imports

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"bober.app/internal/db"
	"bober.app/internal/jsonRespond"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/xuri/excelize/v2"

	"github.com/ledongthuc/pdf"
	openai "github.com/sashabaranov/go-openai"
)

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

// ImportFilesAIResponse – przykładowa odpowiedź handlera
type ImportFilesAIResponse struct {
	Status   string       `json:"status"`
	Tables   []FileTables `json:"tables"`
	PDFFiles []string     `json:"pdf_files"`
	// Tutaj możesz sobie dodać np. pole "LLMResult" jak zaczniesz wołać model.
	LLMResult string `json:"llmResult"`
}

// parseExcelTables parsuje wszystkie arkusze z pliku Excel
func parseExcelTables(r io.Reader, filename string) ([]Table, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			log.Printf("excel close error for %s: %v", filename, cerr)
		}
	}()

	var tables []Table

	sheets := f.GetSheetList()
	for _, sheet := range sheets {
		rows, err := f.GetRows(sheet)
		if err != nil {
			return nil, err
		}

		// Możesz tu dodać logikę przycinania pustych wierszy/kolumn itd.
		tables = append(tables, Table{
			FileName: filename,
			Sheet:    sheet,
			Rows:     rows,
		})
	}

	return tables, nil
}

// parsePdfTables – placeholder, jeśli kiedyś będziesz chciał bawić się w wykrywanie tabel z PDF
// Na razie możesz PDF-y po prostu wysyłać jako pliki do LLM i olać parsowanie po stronie backendu.
func parsePdfTables(r io.Reader, filename string) ([]Table, error) {
	// 1. Zapisujemy upload do pliku tymczasowego,
	//    bo ledongthuc/pdf działa na ścieżce/do pliku.
	tmp, err := os.CreateTemp("", "uploaded-*.pdf")
	if err != nil {
		return nil, fmt.Errorf("could not create temp file: %w", err)
	}
	defer func() {
		tmp.Close()
		if err := os.Remove(tmp.Name()); err != nil {
			log.Printf("could not remove temp file %s: %v", tmp.Name(), err)
		}
	}()

	if _, err := io.Copy(tmp, r); err != nil {
		return nil, fmt.Errorf("could not copy pdf to temp file: %w", err)
	}

	// Musimy zresetować offset na początku pliku,
	// bo pdf.Open zakłada, że zaczyna od 0.
	if _, err := tmp.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("could not seek temp file: %w", err)
	}

	// 2. Otwieramy PDF-a przez ledongthuc/pdf
	f, reader, err := pdf.Open(tmp.Name())
	if err != nil {
		return nil, fmt.Errorf("could not open pdf: %w", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			log.Printf("error closing pdf file %s: %v", filename, cerr)
		}
	}()

	totalPages := reader.NumPage()
	var tables []Table

	for pageIndex := 1; pageIndex <= totalPages; pageIndex++ {
		p := reader.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		plain, err := p.GetPlainText(nil)
		if err != nil {
			return nil, fmt.Errorf("could not get plain text from page %d: %w", pageIndex, err)
		}

		lines := strings.Split(plain, "\n")
		var rows [][]string

		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				continue
			}

			// Na razie: jedna kolumna = cała linia tekstu.
			// Jeśli kiedyś będziesz chciał, możesz tu dodać split po wielu spacjach/tabach.
			rows = append(rows, []string{trimmed})
		}

		if len(rows) == 0 {
			continue
		}

		tables = append(tables, Table{
			FileName: filename,
			Sheet:    fmt.Sprintf("page_%d", pageIndex),
			Rows:     rows,
		})
	}

	return tables, nil
}

// ImportFilesAI – handler przyjmujący wiele plików i wyciągający z nich dane tabel
func ImportFilesAI(c *gin.Context) {
	// 1. DB (jeśli potrzebujesz)
	_, err := db.OpenDB()
	if err != nil {
		jsonRespond.Error(c, http.StatusInternalServerError, "could not connect to database", err)
		return
	}

	// 2. SECRET do JWT
	hmacSecret := []byte(os.Getenv("SECRET"))
	if len(hmacSecret) == 0 {
		jsonRespond.Error(c, http.StatusInternalServerError, "SECRET env is not defined", nil)
		return
	}

	// 3. Token z ciastka
	tokenCookie, err := c.Cookie("token")
	if err != nil {
		jsonRespond.Error(c, http.StatusUnauthorized, "token cookie is not set", err)
		return
	}

	token, err := jwt.Parse(tokenCookie, func(token *jwt.Token) (any, error) {
		return hmacSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !token.Valid {
		jsonRespond.Error(c, http.StatusUnauthorized, "invalid token", err)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		jsonRespond.Error(c, http.StatusInternalServerError, "could not read token claims", nil)
		return
	}

	subVal, ok := claims["sub"].(float64)
	if !ok {
		jsonRespond.Error(c, http.StatusInternalServerError, "sub claim is not a number", nil)
		return
	}

	issuerId := uint(subVal)
	log.Printf("ImportFilesAI called by user %d, claims=%v", issuerId, claims)

	// 4. Pobranie wielu plików z form-data (pole "files")
	form, err := c.MultipartForm()
	if err != nil {
		jsonRespond.Error(c, http.StatusBadRequest, "could not parse multipart form", err)
		return
	}

	fileHeaders, ok := form.File["files"]
	if !ok || len(fileHeaders) == 0 {
		jsonRespond.Error(c, http.StatusBadRequest, "no files uploaded (expected field 'files')", nil)
		return
	}

	var (
		allTables []FileTables
		pdfFiles  []string
	)

	// 5. Iteracja po plikach
	for _, fh := range fileHeaders {
		file, err := fh.Open()
		if err != nil {
			jsonRespond.Error(c, http.StatusBadRequest, "could not open file", err)
			return
		}
		defer file.Close()

		filename := fh.Filename
		lower := strings.ToLower(filename)

		switch {
		case strings.HasSuffix(lower, ".xls"), strings.HasSuffix(lower, ".xlsx"):
			tables, err := parseExcelTables(file, filename)
			if err != nil {
				jsonRespond.Error(c, http.StatusInternalServerError, "could not parse excel file", err)
				return
			}
			if len(tables) > 0 {
				allTables = append(allTables, FileTables{
					FileName: filename,
					Tables:   tables,
				})
			}
		case strings.HasSuffix(lower, ".pdf"):
			tables, err := parsePdfTables(file, filename)
			if err != nil {
				jsonRespond.Error(c, http.StatusInternalServerError, "could not parse pdf file", err)
				return
			}

			if len(tables) > 0 {
				allTables = append(allTables, FileTables{
					FileName: filename,
					Tables:   tables,
				})
			}

			pdfFiles = append(pdfFiles, filename)

		default:
			// Typ pliku nieobsługiwany – możesz:
			// - olać (jak tu),
			// - albo wywalić błąd.
			log.Printf("skipping unsupported file type: %s", filename)
		}
	}

	// 6. Przygotowanie payloadu pod LLM (jeśli chcesz)
	// Możesz sobie to wysłać do modelu jako JSON.
	llmInput := struct {
		UserID   uint         `json:"user_id"`
		Tables   []FileTables `json:"tables"`
		PDFFiles []string     `json:"pdf_files"`
	}{
		UserID:   issuerId,
		Tables:   allTables,
		PDFFiles: pdfFiles,
	}

	llmPayload, err := json.Marshal(llmInput)
	if err != nil {
		jsonRespond.Error(c, http.StatusInternalServerError, "could not marshal llm payload", err)
		return
	}

	// TODO: w tym miejscu wołasz swoje LLM:
	// result, err := callLLM(ctx, llmPayload)
	// i potem dorzucasz to do odpowiedzi.

	oaitoken := os.Getenv("OPEN_AI_KEY")
	if oaitoken == "" {
		jsonRespond.Error(c, http.StatusInternalServerError, "open ai token is not defined", nil)
		return
	}
	dataStructure, err := os.ReadFile("test.json")
	if err != nil {
		jsonRespond.Error(c, http.StatusInternalServerError, "could not read json file", nil)
		return
	}

	client := openai.NewClient(oaitoken)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleUser,
					Content: "Jesteś analizatorem ogromnej ilości danych z różnych baz danych wyeksportowanych do plików pdf lub xlsl które musisz zwrócić do jsona" +
						"Każda dana jest podzielona na osobne pliki pdf oraz excel które trzeba przerobić na poniższy obiekt JSON który zawiera te dane" +
						"Poniżej jest wypisane jak powinna wyglądać struktura JSON:" +
						string(dataStructure),
				},
				{
					Role: openai.ChatMessageRoleUser,
					Content: "Oto dane do analizy:" +
						string(llmPayload),
				},
			},
			Temperature: 0.1,
		},
	)

	if err != nil {
		jsonRespond.Error(c, http.StatusInternalServerError, "error when sending AI message", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)

	analysis := resp.Choices[0].Message.Content

	// 7. Na razie zwracamy „surowe” dane, żebyś widział co backend wypluwa.
	jsonRespond.SendJSON(c, ImportFilesAIResponse{
		Status:    "OK",
		Tables:    allTables,
		PDFFiles:  pdfFiles,
		LLMResult: analysis, // jak już będziesz miał
	})

	// Dla debugowania, możesz chwilowo logować:
	log.Printf("LLM payload preview: %s", string(llmPayload))
}
