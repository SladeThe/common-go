package io

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

const bufferSize = 1 << 20

func errEmptyPath() error {
	return errors.New("path is empty")
}

func errNotFile(path string) error {
	return fmt.Errorf("not a file: %s", path)
}

func errNotDir(path string) error {
	return fmt.Errorf("not a directory: %s", path)
}

func exists(path string) (bool, error) {
	if len(path) <= 0 {
		return false, errEmptyPath()
	}

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func isDir(info os.FileInfo) bool {
	return info.IsDir()
}

func isFile(info os.FileInfo) bool {
	return info.Mode().IsRegular()
}

func isSymlink(info os.FileInfo) bool {
	return info.Mode()&os.ModeSymlink != 0
}

func isFileOrDir(path string, dir bool) (bool, error) {
	if len(path) <= 0 {
		return false, errEmptyPath()
	}

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return dir && isDir(info) || !dir && isFile(info), nil
}

func isEmpty(path string) (bool, error) {
	if len(path) <= 0 {
		return false, errEmptyPath()
	}

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if isDir(info) {
		dir, err := os.Open(path)
		if err != nil {
			return false, err
		}

		_, err = dir.Readdirnames(1)

		if err := dir.Close(); err != nil {
			return false, err
		}

		if errors.Is(err, EOF) {
			return true, nil
		}

		return false, err
	}

	if isFile(info) {
		return info.Size() == 0, nil
	}

	return false, nil
}

func parent(path string) (string, error) {
	if len(path) <= 0 {
		return "", errEmptyPath()
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	parent := filepath.Dir(abs)
	if len(parent) <= 0 {
		return "", fmt.Errorf("parent is empty: %s", path)
	}

	parent, err = filepath.Abs(parent)
	if err != nil {
		return "", err
	}

	if parent == abs || parent == path {
		return "", fmt.Errorf("parent equals to absolute path: %s", path)
	}

	return parent, nil
}

func isParentDir(path string) (bool, error) {
	parent, err := parent(path)
	if err != nil {
		return false, err
	}

	if isDir, err := IsDir(parent); err != nil {
		return false, err
	} else {
		return isDir, nil
	}
}

func checkFileOrDir(path string, dir bool) (*FileInfo, error) {
	if len(path) <= 0 {
		return nil, errEmptyPath()
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	{
		info := &FileInfo{FileInfo: info, FullPath: path}

		if dir {
			if !isDir(info) {
				return info, errNotDir(path)
			}
		} else {
			if !isFile(info) {
				return info, errNotFile(path)
			}
		}

		return info, nil
	}
}

func closeMany(closers ...io.Closer) error {
	var firstErr error

	for _, c := range closers {
		if c != nil {
			if err := c.Close(); err != nil && firstErr == nil {
				firstErr = err
			}
		}
	}

	return firstErr
}

func closeQuietly(closers ...io.Closer) {
	for _, c := range closers {
		if c != nil {
			_ = c.Close()
		}
	}
}

func closeOr(handleErr func(err error), closers ...io.Closer) {
	for _, c := range closers {
		if c != nil {
			if err := c.Close(); err != nil {
				handleErr(err)
			}
		}
	}
}

func readChunk(r io.Reader, buf []byte) (count int, err error) {
	for len(buf) > 0 {
		n, err := r.Read(buf)
		if n > 0 {
			count += n
		}
		if err != nil {
			return count, err
		}
		buf = buf[n:]
	}

	return
}

func readersContentEqual(r1, r2 io.Reader) (bool, error) {
	buf1 := make([]byte, bufferSize)
	buf2 := make([]byte, bufferSize)

	var wg sync.WaitGroup

	for {
		var (
			n1, n2     int
			err1, err2 error
		)

		wg.Add(2)

		go func() {
			defer wg.Done()
			n1, err1 = readChunk(r1, buf1)
		}()

		go func() {
			defer wg.Done()
			n2, err2 = readChunk(r2, buf2)
		}()

		wg.Wait()

		eof1 := errors.Is(err1, EOF)
		if err1 != nil && !eof1 {
			return false, err1
		}

		eof2 := errors.Is(err2, EOF)
		if err2 != nil && !eof2 {
			return false, err2
		}

		if n1 != n2 || !bytes.Equal(buf1[:n1], buf2[:n2]) {
			return false, nil
		}

		if eof1 && eof2 {
			return true, nil
		}
	}
}

func readClosersContentEqual(r1, r2 io.ReadCloser, close1, close2 bool) (bool, error) {
	var closers [2]io.Closer
	if close1 {
		closers[0] = r1
	}
	if close2 {
		closers[1] = r2
	}

	equals, err := readersContentEqual(r1, r2)
	if err != nil {
		closeQuietly(closers[:]...)
		return equals, err
	}

	if err := closeMany(closers[:]...); err != nil {
		return false, err
	}

	return equals, nil
}

func absolutePathsEqual(path1, path2 string) (bool, error) {
	abs1, err := filepath.Abs(path1)
	if err != nil {
		return false, err
	}

	abs2, err := filepath.Abs(path2)
	if err != nil {
		return false, err
	}

	return abs1 == abs2, nil
}

func pathsEqual(path1, path2 string, diffs *[]Diff) (bool, error) {
	if equal, err := absolutePathsEqual(path1, path2); err != nil {
		return false, err
	} else if equal {
		return true, nil
	}

	info1, err := os.Stat(path1)
	notExists1 := os.IsNotExist(err)
	if err != nil && !notExists1 {
		return false, err
	}

	info2, err := os.Stat(path2)
	notExists2 := os.IsNotExist(err)
	if err != nil && !notExists2 {
		return false, err
	}

	if notExists1 && notExists2 {
		return true, nil
	}

	if notExists1 != notExists2 {
		if diffs != nil {
			if notExists2 {
				*diffs = append(*diffs, Diff{Item1: &FileInfo{FileInfo: info1, FullPath: path1}})
			} else {
				*diffs = append(*diffs, Diff{Item2: &FileInfo{FileInfo: info2, FullPath: path2}})
			}
		}

		return false, nil
	}

	updateDiffs := func() {
		if diffs != nil {
			*diffs = append(*diffs, Diff{
				Item1: &FileInfo{FileInfo: info1, FullPath: path1},
				Item2: &FileInfo{FileInfo: info2, FullPath: path2},
			})
		}
	}

	if isFile(info1) {
		if isFile(info2) {
			return filesEqual(path1, path2, diffs)
		} else if isDir(info2) {
			updateDiffs()
			return false, nil
		} else {
			return false, fmt.Errorf("unsupported path type: %v", path2)
		}
	}

	if isDir(info1) {
		if isFile(info2) {
			updateDiffs()
			return false, nil
		} else if isDir(info2) {
			return dirsEqual(path1, path2, diffs)
		} else {
			return false, fmt.Errorf("unsupported path type: %v", path2)
		}
	}

	return false, fmt.Errorf("unsupported path type: %v", path1)
}

func filesEqual(path1, path2 string, diffs *[]Diff) (bool, error) {
	info1, err := checkFileOrDir(path1, false)
	if err != nil {
		return false, err
	}

	info2, err := checkFileOrDir(path2, false)
	if err != nil {
		return false, err
	}

	if equal, err := absolutePathsEqual(path1, path2); err != nil {
		return false, err
	} else if equal {
		return true, nil
	}

	file1, err := os.Open(path1)
	if err != nil {
		return false, err
	}

	file2, err := os.Open(path2)
	if err != nil {
		return false, err
	}

	if equal, err := readClosersContentEqual(file1, file2, true, true); err != nil {
		return false, err
	} else if equal {
		return true, nil
	}

	if diffs != nil {
		*diffs = append(*diffs, Diff{Item1: info1, Item2: info2})
	}

	return false, nil
}

func dirsEqual(path1, path2 string, diffs *[]Diff) (bool, error) {
	_, err := checkFileOrDir(path1, true)
	if err != nil {
		return false, err
	}

	_, err = checkFileOrDir(path2, true)
	if err != nil {
		return false, err
	}

	if equal, err := absolutePathsEqual(path1, path2); err != nil {
		return false, err
	} else if equal {
		return true, nil
	}

	infos1, err := ioutil.ReadDir(path1)
	if err != nil {
		return false, err
	}

	infos2, err := ioutil.ReadDir(path2)
	if err != nil {
		return false, err
	}

	if diffs == nil && len(infos1) != len(infos2) {
		return false, nil
	}

	var itemDiffs []Diff
	var pos1, pos2 int

	for pos1 < len(infos1) && pos2 < len(infos2) {
		itemInfo1 := infos1[pos1]
		itemInfo2 := infos2[pos2]

		itemPath1 := filepath.Join(path1, itemInfo1.Name())
		itemPath2 := filepath.Join(path2, itemInfo2.Name())

		if itemInfo1.Name() != itemInfo2.Name() {
			if diffs == nil {
				return false, nil
			}

			if itemInfo1.Name() < itemInfo2.Name() {
				itemDiffs = append(itemDiffs, Diff{Item1: &FileInfo{FileInfo: itemInfo1, FullPath: itemPath1}})
				pos1++
			} else {
				itemDiffs = append(itemDiffs, Diff{Item2: &FileInfo{FileInfo: itemInfo2, FullPath: itemPath2}})
				pos2++
			}

			continue
		}

		pos1++
		pos2++

		if equal, err := pathsEqual(itemPath1, itemPath2, &itemDiffs); err != nil {
			return false, err
		} else if !equal {
			if diffs == nil {
				return false, nil
			}

			continue
		}
	}

	if diffs != nil {
		for pos1 < len(infos1) {
			itemInfo1 := infos1[pos1]
			itemPath1 := filepath.Join(path1, itemInfo1.Name())
			itemDiffs = append(itemDiffs, Diff{Item1: &FileInfo{FileInfo: itemInfo1, FullPath: itemPath1}})
			pos1++
		}

		for pos2 < len(infos2) {
			itemInfo2 := infos2[pos2]
			itemPath2 := filepath.Join(path2, itemInfo2.Name())
			itemDiffs = append(itemDiffs, Diff{Item2: &FileInfo{FileInfo: itemInfo2, FullPath: itemPath2}})
			pos2++
		}

		if len(itemDiffs) > 0 {
			*diffs = append(*diffs, itemDiffs...)
			return false, nil
		}
	}

	return true, nil
}
