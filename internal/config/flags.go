package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

var (
	Port *string
	Dir  *string
)

func ParseFlags() {
	Port = flag.String("port", "8080", "Port number")
	Dir = flag.String("dir", "data", "Path to the directory")
	help := flag.Bool("help", false, "Show help")
	flag.Parse()

	if *help {
		fmt.Println(`Simple Storage Service.

**Usage:**
    triple-s [-port <N>] [-dir <S>]  
    triple-s --help

**Options:**
- --help     Show this screen.
- --port N   Port number
- --dir S    Path to the directory`)

		os.Exit(0)
	}

	if err := validateDir(); err != nil {
		log.Fatal(err)
	}

	if err := validatePort(); err != nil {
		log.Fatal(err)
	}
}

func validatePort() error {
	port, err := strconv.Atoi(*Port)
	if err != nil {
		return fmt.Errorf("port should be number")
	}

	if port < 1024 || port > 49151 {
		return fmt.Errorf("invalid port, must be between 1024 and 49151")
	}

	return nil
}

func validateDir() error {
	if *Dir == "internal" || *Dir == "pkg" {
		return fmt.Errorf("forbidden dir")
	}
	return nil
}
