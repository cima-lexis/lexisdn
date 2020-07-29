package config

import (
	"fmt"
	"os"
)

// Config contains env configuration
var Config = struct {
	User     string
	Password string
	ClientID string
	AuthURL  string
	URL      string
}{}

func usage() {
	fmt.Print(`
This commands require following environment variable to be set:
  WEBDROPS_USER			-	webdrops user
  WEBDROPS_PWD			-	webdrops password
  WEBDROPS_CLIENT_ID	-	webdrops client id
  WEBDROPS_AUTH_URL		-	URL for KeyCloak authentication
  WEBDROPS_URL			-	base URL for all webdrops endpoints
`)

	os.Exit(1)

}

func check(name string) {
	if os.Getenv(name) == "" {
		fmt.Println("ERROR: missing env variable " + name)
		usage()
	}
}

func Init() {
	check("WEBDROPS_USER")
	check("WEBDROPS_PWD")
	check("WEBDROPS_CLIENT_ID")
	check("WEBDROPS_AUTH_URL")
	check("WEBDROPS_URL")

	Config.User = os.Getenv("WEBDROPS_USER")
	Config.Password = os.Getenv("WEBDROPS_PWD")
	Config.ClientID = os.Getenv("WEBDROPS_CLIENT_ID")
	Config.AuthURL = os.Getenv("WEBDROPS_AUTH_URL")
	Config.URL = os.Getenv("WEBDROPS_URL")

}
