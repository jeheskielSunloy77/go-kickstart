//go:build e2e

package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var (
	cliBinaryPath string
	moduleRoot    string
)

func TestMain(m *testing.M) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "getwd: %v\n", err)
		os.Exit(1)
	}

	moduleRoot = filepath.Dir(cwd)

	buildDir, err := os.MkdirTemp("", "gokickstart-e2e-bin-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "tempdir: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(buildDir)

	cliBinaryPath = filepath.Join(buildDir, "gokickstart")
	build := exec.Command("go", "build", "-o", cliBinaryPath, "./cmd/gokickstart")
	build.Dir = moduleRoot
	build.Env = os.Environ()
	if output, err := build.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "build cli: %v\n%s", err, string(output))
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestNewCommand_GeneratesProjectAndPassesSmoke(t *testing.T) {
	base := t.TempDir()
	projectName := "demo"
	projectDir := filepath.Join(base, projectName)

	result := runCLI(t, baseEnv(t), "new", projectName, base, "--module", "github.com/acme/demo", "--no-git")
	if result.Err != nil {
		t.Fatalf("expected scaffold to succeed: %v\n%s", result.Err, result.Output)
	}

	assertExists(t, filepath.Join(projectDir, "apps/api"))
	assertExists(t, filepath.Join(projectDir, "apps/web"))
	assertExists(t, filepath.Join(projectDir, "packages/ui"))
	assertExists(t, filepath.Join(projectDir, "apps/api/.env"))
	assertExists(t, filepath.Join(projectDir, "apps/web/.env"))
	assertNotExists(t, filepath.Join(projectDir, ".git"))

	assertFileContains(t, filepath.Join(projectDir, "README.md"), "demo/")
	assertFileContains(t, filepath.Join(projectDir, "package.json"), `"web:test"`)

	smokeEnv := baseEnv(t)
	runCommand(t, projectDir, smokeEnv, "bun", "install")
	runCommand(t, projectDir, smokeEnv, "bun", "run", "api:install")
	runCommand(t, projectDir, smokeEnv, "bun", "run", "openapi:generate")
	runCommand(t, projectDir, smokeEnv, "bun", "run", "emails:generate")
	runCommand(t, projectDir, smokeEnv, "bun", "run", "build")
	runCommand(t, projectDir, smokeEnv, "bun", "run", "test")
}

func TestNewCommand_ExcludesWebAndDocker(t *testing.T) {
	base := t.TempDir()
	projectName := "demo"
	projectDir := filepath.Join(base, projectName)

	result := runCLI(
		t,
		baseEnv(t),
		"new",
		projectName,
		base,
		"--module",
		"github.com/acme/demo",
		"--no-web",
		"--no-docker",
		"--no-git",
	)
	if result.Err != nil {
		t.Fatalf("expected scaffold to succeed: %v\n%s", result.Err, result.Output)
	}

	assertExists(t, filepath.Join(projectDir, "apps/api"))
	assertNotExists(t, filepath.Join(projectDir, "apps/web"))
	assertNotExists(t, filepath.Join(projectDir, "packages/ui"))
	assertNotExists(t, filepath.Join(projectDir, "docker-compose.yml"))

	assertFileNotContains(t, filepath.Join(projectDir, "README.md"), "apps/web")
	assertFileNotContains(t, filepath.Join(projectDir, "AGENTS.md"), "apps/web")
	assertFileNotContains(t, filepath.Join(projectDir, "package.json"), `"web:test"`)
	assertFileNotContains(t, filepath.Join(projectDir, "package.json"), `"ui:shadcn:add"`)
}

func TestNewCommand_FailsWithoutModuleFlag(t *testing.T) {
	base := t.TempDir()

	result := runCLI(t, baseEnv(t), "new", "demo", base, "--no-git")
	if result.Err == nil {
		t.Fatalf("expected missing module flag to fail")
	}
	if !strings.Contains(result.Output, "module path is required") {
		t.Fatalf("expected missing module error, got:\n%s", result.Output)
	}
}

func TestNewCommand_FailsForNonEmptyDestination(t *testing.T) {
	base := filepath.Join(t.TempDir(), "demo")
	if err := os.MkdirAll(base, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(base, "existing.txt"), []byte("occupied"), 0o644); err != nil {
		t.Fatalf("write existing file: %v", err)
	}

	result := runCLI(t, baseEnv(t), "new", "demo", base, "--module", "github.com/acme/demo", "--no-git")
	if result.Err == nil {
		t.Fatalf("expected non-empty destination to fail")
	}
	if !strings.Contains(result.Output, "is not empty") {
		t.Fatalf("expected non-empty destination error, got:\n%s", result.Output)
	}
}

func TestNewCommand_FailsForIncompleteS3Configuration(t *testing.T) {
	base := t.TempDir()

	result := runCLI(
		t,
		baseEnv(t),
		"new",
		"demo",
		base,
		"--module",
		"github.com/acme/demo",
		"--storage",
		"s3",
		"--no-git",
	)
	if result.Err == nil {
		t.Fatalf("expected incomplete s3 configuration to fail")
	}
	if !strings.Contains(result.Output, "all s3 connection details are required") {
		t.Fatalf("expected s3 validation error, got:\n%s", result.Output)
	}
}

type runResult struct {
	Output string
	Err    error
}

func runCLI(t *testing.T, env []string, args ...string) runResult {
	t.Helper()

	cmd := exec.Command(cliBinaryPath, args...)
	cmd.Env = env
	cmd.Dir = moduleRoot
	output, err := cmd.CombinedOutput()

	return runResult{
		Output: string(output),
		Err:    err,
	}
}

func runCommand(t *testing.T, dir string, env []string, name string, args ...string) {
	t.Helper()

	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = env
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %s failed: %v\n%s", name, strings.Join(args, " "), err, string(output))
	}
}

func baseEnv(t *testing.T) []string {
	t.Helper()

	home := newIsolatedHome(t)
	values := map[string]string{
		"HOME":                home,
		"XDG_CONFIG_HOME":     filepath.Join(home, ".config"),
		"GOCACHE":             filepath.Join(home, ".cache", "go-build"),
		"GOMODCACHE":          filepath.Join(home, "go", "pkg", "mod"),
		"GOFLAGS":             "-modcacherw",
		"GIT_CONFIG_NOSYSTEM": "1",
		"GIT_TERMINAL_PROMPT": "0",
	}

	envMap := map[string]string{}
	for _, entry := range os.Environ() {
		parts := strings.SplitN(entry, "=", 2)
		key := parts[0]
		value := ""
		if len(parts) == 2 {
			value = parts[1]
		}
		envMap[key] = value
	}
	for key, value := range values {
		envMap[key] = value
	}

	env := make([]string, 0, len(envMap))
	for key, value := range envMap {
		env = append(env, key+"="+value)
	}

	return env
}

func newIsolatedHome(t *testing.T) string {
	t.Helper()

	dir, err := os.MkdirTemp("", "gokickstart-e2e-home-*")
	if err != nil {
		t.Fatalf("create isolated home: %v", err)
	}

	t.Cleanup(func() {
		makeWritable(dir)
		if err := os.RemoveAll(dir); err != nil && !os.IsNotExist(err) {
			t.Fatalf("remove isolated home %s: %v", dir, err)
		}
	})

	return dir
}

func makeWritable(root string) {
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			_ = os.Chmod(path, 0o755)
			return nil
		}
		_ = os.Chmod(path, 0o644)
		return nil
	})
}

func assertExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected %s to exist: %v", path, err)
	}
}

func assertNotExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected %s not to exist", path)
	}
}

func assertFileContains(t *testing.T, path string, needle string) {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if !strings.Contains(string(content), needle) {
		t.Fatalf("expected %s to contain %q", path, needle)
	}
}

func assertFileNotContains(t *testing.T, path string, needle string) {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if strings.Contains(string(content), needle) {
		t.Fatalf("expected %s not to contain %q", path, needle)
	}
}
