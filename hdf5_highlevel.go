package hdf5

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Type that represents the file tree of the hdf file in the form of map[parent][]children
type FileTree struct {
	tree map[string][]string
}

func (file *File) getFileTree() (*FileTree, error) {
	//Get top level groups
	tree := FileTree{map[string][]string{"": []string{}}}

	objCount, err := file.NumObjects()
	if err != nil {
		return &FileTree{}, fmt.Errorf("could not get num of root objects: %w", err)
	}
	filename := filepath.Base(file.FileName())

	var i uint
	var children []string
	for i = 0; i < objCount; i++ {
		child, err := file.ObjectNameByIndex(i)
		if err != nil {
			return &FileTree{}, fmt.Errorf("could not open child object: %w", err)
		}
		children = append(children, child)
		getChildrenRecursive(file, &tree, filename+"/"+child)
	}
	tree.tree[""] = append(tree.tree[""], filename)
	tree.tree[filename] = children

	return &tree, nil
}

func getChildrenRecursive(file *File, tree *FileTree, parent string) error {
	var children []string
	_, relPath, _ := strings.Cut(parent, "/")
	grp, err := file.OpenGroup("/" + relPath)
	if err != nil {
		// Object is not a group, return empty children
		tree.tree[parent] = []string{}
		return nil
	}
	defer grp.Close()

	objCount, err := grp.NumObjects()
	if err != nil {
		return fmt.Errorf("could not get number of children: %w", err)
	}
	var i uint
	for i = 0; i < objCount; i++ {
		name, err := grp.ObjectNameByIndex(i)
		if err != nil {
			return fmt.Errorf("could net get object name for idx %d: %w", i, err)
		}
		children = append(children, name)
	}
	tree.tree[parent] = children
	for _, child := range children {
		getChildrenRecursive(file, tree, parent+"/"+child)
	}
	return nil
}

func (t *FileTree) Print() {
	for _, child := range t.tree[""] {
		fmt.Printf("%s\n", child)
		t.printRecursive("  |--", child)
	}
}

func (t *FileTree) printRecursive(prefix string, parent string) {
	for _, child := range t.tree[parent] {
		fmt.Printf("%s%s\n", prefix, child)
		t.printRecursive("   "+prefix, parent+"/"+child)
	}
}
