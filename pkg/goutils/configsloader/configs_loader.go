package configsloader

import (
	"github.com/joho/godotenv"
)

type ConfigsLoader struct {
}

func (cl *ConfigsLoader) loadFile(fileName string) error {
	if err := godotenv.Load(fileName); err != nil {
		return err
	}
	return nil
}
