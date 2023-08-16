package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var baseURL string = "https://www.jobkorea.co.kr/recruit/joblist?menucode=local&localorder=1"

type extractedJob struct {
	id       string
	title    string
	location string
	summary  string
}

func main() {
	var jobs []extractedJob
	totalPages := getPages()

	for i := 0; i < totalPages; i++ {
		extractedJobs := getPage(i)
		jobs = append(jobs, extractedJobs...)
	}

	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))
}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkError(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"LINK", "TITLE", "LOCATION", "SUMMARY"}

	wErr := w.Write(headers)
	checkError(wErr)

	for _, job := range jobs {
		jobSlice := []string{"https://www.jobkorea.co.kr/Recruit/GI_Read/"+strings.TrimSpace(job.id)+"?rPageCode=AM&logpath=21", job.title, job.location, job.summary}
		jwErr := w.Write(jobSlice)
		checkError(jwErr)
	}
}

func getPage(page int) []extractedJob {
	var jobs []extractedJob

	pageURL := baseURL + "#anchorGICnt_" + strconv.Itoa(page+1)
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	checkError(err)
	checkStatusCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkError(err)

	searchCard := doc.Find(".devloopArea")

	searchCard.Each(func(i int, s *goquery.Selection) {
		job := extractJob(s)
		jobs = append(jobs, job)
	})

	return jobs
}

func extractJob(s *goquery.Selection) extractedJob {
	info, _ := s.Attr("data-info")
	infoValueList := strings.Split(info, "|")
	id := infoValueList[0]
	title, _ := s.Find(".tplTit .normalLog").Attr("title")
	location := ""
	s.Find(".cell").Each(func(i int, cellInfo *goquery.Selection) {
		if i == 2 {
			location = cellInfo.Text()
		}
	})
	summary := s.Find(".dsc").Text()

	title = cleanString(title)
	location = cleanString(location)
	summary = cleanString(summary)

	return extractedJob{
		id:       id,
		title:    title,
		location: location,
		summary:  summary,
	}
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages() int {
	pages := 0

	res, err := http.Get(baseURL)
	checkError(err)
	checkStatusCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkError(err)

	doc.Find("#dvGIPaging").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages

}

func checkError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkStatusCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}
