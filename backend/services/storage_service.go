package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"transok/backend/utils/common"
)

type StorageService struct {
	ctx         context.Context
	data        map[string]interface{}
	storagePath string
	mu          sync.RWMutex
}

var (
	storage     *StorageService
	storageOnce sync.Once
)

// Storage ensures that only one instance of StorageService is created
func Storage() *StorageService {
	if storage == nil {
		storageOnce.Do(func() {
			storage = &StorageService{
				data: make(map[string]interface{}),
			}
		})
	}
	return storage
}

// Init initializes the storage service and sets storage paths based on the environment
func (s *StorageService) Init(ctx context.Context) error {
	s.ctx = ctx

	// Get storage path based on environment variables
	env := system.GetEnv()
	fmt.Println("Environment:", env)
	if env == "" {
		env = "dev" // Default to development environment
	}

	// Get the base storage path suitable for the current operating system
	basePath := common.GetBasePath()
	fmt.Println("basePath => ", basePath)

	s.storagePath = filepath.Join(basePath, "storage.json")

	// Ensure the storage directory exists
	if err := os.MkdirAll(filepath.Dir(s.storagePath), 0755); err != nil {
		return fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Load existing data first
	if err := s.loadData(); err != nil {
		return err
	}

	// Then check and set default values
	keys := s.GetKeys()
	if !slices.Contains(keys, "language") {
		fmt.Println("Setting language to en")
		s.Set("language", "en")
	}

	if !slices.Contains(keys, "port") {
		fmt.Println("Setting port to 9482")
		s.Set("port", "9482")
	}

	if !slices.Contains(keys, "share-list") {
		fmt.Println("Setting share-list to empty")
		s.Set("share-list", []interface{}{})
	}

	if !slices.Contains(keys, "uname") {
		fmt.Println("Setting uname to transok")
		s.Set("uname", "transok")
	}

	if !slices.Contains(keys, "is-share") {
		fmt.Println("Setting is-share to false")
		s.Set("is-share", false)
	}

	return nil
}

// Set sets a key-value pair
func (s *StorageService) Set(key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
	return s.saveData()
}

// Get retrieves a value
func (s *StorageService) Get(key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, exists := s.data[key]
	return value, exists
}

// Delete removes a key-value pair
func (s *StorageService) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
	return s.saveData()
}

// loadData loads data from file
func (s *StorageService) loadData() error {
	data, err := os.ReadFile(s.storagePath)
	if os.IsNotExist(err) {
		return nil // Return empty data if file does not exist
	}
	if err != nil {
		return fmt.Errorf("failed to read storage file: %w", err)
	}

	return json.Unmarshal(data, &s.data)
}

// saveData saves data to file
func (s *StorageService) saveData() error {
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize data: %w", err)
	}

	return os.WriteFile(s.storagePath, data, 0644)
}

// GetKeys returns all keys in storage
func (s *StorageService) GetKeys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]string, 0, len(s.data))
	for key := range s.data {
		keys = append(keys, key)
	}
	return keys
}

// Clear removes all data from storage
func (s *StorageService) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[string]interface{})
	return s.saveData()
}
