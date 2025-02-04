package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"os"
	"path/filepath"
)

func (cli *CLI) saveDiscussions(provider string, examCode string) {
	fmt.Printf("Saving the discussions of %s.%s to a file\n", provider, examCode)
	fileName := fmt.Sprintf("%s_%s_discussions.txt", provider, examCode)
	EXAMTOPICS_BASE_URL := os.Getenv("EXAMTOPICS_BASE_URL")

	discussions, _ := DiscussionScraper(provider, examCode)

	// 현재 작업 디렉토리 가져오기
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Errorf("failed to get current directory: %v", err)
		return
	}

	// 파일 경로 생성
	filePath := filepath.Join(currentDir, fileName)

	// 파일 생성 (이미 존재하면 덮어쓰기)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Errorf("failed to create file: %v", err)
		return
	}
	defer file.Close()

	// 링크를 파일에 쓰기
	for i, discussion := range discussions {
		_, err := fmt.Fprintf(file, "%d: %s%s\n", i+1, EXAMTOPICS_BASE_URL, discussion)
		if err != nil {
			fmt.Errorf("failed to write to file: %v", err)
			return
		}
	}

	log.Infof("Successfully saved %d discussion links to %s", len(discussions), filePath)
}
