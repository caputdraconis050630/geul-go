package main

import "fmt"

func (cli *CLI) listProviders() {
	fmt.Println("Listing the providers of exams offered by examtopics :)\n\n")

	provider := ProviderScraper()

	fmt.Printf("%-4s   %-40s %15s %-100s\n", "[Index]", "[Provider Name]", "[Num of Exams]", "[Link to Exam List]")

	for index, p := range provider {
		indexPlusDot := fmt.Sprintf("%d.", index+1)
		fmt.Printf("%-5s     %-40s  %-15d %-100s\n", indexPlusDot, p.getProvider(), p.getNumOfExams(), p.getExamListLink())
	}
}
