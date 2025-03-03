package ostools

import (
	"os"
	"path/filepath"
)

type FileNode struct {
	Name     string
	IsDir    bool
	Children []FileNode
}

func NewDirectory(name string) error {
	return nil
}

func BuildTree(target string) ([]FileNode, error) {
	var nodes []FileNode
	entries, err := os.ReadDir(target)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		node := FileNode{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
		}

		if entry.IsDir() {
			childPath := filepath.Join(target, entry.Name())
			children, err := BuildTree(childPath)
			if err != nil {
				return nil, err
			}
			node.Children = children
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}
