package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"text/template"
)

const (
	domain      = "https://xkcd.com"
	lineLength  = 80
	defaultPath = "/info.0.json"
	pathFmt     = "/%d/info.0.json"
	pageSize    = 20
	indexHeader = `Num       Date        Title
------------------------------------------------------------------------------`
	indexFmt = "%-8d  %2s/%02s/%s  %-58s\n"
	tmpl     = `
`
)

var list = flag.Bool("l", false, "Lists the most recent xkcd comics.")
var comicNum = flag.Int("n", 0, "The xkcd comic number to fetch.")

var comicReport = template.Must(template.New("comic_info").
	Funcs(template.FuncMap{"wrapText": wrapText}).
	ParseFiles("templates/comic_info"))

type comicData struct {
	Num        int
	Year       string
	Month      string
	Day        string
	Link       string
	Title      string
	SafeTitle  string `json:"safe_title"`
	Img        string
	Alt        string
	News       string
	Transcript string
}

func main() {
	flag.Parse()
	if *list {
		printIndex(index())
	} else {
		c := fetchComic(*comicNum)
		if c == nil {
			os.Exit(1)
		}

		printComic(c)
	}
}

func index() []comicData {
	mostRecent := fetchComic(0)
	if mostRecent == nil {
		return nil
	}

	comics := make([]comicData, pageSize)
	comics[0] = *mostRecent

	var wg sync.WaitGroup
	wg.Add(pageSize - 1)

	for i := 1; i < pageSize; i++ {
		go func(i int) {
			defer wg.Done()
			comics[i] = *fetchComic(mostRecent.Num - i)
		}(i)
	}
	wg.Wait()

	return comics
}

func fetchComic(comicNum int) *comicData {
	res, err := http.Get(domain + pathFor(comicNum))
	if err != nil {
		fmt.Fprintf(os.Stderr, "There was an error fetching the comic.")
		return nil
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr,
			"The server responded with code %v\nMessage: %v\n",
			res.StatusCode, res.Status)
		return nil
	}

	var comic comicData
	if err := json.NewDecoder(res.Body).Decode(&comic); err != nil {
		fmt.Fprintf(os.Stderr,
			"There was an error parsing the response JSON.\n%v\n", err)
		return nil
	}

	return &comic
}

func printComic(comic *comicData) {
	if err := comicReport.Execute(os.Stdout, comic); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
}

func printIndex(comics []comicData) {
	fmt.Println(indexHeader)
	for _, c := range comics {
		fmt.Printf(indexFmt, c.Num, c.Month, c.Day, c.Year, c.Title)
	}
}

func wrapText(input string) string {
	lines := len(input) / lineLength
	output := make([]byte, lines+len(input))

	// Copy every line except for the last one.
	for l := 0; l < lines; l++ {
		iStart := l * lineLength
		iEnd := iStart + lineLength
		copy([]byte(output[iStart+l:iEnd+l]), input[iStart:iEnd])
		output[iEnd+l] = '\n'
	}

	// Copy the last line (which is not necessary lineLength long).
	lastStart := lines * lineLength
	lastEnd := len(input)
	copy([]byte(output[lastStart+lines:lastEnd+lines]), input[lastStart:lastEnd])

	return string(output)
}

func pathFor(comicNum int) string {
	if comicNum > 0 {
		return fmt.Sprintf(pathFmt, comicNum)
	} else {
		return defaultPath
	}
}
