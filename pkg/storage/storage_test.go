package storage

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_AddProfile(t *testing.T) {
	tests := []struct {
		name       string
		storage    *Storage
		newProfile *Profile
	}{
		{
			name: "add new Profile",
			storage: &Storage{
				DefaultProfile: "default_profile",
				Profiles: []*Profile{
					{Name: "default_profile", Key: "key"},
				},
			},
			newProfile: &Profile{Name: "test", Key: "key"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storage.AddProfile(tt.newProfile)
			assert.Len(t, tt.storage.Profiles, 2)
			assert.Equal(t, tt.newProfile, tt.storage.Profiles[1])
		})
	}
}

func TestStorageExtractProfile(t *testing.T) {

	tests := []struct {
		name     string
		storage  *Storage
		wantName string
		wantErr  bool
		want     *Profile
		errMsg   string
	}{
		{
			name: "existing profile",
			storage: &Storage{
				DefaultProfile: "default_profile",
				Profiles: []*Profile{
					{Name: "default_profile", Key: "key"},
					{Name: "test", Key: "key"},
				},
			},
			wantName: "test",
			wantErr:  false,
			want:     &Profile{Name: "test", Key: "key"},
		},
		{
			name: "not found profile",
			storage: &Storage{
				DefaultProfile: "default_profile",
				Profiles: []*Profile{
					{Name: "default_profile", Key: "key"},
					{Name: "test", Key: "key"},
				},
			},
			wantName: "test2",
			wantErr:  true,
			want:     nil,
			errMsg:   "not found profile with name test2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.storage.ExtractProfile(tt.wantName)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.errMsg)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestStorage_Commit(t *testing.T) {
	// Create a temporary directory for tests
	tempDir := t.TempDir()

	tests := []struct {
		name     string
		storage  *Storage
		wantErr  bool
		validate func(*testing.T, string)
	}{
		{
			name: "successful Commit with DefaultProfile profile",
			storage: &Storage{
				DefaultProfile: "test_profile",
				Profiles:       []*Profile{{Name: "test_profile", Key: "test_key"}},
				StoragePath:    path.Join(tempDir, "config1.json"),
			},
			wantErr: false,
			validate: func(t *testing.T, path string) {
				// Read the Commitd file
				data, err := os.ReadFile(path)
				require.NoError(t, err)

				var Commitd map[string]interface{}
				if err := json.Unmarshal(data, &Commitd); err != nil {
					t.Fatalf("Failed to unmarshal Commitd data: %v", err)
				}

				// Check StoragePath is not in JSON
				if _, exists := Commitd["StoragePath"]; exists {
					t.Error("StoragePath should not be present in Commitd JSON")
				}

				// Check required fields
				if def, ok := Commitd["default_profile"]; !ok {
					t.Error("DefaultProfile field missing from JSON")
				} else if def != "test_profile" {
					t.Errorf("DefaultProfile field = %v, want %v", def, "test_profile")
				}

				if profiles, ok := Commitd["profiles"]; !ok {
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
			name: "successful Commit with multiple profiles",
			storage: &Storage{
				DefaultProfile: "test_profile1",
				Profiles:       []*Profile{{Name: "test_profile1", Key: "test_key1"}, {Name: "test_profile2", Key: "test_key2"}},
				StoragePath:    path.Join(tempDir, "config2.json"),
			},
			wantErr: false,
			validate: func(t *testing.T, path string) {
				data, err := os.ReadFile(path)
				require.NoError(t, err)

				var Commitd map[string]interface{}
				if err := json.Unmarshal(data, &Commitd); err != nil {
					t.Fatalf("Failed to unmarshal Commitd data: %v", err)
				}

				if profiles, ok := Commitd["profiles"]; !ok {
					t.Error("profiles field missing from JSON")
				} else if profs, ok := profiles.([]interface{}); !ok {
					t.Error("profiles is not an array")
				} else if len(profs) != 2 {
					t.Errorf("got %d profiles, want 2", len(profs))
				}
			},
		},
		{
			name: "Commit to invalid path",
			storage: &Storage{
				DefaultProfile: "test_profile",
				Profiles:       []*Profile{{Name: "test_profile", Key: "test_key"}},
				StoragePath:    path.Join(tempDir, "invalid", "config.json"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.storage.Commit()

			// Check error condition
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Commit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we expect success, validate the Commitd file
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
			name:      "successful load with DefaultProfile profile",
			setupJSON: `{"default_profile":"test_profile","profiles":[{"name":"test_profile","key":"test_key"}]}`,
			storage: &Storage{
				StoragePath: path.Join(tempDir, "config1.json"),
			},
			wantErr: false,
			validate: func(t *testing.T, s *Storage) {
				assert.Equal(t, "test_profile", s.DefaultProfile)
				assert.Len(t, s.Profiles, 1)
				assert.Equal(t, "test_profile", s.Profiles[0].Name)
				assert.Equal(t, "test_key", s.Profiles[0].Key)
			},
		},
		{
			name:      "successful load with multiple profiles",
			setupJSON: `{"default_profile":"test_profile1","profiles":[{"name":"test_profile1","key":"test_key1"},{"name":"test_profile2","key":"test_key2"}]}`,
			storage: &Storage{
				StoragePath: path.Join(tempDir, "config2.json"),
			},
			wantErr: false,
			validate: func(t *testing.T, s *Storage) {
				assert.Equal(t, "test_profile1", s.DefaultProfile)
				assert.Len(t, s.Profiles, 2)
				assert.Equal(t, "test_profile1", s.Profiles[0].Name)
				assert.Equal(t, "test_key1", s.Profiles[0].Key)
				assert.Equal(t, "test_profile2", s.Profiles[1].Name)
				assert.Equal(t, "test_key2", s.Profiles[1].Key)
			},
		},
		{
			name:      "invalid json",
			setupJSON: `invalid json content`,
			storage: &Storage{
				StoragePath: path.Join(tempDir, "invalid.json"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test file if JSON content is provided
			if tt.setupJSON != "" {
				err := os.WriteFile(tt.storage.StoragePath, []byte(tt.setupJSON), 0600)
				require.NoError(t, err)
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
