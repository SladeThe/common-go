package io

import (
	"io"
)

const (
	BufferSize = bufferSize

	SeekStart   = io.SeekStart
	SeekCurrent = io.SeekCurrent
	SeekEnd     = io.SeekEnd
)

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

// Returns true iff the path entry exists and is one of:
// 1. A directory with no children.
// 2. A regular file of zero length.
func IsEmpty(path string) (bool, error) {
	return isEmpty(path)
}

func Parent(path string) (string, error) {
	return parent(path)
}

// Returns true iff the parent of this path exists and is a directory.
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

// ----------------------------------------------------------------------------------------------------

func WriteString(w io.Writer, s string) (n int, err error) {
	return io.WriteString(w, s)
}

func ReadAtLeast(r io.Reader, buf []byte, min int) (n int, err error) {
	return io.ReadAtLeast(r, buf, min)
}

func ReadFull(r io.Reader, buf []byte) (n int, err error) {
	return io.ReadFull(r, buf)
}

func CopyN(dst io.Writer, src io.Reader, n int64) (written int64, err error) {
	return io.CopyN(dst, src, n)
}

func Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return io.Copy(dst, src)
}

func CopyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
	return io.CopyBuffer(dst, src, buf)
}

func LimitReader(r io.Reader, n int64) io.Reader {
	return io.LimitReader(r, n)
}
