package scaffold

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type TransformFunc func(path string, content []byte) ([]byte, error)

func RenderFS(source fs.FS, dest string, shouldSkip func(path string) bool, transform TransformFunc) error {
	return fs.WalkDir(source, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == "." {
			return nil
		}
		if shouldSkip != nil && shouldSkip(path) {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		outPath := filepath.Join(dest, stripTemplateSuffix(path))
		if d.IsDir() {
			return os.MkdirAll(outPath, 0o755)
		}
		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return err
		}
		in, err := source.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		out, err := os.OpenFile(outPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			return err
		}
		defer out.Close()
		if transform == nil {
			if _, err := io.Copy(out, in); err != nil {
				return err
			}
			return nil
		}
		data, err := io.ReadAll(in)
		if err != nil {
			return err
		}
		data, err = transform(path, data)
		if err != nil {
			return err
		}
		if _, err := out.Write(data); err != nil {
			return err
		}
		return nil
	})
}

func stripTemplateSuffix(path string) string {
	if strings.HasSuffix(path, ".tmpl") {
		return strings.TrimSuffix(path, ".tmpl")
	}
	return path
}

func DefaultSkip(path string) bool {
	base := filepath.Base(path)
	for _, glob := range DefaultIgnoreGlobs {
		if match, _ := filepath.Match(strings.ReplaceAll(glob, "**/", ""), path); match {
			if strings.HasPrefix(base, ".env.example") {
				continue
			}
			return true
		}
	}
	return false
}
