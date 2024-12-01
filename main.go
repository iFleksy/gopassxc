package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/iFleksy/gopassxc/pkg/client"
	"github.com/iFleksy/gopassxc/pkg/helpers"
	"github.com/iFleksy/gopassxc/pkg/storage"
	log "github.com/sirupsen/logrus"
)

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func main() {
	// log.SetOutput(io.Discard)

	var url string
	var debug bool

	flag.StringVar(&url, "u", "", "URL for search")
	flag.BoolVar(&debug, "d", false, "debug")
	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	configDir, err := helpers.GetStoragePath()
	if err != nil {
		panic(err)

	}
	store := storage.Storage{StoragePath: configDir}
	store.Load()
	profile, err := store.ExtractDefaultProfile()

	socketPath := helpers.GetSocketPath()

	var c client.Client
	isInit := false
	if err != nil {
		c = client.NewClient(socketPath, "", "")
	} else {
		isInit = true
		c = client.NewClient(socketPath, profile.Name, profile.Key)
	}

	err = c.Connect()
	if err != nil {
		panic(err)
	}

	defer c.Disconnect()

	err = c.ChangePublicKeys()
	if err != nil {
		panic(err)
	}

	if !isInit {
		_, err = c.Associate()
		if err != nil {
			panic(err)
		}
		name, key := c.GetAssociatedProfile()
		p := &storage.Profile{Name: name, Key: key}
		store.AddProfile(p)
		store.DefaultProfile = name
		err = store.Commit()
		if err != nil {
			panic(err)
		}
	}

	err = c.TestAssociate()
	if err != nil {
		panic(err)
	}

	e, err := c.GetLogins(url)
	if err != nil {
		panic(err)
	}

	jData, err := JSONMarshal(e)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", string(jData))
}
