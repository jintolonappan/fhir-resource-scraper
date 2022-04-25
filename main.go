package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocarina/gocsv"
)

type FhirResource struct {
	Module       string `csv:"Module"`
	Category     string `csv:"Category"`
	Resource     string `csv:"Resource"`
	ResourceDesc string `csv:"ResourceDesc"`
	Url          string `csv:"Url"`
}

func main() {
	url := "https://www.hl7.org/fhir/resourcelist.html"
	resp, err := http.Get(url)
	checkErr(err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		panic(fmt.Sprint("Expected data not received. Error code was ", resp.Status))
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	checkErr(err)

	resFile, err := os.OpenFile("fhir_resources.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	checkErr(err)
	defer resFile.Close()

	// baseurl is prefixed to each Resource URL in struct
	// <meta name="author" content="http://hl7.org/fhir">
	baseurl, _ := doc.Find("meta[name=author]").Attr("content")
	baseurl += "/"

	var resources []FhirResource
	var modules []string
	var categories []string

	// Module Count to fetch array element
	mc := -1

	// Category count to fetch array element
	cc := 0

	doc.Find("div #tabs-1>table>tbody>tr").Each(func(i int, root *goquery.Selection) {
		var resource FhirResource

		mod, err := root.Find("td.frm-group>div").Html()
		checkErr(err)
		resource.Module = mod
		modules = append(modules, mod)

		root.Find("tr.frm-group>td.frm-category").Each(func(i int, s *goquery.Selection) {
			resource.Category = s.Text()
			categories = append(categories, s.Text())
		})

		root.Find("tr.frm-contents>td.frm-set").Each(func(i int, s *goquery.Selection) {
			s.Find("li a").Each(func(i int, s *goquery.Selection) {
				resname := strings.TrimSpace(s.Text())
				if len(resname) > 2 { // This is the Resource Name
					thisresource := resname
					resource.Module = modules[mc]
					resource.Category = categories[cc]
					resource.Resource = thisresource
					resource.ResourceDesc, _ = s.Attr("title")
					resource.ResourceDesc = strings.ReplaceAll(resource.ResourceDesc, "\n", "")
					resource.Url, _ = s.Attr("href")
					resource.Url = baseurl + resource.Url
					resources = append(resources, resource)
				}
			})
			cc++
		})
		mc++
	})
	for i := range resources {
		fmt.Printf("Module: %s ", resources[i].Module)
		fmt.Printf("Cat: %s ", resources[i].Category)
		fmt.Printf("Res: %s ", resources[i].Resource)
		fmt.Printf("Url: %s\n", resources[i].Url)
	}
	err = gocsv.MarshalFile(&resources, resFile)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic("Encountered error. Exiting..")
	}
}
