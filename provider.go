package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PROVIDER struct {
	Provider     string
	ExamListLink string // e.g. https://www.examtopics.com/exams/amazon/
	NumOfExams   int
}

const ProviderListLink = "https://www.examtopics.com/exams/"
const ExamListBaseLink = "https://www.examtopics.com"

func ProviderScraper() []PROVIDER {
	providers := []PROVIDER{}
	resp, err := http.Get(ProviderListLink)
	if err != nil || resp.StatusCode != 200 {
		panic(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

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
