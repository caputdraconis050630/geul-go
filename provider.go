package main

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PROVIDER struct {
	Provider     string
	ExamListLink string // e.g. https://www.examtopics.com/exams/amazon/
	NumOfExams   int
}

func ProviderScraper() []PROVIDER {
	// env
	ExamListBaseLink := os.Getenv("EXAMTOPICS_BASE_URL")
	ProviderListLink := os.Getenv("EXAMTOPICS_EXAM_URL")

	providers := []PROVIDER{}
	resp, err := http.Get(ProviderListLink)
	HandleError(err)
	HandleStatusCodeError(resp)
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	HandleError(err)

	doc.Find("body > div.sec-spacer > div:nth-child(1) > div:nth-child(2) > div").Each(func(i int, s *goquery.Selection) {
		fullText := s.Find("a").Text()
		parts := strings.Split(fullText, "(")
		provider := strings.TrimSpace(parts[0])
		num, _ := strconv.Atoi(strings.Trim(parts[1], " exams)"))

		link, _ := s.Find("a").Attr("href")

		providers = append(providers, PROVIDER{
			Provider:     provider,
			ExamListLink: ExamListBaseLink + link,
			NumOfExams:   num,
		})
	})

	// DUMMY DATA
	// providers := []PROVIDER{
	// 	{
	// 		Provider:     "Amazon",
	// 		ExamListLink: "https://www.examtopics.com/exams/amazon/",
	// 	},
	// 	{
	// 		Provider:     "Cisco",
	// 		ExamListLink: "https://www.examtopics.com/exams/cisco/",
	// 	},
	// 	{
	// 		Provider:     "CompTIA",
	// 		ExamListLink: "https://www.examtopics.com/exams/comptia/",
	// 	},
	// 	{
	// 		Provider:     "Microsoft",
	// 		ExamListLink: "https://www.examtopics.com/exams/microsoft/",
	// 	},
	// 	{
	// 		Provider:     "VMware",
	// 		ExamListLink: "https://www.examtopics.com/exams/vmware/",
	// 	},
	// }

	return providers
}

func (provider *PROVIDER) getProvider() string {
	return provider.Provider
}

func (provider *PROVIDER) getExamListLink() string {
	return provider.ExamListLink
}

func (provider *PROVIDER) getNumOfExams() int {
	return provider.NumOfExams
}
