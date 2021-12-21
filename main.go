package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	filenameArguments := extractFilenameArguments()
	loadFormulaFile(filenameArguments.FormulaFilename)
	// TODO Call creating-symmetry library
	// TODO Get output stream
	// TODO Put into image
}

// FilenameArguments assume the user provides filenames to create a pattern.
type FilenameArguments struct {
	FormulaFilename     string
	OutputFilename      string
	OutputHeight        int
	OutputWidth         int
	SourceImageFilename string
}

func extractFilenameArguments() *FilenameArguments {
	var sourceImageFilename, outputFilename, outputDimensions string
	formulaFilename := "data/oldformula.yml"
	flag.StringVar(&formulaFilename, "f", "data/oldformula.yml", "See -oldformula")
	flag.StringVar(&formulaFilename, "oldformula", "data/oldformula.yml", "The filename of the oldformula file. Defaults to data/oldformula.yml")

	flag.StringVar(&sourceImageFilename, "in", "", "See -source. Required.")
	flag.StringVar(&sourceImageFilename, "source", "", "Source filename. Required.")

	flag.StringVar(&outputFilename, "out", "", "Output filename. Required.")
	outputDimensions = "200x200"
	flag.StringVar(&outputDimensions, "size", "200x200", "Output size in pixels, separated with an x. Default to 200x200.")
	flag.Parse()

	checkSourceArgument(sourceImageFilename)
	outputWidth, outputHeight := checkOutputArgument(outputFilename, outputDimensions)

	return &FilenameArguments{
		FormulaFilename:     formulaFilename,
		OutputFilename:      outputFilename,
		OutputHeight:        outputHeight,
		OutputWidth:         outputWidth,
		SourceImageFilename: sourceImageFilename,
	}
}

func checkSourceArgument(sourceImageFilename string) {
	if sourceImageFilename == "" {
		log.Fatal("missing source filename")
	}
}

func checkOutputArgument(outputFilename, outputDimensions string) (int, int) {
	if outputFilename == "" {
		log.Fatal("missing output filename")
	}

	outputWidth, widthErr := strconv.Atoi(strings.Split(outputDimensions, "x")[0])
	if widthErr != nil {
		log.Fatal(widthErr)
	}

	outputHeight, heightErr := strconv.Atoi(strings.Split(outputDimensions, "x")[1])
	if heightErr != nil {
		log.Fatal(heightErr)
	}

	return outputWidth, outputHeight
}

func loadFormulaFile(filename string) io.Reader {
	formulaFile, err := os.Open(filename)
	if err != nil {
		println(err.Error())
		log.Fatal(err)
	}
	return formulaFile
}
