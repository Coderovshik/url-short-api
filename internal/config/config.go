package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}

type Config struct {
	Env        string `json:"env" env-default:"local"`
	HTTPServer `json:"http_server"`
	Logger     `json:"logger"`
}

type HTTPServer struct {
	Address     string   `json:"address" env-default:"localhost:8080"`
	Timeout     Duration `json:"timeout" env-default:"4s"`
	IdleTimeout Duration `json:"idle_timeout" env-default:"60s"`
}

type Logger struct {
	Type  string `json:"type" env-required:"true"`
	Level string `json:"level" env-required:"true"`
	Out   string `json:"out" env-required:"true"`
}

func MustLoad(configPath string) (cfg Config) {
	configFile, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("%s does not exist\n", configPath)
		}
		log.Fatalf("failed to open %s\n", configPath)
	}

	contents, err := io.ReadAll(configFile)
	if err != nil {
		log.Fatalf("failed to read %s\n", configPath)
	}

	err = json.Unmarshal(contents, &cfg)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("failed to unmarshal contents of %s\n", configPath)
	}

	return
}
