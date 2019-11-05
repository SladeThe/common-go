package io

import (
	"io"
)

const BufferSize = bufferSize

var (
	ErrShortWrite    = io.ErrShortWrite
	ErrShortBuffer   = io.ErrShortBuffer
	EOF              = io.EOF
	ErrUnexpectedEOF = io.ErrUnexpectedEOF
	ErrNoProgress    = io.ErrNoProgress
)

func Exists(path string) (bool, error) {
	return exists(path)
}

func IsFile(path string) (bool, error) {
	return isFileOrDir(path, false)
}

func IsDir(path string) (bool, error) {
	return isFileOrDir(path, true)
}

func Parent(path string) (string, error) {
	return parent(path)
}

func IsParentDir(path string) (bool, error) {
	return isParentDir(path)
}

func Close(closers ...io.Closer) error {
	return closeMany(closers...)
}

func CloseQuietly(closers ...io.Closer) {
	closeQuietly(closers...)
}

func CloseOr(handleErr func(err error), closers ...io.Closer) {
	closeOr(handleErr, closers...)
}

func CloseOrPanic(closers ...io.Closer) {
	CloseOr(func(err error) {
		panic(err)
	}, closers...)
}

func FilesEqual(path1, path2 string) (equal bool, err error) {
	return filesEqual(path1, path2, nil)
}

func DirsEqual(path1, path2 string) (equal bool, err error) {
	return dirsEqual(path1, path2, nil)
}

func DiffDirs(path1, path2 string) (diffs []Diff, err error) {
	_, err = dirsEqual(path1, path2, &diffs)
	return
}
