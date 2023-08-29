package configsloader

type ConfigsLoader interface {
	LoadFile() error
}
