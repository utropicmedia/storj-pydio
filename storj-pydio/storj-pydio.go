package main

import (
	//Standard Packages

	"fmt"

	"io/ioutil"

	"log"
	"os"
	"bytes"

	"utropicmedia/Pydio_storj_interface/storj"


	"github.com/urfave/cli"

	"github.com/minio/minio-go/v6"
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
	app.Version = "1.0.0"
}

// helper function to flag debug
func setDebug(debugVal bool) {
	gbDEBUG = true
/*
	storj.DEBUG = debugVal
	*/
}

// setCommands sets various command-line options for the app.
func setCommands() {

	app.Commands = []*cli.Command{
		{
			Name:    "parse",
			Aliases: []string{"p"},
			Usage:   "Command to read and parse JSON information about Pydio instance properties and then generate a list of all files with their corresponidng paths ",
			//\narguments-\n\t  fileName [optional] = provide full file name (with complete path), storing Pydio properties if this fileName is not given, then data is read from ./config/pydio_property.json\n\t
			// example = ./storj-pydio p ./config/nextcloud_property.json\n",
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

				////////////////////
				fmt.Println(fullFileName)
				////////////////////

				// Establish connection with pydio and get a new client.
				///////////////////////////////////////////////////////////////////////
				endpoint := "s3.amazonaws.com"
				accessKeyID := "AKIAJLQEMXDN46GO52DA"
				secretAccessKey := "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"
				useSSL := false

				// Initialize minio client object.
				minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
				if err != nil {
					log.Fatalln(err)
				}
				///////////////////////////////////////////////////////////////////////


				// Generate a list of all files with their corresponidng paths
				///////////////////////////////////////////////////////////////////////
				// Create a done channel to control 'ListObjects' go routine.
				doneCh := make(chan struct{})

				// Indicate to our routine to exit cleanly upon return.
				defer close(doneCh)

				isRecursive := true

				objectCh := minioClient.ListObjects("qazx123", "", isRecursive, doneCh)
				for object1 := range objectCh {
					if object1.Err != nil {
						fmt.Println(object1.Err)
						/*
						return
						*/
					}
					fmt.Println(object1.Key)
				}
				///////////////////////////////////////////////////////////////////////
				return err
			},
		},

		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "Command to read and parse JSON information about Storj network and upload sample JSON data",
			//\n arguments- 1. fileName [optional] = provide full file name (with complete path), storing Storj configuration information if this fileName is not given, then data is read from ./config/storj_config.json
			//example = ./storj-nextcloud t ./config/storj_config.json\n\n\n",
			Action: func(cliContext *cli.Context) error {

				// Default Storj configuration file name.
				var fullFileName = storjConfigFile
				var foundFirstFileName = false
				var foundSecondFileName = false
				var keyValue string
				var restrict string

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
				// Sample test file and data to be uploaded
				sampleFileName := "testfile"
				testFile := []byte("Hello Storj")
				var fileName string
				if gbDEBUG {
					fileName = "./testFile_.txt"
					data := []byte(testFile)
					err := ioutil.WriteFile(fileName, data, 0644)
					if err != nil {
						fmt.Println("Error while writting to file ")
					}
				}
				fileName = sampleFileName + ".txt"

				//var fileNamesDEBUG []string
				// Connect to storj network.
				var test_buf=bytes.NewBuffer(testFile)
				_,err:= storj.ConnectStorjReadUploadData(fullFileName,test_buf,fileName, keyValue, restrict)
				//
				fmt.Println("\nUpload \"testdata\" on Storj: Successful!")
				return err

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
				fmt.Print(fullFileNameStorj)


				var fullFileNamePydio = pydioConfigFile

				var keyValue string
				var restrict string
				var filePath string
				var fileNamesDEBUG []string
				// process arguments - Reading fileName from the command line.
				var foundFirstFileName = false
				var foundSecondFileName = false
				var foundThirdFileName = false
				var foundFilePath = false
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
										if !foundFilePath {
											filePath = cliContext.Args().Slice()[i]
											/*
											foundFilePath = true
											*/
											foundFilePath = false
										} else {
											restrict = cliContext.Args().Slice()[i]
										}
									}
								}
							}
						}
					}
				}

				if !foundFilePath {
					/*
					log.Fatal("Please enter the file/root address for back-up. Terminating Application...")
					*/
				}

				/////////////////////////////
				fmt.Println(fullFileNamePydio, keyValue, restrict, filePath, fileNamesDEBUG)
				/////////////////////////////

				// Establish connection with Pydio and get a new Pydio Client
				///////////////////////////////////////////////////////////////////////
				endpoint := "s3.amazonaws.com"
				accessKeyID := "AKIAJLQEMXDN46GO52DA"
				secretAccessKey := "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"
				useSSL := false

				// Initialize minio client object.
				minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
				if err != nil {
					log.Fatalln(err)
				}
				///////////////////////////////////////////////////////////////////////

				// Create a list of all files with their respective paths for transfer to Storj


				// Generate a list of all files with their corresponidng paths
				///////////////////////////////////////////////////////////////////////
				// Create a done channel to control 'ListObjects' go routine.
				doneCh := make(chan struct{})

				// Indicate to our routine to exit cleanly upon return.
				defer close(doneCh)

				isRecursive := true

				objectCh := minioClient.ListObjects("qazx123", "", isRecursive, doneCh)
				for object1 := range objectCh {
					if object1.Err != nil {
						fmt.Println(object1.Err)
						/*
						return
						*/
					}
					fmt.Println(object1.Key)
					fmt.Println("---------------------")

					object, err := minioClient.GetObject("qazx123", object1.Key, minio.GetObjectOptions{})
					if err != nil {
						fmt.Println(err)
						/*
						return
						*/
					}

					// Transfer required data from pydio to storj

					////////////////////////////
					//scope ,err:= storj.ConnectStorjReadUploadData(fullFileName,test_buf,fileName, keyValue, restrict)
					if object1.Key != ".pydio" {
						scope , _ := storj.ConnectStorjReadUploadData(fullFileNameStorj, object, object1.Key, keyValue, restrict)
					////////////////////////////


					///////////////////////////////////////////////////////////////////////







				//Display Paths of all files to be uploaded

				// Create connection with the required storj nextwork
/*
				ctx, uplinkStorj, project, bucket, storjConfig, scope, err := storj.ConnectStorjReadUploadData(fullFileNameStorj, keyValue, restrict)
*/

				// Transfer required data from pydio to storj

				////////////////////////////
				//scope ,err:= storj.ConnectStorjReadUploadData(fullFileName,test_buf,fileName, keyValue, restrict)
				//scope ,err:= storj.ConnectStorjReadUploadData(fullFileNameStorj, object, object1.Key, keyValue, restrict)
				////////////////////////////

/*
				// close connection with storj nextwork
				storj.CloseProject(uplinkStorj, project, bucket)
*/

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
