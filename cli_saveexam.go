package main

import (
	"github.com/gofiber/fiber/v2/log"
)

func (cli *CLI) saveExam(importFilename string, exportFilename string) {
	log.Info("Saving exam...")

	// importFilename 으로부터 discussions links 긁어오기
	links, err := ReadFile(importFilename)
	if err != nil {
		panic(err)
	}
	//for _, link := range links {
	//	fmt.Println(link)
	//}

	// Scrape
	examSet, err := ExamScraper(links)
	if err != nil {
		panic(err)
	}

	// PDF 파일로 저장
	Export2PDF(examSet, exportFilename)
	log.Infof("Saved the Exam to %s", exportFilename)
}
