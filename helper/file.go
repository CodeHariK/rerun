package helper

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func AddRecursive(watcher *fsnotify.Watcher, directory string) error {
	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			err = watcher.Add(path)
			if err != nil {
				return err
			}
			fmt.Println("Watching:", path)
		}
		return nil
	})
}

// TreeNode represents a node in the directory tree
type TreeNode struct {
	Name     string      `json:"name"`
	IsDir    bool        `json:"is_dir"`
	Children []*TreeNode `json:"children,omitempty"`
}

// Tree function to build the directory tree structure
func Tree(root string) (*TreeNode, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}

	node := &TreeNode{
		Name:  info.Name(),
		IsDir: info.IsDir(),
	}

	if info.IsDir() {
		files, err := os.ReadDir(root)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			childNode, err := Tree(filepath.Join(root, file.Name()))
			if err != nil {
				return nil, err
			}
			node.Children = append(node.Children, childNode)
		}
	}

	return node, nil
}

func Pwd(directory string) string {
	tree, err := Tree(directory)
	if err != nil {
		panic(err)
	}

	jsonData, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(jsonData)
}
