package command

import (
	"bytes"
	"context"
	"fmt"
	creatingsymmetry "github.com/Chadius/creating-symmetry"
	"github.com/chadius/image-transform-server/rpc/transform/github.com/chadius/image_transform_server"
	"net/http"
	"os"
)

// Processor processes commands.
type Processor struct {
	localImageTransformer  creatingsymmetry.TransformerStrategy
	remoteImageTransformer image_transform_server.ImageTransformer
}

// ProcessArgumentsToTransformImage will read the given arguments and transform the image.
func (p Processor) ProcessArgumentsToTransformImage(args *TransformArguments) error {
	if args.UseServerURL {
		return p.useServerToTransformImage(args)
	}

	return p.useLocalPackageToTransformImage(args)
}

func (p Processor) useLocalPackageToTransformImage(args *TransformArguments) error {
	inputImageReader := bytes.NewBuffer(args.InputImageData)
	formulaReader := bytes.NewBuffer(args.FormulaData)
	outputSettingsReader := bytes.NewBuffer(args.OutputSettingsData)

	return p.localImageTransformer.ApplyFormulaToTransformImage(
		inputImageReader,
		formulaReader,
		outputSettingsReader,
		args.OutputImageData,
	)
}

func (p Processor) useServerToTransformImage(args *TransformArguments) error {
	inputData := image_transform_server.DataStreams{
		InputImage:     args.InputImageData,
		FormulaData:    args.FormulaData,
		OutputSettings: args.OutputSettingsData,
	}
	imageInfo, err := p.remoteImageTransformer.Transform(context.Background(), &inputData)
	if err != nil {
		fmt.Printf("Error from the server: %v", err)
		os.Exit(1)
	}

	args.OutputImageData.Write(imageInfo.GetImageData())
	return nil
}

// GetRemoteTransformer is a getter
func (p Processor) GetRemoteTransformer() image_transform_server.ImageTransformer {
	return p.remoteImageTransformer
}

// GetLocalTransformer is a getter
func (p Processor) GetLocalTransformer() creatingsymmetry.TransformerStrategy {
	return p.localImageTransformer
}

// NewCommandProcessor returns a new CommandProcessor with the given arguments.
func NewCommandProcessor(newRemoteImageTransformer image_transform_server.ImageTransformer, newLocalImageTransformer creatingsymmetry.TransformerStrategy) *Processor {
	if newRemoteImageTransformer == nil {
		newRemoteImageTransformer = image_transform_server.NewImageTransformerProtobufClient("http://localhost:8080", &http.Client{})
	}

	if newLocalImageTransformer == nil {
		newLocalImageTransformer = &creatingsymmetry.FileTransformer{}
	}

	return &Processor{
		localImageTransformer:  newLocalImageTransformer,
		remoteImageTransformer: newRemoteImageTransformer,
	}
}

// TransformArguments holds the arguments scanned from the command line.
type TransformArguments struct {
	InputImageData     []byte
	FormulaData        []byte
	OutputSettingsData []byte
	OutputImageData    *bytes.Buffer
	ServerURL          string
	UseServerURL       bool
}
