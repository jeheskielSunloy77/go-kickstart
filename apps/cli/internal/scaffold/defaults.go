package scaffold

func DefaultConfig() ScaffoldConfiguration {
	return ScaffoldConfiguration{
		ProjectName:  "my-app",
		Destination:  "",
		ModulePath:   "github.com/yourorg/my-app",
		IncludeWeb:   true,
		DatabaseType: DatabasePostgres,
		DBConnection: DBConnection{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "postgres",
			Name:     "app",
			SSLMode:  "disable",
		},
		PackageManager: PackageBun,
		IncludeDocker:  true,
		InitGit:        true,
		Storage: StorageConfig{
			Type: StorageLocal,
			Local: &LocalStorageConfig{
				Path: "storage",
			},
		},
		UseDefaults: true,
	}
}
