package main

import (
	//Standard Packages

	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	"log"
	"os"

	"utropicmedia/pydio_storj_interface/pydio"
	"utropicmedia/pydio_storj_interface/storj"

	"github.com/minio/minio-go"
	"github.com/urfave/cli"
)

var gbDEBUG = false

const pydioConfigFile = "./config/pydio_property.json"
const storjConfigFile = "./config/storj_config.json"

// Create command-line tool to read from CLI.
var app = cli.NewApp()

// SetAppInfo sets information about the command-line application.
func setAppInfo() {
	app.Name = "Storj Pydio Connector"
	app.Usage = "Backup your File from Pydio Cells to the decentralized Storj network"
	app.Authors = []*cli.Author{{Name: "Satyam Shivam - Utropicmedia", Email: "development@utropicmedia.com"}}
	app.Version = "1.0.1"
}

// helper function to flag debug
func setDebug(debugVal bool) {
	gbDEBUG = true
	storj.DEBUG = debugVal

}

// setCommands sets various command-line options for the app.
func setCommands() {

	app.Commands = []*cli.Command{
		{
			Name:    "parse",
			Aliases: []string{"p"},
			Usage:   "Command to read and parse JSON information about Pydio instance properties and then generate a list of all files with their corresponding paths ",
			//\narguments-\n\t  fileName [optional] = provide full file name (with complete path), storing Pydio properties if this fileName is not given, then data is read from ./config/pydio_property.json\n\t
			// example = ./storj-pydio p ./config/pydio_property.json\n",
			Action: func(cliContext *cli.Context) error {
				var fullFileName = pydioConfigFile

				// process arguments
				if len(cliContext.Args().Slice()) > 0 {
					for i := 0; i < len(cliContext.Args().Slice()); i++ {
						// Incase, debug is provided as argument.
						if cliContext.Args().Slice()[i] == "debug" {
							setDebug(true)
						} else {
							fullFileName = cliContext.Args().Slice()[i]
						}
					}
				}

				// Establish connection with Pydio and get io.Reader implementor.
				pydioReader, err := pydio.ConnectToPydio(fullFileName)
				if err != nil {
					log.Fatalf("Failed to establish connection with Pydio: %s\n", err)
				}

				// Create a done channel to control 'ListObjects' go routine.
				doneCh := make(chan struct{})

				// Indicate to our routine to exit cleanly upon return.
				defer close(doneCh)

				isRecursive := true

				fmt.Println("\nReading All files from the Pydio Cells...")
				// ListObjects lists all objects from the specified bucket.
				objectCh := pydioReader.Client.ListObjects(pydioReader.BucketName, "", isRecursive, doneCh)
				for object := range objectCh {
					if object.Err != nil {
						log.Fatal("Object Information Error: ", object.Err)
					}

					if object.Key == ".pydio" {
						continue
					}
					fmt.Println(object.Key)
				}

				fmt.Println("Reading ALL files from the Pydio Cells Bucket...Complete!")
				return err
			},
		},

		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "Command to read and parse JSON information about Storj network and upload sample data",
			//\n arguments- 1. fileName [optional] = provide full file name (with complete path), storing Storj configuration information if this fileName is not given, then data is read from ./config/storj_config.json
			//example = ./storj-pydio t ./config/storj_config.json\n\n\n",
			Action: func(cliContext *cli.Context) error {

				// Default Storj configuration file name.
				var fullFileName = storjConfigFile
				var foundFirstFileName = false
				var foundSecondFileName = false
				var keyValue string
				var restrict string
				gbDEBUG = false
				// process arguments
				if len(cliContext.Args().Slice()) > 0 {
					for i := 0; i < len(cliContext.Args().Slice()); i++ {

						// Incase, debug is provided as argument.
						if cliContext.Args().Slice()[i] == "debug" {
							setDebug(true)
						} else {
							if !foundFirstFileName {
								fullFileName = cliContext.Args().Slice()[i]
								foundFirstFileName = true
							} else {
								if !foundSecondFileName {

									keyValue = cliContext.Args().Slice()[i]
									foundSecondFileName = true
								} else {
									restrict = cliContext.Args().Slice()[i]
								}
							}
						}
					}
				}
				// Sample data to be uploaded with sample data name
				fileName := "testdata"
				testData := "test"
				data := []byte(testData)
				if gbDEBUG {
					t := time.Now()
					time := t.Format("2006-01-02")
					fileName = "uploaddata_" + time + ".txt"
					err := ioutil.WriteFile(fileName, data, 0644)
					if err != nil {
						fmt.Println("Error while writting to file: ", err)
					}
				}
				reader := bytes.NewBuffer(data)
				var fileNamesDEBUG []string

				// Connect to storj network.
				ctx, uplink, project, bucket, storjConfig, _, errr := storj.ConnectStorjReadUploadData(fullFileName, keyValue, restrict)

				// Upload sample data on storj network.
				fileNamesDEBUG = storj.ConnectUpload(ctx, bucket, reader, fileName, fileNamesDEBUG, storjConfig, errr)

				if errr != nil {
					return errr
				}

				// Close storj project.
				storj.CloseProject(uplink, project, bucket)
				//
				fmt.Println("\nUpload \"testdata\" on Storj: Successful!")
				return errr
			},
		},
		{
			Name:    "store",
			Aliases: []string{"s"},
			Usage:   "Command to connect and transfer file(s)/folder(s) from a desired Pydio Cells account to given Storj Bucket.",
			//\n    arguments-\n      1. fileName [optional] = provide full file name (with complete path),
			// storing pydio properties in JSON format\n   if this fileName is not given,
			// then data is read from ./config/pydio_property.json\n
			// 2. fileName [optional] = provide full file name (with complete path), storing Storj
			// configuration in JSON format\n     if this fileName is not given, then
			// data is read from ./config/storj_config.json\n
			// example = ./storj-pydio.exe s ./config/pydio_property.json ./config/storj_config.json fileName/DirectoryName\n"
			Action: func(cliContext *cli.Context) error {

				// Default configuration file names.
				var fullFileNameStorj = storjConfigFile

				var fullFileNamePydio = pydioConfigFile

				var keyValue string
				var restrict string
				var fileNamesDEBUG []string
				// process arguments - Reading fileName from the command line.
				// process arguments - Reading fileName from the command line.
				var foundFirstFileName = false
				var foundSecondFileName = false
				var foundThirdFileName = false
				if len(cliContext.Args().Slice()) > 0 {
					for i := 0; i < len(cliContext.Args().Slice()); i++ {
						// Incase debug is provided as argument.
						if cliContext.Args().Slice()[i] == "debug" {
							setDebug(true)
						} else {
							if !foundFirstFileName {
								fullFileNamePydio = cliContext.Args().Slice()[i]
								foundFirstFileName = true
							} else {
								if !foundSecondFileName {
									fullFileNameStorj = cliContext.Args().Slice()[i]
									foundSecondFileName = true
								} else {
									if !foundThirdFileName {
										keyValue = cliContext.Args().Slice()[i]
										foundThirdFileName = true
									} else {
										restrict = cliContext.Args().Slice()[i]
									}
								}
							}
						}
					}
				}
				// Establish connection with Pydio and get io.Reader implementor.
				pydioReader, err := pydio.ConnectToPydio(fullFileNamePydio)
				if err != nil {
					log.Fatalf("Failed to establish connection with Pydio: %s\n", err)
				}

				// Create a done channel to control 'ListObjects' go routine.
				doneCh := make(chan struct{})

				// Indicate to our routine to exit cleanly upon return.
				defer close(doneCh)

				isRecursive := true

				// ListObjects lists all objects from the specified bucket.
				objectCh := pydioReader.Client.ListObjects(pydioReader.BucketName, "", isRecursive, doneCh)

				// Connect to storj network and it returns context, uplink, project, bucket and storj configration.
				ctx, uplink, project, bucket, storjConfig, scope, errr := storj.ConnectStorjReadUploadData(fullFileNameStorj, keyValue, restrict)
				if errr != nil {
					log.Fatal(err)
				}

				storePath := make([]string, 0)

				t := time.Now()
				timeNow := t.Format("2006-01-02_15_04_05")
				for object := range objectCh {
					if object.Err != nil {
						log.Fatal("Object Information Error", object.Err)
					}

					if object.Key == ".pydio" {
						continue
					}

					fmt.Println("\nReading content from the file:", object.Key)
					// GetObject function returns an seekable, readable object.
					objectReader, err := pydioReader.Client.GetObject(pydioReader.BucketName, object.Key, minio.GetObjectOptions{})
					if err != nil {
						log.Fatal(err)
					}

					pydioPath := pydioReader.BucketName + "_" + timeNow + "/" + object.Key

					storePath = append(storePath, pydioPath)
					// Upload Pydio object on storj Network with file name.
					storj.ConnectUpload(ctx, bucket, objectReader, pydioPath, fileNamesDEBUG, storjConfig, errr)
					if errr != nil {
						log.Fatal(errr)
					}

				}
				// Debug the storj data.
				storj.Debug(ctx, bucket, storjConfig.UploadPath, storePath)

				// Close the storj project.
				storj.CloseProject(uplink, project, bucket)
				fmt.Println(" ")
				if keyValue == "key" {
					if restrict == "restrict" {
						fmt.Println("Restricted Serialized Scope Key: ", scope)
						fmt.Println(" ")
					} else {
						fmt.Println("Serialized Scope Key: ", scope)
						fmt.Println(" ")
					}
				}

				return err
			},
		},
	}
}

func main() {

	// Show application information on screen
	setAppInfo()
	// Get command entered by user on cli
	setCommands()
	// Get detailed information for debugging
	setDebug(false)

	err := app.Run(os.Args)

	if err != nil {
		log.Fatalf("app.Run: %s", err)
	}
}
