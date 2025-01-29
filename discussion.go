package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/gofiber/fiber/v2/log"
)

const maxConcurrent = 5

func DiscussionScraper(provider string, examCode string) ([]string, error) {
	pageNum := GetNumOfDiscussionPage(provider)
	providerDiscussionURL := os.Getenv("EXAMTOPICS_DISCUSSION_URL") + provider + "/"

	// 워커 풀 생성
	jobs := make(chan string, pageNum)
	results := make(chan []string, pageNum)

	// chromedp 컨텍스트 생성 (재사용)
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// 워커 고루틴 시작
	var wg sync.WaitGroup
	for i := 0; i < maxConcurrent; i++ {
		wg.Add(1)
		go worker(ctx, &wg, jobs, results)
	}

	// 작업 전송
	for i := 1; i <= pageNum; i++ {
		requestURL := providerDiscussionURL + strconv.Itoa(i) + "/"
		jobs <- requestURL
	}
	close(jobs)

	// 결과 수집
	go func() {
		wg.Wait()
		close(results)
	}()

	var discussions []string
	for links := range results {
		discussions = append(discussions, links...)
	}

	return discussions, nil
}

func worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan string, results chan<- []string) {
	defer wg.Done()
	for url := range jobs {
		links, err := GetDiscussionLink(ctx, url)
		if err != nil {
			log.Errorf("Error getting discussion links for %s: %v", url, err)
			continue
		}
		results <- links
	}
}

func GetDiscussionLink(ctx context.Context, requestURL string) ([]string, error) {
	var links []string
	err := chromedp.Run(ctx,
		chromedp.Navigate(requestURL),
		chromedp.WaitVisible("body > div.sec-spacer > div > div:nth-child(3) > div > div", chromedp.ByQuery),
		chromedp.Evaluate(`
            (function() {
                var links = [];
                var parentDivs = document.querySelectorAll("body > div.sec-spacer > div > div:nth-child(3) > div > div");
                parentDivs.forEach(function(parentDiv) {
                    var childDivs = parentDiv.querySelectorAll("div");
                    childDivs.forEach(function(childDiv) {
                        var a = childDiv.querySelector("div.col-7.col-md-6.discussion-column.discussion-title > div > h2 > a");
                        if (a && a.href) {
                            links.push(a.href);
                        }
                    });
                });
                return links;
            })()
        `, &links),
	)
	if err != nil {
		return nil, err
	}
	return links, nil
}

func GetNumOfDiscussionPage(provider string) int {
	log.Info("Getting the number of discussion pages for " + provider)

	requestURL := os.Getenv("EXAMTOPICS_DISCUSSION_URL") + provider + "/"
	log.Info("Request URL: " + requestURL)
	resp, err := http.Get(requestURL)
	HandleError(err)
	HandleStatusCodeError(resp)
	defer resp.Body.Close()

	log.Info("Parsing the response body")
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	HandleError(err)

	log.Info("Getting the total number of pages")
	totalPage, err := strconv.Atoi(doc.Find("body > div.sec-spacer > div > div.action-row-container.mb-4 > div > span > span.discussion-list-page-indicator > strong:nth-child(3)").Text())
	HandleError(err)

	return totalPage
}
