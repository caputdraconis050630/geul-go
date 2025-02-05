package main

import (
	"bufio"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"log"
	"net/http"
	"os"
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
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func Export2PDF(exams []Exam, filename string) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8Font("NotoSans", "", "NotoSans-Regular.ttf")
	pdf.AddUTF8Font("NotoSans", "B", "NotoSans-Bold.ttf")

	for i, exam := range exams {
		pdf.AddPage()

		// 문제 번호
		pdf.SetFont("NotoSans", "B", 14)
		pdf.Cell(0, 10, fmt.Sprintf("Q %d", i+1))
		pdf.Ln(10)

		// 문제
		pdf.SetFont("NotoSans", "B", 12)
		pdf.MultiCell(0, 7, exam.ExamQuestion, "", "", false)
		pdf.Ln(5)

		// 선지
		pdf.SetFont("NotoSans", "", 10)
		choices := strings.Split(exam.ExamChoices, "\n") // TODO: sep 수정
		for _, choice := range choices {
			pdf.MultiCell(0, 5, choice, "", "", false)
		}
		pdf.Ln(5)

		// 정답
		pdf.SetFont("NotoSans", "B", 10)
		pdf.Cell(0, 7, fmt.Sprintf("Most Voted: %s", exam.ExamMostVoted))
	}

	err := pdf.OutputFileAndClose("exams.pdf")
	if err != nil {
		panic(err)
	}

	fmt.Println("PDF 파일이 성공적으로 생성되었습니다: exams.pdf")
}
