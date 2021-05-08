package lib

import (
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type extractedJob struct {
	link		string
	title		string
	location	string
	salary		string
	summary		string
}

// Scrape Indeed by a term 📄
func Scrape(term string) {
	var baseURL string = "https://kr.indeed.com/jobs?q="+ term + "&limit=50"
	var jobs []extractedJob
	c := make(chan []extractedJob)
	totalPages := getPages(baseURL)

	for i := 0; i < totalPages; i++ {
		go getPage(i, baseURL, c)
	}

	for i:=0; i<totalPages; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)
	}

	writeJobs(jobs)
}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv"); checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"ID", "Title", "Location", "Salary", "Summary"}

	wErr := w.Write(headers); checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{"https://kr.indeed.com/viewjob?jk="+job.link, job.title, job.location, job.salary, job.summary}
		jwErr := w.Write(jobSlice); checkErr(jwErr)
	}
}

func getPage(page int, url string, mainChannel chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
	PageURL := url + "&start=" + strconv.Itoa(page*50)
	fmt.Println("requesting... ", PageURL)
	res, err := http.Get(PageURL); checkErr(err); checkCode(res)

	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body); checkErr(err)

	searchCards := doc.Find(".jobsearch-SerpJobCard")

	searchCards.Each(func(i int, card *goquery.Selection){
		go extractJob(card, c)
	})

	for i:=0; i<searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}

	mainChannel <- jobs
}

func extractJob(card *goquery.Selection, c chan <- extractedJob) {
	id, _ := card.Attr("data-jk")
	title := CleanString(card.Find(".title>a").Text())
	location := CleanString(card.Find(".sjcl").Text())
	salary := CleanString(card.Find(".salaryText").Text())
	summary := CleanString(card.Find(".summary").Text())
	c <- extractedJob{
		link: 		id,
		title: 		title,
		location: 	location,
		salary: 	salary,
		summary: 	summary}
}

func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str))," ")
}

func getPages(baseURL string) int {
	pages := 0
	res, err := http.Get(baseURL)

	checkErr(err)
	checkCode(res)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages
}

func checkErr(err error) {
	if err != nil {
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("request failed with Status: ", res.StatusCode)
	}
}