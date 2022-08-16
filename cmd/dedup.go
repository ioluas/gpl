package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// dedupCmd represents the dedup command
var dedupCmd = &cobra.Command{
	Use:   "dedup",
	Short: "De-Duplicate lines from stdin",
	Long:  `just pipe text to it`,
	Run:   dedup,
}

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "De-Deduplicate lines in files passed from command arguments",
	Long: `
	Iterates over files in arguments and prints out uniques lines with count.
	
	Example usage:
		gpl dedup file text.txt text2.txt`,
	Run: fromFile,
}

func init() {
	rootCmd.AddCommand(dedupCmd)

	dedupCmd.AddCommand(fileCmd)
	fileCmd.Flags().BoolP("write-result", "w", false, "Write deduplicated lines to file")
}

func dedup(cmd *cobra.Command, args []string) {
	counts := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		counts[input.Text()]++
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

var writeToFile bool = false

func checkWriteFlag(f *pflag.Flag) {
	if f.Name == "write-result" {
		writeToFile = true
	}
}

func fromFile(cmd *cobra.Command, args []string) {
	wg := sync.WaitGroup{}
	cmd.Flags().Visit(checkWriteFlag)
	for _, file := range args {
		wg.Add(1)
		go countLines(file, &wg)
	}
	wg.Wait()
}

func countLines(file string, wg *sync.WaitGroup) {
	f, err := os.Open(file)
	if err != nil {
		wg.Done()
		log.Printf("failed to open file handle: %s %v", f.Name(), err)
		return
	}
	defer func() {
		err = f.Close()
		if err != nil {
			log.Printf("failed to close file handle: %s %v", f.Name(), err)
		}
	}()
	counts := make(map[string]int)
	fmt.Printf("\nFile:\t%s\n", f.Name())
	in := bufio.NewScanner(f)
	for in.Scan() {
		counts[in.Text()]++
	}
	for line, n := range counts {
		fmt.Printf("%d\t%s\n", n, line)
	}
	if writeToFile {
		outName := f.Name() + ".dedup"
		log.Println("Creating file: ", outName)
		df, err := os.Create(outName)
		if err != nil {
			log.Printf("failed to create out file: %s %v", outName, err)
		}
		defer func(df *os.File) {
			log.Println("Closing file: ", df.Name())
			_ = df.Close()
		}(df)
		w := bufio.NewWriter(df)
		for l := range counts {
			n, err := w.WriteString(l + "\n")
			if err != nil {
				log.Printf("failed to write string to buffer: %s %v", f.Name(), err)
			}
			log.Println("Wrote bytes: ", n)
		}
		err = w.Flush()
		if err != nil {
			log.Printf("failed to flush file: %s %v", f.Name(), err)
		}
		log.Println("Flushed writer buffer")
	}
	wg.Done()
}
