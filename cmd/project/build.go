package project

import (
	"fmt"
	"gwb/cmd/ostools"
	"os"
	"path/filepath"
	"strings"
)

var (
	fiberTemplate    = "/cmd/project/template/fiber"
	ginTemplate      = "/cmd/project/template/gin"
	microsrvTemplate = "/cmd/project/template/microservice"
)

type templates struct {
	target string
	path   string
}

func newTemplate(target string) *templates {
	return &templates{target: target}
}

func (tpl *templates) generate() error {
	if err := tpl.pathBuild(); err != nil {
		return err
	}
	fileNodes, err := ostools.BuildTree(tpl.path)
	if err != nil {
		return err
	}

	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	basePath := filepath.Dir(execPath)
	if err := tpl.create(fileNodes, basePath); err != nil {
		return err
	}
	return nil
}

func (tpl *templates) create(fileNodes []ostools.FileNode, basePath string) error {
	for _, fn := range fileNodes {
		path := filepath.Join(basePath, fn.Name)

		if fn.IsDir {
			// Create directory
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", path, err)
			}

			// Recursively create children
			if err := tpl.create(fn.Children, path); err != nil {
				return err
			}
		} else {
			out, err := ostools.NewFile(tpl.path, name)
			if err != nil {
				return err
			}

			realPath := strings.ReplaceAll(path, ".tp", ".go")
			if err := os.WriteFile(realPath, []byte(out), 0644); err != nil {
				return fmt.Errorf("failed to create file %s: %w", path, err)
			}
		}
	}
	return nil
}

func (tpl *templates) pathBuild() error {
	switch tpl.target {
	case "fiber":
		execPath, err := os.Executable()
		if err != nil {
			return err
		}
		tpl.path = filepath.Dir(execPath) + fiberTemplate
	case "gin":
		execPath, err := os.Executable()
		if err != nil {
			return err
		}
		tpl.path = filepath.Dir(execPath) + ginTemplate
	case "microservice":
		execPath, err := os.Executable()
		if err != nil {
			return err
		}
		tpl.path = filepath.Dir(execPath) + microsrvTemplate
	}
	return nil
}
