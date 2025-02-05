package main

import (
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

type Exam struct {
	ExamName      string
	ExamLink      string
	ExamQuestion  string
	ExamChoices   string
	ExamMostVoted string // 정답으로 간주
}

func ExamListScraper(provider string) []Exam {
	exams := []Exam{}

	BaseLink := os.Getenv("EXAMTOPICS_BASE_URL")
	ProviderListLink := os.Getenv("EXAMTOPICS_EXAM_URL") + provider + "/"

	resp, err := http.Get(ProviderListLink)
	HandleError(err)
	HandleStatusCodeError(resp)

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	HandleError(err)

	doc.Find("body > div.sec-spacer > div:nth-child(1) > div:nth-child(2) > div > ul > li").Each(func(i int, s *goquery.Selection) {
		examName := s.Find("span").Text()
		examLink, _ := s.Find("a").Attr("href")

		exams = append(exams, Exam{
			ExamName: examName,
			ExamLink: BaseLink + examLink,
		})
	})

	return exams
}
