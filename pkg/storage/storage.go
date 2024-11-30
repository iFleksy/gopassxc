package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

type Profile struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type Storage struct {
	Default     string     `json:"default"`
	Profiles    []*Profile `json:"profiles"`
	StoragePath string     `json:"-"`
}

func (s *Storage) checkPath() error {
	if _, err := os.Stat(s.StoragePath); os.IsNotExist(err) {
		return err
	}
	return nil
}

func (s *Storage) Load() error {
	if err := s.checkPath(); err != nil {
		return err
	}

	content, err := os.ReadFile(s.StoragePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, s)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Save() error {

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
	return nil, fmt.Errorf("not found profile with name %s", name)
}

func (s *Storage) ExtractDefaultProfile() (*Profile, error) {
	return s.ExtractProfile(s.Default)
}
