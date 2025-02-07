package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct{}

// This function will run the CLI
func (cli *CLI) Run() {
	cli.validateArgs(os.Args)

	listProvidersCmd := flag.NewFlagSet("listproviders", flag.ExitOnError)
	listExamsCmd := flag.NewFlagSet("listexams", flag.ExitOnError)
	saveDiscussionsCmd := flag.NewFlagSet("savediscussions", flag.ExitOnError)
	saveExamCmd := flag.NewFlagSet("saveexam", flag.ExitOnError)

	listExamsProvider := listExamsCmd.String("provider", "", "The provider of the exams")
	saveDiscussionsProvider := saveDiscussionsCmd.String("provider", "", "The provider of the exams")
	saveDiscussionsExam := saveDiscussionsCmd.String("exam", "", "The exam code")
	saveExamImport := saveExamCmd.String("import", "", "The filename to get the exam discussion links(output of savediscussions feature)")
	saveExamExport := saveExamCmd.String("export", "", "The filename to save the exam to")

	switch os.Args[1] {
	case "listproviders":
		err := listProvidersCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listexams":
		err := listExamsCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "savediscussions":
		err := saveDiscussionsCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "saveexam":
		err := saveExamCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if listProvidersCmd.Parsed() {
		cli.listProviders()
	}

	if listExamsCmd.Parsed() {
		if *listExamsProvider == "" {
			listExamsCmd.Usage()
			os.Exit(1)
		}
		cli.listExams(*listExamsProvider)
	}

	if saveDiscussionsCmd.Parsed() {
		if *saveDiscussionsProvider == "" || *saveDiscussionsExam == "" {
			saveDiscussionsCmd.Usage()
			os.Exit(1)
		}
		cli.saveDiscussions(*saveDiscussionsProvider, *saveDiscussionsExam)
	}

	if saveExamCmd.Parsed() {
		if *saveExamImport == "" || *saveExamExport == "" {
			saveExamCmd.Usage()
			os.Exit(1)
		}
		cli.saveExam(*saveExamImport, *saveExamExport)
	}
}

// This function will print the usage of the CLI
func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  listproviders - List the providers of exams offered by examtopics.")
	fmt.Println("  listexams -provider PROVIDER - List the exams offered by the specific provider.")
	fmt.Println("  savediscussions -provider PROVIDER -exam EXAMCODE - Save the discussion links of the exam to a text file.")
	fmt.Println("  saveexam -provider PROVIDER -exam EXAM_CODE -import FILENAME -export FILENAME - Save the exam to a pdf file.")
}

// This function will list the providers of exams offered by examtopics
func (cli *CLI) validateArgs(args []string) {
	if len(args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}
