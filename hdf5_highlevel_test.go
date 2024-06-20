package hdf5

import (
	"testing"
)

func TestGetFileTree(t *testing.T) {
	filename := "testdata/HS01919-04C0_MM1044_20240112_1159.h5"
	f, err := OpenFile(filename, F_ACC_RDONLY)
	if err != nil {
		panic(err)
	}
	tree, err := f.getFileTree()
	if err != nil {
		panic(err)
	}
	print(tree.tree, "\n")
	tree.Print()
}
