package main

// Import required packages
import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gookit/color"
)

var pageNum []string
var paste [][]string
var dirName, userName, url string
var directory, quiet, user, dontVerbose, help, version bool

func main() {
	args := os.Args[1:]

	for _, s := range args {
		switch s {
		case "-d":
			directory = true
		case "--dir":
			directory = true
		case "-q":
			quiet = true
		case "--quiet":
			quiet = true
		case "-u":
			user = true
		case "--user":
			user = true
		case "-v":
			version = true
		case "--version":
			version = true
		case "-h":
			help = true
		case "--help":
			help = true
		}
	}

	if version {
		color.Style{color.FgWhite, color.OpBold, color.BgRed}.Printf(" Pascra v1.0 \n")
		os.Exit(1)
	}

	if help {
		clearScreen()
		color.Style{color.FgWhite, color.OpBold, color.BgRed}.Printf("\n Pascra v1.0-alpha written by z89  ")
		color.Style{color.FgWhite, color.OpBold, color.BgLightMagenta}.Printf(" Usage: \n\n\n")

		fmt.Printf("		-h, --help		show this help file\n\n")
		fmt.Printf("		-v, --version		do not display directory overwrite warning\n\n")
		fmt.Printf("		-d, --dir		specify the directory for the downloaded pastes\n\n")
		fmt.Printf("		-u, --user		select the user from pastebin.com to download from\n\n")
		fmt.Printf("		-q, --quiet		do not display directory overwrite warning and show limited features\n\n")

		os.Exit(1)
	}

	if user {
		index := indexOf("-u", args) + 1
		userName = args[index]
		url = "https://pastebin.com/u/" + userName

	} else if !user {
		color.Style{color.FgWhite, color.OpBold, color.BgLightMagenta}.Printf(" ERROR! ")
		color.Style{color.FgWhite, color.OpBold}.Printf(" This tool must have a user specified to download their pastebin! \n\n")

		os.Exit(1)
	}

	if directory {
		index := indexOf("-d", args) + 1
		dirName = args[index]
	} else if !directory {
		dirName = "pastebin.com"
	}

	if quiet {
		connectToPastebin()
	} else if !quiet {
		if directory {
			clearScreen() // Clear the CLI

			color.Style{color.FgWhite, color.OpBold, color.BgRed}.Printf(" WARNING! ")
			color.Style{color.FgWhite, color.OpBold}.Printf(" THIS TOOL WILL OVERWRITE YOUR CHOSEN DIRECTORY!\n\n")

			c := askForConfirmation("Do you really want to overwrite your chosen directory '" + dirName + "'?")
			if !c {
				os.Exit(1)
			} else {
				connectToPastebin()
			}

		} else if !directory {
			clearScreen()
			color.Style{color.FgWhite, color.OpBold, color.BgRed}.Printf(" WARNING! ")
			color.Style{color.FgWhite, color.OpBold}.Printf(" THIS TOOL WILL OVERWRITE THE DEFAULT DIRECTORY!\n\n")

			c := askForConfirmation("Do you really want to overwrite the 'pastebin.com' directory?")
			if !c {
				os.Exit(1)
			} else {
				connectToPastebin()
			}
		}
	}

	c := colly.NewCollector()
	t := colly.NewCollector()

	c.OnHTML(".pagination div a", func(e *colly.HTMLElement) {
		pageNum = append(pageNum, e.Text)
	})

	c.OnError(func(r *colly.Response, err error) {
		if err.Error() == "Not Found" {
			color.Style{color.FgWhite, color.OpBold, color.BgRed}.Printf("\n 404 NOT FOUND! ")
			color.Style{color.FgWhite, color.OpBold}.Printf(" The data you supplied resulted in a 404 not found, the program will exit!\n\n")

			os.Exit(1)
		}
	})

	c.Visit(url)

	count := 0
	index := 0
	t.OnHTML("tr td a[href]", func(e *colly.HTMLElement) {
		count++
		if count%2 == 1 {
			paste = append(paste, []string{e.Text, e.Attr("href")})
			index++
		}
	})

	if len(pageNum) > 0 {
		pageNum := pageNum[:len(pageNum)-1]

		for i := 1; i <= len(pageNum); i++ {
			t.Visit(url + "/" + strconv.Itoa(i))
		}

	} else {
		t.Visit(url)
	}

	startTime := time.Now()

	os.RemoveAll(dirName)
	os.MkdirAll(dirName, 0770)

	dotCounter := 0
	for index, s := range paste {
		filename := "./" + dirName + "/" + strconv.Itoa(index+1) + " - " + s[0] + ".txt"
		url := "https://pastebin.com/raw" + s[1]
		err := wget(url, filename)

		if err == nil {

			if quiet {
				clearScreen()
				if dotCounter%2 == 0 {
					color.Style{color.FgYellow, color.OpBold, color.OpReverse}.Printf(" Downloading pastes.. ")
				} else {
					color.Style{color.FgYellow, color.OpBold, color.OpReverse}.Printf(" Downloading pastes.  ")
				}
				fmt.Printf("  ")
				color.Style{color.FgRed, color.OpBold}.Printf(" File " + strconv.Itoa(index+1) + "/" + strconv.Itoa(len(paste)) + " ")

				dotCounter++
			} else if !quiet {
				clearScreen()
				color.Style{color.Cyan, color.OpUnderscore, color.OpBold}.Printf("Pastebin Crawler Statistics:\n\n")
				color.Style{color.White, color.OpBold}.Printf("Status:     ")
				color.Style{color.FgYellow, color.OpBold, color.OpReverse}.Printf(" In Progress \n\n")
				color.Style{color.White, color.OpBold}.Printf("Files:      ")
				color.Style{color.FgRed, color.OpBold}.Printf(" " + strconv.Itoa(index+1) + "/" + strconv.Itoa(len(paste)) + " \n")
				color.Style{color.White, color.OpBold}.Printf("Filename:   ")
				color.Style{color.FgRed, color.OpBold}.Printf(" \"" + filename[len(dirName)+3:] + "\" \n")
				color.Style{color.White, color.OpBold}.Printf("Download:   ")
				color.Style{color.FgRed, color.OpBold}.Printf(" " + url + " \n")
				color.Style{color.White, color.OpBold}.Printf("User:       ")
				color.Style{color.FgRed, color.OpBold}.Printf(" " + userName + " \n")
				color.Style{color.White, color.OpBold}.Printf("Directory:  ")
				color.Style{color.FgRed, color.OpBold}.Printf(" ./" + dirName + "  \n\n")

				now := time.Now()

				color.Style{color.Yellow, color.OpBold}.Printf("Elapsed duration: %v\n", now)
				color.Style{color.Green, color.OpBold}.Printf("Starting Time:    %v\n", startTime)
			}

		} 
	}

	if quiet {
		clearScreen()
		color.Style{color.FgGreen, color.OpBold}.Printf(" Download completed! Files are in the '" + dirName + "' directory\n\n")
	} else if !quiet {
		clearScreen()
		color.Style{color.Cyan, color.OpUnderscore, color.OpBold}.Printf("Pascra Statistics:\n\n")
		color.Style{color.White, color.OpBold}.Printf("Status:     ")
		color.Style{color.FgGreen, color.OpBold, color.OpReverse}.Printf(" Completed \n\n")
		color.Style{color.White, color.OpBold}.Printf("Files:      ")
		color.Style{color.FgGreen, color.OpBold}.Printf(" " + strconv.Itoa(len(paste)) + "/" + strconv.Itoa(len(paste)) + " \n")
		color.Style{color.White, color.OpBold}.Printf("Filename:   ")
		color.Style{color.FgGreen, color.OpBold}.Printf(" ~ \n")
		color.Style{color.White, color.OpBold}.Printf("Download:   ")
		color.Style{color.FgGreen, color.OpBold}.Printf(" ~ \n")
		color.Style{color.White, color.OpBold}.Printf("User:       ")
		color.Style{color.FgGreen, color.OpBold}.Printf(" " + userName + " \n")
		color.Style{color.White, color.OpBold}.Printf("Directory:  ")
		color.Style{color.FgGreen, color.OpBold}.Printf(" ./" + dirName + "  \n\n")

		now := time.Now()

		color.Style{color.Yellow, color.OpBold}.Printf("Elapsed duration: %v\n", now)
		color.Style{color.Green, color.OpBold}.Printf("Starting Time:    %v\n", startTime)
	}

}

func connectToPastebin() {
	clearScreen()
	color.Style{color.FgWhite, color.OpBold}.Printf("Connecting to pastebin.com...\n")
}

func clearScreen() {
	fmt.Printf("\033[H\033[2J")
}

func wget(url, filepath string) error {
	cmdArgs := []string{url, "--content-on-error", "-N", "-q", "-O", filepath}

	cmd := exec.Command("wget", "-V")

	err := cmd.Run()
	if err != nil {
		color.Style{color.FgWhite, color.OpBold}.Printf("Wget seems to not be installed?\n")
		os.Exit(1)
		return err
		
	} else {
		cmd := exec.Command("wget", cmdArgs...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		return cmd.Run()
	}
	
}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		color.Style{color.FgWhite, color.OpBold, color.BgRed}.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
