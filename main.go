package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/PuerkitoBio/goquery"
	_ "go/printer"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

//  https://www.youtube.com/watch?v=mS74M-rnc90 - Tutorial WebScapper Ebay

func getHtml(url string) *http.Response {
	response, err := http.Get(url)
	checkErr(err)

	if response.StatusCode > 400 {
		fmt.Println("Status Code:", response.StatusCode)
	}

	if response.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", response.StatusCode, response.Status)
	}
	return response
}
func saveToCSV(result []string) {
	file, error := os.Create("results.csv")
	checkErr(error)
	writer := csv.NewWriter(file)
	writer.Write(result)
	writer.Flush()
}

func scrapeHtml(doc *goquery.Document) []string {
	var result []string

	doc.Find(".product-wrapper").Each(func(index int, item *goquery.Selection) {
		a := item.Find("h2 a")
		price := item.Find(".price.small")
		title := strings.TrimSpace(strings.Trim(a.Text(), "Marktverfügbarkeit prüfen"+"Warenkorb"))
		priceValue := strings.TrimSpace(strings.Trim(price.Text(), "Warenkorb"))
		fmt.Println("Titel:")
		fmt.Printf("%+v \n", title)
		fmt.Println("Preis:")
		fmt.Printf("%+v \n", priceValue)
		result = append(result, title, priceValue)

	})
	saveToCSV(result)
	return result
}

func main() {

	ClearTerminal()
	PrintMenue()
	for true {
		ExecuteCommand()
	}

}
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func loopOver2(page int, command string) []string {

	var result []string
	url := fmt.Sprintf("https://www.mediamarkt.ch/de/search.html?query=%s&searchProfile=onlineshop&channel=mmchde&page=%d", command, page)

	response := getHtml(url)
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)

	fmt.Println(&doc)
	checkErr(err)

	result = append(result, scrapeHtml(doc)...)
	return result

}

func loopOver(command string) []string {
	var result []string
	//Markenbezeichnungen in "command" generienen teilweise Unterseiten welche nicht gelesen werden können
	go loopOver2(1, command)
	go loopOver2(2, command)
	go loopOver2(3, command)
	go loopOver2(4, command)
	go loopOver2(5, command)
	return result
}

func PrintMenue() {
	fmt.Println(`
##########################################
#******** Webscapper Mediamarkt **********
#********* Wähle eine Option *************
# 1. Suchen
# 2. Export (Beta)
#
# q. Beenden
##########################################
##########################################
##########################################
`)
}

func askForCommand() string {
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	response = strings.TrimSpace(response)
	return response
}

func ExecuteCommand() {
	command := askForCommand()
	parseCommand(command)
}

func ClearTerminal() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func toInt(info string) int {
	aInt, _ := strconv.Atoi(info)
	return aInt
}

func toStr(info int) string {
	aStr := strconv.Itoa(info)
	return aStr
}

func parseCommand(input string) {
	switch {
	case input == "1":
		ClearTerminal()
		fmt.Println("Suchwort eingeben:")
		command := askForCommand()
		loopOver(command)
		break
	case input == "2":
		ClearTerminal()
		PrintMenue()
		break
	}
}
