package main

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type FhirResource struct {
	Module       string
	Category     string
	Resource     string
	ResourceDesc string
	Url          string
}

func main() {

	fmt.Println("Testing")
	url := "https://www.hl7.org/fhir/resourcelist.html"
	fmt.Println(url)

	resp, err := http.Get(url)

	checkErr(err)

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Print("Couldnt fetch data from web. ")
		fmt.Println("Error code was", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	checkErr(err)

	// content, err := doc.Find("div #tabs-1").Html()
	// checkerr(err)
	// fmt.Println(content)

	var resources []FhirResource
	doc.Find("div #tabs-1 .frm-group .rotate div").Each(func(i int, s *goquery.Selection) {
		// modules, err := s.Find("td").Html()
		// checkerr(err)
		// fmt.Println(modules)
		fmt.Print(s.Text())
		var resource FhirResource
		resource.Module = s.Text()

		resources = append(resources, resource)
	})
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic("Encountered error. Exiting..")
	}
}
