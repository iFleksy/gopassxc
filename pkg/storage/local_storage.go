package storage

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type LocalStorage struct {
	DefaultProfile string     `json:"default_profile"`
	Profiles       []*Profile `json:"profiles"`
	StoragePath    string     `json:"-"`
}

func NewLocalStorage(storagePath string) *LocalStorage {
	storage := &LocalStorage{
		StoragePath: storagePath,
	}
	storage.Load()
	return storage
}

func (s *LocalStorage) Load() error {
	if _, err := os.Stat(s.StoragePath); os.IsNotExist(err) {
		return nil
	}

	fp, err := os.Open(s.StoragePath)

	if err != nil {
		return errors.Errorf("Error opening file")
	}

	defer fp.Close()

	decoder := json.NewDecoder(fp)

	err = decoder.Decode(s)
	if err != nil {
		return err
	}
	return nil
}

func (s *LocalStorage) Commit() error {
	log.Info("Try to update storage file")
	content, err := json.Marshal(s)

	if err != nil {
		return err
	}

	log.Info("Writer content to storage file" + s.StoragePath)

	err = os.WriteFile(s.StoragePath, content, 0600)
	if err != nil {
		return err
	}
	log.Info("Storage file updated")
	return nil
}

func (s *LocalStorage) AddProfile(p *Profile) {
	s.Profiles = append(s.Profiles, p)
}

func (s *LocalStorage) ExtractProfile(name string) (*Profile, error) {
	for _, p := range s.Profiles {
		if p.Name == name {
			return p, nil
		}
	}
	return nil, errors.Errorf("not found profile with name %s", name)
}

func (s *LocalStorage) ExtractDefaultProfile() (*Profile, error) {
	return s.ExtractProfile(s.DefaultProfile)
}
