package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/gocolly/colly"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

// define paste struct
type paste struct {
	title string
	url   string
}

func createDirectory(name string, subDirectory bool) string {
	// don't attempt to delete subdirectories
	if !subDirectory {
		directory, err := os.Stat(name)
		if err == nil {
			fmt.Print(name + "/ already exists, are you sure you want to overwrite it? (y/N): ")

			var response string

			fmt.Scanln(&response)
			fmt.Println(" ")

			if response != "y" {
				os.Exit(0)
			}

			err := os.RemoveAll(directory.Name())

			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
		}
	}

	err := os.Mkdir(name, 0777)

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	return name
}

func getPastes(c *colly.Collector) *[]paste {
	var pastes []paste

	// find all pastes on the page
	c.OnHTML(".maintable tbody tr td:first-child a[href]", func(e *colly.HTMLElement) {
		// create a paste struct from fetched data
		record := paste{title: strings.ReplaceAll(e.Text, "/", "-"), url: "https://pastebin.com/raw" + e.Attr("href")}

		pastes = append(pastes, record)
	})

	return &pastes
}

func downloadFile(bar *progressbar.ProgressBar, filepath string, url string, wg *sync.WaitGroup) {
	// fetch the paste
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	// check response is valid
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 429 {
			fmt.Println("\npastebin is blocking your requests, try increasing your delay")
		} else {
			fmt.Printf("bad status: %v", resp.Status)
		}

		os.Exit(0)
	}

	// create a file to write to
	out, err := os.Create(filepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	defer out.Close()

	// write the response to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// increment progress bar
	bar.Add(1)

	// decrement wait group
	wg.Done()
}

func downloadPage(synchronous bool, path string, page string, user string, interval time.Duration, wg *sync.WaitGroup) int {
	collector := colly.NewCollector(
		colly.AllowedDomains("pastebin.com"),
	)

	// get all pastes from the page
	pastes := getPastes(collector)

	// fetch the pastes found
	err := collector.Visit("https://pastebin.com/u/" + user + "/" + page)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// set up progress bar
	bar := progressbar.NewOptions(len(*pastes),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetDescription(user+" : [blue]"+page+"[reset]"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[white]#",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionOnCompletion(func() {
			fmt.Println(" ")
		}),
	)

	// download each paste
	for _, paste := range *pastes {
		wg.Add(1)

		// set the delay between requests
		time.Sleep(interval * time.Millisecond)

		// check if user wants to download synchronously
		if synchronous {
			downloadFile(bar, path+paste.title+".txt", paste.url, wg)
		} else {
			go downloadFile(bar, path+paste.title+".txt", paste.url, wg)
		}
	}

	return len(*pastes)
}

func Fetch() *cobra.Command {
	var paginationResult []string
	var rootDir string

	command := &cobra.Command{
		Use:   "fetch",
		Short: "Fetch command to download pastes",
		Long:  "fetch - command to download pastes from pastebin.com",
		Run: func(cmd *cobra.Command, args []string) {

			// get user flags
			user, _ := cmd.Flags().GetString("user")
			dir, _ := cmd.Flags().GetString("dir")
			pages, _ := cmd.Flags().GetStringSlice("pages")
			delay, _ := cmd.Flags().GetInt64("delay")
			synchronous, _ := cmd.Flags().GetBool("synchronous")

			interval := time.Duration(delay)

			// if user is not specified, throw help menu
			if len(user) == 0 {
				cmd.Help()
				os.Exit(0)
			}

			collector := colly.NewCollector(
				colly.AllowedDomains("pastebin.com"),
			)

			// check if user has no pastes
			collector.OnHTML(".notice", func(e *colly.HTMLElement) {
				fmt.Println("user has no pastes")
				os.Exit(0)
			})

			// get pagination numbers for user
			collector.OnHTML(".pagination div a", func(e *colly.HTMLElement) {
				paginationResult = append(paginationResult, e.Attr("data-page"))
			})

			// fetch the pagination numbers
			err := collector.Visit("https://pastebin.com/u/" + user)

			// check if user exists
			if err != nil {
				if err.Error() == "Not Found" {
					fmt.Println("user was not found")
				} else {
					fmt.Println(err)
				}
				os.Exit(0)
			}

			var path string
			pageTarget := 1
			var scope int

			// if the user has multiple pages
			if !(len(paginationResult) == 0) {
				// get the last page number & convert to int
				converted, err := strconv.Atoi(paginationResult[len(paginationResult)-1])

				if err != nil {
					fmt.Println(err)
					os.Exit(0)
				}

				// crawled pagination numbers are 1 less than actual pages
				pageTarget = converted + 1
			}

			// if pages flag was specified
			for _, page := range pages {
				// convert page to int
				pageInt, err := strconv.Atoi(page)

				if err != nil {
					fmt.Println(err)
					os.Exit(0)
				}

				// if page flag is 0 or a negative number
				if pageInt < 1 {
					fmt.Println("page number must be above 0")
					os.Exit(0)
				}

				// check if page is greater than available pages for user
				if pageInt > pageTarget {
					if pageTarget == 1 {
						fmt.Println(user + " only has " + strconv.Itoa(pageTarget) + " page")
					} else {
						fmt.Println(user + " only has " + strconv.Itoa(pageTarget) + " pages")
					}

					os.Exit(0)
				}
			}

			// create directory with custom or default name
			if len(dir) == 0 {
				rootDir = createDirectory(user, false)
			} else {
				rootDir = createDirectory(dir, false)
			}

			// if pages flag was not specified
			if len(pages) == 0 {
				// set scope to the last page
				scope = pageTarget
			} else {
				// set scope to length of pages flag
				scope = len(pages)
			}

			startTime := time.Now()
			wg := new(sync.WaitGroup)
			var total int

			// loop for the specifc number of pages
			for i := 1; i <= scope; i++ {
				var page string

				// if pages flag was specified
				if len(pages) != 0 {
					page = pages[i-1]                       // set page to the page flag
					path = rootDir + "/" + pages[i-1] + "/" // create the path for the page directory
					createDirectory(path, true)             // create the page directory
				} else {
					page = strconv.FormatInt(int64(i), 10)                       // set the page to the current iteration
					path = rootDir + "/" + strconv.FormatInt(int64(i), 10) + "/" // create the path for the page directory
					createDirectory(path, true)                                  // create the page directory
				}

				// download pastes from page
				total += downloadPage(synchronous, path, page, user, interval, wg)

				// wait for all pastes to be downloaded before moving on to next page
				wg.Wait()
			}

			fmt.Printf("\nstatistics   ")
			fmt.Printf("\n----------------\n")
			w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
			fmt.Fprintf(w, "username:\t %v\n", user)
			fmt.Fprintf(w, "total pastes:\t %v\n", total)
			fmt.Fprintf(w, "time elapsed:\t %v\n", time.Since(startTime).Round(time.Second))
			w.Flush()
		},
	}

	// user flags
	command.PersistentFlags().String(
		"user",
		"",
		"user to download pastes from (required)",
	)

	command.MarkFlagRequired("user")

	command.PersistentFlags().String(
		"dir",
		"",
		"directory to store downloaded pastes",
	)

	command.PersistentFlags().StringSlice(
		"pages",
		[]string{},
		"download specifc pages from user",
	)

	command.PersistentFlags().Int64(
		"delay",
		250,
		"time delay between downloads to prevent too many requests (milliseconds)",
	)

	command.PersistentFlags().Bool(
		"synchronous",
		false,
		"turn off concurrent downloading of pastes (slows performance)",
	)

	return command
}
