package scaffold

import "fmt"

func EnvOverridesFromConfig(cfg ScaffoldConfiguration) map[string]map[string]string {
	key := "apps/api/.env.example"
	overrides := map[string]map[string]string{}

	api := map[string]string{}
	api["API_PRIMARY.APP_NAME"] = fmt.Sprintf("\"%s\"", cfg.ProjectName)
	api["API_DATABASE.HOST"] = fmt.Sprintf("\"%s\"", cfg.DBConnection.Host)
	api["API_DATABASE.PORT"] = fmt.Sprintf("\"%s\"", cfg.DBConnection.Port)
	api["API_DATABASE.USER"] = fmt.Sprintf("\"%s\"", cfg.DBConnection.User)
	api["API_DATABASE.PASSWORD"] = fmt.Sprintf("\"%s\"", cfg.DBConnection.Password)
	api["API_DATABASE.NAME"] = fmt.Sprintf("\"%s\"", cfg.DBConnection.Name)
	api["API_DATABASE.SSL_MODE"] = fmt.Sprintf("\"%s\"", cfg.DBConnection.SSLMode)

	api["API_FILE_STORAGE.PROVIDER"] = fmt.Sprintf("\"%s\"", cfg.Storage.Type)
	if cfg.Storage.Type == StorageLocal && cfg.Storage.Local != nil {
		api["API_FILE_STORAGE.LOCAL.BASE_DIR"] = fmt.Sprintf("\"%s\"", cfg.Storage.Local.Path)
	}
	if cfg.Storage.Type == StorageS3 && cfg.Storage.S3 != nil {
		api["API_FILE_STORAGE.S3.BUCKET"] = fmt.Sprintf("\"%s\"", cfg.Storage.S3.Bucket)
		api["API_FILE_STORAGE.S3.REGION"] = fmt.Sprintf("\"%s\"", cfg.Storage.S3.Region)
		api["API_FILE_STORAGE.S3.ENDPOINT"] = fmt.Sprintf("\"%s\"", cfg.Storage.S3.Endpoint)
		api["API_FILE_STORAGE.S3.ACCESS_KEY_ID"] = fmt.Sprintf("\"%s\"", cfg.Storage.S3.AccessKey)
		api["API_FILE_STORAGE.S3.SECRET_ACCESS_KEY"] = fmt.Sprintf("\"%s\"", cfg.Storage.S3.SecretKey)
	}
	overrides[key] = api

	return overrides
}
