package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestStorage_AddProfile(t *testing.T) {
	tests := []struct {
		name       string
		storage    *Storage
		newProfile *Profile
		validate   func(*testing.T, *Storage)
	}{
		{
			name: "add new Profile",
			storage: &Storage{
				Default: "default",
				Profiles: []*Profile{
					{Name: "default", Key: "key"},
				},
			},
			newProfile: &Profile{Name: "test", Key: "key"},
			validate: func(t *testing.T, s *Storage) {
				if len(s.Profiles) != 2 {
					t.Fatalf("got %d profiles, want 2", len(s.Profiles))
				}
				if s.Profiles[1].Name != "test" {
					t.Fatalf("got %s, want %s", s.Profiles[1].Name, "test")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storage.AddProfile(tt.newProfile)
			tt.validate(t, tt.storage)
		})
	}
}

func TestStorageExtractProfile(t *testing.T) {

	tests := []struct {
		name     string
		storage  *Storage
		wantName string
		wantErr  bool
		validate func(*testing.T, *Profile)
	}{
		{
			name: "existing profile",
			storage: &Storage{
				Default: "default",
				Profiles: []*Profile{
					{Name: "default", Key: "key"},
					{Name: "test", Key: "key"},
				},
			},
			wantName: "test",
			wantErr:  false,
			validate: func(t *testing.T, p *Profile) {
				if p == nil {
					t.Fatalf("got nil, want %s", "test")
				}
				if p.Name != "test" {
					t.Fatalf("got %s, want %s", p.Name, "test")
				}
			},
		},
		{
			name: "not found profile",
			storage: &Storage{
				Default: "default",
				Profiles: []*Profile{
					{Name: "default", Key: "key"},
					{Name: "test", Key: "key"},
				},
			},
			wantName: "test2",
			wantErr:  true,
			validate: func(t *testing.T, p *Profile) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.storage.ExtractProfile(tt.wantName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.ExtractProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.validate(t, got)
		})
	}
}

func TestStorage_Save(t *testing.T) {
	// Create a temporary directory for tests
	tempDir := t.TempDir()

	tests := []struct {
		name     string
		storage  *Storage
		wantErr  bool
		validate func(*testing.T, string)
	}{
		{
			name: "successful save with default profile",
			storage: &Storage{
				Default:     "test_profile",
				Profiles:    []*Profile{{Name: "test_profile", Key: "test_key"}},
				StoragePath: fmt.Sprintf("%s/config1.json", tempDir),
			},
			wantErr: false,
			validate: func(t *testing.T, path string) {
				// Read the saved file
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("Failed to read saved file: %v", err)
				}

				var saved map[string]interface{}
				if err := json.Unmarshal(data, &saved); err != nil {
					t.Fatalf("Failed to unmarshal saved data: %v", err)
				}

				// Check StoragePath is not in JSON
				if _, exists := saved["StoragePath"]; exists {
					t.Error("StoragePath should not be present in saved JSON")
				}

				// Check required fields
				if def, ok := saved["default"]; !ok {
					t.Error("default field missing from JSON")
				} else if def != "test_profile" {
					t.Errorf("default field = %v, want %v", def, "test_profile")
				}

				if profiles, ok := saved["profiles"]; !ok {
					t.Error("profiles field missing from JSON")
				} else if profs, ok := profiles.([]interface{}); !ok {
					t.Error("profiles is not an array")
				} else if len(profs) != 1 {
					t.Errorf("got %d profiles, want 1", len(profs))
				}

				// Check file permissions
				info, err := os.Stat(path)
				if err != nil {
					t.Fatalf("Failed to get file info: %v", err)
				}
				if info.Mode().Perm() != 0600 {
					t.Errorf("Wrong file permissions, got %v, want %v", info.Mode().Perm(), 0600)
				}
			},
		},
		{
			name: "successful save with multiple profiles",
			storage: &Storage{
				Default: "profile1",
				Profiles: []*Profile{
					{Name: "profile1", Key: "key1"},
					{Name: "profile2", Key: "key2"},
				},
				StoragePath: fmt.Sprintf("%s/config2.json", tempDir),
			},
			wantErr: false,
			validate: func(t *testing.T, path string) {
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("Failed to read saved file: %v", err)
				}

				var saved map[string]interface{}
				if err := json.Unmarshal(data, &saved); err != nil {
					t.Fatalf("Failed to unmarshal saved data: %v", err)
				}

				if profiles, ok := saved["profiles"]; !ok {
					t.Error("profiles field missing from JSON")
				} else if profs, ok := profiles.([]interface{}); !ok {
					t.Error("profiles is not an array")
				} else if len(profs) != 2 {
					t.Errorf("got %d profiles, want 2", len(profs))
				}
			},
		},
		{
			name: "save to invalid path",
			storage: &Storage{
				Default:     "test_profile",
				Profiles:    []*Profile{{Name: "test_profile", Key: "test_key"}},
				StoragePath: fmt.Sprintf("%s/nonexistent/config.json", tempDir),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.storage.Save()

			// Check error condition
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we expect success, validate the saved file
			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, tt.storage.StoragePath)
			}
		})
	}
}

func TestStorage_Load(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name      string
		setupJSON string
		storage   *Storage
		wantErr   bool
		validate  func(*testing.T, *Storage)
	}{
		{
			name: "successful load with default profile",
			setupJSON: `{
				"default": "test_profile",
				"profiles": [
					{"name": "test_profile", "key": "test_key"}
				]
			}`,
			storage: &Storage{},
			wantErr: false,
			validate: func(t *testing.T, s *Storage) {
				if s.Default != "test_profile" {
					t.Errorf("Default = %v, want %v", s.Default, "test_profile")
				}
				if len(s.Profiles) != 1 {
					t.Errorf("got %d profiles, want 1", len(s.Profiles))
					return
				}
				if s.Profiles[0].Name != "test_profile" {
					t.Errorf("Profile name = %v, want %v", s.Profiles[0].Name, "test_profile")
				}
				if s.Profiles[0].Key != "test_key" {
					t.Errorf("Profile key = %v, want %v", s.Profiles[0].Key, "test_key")
				}
			},
		},
		{
			name: "successful load with multiple profiles",
			setupJSON: `{
				"default": "profile1",
				"profiles": [
					{"name": "profile1", "key": "key1"},
					{"name": "profile2", "key": "key2"}
				]
			}`,
			storage: &Storage{},
			wantErr: false,
			validate: func(t *testing.T, s *Storage) {
				if s.Default != "profile1" {
					t.Errorf("Default = %v, want %v", s.Default, "profile1")
				}
				if len(s.Profiles) != 2 {
					t.Errorf("got %d profiles, want 2", len(s.Profiles))
					return
				}
				if s.Profiles[0].Name != "profile1" || s.Profiles[1].Name != "profile2" {
					t.Errorf("Profile names incorrect, got %v and %v, want profile1 and profile2",
						s.Profiles[0].Name, s.Profiles[1].Name)
				}
			},
		},
		{
			name:      "file not found",
			storage:   &Storage{},
			setupJSON: "",
			wantErr:   true,
		},
		{
			name: "invalid json",
			setupJSON: `{
				"default": "test_profile",
				"profiles": [
					{"name": "test_profile", "key": "test_key"
				]
			}`,
			storage: &Storage{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test file if JSON content is provided
			tt.storage.StoragePath = fmt.Sprintf("%s/%s.json", tempDir, tt.name)
			if tt.setupJSON != "" {
				err := os.WriteFile(tt.storage.StoragePath, []byte(tt.setupJSON), 0600)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			}

			// Run the test
			err := tt.storage.Load()

			// Check error condition
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Validate the loaded data if no error was expected
			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, tt.storage)
			}
		})
	}
}
