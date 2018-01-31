package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/codegangsta/cli"

	"github.com/fatih/color"
)

func main() {
	app := cli.NewApp()
	app.Name = "gopeoplesay"
	app.Usage = "Command Line Client for https://dopeoplesay.com"
	app.ArgsUsage = "[keywords]"
	app.HideHelp = true

	app.Action = func(c *cli.Context) {
		keywords := c.Args()

		if len(keywords) == 0 {
			cli.ShowAppHelp(c)
			os.Exit(1)
		}
		url := buildURL(keywords)
		crawl(url, keywords)
	}

	app.Run(os.Args)
}

func buildURL(keywords []string) string {
	return fmt.Sprintf("https://dopeoplesay.com/q/%s", url.QueryEscape(strings.Join(keywords, "-")))
}

func crawl(url string, keywords []string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("div.container div.example div p").Each(func(_ int, line *goquery.Selection) {
		text := line.Text()
		fmt.Println(highlightWords(text, keywords))
		fmt.Println("")
	})
}

func highlightWords(message string, keywords []string) string {
	words := []string{}
	for _, word := range keywords {
		words = append(words, regexp.QuoteMeta(word))
	}

	pattern := regexp.MustCompile(strings.Join(words, "|"))
	return pattern.ReplaceAllStringFunc(message, func(s string) string {
		return color.YellowString(s)
	})
}
