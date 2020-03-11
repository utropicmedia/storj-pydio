// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package pydio

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/minio/minio-go"
	"github.com/pydio/cells-client/rest"
	cells_sdk "github.com/pydio/cells-sdk-go"
)

// DEBUG allows more detailed working to be exposed through the terminal.
var DEBUG = false

// ConfigPydio defines the variables and types.
type ConfigPydio struct {
	PydioUrl        string `json:"pydioUrl"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	EndPoint        string `json:"endpoint"`
	AccessKeyID     string `json:"accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey"`
	BucketName      string `json:"bucketname"`
}

// PydioReader implements an io.Reader interface
type PydioReader struct {
	List       <-chan minio.ObjectInfo
	Client     *minio.Client
	BucketName string
}

// LoadPydioProperty reads and parses the JSON file.
// that contain a Pydio instance's property.
// and returns all the properties as an object.
func LoadPydioProperty(fullFileName string) (ConfigPydio, error) { // fullFileName for fetching Pydio credentials from  given JSON filename.
	var configPydio ConfigPydio
	// Open and read the file
	fileHandle, err := os.Open(fullFileName)
	if err != nil {
		return configPydio, err
	}
	defer fileHandle.Close()

	jsonParser := json.NewDecoder(fileHandle)
	jsonParser.Decode(&configPydio)

	// Display read information.
	fmt.Println("Read Pydio configuration from the ", fullFileName, " file")
	fmt.Println("PydioUrl\t: ", configPydio.PydioUrl)
	fmt.Println("Username\t: ", configPydio.Username)
	fmt.Println("Password\t: ", configPydio.Password)
	fmt.Println("End Point\t: ", configPydio.EndPoint)
	fmt.Println("Bucket Name\t: ", configPydio.BucketName)
	return configPydio, nil
}

// ConnectToPydio will connect to a Pydio instance,
// based on the read property from an external file.
// It returns a reference to an io.Reader with Pydio instance information
func ConnectToPydio(fullFileName string) (*PydioReader, error) { // fullFileName for fetching Pydio credentials from given JSON filename.
	// Read Pydio instance's properties from an external file.
	configPydio, err := LoadPydioProperty(fullFileName)
	//
	if err != nil {
		log.Printf("LoadPydioProperty: %s\n", err)
		return nil, err
	}

	fmt.Println("\nConnecting to Pydio...")

	// Authentication of Pydio user account
	err = newPydioClient(configPydio.PydioUrl, configPydio.Username, configPydio.Password)
	if err != nil {
		return nil, err
	}

	// Inform about successful connection.
	fmt.Println("Successfully connected to Pydio!")

	splitString := strings.Split(configPydio.PydioUrl, ":")
	var secure bool
	if splitString[0] == "https" {
		secure = true
	} else {
		secure = false
	}
	// Initialize minio client object.
	minioClient, err := minio.New(configPydio.EndPoint, configPydio.AccessKeyID, configPydio.SecretAccessKey, secure)
	if err != nil {
		log.Fatal(err)
	}

	// Return Pydio connection client, bucket name.
	return &PydioReader{Client: minioClient, BucketName: configPydio.BucketName}, nil
}

// Create function for authenticate pydio user account
func newPydioClient(urlPydio string, user string, password string) error {
	conf := &cells_sdk.SdkConfig{}
	conf.Url = urlPydio
	conf.User = user
	conf.Password = password

	// Insure values are legal
	conf.Url = strings.Trim(conf.Url, " ")
	if len(conf.Url) == 0 {
		return fmt.Errorf("Field cannot be empty!")
	}

	// Parse the url of pydio
	parseUrl, err := url.Parse(conf.Url)
	if err != nil || parseUrl == nil || parseUrl.Scheme == "" || parseUrl.Host == "" {
		return fmt.Errorf("Provide a valid URL")
	}

	if err != nil {
		return fmt.Errorf("URL %s is not valid: %s", conf.Url, err.Error())
	}

	// Test a simple PING with this config before saving!
	rest.DefaultConfig = conf

	_, _, err = rest.GetApiClient()
	if err != nil {
		return err
	}
	return nil
}
