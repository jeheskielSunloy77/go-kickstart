package prompts

import (
	"os"

	"github.com/charmbracelet/huh"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/scaffold"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/ui"
)

func StorageFlow(cfg *scaffold.ScaffoldConfiguration) error {
	choice := string(cfg.Storage.Type)
	if choice == "" {
		choice = string(scaffold.StorageLocal)
	}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(ui.StorageTitle).
				Description(ui.StorageDescription).
				Options(
					huh.NewOption(ui.StorageLocalLabel, string(scaffold.StorageLocal)),
					huh.NewOption(ui.StorageS3Label, string(scaffold.StorageS3)),
				).Value(&choice),
		),
	)
	form.WithTheme(ui.HuhTheme())
	form.WithWidth(80)
	form.WithHeight(10)
	form.WithOutput(os.Stdout)
	if err := form.Run(); err != nil {
		return err
	}

	cfg.Storage.Type = scaffold.StorageType(choice)
	if cfg.Storage.Type == scaffold.StorageLocal {
		if cfg.Storage.Local == nil {
			cfg.Storage.Local = &scaffold.LocalStorageConfig{Path: "storage"}
		}
		localForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title(ui.LocalPathTitle).
					Description(ui.LocalPathDesc).
					Value(&cfg.Storage.Local.Path),
			),
		)
		localForm.WithTheme(ui.HuhTheme())
		localForm.WithWidth(80)
		localForm.WithHeight(8)
		localForm.WithOutput(os.Stdout)
		return localForm.Run()
	}

	if cfg.Storage.Type == scaffold.StorageS3 {
		if cfg.Storage.S3 == nil {
			cfg.Storage.S3 = &scaffold.S3Config{}
		}
		s3Form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title(ui.S3EndpointTitle).Value(&cfg.Storage.S3.Endpoint),
				huh.NewInput().Title(ui.S3RegionTitle).Value(&cfg.Storage.S3.Region),
				huh.NewInput().Title(ui.S3BucketTitle).Value(&cfg.Storage.S3.Bucket),
				huh.NewInput().Title(ui.S3AccessTitle).Value(&cfg.Storage.S3.AccessKey),
				huh.NewInput().Title(ui.S3SecretTitle).Value(&cfg.Storage.S3.SecretKey),
			),
		)
		s3Form.WithTheme(ui.HuhTheme())
		s3Form.WithWidth(80)
		s3Form.WithHeight(16)
		s3Form.WithOutput(os.Stdout)
		return s3Form.Run()
	}

	return nil
}
