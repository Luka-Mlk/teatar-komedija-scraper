package main

import (
	"fmt"
	"log"
	"math/rand"
	"teatarScraper/pkg/genAgent"
	model "teatarScraper/pkg/models"
	"time"

	"github.com/gocolly/colly"
)

func main() {
	initScrape := colly.NewCollector()
	secondScrape := initScrape.Clone()

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
		content := e.ChildText(":nth-child(2) > :nth-child(1) >:nth-child(1) > :nth-child(1) > :nth-child(1) > :nth-child(2)")
		if title == "" {
			title = "\nEmpty string bad HTML"

		}
		if content == "" {
			content = "\nEmpty string bad HTML"
		}
		event := model.Event{
			Title: title,
			Content: content,
		}
		fmt.Println(event)
	})

	initScrape.Visit("https://teatarkomedija.mk/portfolio-category/репертоар")
}