package scaffold

type DatabaseType string

type PackageManager string

type StorageType string

const (
	DatabasePostgres DatabaseType   = "postgres"
	PackageBun       PackageManager = "bun"
	StorageLocal     StorageType    = "local"
	StorageS3        StorageType    = "s3"
)

type DBConnection struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type S3Config struct {
	Endpoint  string
	Region    string
	Bucket    string
	AccessKey string
	SecretKey string
}

type LocalStorageConfig struct {
	Path string
}

type StorageConfig struct {
	Type  StorageType
	S3    *S3Config
	Local *LocalStorageConfig
}

type ScaffoldConfiguration struct {
	ProjectName    string
	Destination    string
	ModulePath     string
	IncludeWeb     bool
	DatabaseType   DatabaseType
	DBConnection   DBConnection
	PackageManager PackageManager
	IncludeDocker  bool
	InitGit        bool
	Storage        StorageConfig
	UseDefaults    bool
}
