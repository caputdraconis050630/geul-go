package main

import "net/http"

type Exam struct {
	ExamName string
	ExamLink string
}

func (exam *Exam) String() string {
	return exam.ExamName
}

func (exam *Exam) Link() string {
	return exam.ExamLink
}

func (exam *Exam) SetName(name string) {
	exam.ExamName = name
}

func (exam *Exam) SetLink(link string) {
	exam.ExamLink = link
}

func GetExams(provider string) []Exam {
	exams := []Exam{}

	resp, err := http.Get(provider)
	if err != nil || resp.StatusCode != 200 {
		panic(err)
	}

	defer resp.Body.Close()

	return exams
}
