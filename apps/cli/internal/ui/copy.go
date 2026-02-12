package ui

const (
	WelcomeTitle     = "🧪 Boot Sequence"
	WelcomeNextLabel = "jack in // continue"

	FlowTitle       = "Pick your loadout"
	FlowDescription = "Choose how hard you want to min-max this setup."
	FlowBasicLabel  = "🚀 Speedrun (defaults)"
	FlowAdvLabel    = "🛠️  Full control (advanced)"

	ProjectNameTitle       = "Project name"
	ProjectNameDescription = "Name your repo. Keep it clean, future-you is watching."
	ModuleTitle            = "Go module path"
	ModuleDescription      = "Package identity unlocked."

	DestinationTitle       = "Base destination path"
	DestinationDescription = "Where the repo spawns. Final path is <base>/<project-name>."

	IncludeWebTitle       = "Include web app"
	IncludeWebDescription = "Ship full-stack mode with Vite + React."
	IncludeDockerTitle    = "Include Docker Compose"
	IncludeDockerDesc     = "One command stack spin-up energy."
	IncludeGitTitle       = "Initialize git repository"
	IncludeGitDescription = "Start with commit history from frame 1."

	DBHostTitle = "Postgres host"
	DBHostDesc  = "Database coordinates. localhost is a classic."
	DBPortTitle = "Postgres port"
	DBPortDesc  = "Usually 5432 unless your setup is spicy."
	DBUserTitle = "Postgres user"
	DBUserDesc  = "The DB hero account."
	DBPassTitle = "Postgres password"
	DBPassDesc  = "Type carefully. This one bites."
	DBNameTitle = "Postgres database"
	DBNameDesc  = "Primary schema arena."
	DBSSLTitle  = "Postgres SSL mode"
	DBSSLDesc   = "dev tip: disable locally, require in prod."

	StorageTitle       = "Storage provider"
	StorageDescription = "Pick where files live after upload."
	StorageLocalLabel  = "💾 Local disk"
	StorageS3Label     = "☁️ S3-compatible"
	LocalPathTitle     = "Local storage base path"
	LocalPathDesc      = "Relative path inside your API app."
	S3EndpointTitle    = "S3 endpoint"
	S3RegionTitle      = "S3 region"
	S3BucketTitle      = "S3 bucket"
	S3AccessTitle      = "S3 access key"
	S3SecretTitle      = "S3 secret key"

	ReviewTitle          = "Final boss check"
	ReviewActionTitle    = "What do you think?"
	ReviewGenerateLabel  = "🚀 Looks good, lets get busy!"
	ReviewEditLabel      = "↩️  Changed my mind"
	ReviewCancelLabel    = "🛑 Abort mission"
	ReviewSummaryHeading = "Config snapshot"

	OverwriteTitle = "Destination is not empty. Continue anyway?"
	OverwriteAff   = "YOLO"
	OverwriteNeg   = "Nope"
)

func ContributionLine() string {
	return "Contribute: PRs and spicy bug reports welcome at the repo."
}
