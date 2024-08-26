package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/scrapli/scrapligo/driver/options"
	"github.com/scrapli/scrapligo/platform"
)

func main() {

	//Setup logging
	var channelLog bytes.Buffer

	//Defered anonymous function to print full output. Comment out to ignore output. (full io)
	defer func() {
		b := make([]byte, channelLog.Len())
		_, _ = channelLog.Read(b)
		fmt.Printf("\n\n\n\nChannel log output:\n%s", b)
	}()

	//Setup scrapli platform
	p, err := platform.NewPlatform(
		"fs_s3400", //Custom platform
		"10.0.0.1", //Host address
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername("admin"),     //Login username
		options.WithAuthPassword("Adm1n"),     //Login password
		options.WithTransportType("standard"), //Uses crypto/ssh https://github.com/scrapli/scrapligo/blob/main/driver/options/generic.go#L14
		options.WithChannelLog(&channelLog),   //Return channel for logging
	)
	if err != nil {
		log.Fatal(err)
	}

	//Setup scrapli driver
	d, err := p.GetNetworkDriver()
	if err != nil {
		log.Fatal(err)
	}

	//Open scrapli connection, defer close
	err = d.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()

	//Send command
	commandR, err := d.SendCommand("show ver")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(commandR.Result)
}
