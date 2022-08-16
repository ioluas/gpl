package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"sync"

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

type response struct {
	Code int
	Body string
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}

func fetch(cmd *cobra.Command, args []string) {
	urls := os.Args[2:]
	if len(urls) < 1 {
		log.Fatal("at least one url is expected as argument")
	}
	wg := sync.WaitGroup{}
	options := cookiejar.Options{PublicSuffixList: publicsuffix.List}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}
	client := http.Client{Jar: jar}
	log.Printf("URLs to fetch: %v\n", urls)
	for _, url := range urls {
		wg.Add(1)
		go fetchUrl(&client, url, &wg)
	}
	wg.Wait()
}

func fetchUrl(c *http.Client, url string, wg *sync.WaitGroup) {
	res, err := c.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	resp := response{Code: res.StatusCode}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	resp.Body = string(b)
	fmt.Printf("Url:\t%s\n", url)
	fmt.Printf("StatusCode:\t%d\n", resp.Code)
	fmt.Printf("Length:\t%s\n", units.HumanSize(float64(len(resp.Body))))
	fmt.Print("---------------\n")
	wg.Done()
}
