package configsloader

import (
	"github.com/joho/godotenv"
)

type configsEnvLoader struct {
}

func (cl *configsEnvLoader) LoadFile() error {
	if err := godotenv.Load(".env"); err != nil {
		return err
	}
	return nil
}
