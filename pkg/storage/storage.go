package storage

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

type Profile struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type Storage struct {
	DefaultProfile string     `json:"default_profile"`
	Profiles       []*Profile `json:"profiles"`
	StoragePath    string     `json:"-"`
}

func (s *Storage) Load() error {
	if _, err := os.Stat(s.StoragePath); os.IsNotExist(err) {
		return err
	}

	fp, err := os.Open(s.StoragePath)
	if err != nil {
		return errors.Errorf("Error opening file: %s", err)
	}
	defer fp.Close()

	decoder := json.NewDecoder(fp)

	err = decoder.Decode(s)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Commit() error {

	content, err := json.Marshal(s)
	if err != nil {
		return err
	}

	err = os.WriteFile(s.StoragePath, content, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) AddProfile(p *Profile) {
	s.Profiles = append(s.Profiles, p)
}

func (s *Storage) ExtractProfile(name string) (*Profile, error) {
	for _, p := range s.Profiles {
		if p.Name == name {
			return p, nil
		}
	}
	return nil, errors.Errorf("not found profile with name %s", name)
}

func (s *Storage) ExtractDefaultProfile() (*Profile, error) {
	return s.ExtractProfile(s.DefaultProfile)
}
