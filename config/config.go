package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	configDir  = ".ackctl"
	configFile = "config.json"
)

type aliyunCliConfig struct {
	Current  string             `json:"current"`
	Profiles []aliyunCliProfile `json:"profiles"`
}

type aliyunCliProfile struct {
	Name            string `json:"name"`
	Mode            string `json:"mode"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
}

type AkConfig struct {
	AccessKeyId     string
	AccessKeySecret string
}

func MustLoadConfig() *AkConfig {
	c, err := loadOwnConfig()
	if err != nil {
		c, err = loadAliyunCliConfig()
		if err != nil {
			log.Fatal("No config available. Try 'ackctl configure'.")
		}
		fmt.Println("Using config from Aliyun CLI.")
	} else {
		fmt.Println("Using config for ackctl")
	}

	return &AkConfig{
		AccessKeyId:     c.AccessKeyId,
		AccessKeySecret: c.AccessKeySecret,
	}
}

func loadAliyunCliConfig() (*aliyunCliProfile, error) {
	path := getAliyunCLIConfigFilePath()

	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := aliyunCliConfig{}
	err = json.Unmarshal(f, &config)
	if err != nil {
		return nil, err
	}

	for _, profile := range config.Profiles {
		if profile.Mode == "AK" && config.Current == profile.Name {
			return &profile, nil
		}
	}

	return nil, fmt.Errorf("no access key found in Aliyun CLI config")
}

func getAliyunCLIConfigFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error to get user home dir: %v", err)
	}
	return filepath.Join(homeDir, ".aliyun/config.json")
}

func getConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error to get user home dir: %v", err)
	}
	return filepath.Join(homeDir, configDir)
}

func getConfigFilePath() string {
	return filepath.Join(getConfigDir(), configFile)
}

func ensureConfigFile() {
	path := getConfigFilePath()
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(getConfigDir(), 0755); err != nil {
				log.Fatal(err)
			}

			text, _ := json.Marshal(aliyunCliProfile{})
			if err := ioutil.WriteFile(getConfigFilePath(), text, 0600); err != nil {
				log.Fatal(err)
			}

			return
		} else {
			log.Fatal(err)
		}
	}
}

func loadOwnConfig() (*aliyunCliProfile, error) {
	t, err := ioutil.ReadFile(getConfigFilePath())
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var c aliyunCliProfile
	if err := json.Unmarshal(t, &c); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &c, nil
}

func UpdateConfigFile(input AkConfig) error {
	ensureConfigFile()

	var c *aliyunCliProfile
	c, err := loadOwnConfig()
	if err != nil {
		c = &aliyunCliProfile{}
	}

	c.AccessKeyId = input.AccessKeyId
	c.AccessKeySecret = input.AccessKeySecret

	t, _ := json.Marshal(c)
	if err := ioutil.WriteFile(getConfigFilePath(), t, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
