package main

import (
	"goquery/scrapper"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

const fileName string = "jobs.csv"

func handleHome(c echo.Context) error {
	return c.File("home.html")
}

func handleScope(c echo.Context) error {
	defer os.Remove(fileName)
	term, err := strconv.Atoi(strings.ToLower(scrapper.CleanString(c.FormValue("term"))))
	scrapper.CheckError(err)
	scrapper.Scrape(term)
	return c.Attachment(fileName, fileName)
}

func main() {
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScope)
	e.Logger.Fatal(e.Start(":1323"))
}
