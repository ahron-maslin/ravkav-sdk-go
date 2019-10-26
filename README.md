# ravkav-sdk-go


<p align="center">
    <img width="150" alt="ravkav-sdk-go" src="https://github.com/ybaruchel/ravkav-sdk-go/blob/master/assets/logo.png">
</p>


[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/ybaruchel/ravkav-sdk-go/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/ybaruchel/ravkav-sdk-go)](https://goreportcard.com/report/github.com/ybaruchel/ravkav-sdk-go)
[![Go Doc](https://godoc.org/github.com/ybaruchel/ravkav-sdk-go?status.svg)](https://godoc.org/github.com/ybaruchel/ravkav-sdk-go)

## Overview
The most simple [ravkav card](https://en.wikipedia.org/wiki/Rav-Kav) reader implemented in golang

## Download
To download this package, run:
```
go get github.com/ybaruchel/ravkav-sdk-go...
```

## Example
```go
package main

import (
    "fmt"
    "github.com/ybaruchel/ravkav-sdk-go/card"
    "github.com/ybaruchel/ravkav-sdk-go/reader"
)

func main() {
    ravkavReader := reader.NewReader()
    availableReaders, err := ravkavReader.ListReaders()
    if err != nil {
    	panic("can't find available card readers")
    }
    err = ravkavReader.Connect(availableReaders[0]) // Connect to first available reader
    if err != nil {
    	panic("error when trying to connect to reader")
    }
    
    defer func() {
    	if ravkavReader.Disconnect() != nil {
    		panic("error when trying to disconnect from reader")
    	}
    }()
    
    c := card.NewByReader(ravkavReader) // Get new card instance
    err = c.Populate()                  // Populate the card instance with physical card records
    if err != nil {
    	panic(err)
    }
    err = c.Normalize() // Normalize card records to human readable format
    if err != nil {
    	panic(err)
    }
    jsonOutput, err := c.Output().JSON() // Get JSON string representation of the card
    if err != nil {
    	panic("error getting card output")
    }
    fmt.Println(jsonOutput)
}
```

## License
This work is published under the MIT license.

Please see the `LICENSE` file for details.
