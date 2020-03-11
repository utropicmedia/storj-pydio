# storj-pydio
### Developed using libuplink version : v0.34.6

To build from scratch, [install Go](https://golang.org/doc/install#install).

```
$ go get -u github.com/minio/minio-go
$ go get -u github.com/pydio/cells-client/rest
$ go get -u github.com/pydio/cells-sdk-go
$ go get -u github.com/urfave/cli
$ go get -u storj.io/storj/lib/uplink
$ go get -u ./...
```

## Set-up Files
* Create a `pydio_property.json` file, with following contents about a Pydio instance:
    * pydioUrl :- URL of pydio on the server.
    * username :- User Name of Pydio
    * password :- Password of Pydio
    * endpoint :- Amazon AWS S3 End point
    * accessKeyID :- Access Key ID created in Amazon AWS S3
    * secretAccessKey :- Secret Access Key created in Amazon AWS S3
    * bucketname :- Amazon AWS S3 Bucket Name


```json
    { 
        "pydioUrl": "IPAddressPortToPydioServer",
        "username": "pydioUserName",
        "password": "pydioPassword",
        "endpoint": "s3.amazonaws.com",
        "accessKeyID": "amazonawsS3AccessKey",
        "secretAccessKey":"amazonawsS3SecretAccessKey",
        "bucketname":"amazonawsS3BucketName"
    }
```

* Create a `storj_config.json` file, with Storj network's configuration information in JSON format:
    * apiKey :- API key created in Storj satellite gui
    * satelliteURL :- Storj Satellite URL
    * encryptionPassphrase :- Storj Encryption Passphrase.
    * bucketName :- Split file into given size before uploading.
    * uploadPath :- Path on Storj Bucket to store data (optional) or "/"
    * serializedScope:- Serialized Scope Key shared while uploading data used to access bucket without API key
    * disallowReads:- Set true to create serialized scope key with restricted read access
    * disallowWrites:- Set true to create serialized scope key with restricted write access
    * disallowDeletes:- Set true to create serialized scope key with restricted delete access

```json
    { 
        "apikey":     "change-me-to-the-api-key-created-in-satellite-gui",
        "satelliteURL":  "us-central-1.tardigrade.io:7777",
        "bucketName":     "change-me-to-desired-bucket-name",
        "uploadPath": "optionalpath/requiredfilename",
        "encryptionpassphrase": "you'll never guess this",
        "serializedScope": "change-me-to-the-api-key-created-in-encryption-access-apiKey",
        "disallowReads": "true/false-to-disallow-reads",
        "disallowWrites": "true/false-to-disallow-writes",
        "disallowDeletes": "true/false-to-disallow-deletes"
    }
```

* Store both these files in a `config` folder.  Filename command-line arguments are optional.  defualt locations are used.

## Run the command-line tool

**NOTE**: The following commands operate

* Get help
```
$ storj-pydio -h
```

* Check version
```
$ storj-pydio -v
```

* Read files' data from desired Pydio instance and upload it to given Storj network bucket using Serialized Scope Key.  [note: filename arguments are optional.  default locations are used.]
```
$ storj-pydio store ./config/pydio_property.json ./config/storj_config.json  
```

* Read files' data from desired Pydio instance and upload it to given Storj network bucket API key and EncryptionPassPhrase from storj_config.json and creates an unrestricted shareable Serialized Scope Key.  [note: filename arguments are optional. default locations are used.]
```
$ storj-pydio store ./config/pydio_property.json ./config/storj_config.json key
```

* Read files' data from desired Pydio instance and upload it to given Storj network bucket API key and EncryptionPassPhrase from storj_config.json and creates a restricted shareable Serialized Scope Key.  [note: filename arguments are optional. default locations are used. `restrict` can only be used with `key`]
```
$ storj-pydio store ./config/pydio_property.json ./config/storj_config.json key restrict
```

* Read files' data in `debug` mode from desired Pydio instance and upload it to given Storj network bucket.  [note: filename arguments are optional.  default locations are used. Make sure `debug` folder already exist in project folder.]
```
$ storj-pydio store debug ./config/pydio_property.json ./config/storj_config.json  
```

* Read Pydio instance property from a desired JSON file and display all its files' data
```
$ storj-pydio parse   
```

* Read Pydio instance property in `debug` mode from a desired JSON file and display all its files' data
```
$ storj-pydio parse debug 
```

* Read and parse Storj network's configuration, in JSON format, from a desired file and upload a sample object
```
$ storj-pydio test 
```
* Read and parse Storj network's configuration, in JSON format, from a desired file and upload a sample object in `debug` mode
```
$ storj-pydio test debug 
```