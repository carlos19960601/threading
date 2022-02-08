package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/spf13/viper"
	"github.com/zengqiang96/threading/common"
	"github.com/zengqiang96/threading/plotter"
	"github.com/zengqiang96/threading/thread"
)

var config common.Config

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	if err = viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("Fatal unmarshal config: %w \n", err))
	}
}

const MAX_COMPUTING_TIME_PER_FRAME = 20 // ms

func main() {
	file, err := os.OpenFile(config.SourceImage, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file failed, err: %v", err)
	}
	defer file.Close()
	sourceImage, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("image decode failed, err: %v", err)
	}

	svgPlotter := &plotter.PlotterSVG{}
	threadComputer := thread.NewThreadComputer(sourceImage, config)
	threadPlotter := thread.NewThreadPlotter(svgPlotter, threadComputer)

	for threadComputer.ComputeNextSegments(MAX_COMPUTING_TIME_PER_FRAME) {
		threadPlotter.Plot()
	}
	threadPlotter.Plotter.Finalize()
}
