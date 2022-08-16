package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"

	"github.com/docker/go-units"
	"github.com/spf13/cobra"
	"golang.org/x/net/publicsuffix"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetch urls from cli input",
	Long:  ``,
	Run:   fetch,
}

type fetchResult struct {
	Code      int
	Url, Body string
	time      float64
}

func init() {
	rootCmd.AddCommand(fetchCmd)
	fetchCmd.PersistentFlags().Bool("use-channels", false, "Enable use of channels in fetch go routines")
	fetchCmd.PersistentFlags().StringSlice("url", []string{}, "Slice of Urls to fetch")
}

func fetch(cmd *cobra.Command, args []string) {
	start := time.Now()
	urls, err := cmd.PersistentFlags().GetStringSlice("url")
	if err != nil || len(urls) < 1 {
		log.Fatalf("error reading url flag: %v", err)
	}
	useChannels, err := cmd.PersistentFlags().GetBool("use-channels")
	if err != nil {
		log.Printf("error reading use-channels flag: %v", err)
		useChannels = false
	}

	// create client wth cookie jar
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}
	client := http.Client{Jar: jar}

	fmt.Printf("URLs to fetch: %v\n", urls)
	fmt.Print("---------------\n")

	if useChannels {
		ch := make(chan fetchResult)
		defer close(ch)
		for _, url := range urls {
			go fetchUrl(&client, url, nil, ch)
		}
		for _, _ = range urls {
			res := <-ch
			fmt.Printf("Url:\t%s\n", res.Url)
			fmt.Printf("StatusCode:\t%d\n", res.Code)
			fmt.Printf("Length:\t%s\n", units.HumanSize(float64(len(res.Body))))
			fmt.Printf("Fetch timei:\t%.3fs\n", res.time)
			fmt.Print("---------------\n")
		}
	} else {
		wg := sync.WaitGroup{}
		for _, url := range urls {
			wg.Add(1)
			go fetchUrl(&client, url, &wg, nil)
		}
		wg.Wait()
	}
	fmt.Printf("Command time: %.3fs\n", time.Since(start).Seconds())
}

func fetchUrl(c *http.Client, url string, wg *sync.WaitGroup, ch chan<- fetchResult) {
	start := time.Now()
	res, err := c.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	resp := fetchResult{Code: res.StatusCode, Url: url}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	_ = res.Body.Close()
	resp.Body = string(b)
	resp.time = time.Since(start).Seconds()
	if ch != nil {
		ch <- resp
	} else if wg != nil {
		fmt.Printf("Url:\t%s\n", url)
		fmt.Printf("StatusCode:\t%d\n", resp.Code)
		fmt.Printf("Length:\t%s\n", units.HumanSize(float64(len(resp.Body))))
		fmt.Printf("Fetch time:\t%.3fs\n", time.Since(start).Seconds())
		fmt.Print("---------------\n")
		wg.Done()
	}
}
