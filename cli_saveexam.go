package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
)

func (cli *CLI) saveExam(provider string, examCode string, importFilename string, exportFilename string) {
	fmt.Println("Saving the %s.%s exam to a %s.pdf file", provider, examCode, exportFilename)

	// importFilename 으로부터 discussions links 긁어오기
	links, err := ReadFile(importFilename)
	if err != nil {
		panic(err)
	}

	// Main Logic
	examSet, err := ExamScraper(links)
	if err != nil {
		panic(err)
	}

	// PDF 파일로 저장
	Export2PDF(examSet, exportFilename)
	log.Info("Saved the %s exam to a %s.pdf file", provider, exportFilename)
}
