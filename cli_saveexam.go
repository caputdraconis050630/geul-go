package main

import "fmt"

func (cli *CLI) saveExam(provider string, examCode string, filename string) {
	fmt.Println("Saving the %s.%s exam to a %s.pdf file", provider, examCode, filename)
}
