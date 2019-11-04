package io

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirsEqual_SamePath(t *testing.T) {
	equal, err := DirsEqual(diffPath("same"), diffPath("same"))
	assert.Nil(t, err)
	assert.True(t, equal)
}

func TestDirsEqualA_Equal(t *testing.T) {
	equal, err := DirsEqual(diffPath("a1"), diffPath("a2"))
	assert.Nil(t, err)
	assert.True(t, equal)
}

func TestDirsEqualB_DifferentContent(t *testing.T) {
	equal, err := DirsEqual(diffPath("b1"), diffPath("b2"))
	assert.Nil(t, err)
	assert.False(t, equal)
}

func TestDirsEqualC_DifferentContentStructure(t *testing.T) {
	equal, err := DirsEqual(diffPath("c1"), diffPath("c2"))
	assert.Nil(t, err)
	assert.False(t, equal)
}

func TestDiffDirs_SamePath(t *testing.T) {
	diffs, err := DiffDirs(diffPath("same"), diffPath("same"))
	assert.Nil(t, err)
	assert.Zero(t, len(diffs))
}

func TestDiffDirsA_Equal(t *testing.T) {
	diffs, err := DiffDirs(diffPath("a1"), diffPath("a2"))
	assert.Nil(t, err)
	assert.Zero(t, len(diffs))
}

func TestDiffDirsB_DifferentContent(t *testing.T) {
	diffs, err := DiffDirs(diffPath("b1"), diffPath("b2"))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(diffs))

	diff := diffs[0]
	assert.Equal(t, "0.bin", diff.Item1.Name())
	assert.Equal(t, "0.bin", diff.Item2.Name())
	assert.Equal(t, diffPath("b1/diff/diff/0.bin"), diff.Item1.FullPath)
	assert.Equal(t, diffPath("b2/diff/diff/0.bin"), diff.Item2.FullPath)
	assert.True(t, diff.Item1.IsFile())
	assert.True(t, diff.Item2.IsFile())
	assert.False(t, diff.Item1.IsDir())
	assert.False(t, diff.Item2.IsDir())
}

func TestDiffDirsC_DifferentContentStructure(t *testing.T) {
	diffs, err := DiffDirs(diffPath("c1"), diffPath("c2"))
	assert.Nil(t, err)
	assert.Equal(t, 15, len(diffs))

	{
		diff := diffs[0]
		assert.Equal(t, diffPath("c1/s0"), diff.Item1.FullPath)
		assert.Nil(t, diff.Item2)
		assert.False(t, diff.Item1.IsFile())
		assert.True(t, diff.Item1.IsDir())
	}

	{
		diff := diffs[1]
		assert.Nil(t, diff.Item1)
		assert.Equal(t, diffPath("c2/s1/S2"), diff.Item2.FullPath)
		assert.False(t, diff.Item2.IsFile())
		assert.True(t, diff.Item2.IsDir())
	}

	{
		diff := diffs[2]
		assert.Equal(t, "0.bin", diff.Item1.Name())
		assert.Equal(t, "0.bin", diff.Item2.Name())
		assert.Equal(t, diffPath("c1/s1/s1/0.bin"), diff.Item1.FullPath)
		assert.Equal(t, diffPath("c2/s1/s1/0.bin"), diff.Item2.FullPath)
		assert.True(t, diff.Item1.IsFile())
		assert.True(t, diff.Item2.IsFile())
		assert.False(t, diff.Item1.IsDir())
		assert.False(t, diff.Item2.IsDir())
	}

	{
		diff := diffs[3]
		assert.Equal(t, diffPath("c1/s1/s2"), diff.Item1.FullPath)
		assert.Nil(t, diff.Item2)
		assert.False(t, diff.Item1.IsFile())
		assert.True(t, diff.Item1.IsDir())
	}

	{
		diff := diffs[4]
		assert.Nil(t, diff.Item1)
		assert.Equal(t, diffPath("c2/s2"), diff.Item2.FullPath)
		assert.False(t, diff.Item2.IsFile())
		assert.True(t, diff.Item2.IsDir())
	}

	{
		diff := diffs[5]
		assert.Equal(t, diffPath("c1/s3"), diff.Item1.FullPath)
		assert.Nil(t, diff.Item2)
		assert.False(t, diff.Item1.IsFile())
		assert.True(t, diff.Item1.IsDir())
	}

	{
		diff := diffs[6]
		assert.Nil(t, diff.Item1)
		assert.Equal(t, diffPath("c2/s4"), diff.Item2.FullPath)
		assert.False(t, diff.Item2.IsFile())
		assert.True(t, diff.Item2.IsDir())
	}

	{
		diff := diffs[7]
		assert.Nil(t, diff.Item1)
		assert.Equal(t, diffPath("c2/s5/aa"), diff.Item2.FullPath)
	}

	{
		diff := diffs[8]
		assert.Equal(t, diffPath("c1/s5/s0/0.bin"), diff.Item1.FullPath)
		assert.Nil(t, diff.Item2)
	}

	{
		diff := diffs[9]
		assert.Nil(t, diff.Item1)
		assert.Equal(t, diffPath("c2/s5/s0/1.bin"), diff.Item2.FullPath)
	}

	{
		diff := diffs[10]
		assert.Equal(t, diffPath("c1/s5/s2/0.bin"), diff.Item1.FullPath)
		assert.Equal(t, diffPath("c2/s5/s2/0.bin"), diff.Item2.FullPath)
		assert.True(t, diff.Item1.IsFile())
		assert.True(t, diff.Item2.IsFile())
	}

	{
		diff := diffs[11]
		assert.Nil(t, diff.Item1)
		assert.Equal(t, diffPath("c2/s5/s2/1.bin"), diff.Item2.FullPath)
		assert.True(t, diff.Item2.IsFile())
	}

	{
		diff := diffs[12]
		assert.Equal(t, diffPath("c1/s5/s3/0.bin"), diff.Item1.FullPath)
		assert.Equal(t, diffPath("c2/s5/s3/0.bin"), diff.Item2.FullPath)
		assert.True(t, diff.Item1.IsFile())
		assert.True(t, diff.Item2.IsDir())
	}

	{
		diff := diffs[13]
		assert.Nil(t, diff.Item1)
		assert.Equal(t, diffPath("c2/s5/zz"), diff.Item2.FullPath)
	}

	{
		diff := diffs[14]
		assert.Nil(t, diff.Item1)
		assert.Equal(t, diffPath("c2/s6"), diff.Item2.FullPath)
	}
}

func testPath() string {
	for _, path := range []string{"io/test", "test"} {
		if ok, err := IsDir(path); err != nil {
			panic(err)
		} else if ok {
			return path
		}
	}

	panic(errors.New("failed to find test path"))
}

func diffPath(path string) string {
	return filepath.Join(testPath(), "diffs", path)
}
