package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/iFleksy/gopassxc/pkg/client/impl"
	"github.com/sirupsen/logrus"
)

const (
	SOCKET_PATH       = "/run/user/1000/org.keepassxc.KeePassXC.BrowserServer"
	AssociatedName    = "fsd"
	identificationKey = "fxDlZEjN7wI/Mncr6iJIdYccUwS0frUn"

	ConfigPath = "/tmp/gokeepass.json"

	// AssociatedName    = ""
	// identificationKey = ""
)

type ConfigData struct {
	AssociatedName    string `json:"associated_name"`
	IdentificationKey string `json:"identification_Key"`
}

func saveConfig(c ConfigData) error {
	_, err := os.Stat(ConfigPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var file *os.File
	if os.IsNotExist(err) {
		file, err = os.Create(ConfigPath)
		if err != nil {
			return err
		}
	}
	defer file.Close()

	as_json, _ := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	_, err = file.Write(as_json)
	return err
}

func loadConfig() (ConfigData, error) {
	var c ConfigData
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		return c, err
	}

	f, err := os.Open(ConfigPath)
	if err != nil {
		return c, err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return c, err
	}
	json.Unmarshal(data, &c)
	return c, err
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	c, err := loadConfig()
	isInited := false
	var client impl.Client
	if err != nil {
		fmt.Println("Config not found, creating new client")
		client = impl.New(impl.ClientOPTS{
			SocketPath: SOCKET_PATH,
		})
	} else {
		isInited = true
		client = impl.New(impl.ClientOPTS{
			SocketPath:        SOCKET_PATH,
			AssociatedName:    &c.AssociatedName,
			IdentificationKey: &c.IdentificationKey,
		})
	}

	err = client.Connect()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer client.Disconnect()

	pubKey, err := client.ChangePublicKeys()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Pub key: %s\n", pubKey)

	if !isInited {
		aid, akey, err := client.Associate()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		c := ConfigData{
			AssociatedName:    aid,
			IdentificationKey: akey,
		}
		saveConfig(c)
	} else {
		err = client.TestAssociate()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	resp, err := client.GetDBHash()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(resp)

	err = client.GetLogins("http://dot.net")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
