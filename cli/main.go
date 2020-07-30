package main

import (
	"fmt"
	"os"
	"time"

	"github.com/cima-lexis/lexisdn/config"
	"github.com/cima-lexis/lexisdn/fetcher"
	"github.com/cima-lexis/lexisdn/webdrops"
)

func main() {
	config.Init()

	sess := webdrops.Session{}
	err := sess.Login()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("Error during login: %w", err))
		os.Exit(1)
	}

	err = fetcher.RisicoSensorsMaps(sess, time.Date(2020, 7, 28, 9, 0, 0, 0, time.Local))
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("Error fetching wunderground observations for WRFDA: %w", err))
		os.Exit(1)
	}
}
