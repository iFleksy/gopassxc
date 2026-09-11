package helper

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/iFleksy/gopassxc/pkg/client/impl"
	"github.com/iFleksy/gopassxc/pkg/utils"
)

type Config struct {
	AssociatedName    string `json:"associated_name"`
	IdentificationKey string `json:"identification_Key"`
}

var ConfigNotFoundError = errors.New("Config file not found. Please run 'gopassxc init' to create a new configuration")

func BuildConfig() (Config, error) {
	ConfigPath := utils.GetConfigPath()
	var c Config
	f, err := os.Open(ConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return c, ConfigNotFoundError
		}
		return c, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(data, &c)
	return c, err
}

func BuildClient(config Config) (*impl.Client, error) {
	var client *impl.Client
	client, err := impl.New(impl.ClientOPTS{
		SocketPath:        utils.GetSocketPath(),
		AssociatedName:    &config.AssociatedName,
		IdentificationKey: &config.IdentificationKey,
		TriggerUnlock:     impl.ValueToPointer(true), // Set to nil to use default value
	})
	return client, err
}

func BuildClientWithConfig() (*impl.Client, error) {
	return impl.New(impl.ClientOPTS{
		SocketPath: utils.GetSocketPath(),
	})
}

func BuildWithConfig() (*impl.Client, error) {
	config, err := BuildConfig()
	if err != nil {
		return nil, err
	}
	client, err := BuildClient(config)
	return client, err
}

func PrepareClient(ctx context.Context, client *impl.Client) (*impl.Client, error) {
	err := client.Connect()
	if err != nil {
		return client, err
	}
	err = client.ChangePublicKeys(ctx)
	return client, err
}
