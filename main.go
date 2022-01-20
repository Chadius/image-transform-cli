package main

import (
	"bytes"
	"flag"
	"fmt"
	creatingsymmetry "github.com/Chadius/creating-symmetry"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	filenameArguments := extractFilenameArguments()

	inputImageDataByteStream := new(bytes.Buffer)
	png.Encode(inputImageDataByteStream, openSourceImage(filenameArguments))

	formulaDataByteStream, formulaErr := os.Open(filenameArguments.FormulaFilename)
	if formulaErr != nil {
		log.Fatal(formulaErr)
	}

	outputSettingsJSONByteStream := []byte(
		fmt.Sprintf(
			"{'output_width':%d,'output_height':%d}",
			filenameArguments.OutputWidth,
			filenameArguments.OutputHeight,
		),
	)
	outputSettingsDataByteStream := bytes.NewReader(outputSettingsJSONByteStream)

	var output bytes.Buffer

	// TODO Put into Args
	// TODO Build new client object
	// TODO pass onto new client object and get the output
	defaultTransformer := creatingsymmetry.FileTransformer{}
	defaultTransformer.ApplyFormulaToTransformImage(inputImageDataByteStream, formulaDataByteStream, outputSettingsDataByteStream, &output)

	outputReader := bytes.NewReader(output.Bytes())
	outputImage, _ := png.Decode(outputReader)
	outputToFile(filenameArguments.OutputFilename, outputImage)
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

func openSourceImage(filenameArguments *FilenameArguments) image.Image {
	reader, err := os.Open(filenameArguments.SourceImageFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	colorSourceImage, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	return colorSourceImage
}

func outputToFile(outputFilename string, outputImage image.Image) {
	err := os.MkdirAll(filepath.Dir(outputFilename), 0777)
	if err != nil {
		panic(err)
	}
	outputImageFile, err := os.Create(outputFilename)
	if err != nil {
		panic(err)
	}
	defer outputImageFile.Close()
	png.Encode(outputImageFile, outputImage)
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
