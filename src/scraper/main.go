package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"teatarScraper/pkg/genAgent"
	model "teatarScraper/pkg/models"
	"time"

	"github.com/gocolly/colly"
)

func main() {
	initScrape := colly.NewCollector()
	secondScrape := initScrape.Clone()
	var Events []model.Event

	initScrape.OnRequest(func(r * colly.Request){
		userAgent := genAgent.GenBiasAgent()
		r.Headers.Set("User-Agent", userAgent)
	})

	initScrape.OnError(func(_ * colly.Response, err error){
		log.Println("Error:", err)
	})

	// initScrape.OnHTML("#qodef-page-content > :nth-child(1) > :nth-child(1) > :nth-child(1) > :nth-child(1) > :nth-child(1)", func(e * colly.HTMLElement){
	// 	randNum := rand.Intn(10) + 5
	// 	time.Sleep(time.Duration(randNum) * time.Second)
	// 	fmt.Println(e.ChildAttr("article > a", "href"))
	// 	// fmt.Println(e.ChildText("article > :nth-child(1) > :nth-child(2) > h4 > a"))
	// })

	initScrape.OnHTML("article", func(e * colly.HTMLElement){
		randNum := rand.Intn(10) + 5
		time.Sleep(time.Duration(randNum) * time.Second)
		url := e.ChildAttr("a", "href")
		secondScrape.Visit(url)
	})

	secondScrape.OnRequest(func(r * colly.Request){
		userAgent := genAgent.GenBiasAgent()
		r.Headers.Set("User-Agent", userAgent)
	})

	secondScrape.OnError(func(_ * colly.Response, err error){
		log.Println("Error:", err)
	})

	secondScrape.OnHTML("#qodef-page-outer", func(e * colly.HTMLElement){
		title := e.ChildText(":nth-child(1) > :nth-child(1) > :nth-child(1) > h5")
		content := ""
		contentSelectors := []string{
			"#qodef-page-content > div > div > div > section.elementor-section.elementor-top-section.elementor-element.elementor-element-2192be3.elementor-section-full_width.qodef-elementor-content-grid.elementor-section-height-default.elementor-section-height-default > div > div > div > section > div > div.elementor-column.elementor-col-50.elementor-inner-column.elementor-element.elementor-element-ab57097 > div > div > div",
			"#qodef-page-content > div > div > div > section.elementor-section.elementor-top-section.elementor-element.elementor-element-f5c67d6.elementor-section-boxed.elementor-section-height-default.elementor-section-height-default.qodef-elementor-content-no > div",
			"#qodef-page-content > div > div > div > section.elementor-section.elementor-top-section.elementor-element.elementor-element-95b1ede.elementor-section-boxed.elementor-section-height-default.elementor-section-height-default.qodef-elementor-content-no > div > div > div > div > div",
		}
		for _, selector := range contentSelectors {
			if content == "" {
				content = e.ChildText(selector)
			}
		}
		if title == "" { 
			title = e.ChildText("h3")
		}
		event := model.Event{
			Title: title,
			Content: strings.ReplaceAll(strings.ReplaceAll(content, "\t", ""), "\n", " "),
		}
		fmt.Println(event)
		Events = append(Events, event)
	})

	initScrape.Visit("https://teatarkomedija.mk/portfolio-category/репертоар")

	fp, _ := filepath.Abs("")
	file, err := os.Create(fp + "/src/json/teatar.json")
	if err != nil {
		log.Println("Error creating file", err)
	}
	defer file.Close()

	writer := json.NewEncoder(file)
	writer.Encode(Events)
}