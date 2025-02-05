package main

import (
	"github.com/gofiber/fiber/v2/log"
)

func (cli *CLI) saveExam(provider string, examCode string, importFilename string, exportFilename string) {
	log.Info("Saving exam...")

	// importFilename 으로부터 discussions links 긁어오기
	links, err := ReadFile(importFilename)
	if err != nil {
		panic(err)
	}

	// Scrape
	examSet, err := ExamScraper(links)
	if err != nil {
		panic(err)
	}

	// PDF 파일로 저장
	Export2PDF(examSet, exportFilename)
	log.Info("Saved the %s exam to a %s.pdf file", provider, exportFilename)
}
