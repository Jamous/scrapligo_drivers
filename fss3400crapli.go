package drivers

import (
	"bufio"
	"fmt"
	"nettools/pkg/logging"
	"strings"

	"github.com/scrapli/scrapligo/driver/network"
	"github.com/scrapli/scrapligo/driver/options"
	"github.com/scrapli/scrapligo/platform"
)

// Functions
func Fss3400SetupScrapli(host NodeDetails) (*network.Driver, string, error) {
	//Log attempted connection
	hostm := fmt.Errorf("drivers/Fss3400SetupScrapli attempting to connect to host %s", host.Address)
	logging.Debug(hostm)

	//Scrapli code
	p, err := platform.NewPlatform(
		//Custom platform
		"fs_s3400",
		host.Address,
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername(host.Username),
		options.WithAuthPassword(host.Password),
		options.WithTransportType("standard"), //Uses crypto/ssh https://github.com/scrapli/scrapligo/blob/main/driver/options/generic.go#L14
	)

	if err != nil {
		errm := fmt.Errorf("drivers/Fss3400SetupScrapli failed to create platform on host %s; error: %+v\n\n", host.Address, err)
		return nil, "", errm
	}

	d, err := p.GetNetworkDriver()
	if err != nil {
		errm := fmt.Errorf("Fss3400SetupScrapli failed to fetch network driver from the platfor mon host %s; error: %+v\n\n", host.Address, err)
		return nil, "", errm
	}

	err = d.Open()
	if err != nil {
		errm := fmt.Errorf("Fss3400SetupScrapli failed to open driver on host %s, error: %+v\n\n", host.Address, err)
		return nil, "", errm
	} else {
		hostm := fmt.Errorf("Fss3400SetupScrapli connected to host %s", host.Address)
		logging.Debug(hostm)
	}

	//Get hostname
	hnameR, err := d.SendCommand("show run | i hostname")
	if err != nil {
		errm := fmt.Errorf("Fss3400SetupScrapli failed to send command on host %s; error: %+v\n\n", host.Address, err)
		logging.Warn(errm)
	} else {
		cmdm := fmt.Errorf("Host %s, Input: %s, Result%s", host.Address, string(hnameR.Input), string(hnameR.Result))
		logging.Debug(cmdm)
	}

	//Get clean hostname - parse through results line by line
	hostname := ""
	scanner := bufio.NewScanner(strings.NewReader(hnameR.Result))
	for scanner.Scan() {
		line := scanner.Text()
		//Ignore values smaller than 10
		if len(line) >= 10 {
			//If the string starts with "hostname"
			if line[:8] == "hostname" {
				hostname = line[9:]
			}
		}
	}

	logm := fmt.Errorf("Connected to %s", hostname)
	logging.Info(logm)
	fmt.Println(logm)

	return d, hostname, nil
}
