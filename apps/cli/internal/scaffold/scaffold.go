package scaffold

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/templates"
)

const (
	TemplateModulePath  = "github.com/jeheskielSunloy77/go-kickstart"
	TemplateProjectName = "go-kickstart"
)

func ScaffoldProject(cfg ScaffoldConfiguration, allowOverwrite bool) error {
	sub, err := fs.Sub(templates.MonorepoFS, "monorepo")
	if err != nil {
		return err
	}
	return ScaffoldFromFS(cfg, allowOverwrite, sub, nil)
}

func ScaffoldFromFS(cfg ScaffoldConfiguration, allowOverwrite bool, source fs.FS, envOverrides map[string]map[string]string) error {
	if err := EnsureSafeDestination(cfg.Destination, allowOverwrite); err != nil {
		return err
	}
	replacements := map[string]string{
		"{{PROJECT_NAME}}":  cfg.ProjectName,
		"{{MODULE_PATH}}":   cfg.ModulePath,
		TemplateModulePath:  cfg.ModulePath,
		TemplateProjectName: cfg.ProjectName,
	}
	skip := combineSkips(DefaultSkip, ShouldSkipForConfig(cfg))
	transform := func(path string, content []byte) ([]byte, error) {
		return []byte(ReplaceTokens(string(content), replacements)), nil
	}
	if err := RenderFS(source, cfg.Destination, skip, transform); err != nil {
		return err
	}
	if envOverrides == nil {
		envOverrides = EnvOverridesFromConfig(cfg)
	}
	if err := generateEnvFiles(cfg.Destination, envOverrides); err != nil {
		return err
	}
	if cfg.InitGit {
		if err := InitGitRepo(cfg.Destination); err != nil {
			return err
		}
	}
	return nil
}

func combineSkips(skips ...func(string) bool) func(string) bool {
	return func(path string) bool {
		for _, fn := range skips {
			if fn != nil && fn(path) {
				return true
			}
		}
		return false
	}
}

func generateEnvFiles(root string, overrides map[string]map[string]string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(info.Name(), ".env.example") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			target := strings.TrimSuffix(path, ".example")
			key := strings.TrimPrefix(path, root+string(os.PathSeparator))
			merged := MergeEnvExample(string(content), nil)
			if overrides != nil {
				if o, ok := overrides[key]; ok {
					merged = MergeEnvExample(string(content), o)
				}
			}
			if err := os.WriteFile(target, []byte(merged), info.Mode()); err != nil {
				return err
			}
		}
		return nil
	})
}
