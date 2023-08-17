package scrapper

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

type extractedJob struct {
	id       string
	title    string
	location string
	summary  string
}

//Scrape jobkorea
func Scrape(totalPagesTerm int) {
	var baseURL string = "https://www.jobkorea.co.kr/recruit/joblist?menucode=local&localorder=1"
	var c = make(chan []extractedJob)

	var jobs []extractedJob
	totalPages := getPages(baseURL)
	if totalPagesTerm > 0 {
		totalPages = totalPagesTerm
	}

	for i := 0; i < totalPages; i++ {
		go getPage(i, baseURL, c)
	}
	
	for i :=0; i< totalPages; i++ {
		extractedJobs := <- c
		jobs = append(jobs, extractedJobs...)
	}

	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))
}

func getPage(page int, url string, mainC chan<- []extractedJob){
	var c = make(chan extractedJob)

	var jobs []extractedJob

	pageURL := url + "#anchorGICnt_" + strconv.Itoa(page+1)
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	CheckError(err)
	checkStatusCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	CheckError(err)

	searchCard := doc.Find(".devloopArea")

	searchCard.Each(func(i int, s *goquery.Selection) {
		go extractJob(s, c)
	})
	for i:=0;i<searchCard.Length();i++{
		job := <-c
		jobs = append(jobs, job)
	}
	mainC <- jobs
}

func extractJob(s *goquery.Selection, c chan<- extractedJob){
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

	title = CleanString(title)
	location = CleanString(location)
	summary = CleanString(summary)

	c <- extractedJob{
		id:       id,
		title:    title,
		location: location,
		summary:  summary,
	}
}

//CleanString cleans a string
func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages(url string) int {
	pages := 0

	res, err := http.Get(url)
	CheckError(err)
	checkStatusCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	CheckError(err)

	doc.Find("#dvGIPaging").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages

}

func writeJobs(jobs []extractedJob) {
	c := make(chan []string)
	file, err := os.Create("jobs.csv")
	CheckError(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"LINK", "TITLE", "LOCATION", "SUMMARY"}

	wErr := w.Write(headers)
	CheckError(wErr)

	for _, job := range jobs {
		go writeJob(job, c)
	}
	jobSlices := [][]string{}
	for i:=0;i<len(jobs);i++{
		jobSlice := <- c
		jobSlices = append(jobSlices, jobSlice)
	}
	jwErr := w.WriteAll(jobSlices)
	CheckError(jwErr)
}

func writeJob(job extractedJob, c chan []string) {
	c <- []string{"https://www.jobkorea.co.kr/Recruit/GI_Read/"+strings.TrimSpace(job.id)+"?rPageCode=AM&logpath=21", job.title, job.location, job.summary}
}

//CheckError check a error
func CheckError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkStatusCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}
