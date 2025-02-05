/*
 * cli_listexams.go
 * This file will contain the code for listing the exams offered by a specific provider
 */

package main

import "fmt"

func (cli *CLI) listExams(provider string) {
	fmt.Printf("Listing the exams offered by the provider: %s\n\n", provider)

	// TODO: Implement the code to validate that provider is valid

	exams := ExamListScraper(provider)

	fmt.Printf("%-4s  %-100s %s\n", "[Index]", "[Exam Name]", "[Exam Link]")

	for index, exam := range exams {
		indexPlusDot := fmt.Sprintf("%d.", index+1)
		fmt.Printf("%-4s      %-100s %s\n", indexPlusDot, exam.ExamName, exam.ExamLink)
	}
}
