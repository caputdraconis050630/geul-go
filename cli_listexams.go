/*
 * cli_listexams.go
 * This file will contain the code for listing the exams offered by a specific provider
 */

package main

import "fmt"

func (cli *CLI) listExams(provider string) {
	fmt.Println("Listing the exams offered by the provider: ", provider)

}
