package io

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

var junctionImpl = func() func(directory, junction string) error {
	if runtime.GOOS == "windows" {
		// Ignore error, because mklink returns 1 after printing help.
		output, _ := exec.Command("cmd", "/c", "mklink", "/?").Output()

		if strings.Contains(string(output), " /J ") {
			return func(directory, junction string) error {
				_, err := exec.Command("cmd", "/c", "mklink", "/J", junction, directory).CombinedOutput()
				return err
			}
		}
	}

	// TODO see os_windows_test.go: createMountPoint for more impl options

	return func(directory, junction string) error {
		return fmt.Errorf("junctions are not supported on this OS/ARCH")
	}
}()

func Junction(directory, junction string) error {
	if isDir, err := IsDir(directory); err != nil {
		return err
	} else if !isDir {
		return fmt.Errorf("not a directory: %s", directory)
	}

	if exists, err := Exists(junction); err != nil {
		return err
	} else if exists {
		return fmt.Errorf("must not exist: %s", junction)
	}

	return junctionImpl(directory, junction)
}
