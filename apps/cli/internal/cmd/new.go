package cmd

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/prompts"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/scaffold"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/ui"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/validate"
	"github.com/spf13/cobra"
)

var interactive bool
var flags newFlags

type newFlags struct {
	name       string
	modulePath string
	web        bool
	noWeb      bool
	docker     bool
	noDocker   bool
	git        bool
	noGit      bool
	db         string
	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbName     string
	dbSSLMode  string
	pkg        string
	storage    string
	s3Endpoint string
	s3Region   string
	s3Bucket   string
	s3Access   string
	s3Secret   string
}

func init() {
	newCmd := &cobra.Command{
		Use:   "new [name] [path]",
		Short: "Create a new project",
		Args:  cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if interactive || len(args) == 0 {
				return runInteractive()
			}
			cfg, err := configFromFlags(args, flags)
			if err != nil {
				return err
			}
			return runNonInteractive(cfg)
		},
	}
	newCmd.Flags().BoolVar(&interactive, "interactive", false, "force interactive wizard")
	newCmd.Flags().StringVar(&flags.name, "name", "", "project name")
	newCmd.Flags().StringVar(&flags.modulePath, "module", "", "go module path")
	newCmd.Flags().BoolVar(&flags.web, "web", true, "include web app")
	newCmd.Flags().BoolVar(&flags.noWeb, "no-web", false, "exclude web app")
	newCmd.Flags().BoolVar(&flags.docker, "docker", true, "include docker compose")
	newCmd.Flags().BoolVar(&flags.noDocker, "no-docker", false, "exclude docker compose")
	newCmd.Flags().BoolVar(&flags.git, "git", true, "initialize git repository")
	newCmd.Flags().BoolVar(&flags.noGit, "no-git", false, "do not initialize git repository")
	newCmd.Flags().StringVar(&flags.db, "db", "postgres", "database type (postgres)")
	newCmd.Flags().StringVar(&flags.dbHost, "db-host", "", "database host")
	newCmd.Flags().StringVar(&flags.dbPort, "db-port", "", "database port")
	newCmd.Flags().StringVar(&flags.dbUser, "db-user", "", "database user")
	newCmd.Flags().StringVar(&flags.dbPassword, "db-password", "", "database password")
	newCmd.Flags().StringVar(&flags.dbName, "db-name", "", "database name")
	newCmd.Flags().StringVar(&flags.dbSSLMode, "db-ssl-mode", "", "database ssl mode")
	newCmd.Flags().StringVar(&flags.pkg, "pkg", "bun", "package manager (bun)")
	newCmd.Flags().StringVar(&flags.storage, "storage", "local", "storage type (local|s3)")
	newCmd.Flags().StringVar(&flags.s3Endpoint, "s3-endpoint", "", "s3 endpoint")
	newCmd.Flags().StringVar(&flags.s3Region, "s3-region", "", "s3 region")
	newCmd.Flags().StringVar(&flags.s3Bucket, "s3-bucket", "", "s3 bucket")
	newCmd.Flags().StringVar(&flags.s3Access, "s3-access-key", "", "s3 access key")
	newCmd.Flags().StringVar(&flags.s3Secret, "s3-secret-key", "", "s3 secret key")
	rootCmd.AddCommand(newCmd)
}

func runInteractive() error {
	cfg := scaffold.DefaultConfig()
	flow := prompts.FlowBasic

	for {
		choice, err := prompts.ChooseFlow()
		if err != nil {
			return err
		}
		flow = choice
		cfg.UseDefaults = flow == prompts.FlowBasic

		result, err := prompts.BasicFlow(cfg)
		if err != nil {
			return err
		}
		cfg = result

		if flow == prompts.FlowAdvanced {
			cfg.UseDefaults = false
			if err := prompts.ComponentsFlow(&cfg); err != nil {
				return err
			}
			if err := prompts.DatabaseFlow(&cfg); err != nil {
				return err
			}
			if err := prompts.StorageFlow(&cfg); err != nil {
				return err
			}
		}

		action, err := prompts.ReviewConfig(cfg)
		if err != nil {
			return err
		}
		switch action {
		case prompts.ReviewEdit:
			continue
		case prompts.ReviewCancel:
			return errors.New("cancelled")
		case prompts.ReviewGenerate:
			// proceed
		}
		break
	}

	if err := validate.ProjectName(cfg.ProjectName); err != nil {
		return err
	}
	if err := validate.ModulePath(cfg.ModulePath); err != nil {
		return err
	}
	dest, err := validate.ResolveProjectDestination(cfg.Destination, cfg.ProjectName)
	if err != nil {
		return err
	}
	cfg.Destination = dest

	nonEmpty, err := validate.IsNonEmptyDir(dest)
	if err != nil {
		return err
	}
	allowOverwrite := false
	if nonEmpty {
		confirm := false
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().Title("Destination is not empty. Continue?").Value(&confirm),
			),
		)
		if err := form.Run(); err != nil {
			return err
		}
		if !confirm {
			return errors.New("cancelled")
		}
		allowOverwrite = true
	}

	err = ui.RunWithSpinner("Generating project...", func() error {
		return scaffold.ScaffoldProject(cfg, allowOverwrite)
	})
	if err != nil {
		return err
	}

	ui.PrintSummary(cfg.Destination)
	return nil
}

func configFromFlags(args []string, flags newFlags) (scaffold.ScaffoldConfiguration, error) {
	cfg := scaffold.DefaultConfig()
	cfg.UseDefaults = false

	if len(args) > 0 {
		cfg.ProjectName = args[0]
	} else if flags.name != "" {
		cfg.ProjectName = flags.name
	}
	if cfg.ProjectName == "" {
		return cfg, errors.New("project name is required (arg or --name)")
	}

	if len(args) > 1 {
		cfg.Destination = args[1]
	}

	if flags.modulePath == "" {
		return cfg, errors.New("module path is required (--module)")
	}
	cfg.ModulePath = flags.modulePath

	cfg.IncludeWeb = flags.web
	if flags.noWeb {
		cfg.IncludeWeb = false
	}
	cfg.IncludeDocker = flags.docker
	if flags.noDocker {
		cfg.IncludeDocker = false
	}
	cfg.InitGit = flags.git
	if flags.noGit {
		cfg.InitGit = false
	}

	if flags.db != "" {
		cfg.DatabaseType = scaffold.DatabaseType(flags.db)
	}
	if flags.dbHost != "" {
		cfg.DBConnection.Host = flags.dbHost
	}
	if flags.dbPort != "" {
		cfg.DBConnection.Port = flags.dbPort
	}
	if flags.dbUser != "" {
		cfg.DBConnection.User = flags.dbUser
	}
	if flags.dbPassword != "" {
		cfg.DBConnection.Password = flags.dbPassword
	}
	if flags.dbName != "" {
		cfg.DBConnection.Name = flags.dbName
	}
	if flags.dbSSLMode != "" {
		cfg.DBConnection.SSLMode = flags.dbSSLMode
	}

	cfg.PackageManager = scaffold.PackageManager(flags.pkg)
	cfg.Storage.Type = scaffold.StorageType(flags.storage)

	if cfg.Storage.Type == scaffold.StorageS3 {
		if flags.s3Endpoint == "" || flags.s3Region == "" || flags.s3Bucket == "" || flags.s3Access == "" || flags.s3Secret == "" {
			return cfg, errors.New("s3 storage selected: all s3 connection details are required")
		}
		cfg.Storage.S3 = &scaffold.S3Config{
			Endpoint:  flags.s3Endpoint,
			Region:    flags.s3Region,
			Bucket:    flags.s3Bucket,
			AccessKey: flags.s3Access,
			SecretKey: flags.s3Secret,
		}
	}

	if cfg.Storage.Type == scaffold.StorageLocal && cfg.Storage.Local == nil {
		cfg.Storage.Local = &scaffold.LocalStorageConfig{Path: "storage"}
	}

	if err := validate.ProjectName(cfg.ProjectName); err != nil {
		return cfg, err
	}
	if err := validate.ModulePath(cfg.ModulePath); err != nil {
		return cfg, err
	}

	dest, err := validate.ResolveProjectDestination(cfg.Destination, cfg.ProjectName)
	if err != nil {
		return cfg, err
	}
	cfg.Destination = dest

	return cfg, nil
}

func runNonInteractive(cfg scaffold.ScaffoldConfiguration) error {
	nonEmpty, err := validate.IsNonEmptyDir(cfg.Destination)
	if err != nil {
		return err
	}
	if nonEmpty {
		return fmt.Errorf("destination %s is not empty", cfg.Destination)
	}
	if err := scaffold.ScaffoldProject(cfg, false); err != nil {
		return err
	}
	ui.PrintSummary(cfg.Destination)
	return nil
}
