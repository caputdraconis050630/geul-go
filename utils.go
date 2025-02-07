package main

import (
	"bufio"
	"fmt"
	"github.com/jaytaylor/html2text"
	"github.com/jung-kurt/gofpdf"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func HandleStatusCodeError(resp *http.Response) {
	if resp.StatusCode != 200 {
		log.Fatal("Failed to get the page. Status code: ", resp.StatusCode)
	}
}

func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// ReadFile 주어진 파일로부터 데이터를 읽어와서, 문자열 배열로 리턴
func ReadFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url, _ := extractURL(scanner.Text()) // URL만 추출
		lines = append(lines, url)
	}

	return lines, scanner.Err()
}

func cleanText(input string) string {
	// HTML을 텍스트로 변환
	text, err := html2text.FromString(input, html2text.Options{PrettyTables: true})
	if err != nil {
		return input // 에러 발생 시 원본 텍스트 반환
	}

	// 앞뒤 공백 제거 및 연속된 공백을 하나로 줄임
	text = strings.TrimSpace(text)
	text = strings.Join(strings.Fields(text), " ")

	return text
}

func Export2PDF(exams []Exam, filename string) {
	envBoldTTF := os.Getenv("FONT_FILENAME_BOLD")
	envRegularTTF := os.Getenv("FONT_FILENAME_REGULAR")
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8Font("NotoSans", "", envRegularTTF)
	pdf.AddUTF8Font("NotoSans", "B", envBoldTTF)

	for i, exam := range exams {
		pdf.AddPage()

		// 문제 번호
		pdf.SetFont("NotoSans", "B", 14)
		pdf.Cell(0, 10, fmt.Sprintf("Q%d", i+1))
		pdf.Ln(10)

		// 문제
		pdf.SetFont("NotoSans", "B", 12)
		cleanedQuestion := cleanText(exam.ExamQuestion)
		pdf.MultiCell(0, 7, cleanedQuestion, "", "", false)
		pdf.Ln(5)

		// 선지
		pdf.SetFont("NotoSans", "", 10)
		choices := exam.ExamChoices
		for _, choice := range choices {
			cleanedChoice := cleanText(choice)
			pdf.MultiCell(0, 5, cleanedChoice, "", "", false)
		}
		pdf.Ln(5)

		// 정답
		pdf.SetFont("NotoSans", "B", 10)
		pdf.Cell(0, 7, fmt.Sprintf("Answer: %s", cleanText(exam.ExamAnswer)))
		pdf.Ln(5)

		// Discussion Link
		pdf.SetFont("NotoSans", "B", 6)
		pdf.Cell(0, 7, fmt.Sprintf("Reference Link: %s", exam.ExamLink))
	}

	err := pdf.OutputFileAndClose(filename)
	if err != nil {
		panic(err)
	}

	fmt.Println("file is successfully saved:" + filename)
}

func extractURL(input string) (string, error) {
	// 정규표현식 패턴 정의
	pattern := `\d+:\s*(https?://\S+)`

	// 정규표현식 컴파일
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("REGEX expression error: %v", err)
	}

	// 매칭 실행
	matches := re.FindStringSubmatch(input)
	if len(matches) < 2 {
		return "", fmt.Errorf("Failed to extract URL from %s", input)
	}

	// 첫 번째 캡처 그룹 (URL) 반환
	return matches[1], nil

}
