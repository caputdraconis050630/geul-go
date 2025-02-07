package main

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/gofiber/fiber/v2/log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Exam struct {
	ExamName      string
	ExamLink      string
	ExamQuestion  string
	ExamChoices   []string
	ExamAnswer    string // Examtopics에서 제공하는 답안
	ExamMostVoted string // 사용자들이 정답이라고 투표한 답안
}

func ExamScraper(links []string) ([]Exam, error) {
	log.Infof("Starting ExamScrapper...")
	exams := []Exam{}

	jobs := make(chan string, maxConcurrent*2)  // 크롤링할 링크 전달 채널
	results := make(chan Exam, maxConcurrent*5) // 크롤링한 결과 전달 채널

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	var wg sync.WaitGroup
	for i := 0; i < maxConcurrent; i++ {
		wg.Add(1) // worker ++
		go func(workerID int) {
			defer wg.Done() // worker -- when all job is done
			log.Infof("Worker %d starting...", workerID)
			examWorker(ctx, jobs, results, workerID)
			log.Infof("Worker %d finished.", workerID)
		}(i)
	}

	go func() {
		for _, link := range links {
			jobs <- link
			log.Debugf("Job added to queue: %s", link)
		}
		close(jobs)
		log.Info("All jobs added to queue, jobs channel closed.")
	}()

	// Combine fetched data
	done := make(chan struct{})
	go func() {
		defer close(done)
		resultCount := 0
		for set := range results {
			exams = append(exams, set)
			resultCount++
			log.Infof("Completed %d exam questions.", resultCount)
		}
	}()

	go func() {
		wg.Wait()
		close(results)
		log.Info("All workers finished, results channel closed.")
	}()

	select {
	case <-done:
		log.Infof("All results collected successfully. Total Questions: %d", len(exams))
	case <-ctx.Done():
		log.Error("Context deadline exceeded")
		return nil, ctx.Err()
	}

	return exams, nil
}

func examWorker(ctx context.Context, jobs <-chan string, results chan<- Exam, workerID int) {
	for {
		select {
		case url, ok := <-jobs:
			if !ok {
				// 채널이 닫혔으면 워커 종료
				log.Infof("Worker %d: No more jobs, exiting", workerID)
				return
			}
			for attempt := 1; attempt <= maxRetries; attempt++ {
				log.Infof("Worker %d processing URL: %s (Attempt %d/%d)", workerID, url, attempt, maxRetries)
				exam, err := GetExamSet(url)
				if err == nil {
					results <- exam
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
		case <-ctx.Done():
			log.Infof("Worker %d: Context cancelled, stopping", workerID)
			return
		}
	}
}

// GetExamSet 파라미터로 받은 링크에 대해서 시험 문제/답안/최다투표 데이터를 스크래핑 후, Exam 구조체로 반환
func GetExamSet(link string) (Exam, error) {
	log.Debugf("Getting exam set for link: %s", link)
	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Get(link)
	if err != nil {
		log.Errorf("HTTP GET request failed for %s: %v", link, err)
		return Exam{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("HTTP request failed with status %d for URL: %s", resp.StatusCode, link)
		return Exam{}, &HTTPError{StatusCode: resp.StatusCode}
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Errorf("Failed to parse HTML for URL %s: %v", link, err)
		return Exam{}, err
	}

	exam := Exam{}
	doc.Find("body > div.sec-spacer.pt-50 > div > div:nth-child(5) > div > div.discussion-header-container > div.question-body.mt-3.pt-3.border-top").Each(func(i int, parentDiv *goquery.Selection) {
		question := parentDiv.Find("p").Text()
		choices := parentDiv.Find("div.question-choices-container ul li").Map(func(i int, s *goquery.Selection) string {
			// span 요소와 텍스트 노드 분리 추출
			clone := s.Clone()
			index := clone.Find("span.multi-choice-letter").Text()
			text := strings.TrimSpace(clone.Text())

			// MostVoted parsing logic
			//if clone.Has("span.badge.badge-success.most-voted-answer-badge").Length() > 0 {
			//	exam.ExamMostVoted = index
			//}
			clone.Remove()

			return index + " " + text
		})
		answer := parentDiv.Find("div.question-choices-container > ul > li.multi-choice-item.correct-hidden > span").Text()

		exam.ExamQuestion = question
		exam.ExamChoices = choices
		exam.ExamAnswer = answer

	})

	log.Infof("Successfully got exam set for URL: %s", link)

	return exam, nil
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
