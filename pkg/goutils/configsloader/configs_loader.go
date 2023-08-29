package configsloader

import (
	"github.com/joho/godotenv"
)

type configsLoader struct {
}

func (cl *configsLoader) LoadFile(fileName string) error {
	if err := godotenv.Load(fileName); err != nil {
		return err
	}
	return nil
}
