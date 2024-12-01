package main

import (
	"encoding/json"
	"fmt"

	"github.com/iFleksy/gopassxc/pkg/client"
	"github.com/iFleksy/gopassxc/pkg/helpers"
	"github.com/iFleksy/gopassxc/pkg/storage"
)

func main() {
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

	assoName, assoKey := c.GetAssociatedProfile()
	fmt.Printf("assoName: %s, assoKey: %s\n", assoName, assoKey)

	e, err := c.GetLogins("https://test.local")
	if err != nil {
		panic(err)
	}

	jdata, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	fmt.Printf("jdata: %s\n", jdata)
}
