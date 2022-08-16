package cmd

import (
	"bufio"
	"image"
	"image/color"
	"image/gif"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/spf13/cobra"
)

// lissajousCmd represents the lissajous command
var lissajousCmd = &cobra.Command{
	Use:   "lissajous",
	Short: "Generates a Lissajous figure encoded as gif",
	Long:  `Generates simple visual effect used in sci-fi films of the 1960s`,
	Run:   lissajous,
}

func init() {
	rootCmd.AddCommand(lissajousCmd)
}

func lissajous(cmd *cobra.Command, args []string) {
	var palette = []color.Color{
		color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 255,
		}, color.RGBA{
			R: 0,
			G: 255,
			B: 0,
			A: 255,
		}, color.RGBA{
			R: 255,
			G: 0,
			B: 0,
			A: 255,
		}, color.RGBA{
			R: 0,
			G: 0,
			B: 255,
			A: 255,
		}, color.RGBA{
			R: 255,
			G: 255,
			B: 0,
			A: 255,
		}, color.RGBA{
			R: 0,
			G: 255,
			B: 255,
			A: 255,
		}, color.RGBA{
			R: 255,
			G: 0,
			B: 255,
			A: 255,
		}, color.RGBA{
			R: 255,
			G: 255,
			B: 255,
			A: 255,
		}}
	const (
		cycles    = 5
		res       = 0.001
		size      = 250
		numFrames = 64
		delay     = 8
	)
	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount: numFrames}
	phase := 0.0
	for i := 0; i < numFrames; i++ {
		colorIdx := uint8(i % len(palette))
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), colorIdx)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	file, err := os.Create(execDir + "/data/lissajous.gif")
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("error when closing output file: %v\n", err)
		}
	}(file)
	buf := bufio.NewWriter(file)
	err = gif.EncodeAll(buf, &anim)
	if err != nil {
		log.Fatal(err)
	}
	err = buf.Flush()
	if err != nil {
		log.Printf("error when flushing data to file: %v\n", err)
	}
}
