package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gofiber/fiber/v2/log"
)

const (
	maxConcurrent = 5
	httpTimeout   = 30 * time.Second
	maxRetries    = 3
	retryDelay    = 5 * time.Second
)

func DiscussionScraper(provider string, examCode string) ([]string, error) {
	log.Infof("Starting DiscussionScraper for provider: %s, examCode: %s", provider, examCode)

	pageNum := GetNumOfDiscussionPage(provider)
	log.Infof("Total number of pages to scrape: %d", pageNum)

	providerDiscussionURL := os.Getenv("EXAMTOPICS_DISCUSSION_URL") + provider + "/"

	jobs := make(chan string, maxConcurrent*2)
	results := make(chan []string, maxConcurrent*5)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	var wg sync.WaitGroup
	for i := 0; i < maxConcurrent; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			log.Infof("Worker %d started", workerID)
			worker(ctx, jobs, results, examCode, workerID)
			log.Infof("Worker %d finished", workerID)
		}(i)
	}

	go func() {
		for i := 1; i <= pageNum; i++ {
			requestURL := providerDiscussionURL + strconv.Itoa(i) + "/"
			jobs <- requestURL
			log.Debugf("Job added to queue: %s", requestURL)
		}
		close(jobs)
		log.Info("All jobs added to queue, jobs channel closed")
	}()

	var discussions []string
	done := make(chan struct{})
	go func() {
		defer close(done)
		resultCount := 0
		for links := range results {
			discussions = append(discussions, links...)
			resultCount++
			log.Infof("Collected results from page %d/%d", resultCount, pageNum)
		}
	}()

	go func() {
		wg.Wait()
		close(results)
		log.Info("All workers finished, results channel closed")
	}()

	select {
	case <-done:
		log.Infof("All results collected successfully. Total links: %d", len(discussions))
	case <-ctx.Done():
		log.Error("Context deadline exceeded")
		return nil, ctx.Err()
	}

	return discussions, nil
}

func worker(ctx context.Context, jobs <-chan string, results chan<- []string, examCode string, workerID int) {
	for url := range jobs {
		for attempt := 1; attempt <= maxRetries; attempt++ {
			log.Infof("Worker %d processing URL: %s (Attempt %d/%d)", workerID, url, attempt, maxRetries)
			links, err := GetDiscussionLink(url, examCode)
			if err == nil {
				results <- links
				log.Infof("Worker %d successfully processed URL: %s", workerID, url)
				break
			}

			log.Warnf("Worker %d: Attempt %d/%d failed for %s: %v",
				workerID, attempt, maxRetries, url, err)

			if attempt == maxRetries {
				log.Errorf("Worker %d: Final attempt failed for %s", workerID, url)
				break
			}

			select {
			case <-time.After(retryDelay):
				log.Infof("Worker %d: Retrying URL %s after delay", workerID, url)
			case <-ctx.Done():
				log.Infof("Worker %d: Context cancelled, stopping", workerID)
				return
			}
		}
	}
}

func GetDiscussionLink(url string, examCode string) ([]string, error) {
	log.Debugf("Fetching discussion links from URL: %s", url)
	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Get(url)
	if err != nil {
		log.Errorf("HTTP GET request failed for %s: %v", url, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("HTTP request failed with status %d for URL: %s", resp.StatusCode, url)
		return nil, &HTTPError{StatusCode: resp.StatusCode}
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Errorf("Failed to parse HTML for URL %s: %v", url, err)
		return nil, err
	}

	var links []string
	doc.Find("body > div.sec-spacer > div > div:nth-child(3) > div > div > div").Each(func(i int, parentDiv *goquery.Selection) {
		childDiv := parentDiv.Find("div.col-7.col-md-6.discussion-column.discussion-title > div > h2 > a")
		title := childDiv.Text()
		if strings.Contains(strings.ToLower(title), strings.ToLower(examCode)) {
			if link, linkExists := childDiv.Attr("href"); linkExists {
				links = append(links, link)
				log.Debugf("Found link: %s", link)
			}
		}

	})

	log.Infof("Retrieved %d links from URL: %s", len(links), url)
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

	log.Info("Total Page: " + strconv.Itoa(totalPage))
	return totalPage
}

type HTTPError struct {
	StatusCode int
}

func (e *HTTPError) Error() string {
	return "HTTP request failed with status: " + strconv.Itoa(e.StatusCode)
}
