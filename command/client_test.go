package command_test

import (
	"bytes"
	"errors"
	creatingsymmetry "github.com/Chadius/creating-symmetry"
	"github.com/chadius/image-transform-server/creatingsymmetryfakes"
	"github.com/chadius/image-transform-server/rpc/transform/github.com/chadius/image_transform_server"
	"github.com/cserrant/image-transform-cli/command"
	"github.com/cserrant/image-transform-cli/imagetransformserverfakes"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestClientCallsServerSuite(t *testing.T) {
	suite.Run(t, new(ClientCallsServerSuite))
}

type ClientCallsServerSuite struct {
	suite.Suite
	inputImageData     []byte
	formulaData        []byte
	outputSettingsData []byte
	expectedImage      []byte

	fakeRemoteImageTransformerClient *imagetransformserverfakes.FakeImageTransformer
	fakeLocalImageTransformerClient  *creatingsymmetryfakes.FakeTransformerStrategy
	fakeTransformedImage             *image_transform_server.Image
}

func (suite *ClientCallsServerSuite) SetupTest() {
	suite.inputImageData = []byte(`images go here`)
	suite.formulaData = []byte(`formula goes here`)
	suite.outputSettingsData = []byte(`outputSettings go here`)

	suite.expectedImage = []byte(`transformed image`)
	suite.fakeTransformedImage = &image_transform_server.Image{
		ImageData: suite.expectedImage,
	}
	suite.fakeRemoteImageTransformerClient = &imagetransformserverfakes.FakeImageTransformer{}
	suite.fakeRemoteImageTransformerClient.TransformReturns(suite.fakeTransformedImage, nil)

	suite.fakeLocalImageTransformerClient = &creatingsymmetryfakes.FakeTransformerStrategy{}
	suite.fakeLocalImageTransformerClient.ApplyFormulaToTransformImageStub = func(inputImageDataByteStream, formulaDataByteStream, outputSettingsDataByteStream io.Reader, output io.Writer) error {
		output.Write(suite.expectedImage)
		return nil
	}
}

func (suite *ClientCallsServerSuite) TestWhenURLIsProvided_ThenCallsServer() {
	// Setup
	commandProcessor := command.NewCommandProcessor(suite.fakeRemoteImageTransformerClient, suite.fakeLocalImageTransformerClient)
	var outputImageData bytes.Buffer

	// Act
	commandProcessor.ProcessArgumentsToTransformImage(&command.TransformArguments{
		InputImageData:     suite.inputImageData,
		FormulaData:        suite.formulaData,
		OutputSettingsData: suite.outputSettingsData,
		OutputImageData:    &outputImageData,
		ServerURL:          "localhost://8000",
		UseServerURL:       true,
	})

	// Require
	require := require.New(suite.T())
	require.Equal(1, suite.fakeRemoteImageTransformerClient.TransformCallCount(), "Client was not called")

	_, dataStreamsUsed := suite.fakeRemoteImageTransformerClient.TransformArgsForCall(0)
	require.Equal(suite.inputImageData, dataStreamsUsed.InputImage)
	require.Equal(suite.formulaData, dataStreamsUsed.FormulaData)
	require.Equal(suite.outputSettingsData, dataStreamsUsed.OutputSettings)
	require.Equal(suite.fakeTransformedImage.GetImageData(), outputImageData.Bytes(), "output image received from mock object is different")
}

func (suite *ClientCallsServerSuite) TestWhenDoNotUseServer_ThenLocalTransformerIsUsed() {
	// Setup
	commandProcessor := command.NewCommandProcessor(suite.fakeRemoteImageTransformerClient, suite.fakeLocalImageTransformerClient)
	var outputImageData bytes.Buffer

	// Act
	commandProcessor.ProcessArgumentsToTransformImage(&command.TransformArguments{
		InputImageData:     suite.inputImageData,
		FormulaData:        suite.formulaData,
		OutputSettingsData: suite.outputSettingsData,
		OutputImageData:    &outputImageData,
		ServerURL:          "",
		UseServerURL:       false,
	})

	// Require
	require := require.New(suite.T())
	require.Equal(1, suite.fakeLocalImageTransformerClient.ApplyFormulaToTransformImageCallCount(), "Client was not called")
	require.Equal(suite.fakeTransformedImage.GetImageData(), outputImageData.Bytes(), "output image received from mock object is different")
}

type InjectTransformerSuite struct {
	suite.Suite
}

func TestInjectTransformerSuite(t *testing.T) {
	suite.Run(t, new(InjectTransformerSuite))
}

func (suite *InjectTransformerSuite) TestWhenRemoteImageTransformerIsInjected_ThenUsesGivenObject() {
	// Setup
	fakeImageTransformerClient := &imagetransformserverfakes.FakeImageTransformer{}

	// Act
	commandProcessor := command.NewCommandProcessor(fakeImageTransformerClient, nil)

	// Assert
	require := require.New(suite.T())
	require.Equal(
		reflect.TypeOf(commandProcessor.GetRemoteTransformer()),
		reflect.TypeOf(fakeImageTransformerClient),
	)
}

func (suite *InjectTransformerSuite) TestWhenNoRemoteImageTransformerIsInjected_ThenUsesDefaultObject() {
	// Setup
	productionClient := image_transform_server.NewImageTransformerProtobufClient("http://localhost:8080", &http.Client{})

	// Act
	commandProcessor := command.NewCommandProcessor(nil, nil)

	// Assert
	require := require.New(suite.T())
	require.Equal(
		reflect.TypeOf(commandProcessor.GetRemoteTransformer()),
		reflect.TypeOf(productionClient),
	)
}

func (suite *InjectTransformerSuite) TestWhenNoLocalImageTransformerIsInjected_ThenUsesDefaultObject() {
	// Setup
	productionClient := &creatingsymmetry.FileTransformer{}

	// Act
	commandProcessor := command.NewCommandProcessor(nil, nil)

	// Assert
	require := require.New(suite.T())
	require.Equal(
		reflect.TypeOf(commandProcessor.GetLocalTransformer()),
		reflect.TypeOf(productionClient),
	)
}

func (suite *InjectTransformerSuite) TestWhenLocalImageTransformerIsInjected_ThenUsesGivenObject() {
	// Setup
	injectedTransformer := &creatingsymmetryfakes.FakeTransformerStrategy{}

	// Act
	commandProcessor := command.NewCommandProcessor(nil, injectedTransformer)

	// Assert
	require := require.New(suite.T())
	require.Equal(
		reflect.TypeOf(commandProcessor.GetLocalTransformer()),
		reflect.TypeOf(injectedTransformer),
	)
}

type LocalPackageErrorReportingSuite struct {
	suite.Suite
	fakeLocalImageTransformerClient *creatingsymmetryfakes.FakeTransformerStrategy
}

func TestLocalPackageErrorReportingSuite(t *testing.T) {
	suite.Run(t, new(LocalPackageErrorReportingSuite))
}

func (suite *LocalPackageErrorReportingSuite) TestWhenLocalPackageReturnsError_ThenReturnError() {
	// Setup
	suite.fakeLocalImageTransformerClient = &creatingsymmetryfakes.FakeTransformerStrategy{}
	suite.fakeLocalImageTransformerClient.ApplyFormulaToTransformImageStub = func(inputImageDataByteStream, formulaDataByteStream, outputSettingsDataByteStream io.Reader, output io.Writer) error {
		return errors.New("irrelevant error")
	}

	var dummyOutputImageData bytes.Buffer

	// Act
	commandProcessor := command.NewCommandProcessor(nil, suite.fakeLocalImageTransformerClient)
	err := commandProcessor.ProcessArgumentsToTransformImage(&command.TransformArguments{
		InputImageData:     nil,
		FormulaData:        nil,
		OutputSettingsData: nil,
		OutputImageData:    &dummyOutputImageData,
		ServerURL:          "",
		UseServerURL:       false,
	})

	// Require
	require := require.New(suite.T())
	require.Error(err, "No error was found")
	require.Equal("irrelevant error", err.Error(), "error message does not match")
}
