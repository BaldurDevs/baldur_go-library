package configsloader

func ConfigsEnvLoaderFactory() ConfigsLoader {
	return &configsEnvLoader{}
}
