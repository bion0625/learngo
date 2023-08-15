package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

var baseURL string = "https://www.jobkorea.co.kr/recruit/joblist?menucode=local&localorder=1"
// var baseURL string = "https://www.jobkorea.co.kr/recruit/joblist?menucode=local&localorder=1#anchorGICnt_1"

func main() {
	totalPages := getPages()

	for i := 0; i < totalPages; i++ {
		getPage(i)
	}
}

func getPage(page int){
	pageURL := baseURL + "#anchorGICnt_" + strconv.Itoa(page + 1)
	fmt.Println("Requesting", pageURL)
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
