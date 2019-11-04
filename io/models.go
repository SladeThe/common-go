package io

import (
	"fmt"
	"os"
)

type FileInfo struct {
	FullPath string
	os.FileInfo
}

func (fi FileInfo) IsFile() bool {
	return isFile(fi)
}

func (fi FileInfo) IsSymlink() bool {
	return isSymlink(fi)
}

func (fi FileInfo) String() string {
	return fmt.Sprintf("{FullPath:%v, FileInfo:%+v}", fi.FullPath, fi.FileInfo)
}

type Diff struct {
	Item1 *FileInfo
	Item2 *FileInfo
}
