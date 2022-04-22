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
	doc.Find("div #tabs-1 .frm-group").Each(func(i int, root *goquery.Selection) {
		var resource FhirResource

		mod, err := root.Find(".frm-group .rotate div").Html()
		checkErr(err)
		resource.Module = mod

		root.Find(".frm-category").Each(func(i int, s *goquery.Selection) {
			resource.Category = s.Text()
			resources = append(resources, resource)
		})
	})

	fmt.Println(resources)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic("Encountered error. Exiting..")
	}
}
