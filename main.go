package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/gocolly/colly"
)

var pageNum []string
var userName, dirName, url string
var directory, user bool

func wget(url, filepath string) error { // Define wget function for downloading & saving files
	cmdArgs := []string{url, "--content-on-error", "-N", "-q", "-O", filepath}

	cmd := exec.Command("wget", "-V")

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Wget seems to not be installed?\n")
		os.Exit(1)
		return err
	} else {
		cmd := exec.Command("wget", cmdArgs...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		return cmd.Run()
	}
}

func indexOf(element string, data []string) int { // Find the index of a certain string in a slice
	for i, o := range data {
		if element == o {
			return i
		}
	}
	return -1
}

func clearScreen() { // Clear stdout screen
	fmt.Printf("\033[H\033[2J")
}

func main() {
	args := os.Args[1:] // Assign all arguments to args not including the first arg

	for _, arg := range args { // Define arguments & display help
		switch arg {
		case "-d", "--dir":
			directory = true
		case "-u", "--user":
			user = true
		case "-v", "--version":
			fmt.Printf("Pascra v1.1.0-release\n")
			os.Exit(1)
		case "-h", "--help":
			clearScreen()

			fmt.Printf("Pascra v1.1.0-release written by z89, Instructions:\n\n\n")
			fmt.Printf("-h, --help		show these help instructions \n\n")
			fmt.Printf("-v, --version		display the program version\n\n")
			fmt.Printf("-d, --dir		specify the directory for the downloaded pastes\n\n")
			fmt.Printf("-u, --user		select the user from pastebin.com to download from\n\n")

			os.Exit(1)
		}
	}

	if user { // If the user exists from stdin arg
		index := indexOf("-u", args) + 1           // Get the args index for username
		userName = args[index]                     // Get username from stdin
		url = "https://pastebin.com/u/" + userName // Define url from username
	} else {
		fmt.Printf("\nERROR!\n")
		fmt.Printf("This tool must have a user specified to download their pastebin! \n\n")
		os.Exit(1)
	}

	if directory { // Select custom or default directory for pastes
		index := indexOf("-d", args) + 1
		dirName = args[index]
	} else {
		dirName = "default"
	}

	pasteCollector := colly.NewCollector()
	numCollector := colly.NewCollector()

	go numCollector.OnHTML(".pagination div a", func(e *colly.HTMLElement) {
		pageNum = append(pageNum, e.Text) // Setup collector to get page numbers and append to pageNum slice
	})

	numCollector.Visit(url) // Vist url and exec collector

	type paste struct {
		title string
		url   string
	}

	pastes := []paste{} // Make a slice of paste{} objects/structs
	pasteCounter := 0   // To count even number of paste titles for filtering

	go pasteCollector.OnHTML("tr td a[href]", func(e *colly.HTMLElement) {
		pasteCounter++
		if pasteCounter%2 == 1 { // For every even result append that paste object to pastes slice
			data := paste{title: e.Text, url: "https://pastebin.com/raw" + e.Attr("href")}
			pastes = append(pastes, data)
		}
	})

	go pasteCollector.OnError(func(r *colly.Response, err error) {
		if err.Error() == "Not Found" { // Usually a username error results in this error
			fmt.Printf("\n404 NOT FOUND!\n")
			fmt.Printf("The data you supplied resulted in a 404 not found, maybe you typed\n")
			fmt.Printf("in the wrong username? The program will now exit!\n\n")

			os.Exit(1)
		}
	})

	if len(pageNum) > 0 {
		pageNum := pageNum[:len(pageNum)-1] // Remove last slice element to remove "oldest" button

		for i := 1; i <= len(pageNum); i++ {
			pasteCollector.Visit(url + "/" + strconv.Itoa(i)) // For each page vist all the pastes for the user
		}
	} else {
		pasteCollector.Visit(url) // User only has one page worth of pastes to visit
	}

	startTime := time.Now() // Start the time for downloading time of pastes

	os.RemoveAll(dirName)      // Remove the current specified directory
	os.MkdirAll(dirName, 0770) // Make directory with correct permissions

	type statistic struct {
		status    string
		files     string
		filename  string
		url       string
		user      string
		directory string
	}

	completed := statistic{ // Completed version of console output to stdout
		status:    "\tCompleted\n",
		files:     "\t\t" + strconv.Itoa(len(pastes)) + "/" + strconv.Itoa(len(pastes)) + "\n",
		filename:  "\t~\n",
		url:       "\t\t~\n",
		user:      "\t\t" + userName + "\n",
		directory: "\t./" + dirName + "",
	}

	for index, paste := range pastes {
		filename := "./" + dirName + "/" + strconv.Itoa(index+1) + " - " + paste.title + ".txt"
		err := wget(paste.url, filename)

		if err == nil {
			clearScreen()

			fmt.Printf("Pastebin Crawler Statistics:\n\n")

			progress := statistic{ // In Progress version of console output to stdout
				status:    "\tIn Progress\n",
				files:     "\t\t" + strconv.Itoa(index+1) + "/" + strconv.Itoa(len(pastes)) + "\n",
				filename:  "\t\"" + filename[len(dirName)+3:] + "\"\n",
				url:       "\t\t" + paste.url + "\n",
				user:      "\t\t" + userName + "\n",
				directory: "\t./" + dirName + "",
			}

			str := " " + fmt.Sprintf("%+v", progress)[1:]
			fmt.Printf("%s", str[:len(str)-1])

			now := time.Now()

			fmt.Printf("\n\nElapsed duration: %v\n", now)
			fmt.Printf("Starting Time:    %v\n", startTime)
		}

	}
	clearScreen()

	fmt.Printf("Pastebin Crawler Statistics:\n\n")

	str := " " + fmt.Sprintf("%+v", completed)[1:]
	fmt.Printf("%s", str[:len(str)-1])

	now := time.Now()

	fmt.Printf("\n\nElapsed duration: %v\n", now)
	fmt.Printf("Starting Time:    %v\n", startTime)
}
