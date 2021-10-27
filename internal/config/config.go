package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	configsFolderPathEnv     = "BENDY_CONFIGS_FOLDER_PATH"
	defaultConfigsFolderPath = "/home/ythosa/Projects/bendy/configs"
	configNameEnv            = "BENDY_CONFIG_NAME"
	defaultConfigName        = "default"
)

type Config struct {
	Index   *Index
	Storage *Storage
}

func newConfig() *Config {
	return &Config{
		Index:   newIndex(),
		Storage: newStorage(),
	}
}

func Get() *Config {
	var (
		once sync.Once
		cfg  *Config
	)

	once.Do(func() {
		if err := initConfigParser(); err != nil {
			logrus.Fatal(err)
		}

		cfg = newConfig()
	})

	return cfg
}

func initConfigParser() error {
	configsFolderPath := os.Getenv(configsFolderPathEnv)
	if configsFolderPath == "" {
		configsFolderPath = defaultConfigsFolderPath
	}

	configName := os.Getenv(configNameEnv)
	if configName == "" {
		configName = defaultConfigName
	}

	viper.AddConfigPath(configsFolderPath)
	viper.SetConfigName(configName)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error while reading config: %w", err)
	}

	return nil
}

type Index struct {
	TempFilesStoragePath string
	MaxOpenFilesCount    int64
}

func newIndex() *Index {
	return &Index{
		TempFilesStoragePath: viper.GetString("index.temp_file_storage_path"),
		MaxOpenFilesCount:    viper.GetInt64("index.max_open_files_count"),
	}
}

type Storage struct {
	IndexStoragePath           string
	IndexingFilesFilenamesPath string
}

func newStorage() *Storage {
	return &Storage{
		IndexStoragePath:           viper.GetString("storage.index_storage_path"),
		IndexingFilesFilenamesPath: viper.GetString("storage.indexing_files_filenames_path"),
	}
}
