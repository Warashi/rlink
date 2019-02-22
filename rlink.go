package main

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"
)

func MkLinks(srcRoot, dstRoot string, relative, force bool) error {
	return filepath.Walk(srcRoot, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(srcRoot, path)
		if err != nil {
			return xerrors.Errorf("failed to filepath.Rel of %s, %s: %w", srcRoot, path, err)
		}

		src, err := filepath.Abs(path)
		if err != nil {
			return xerrors.Errorf("failed to filepath.Abs of %s: %w", path, err)
		}

		dst, err := filepath.Abs(filepath.Join(dstRoot, rel))
		if err != nil {
			return xerrors.Errorf("failed to filepath.Abs of %s: %w", filepath.Join(dstRoot, rel), err)
		}

		if err := mklink(src, dst, relative, force); err != nil {
			return xerrors.Errorf("failed to mklink: %w", err)
		}

		return nil
	})
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func mklink(srcAbs, dstAbs string, relative, force bool) error {
	if err := os.MkdirAll(filepath.Dir(dstAbs), 0755); err != nil {
		return xerrors.Errorf("failed to MkdirAll %s: %w", filepath.Dir(dstAbs), err)
	}

	var path string
	if relative {
		var err error
		path, err = filepath.Rel(filepath.Dir(dstAbs), srcAbs)
		if err != nil {
			return xerrors.Errorf("failed to filepath.Rel of %s, %s: %w", dstAbs, srcAbs, err)
		}
	} else {
		path = srcAbs
	}

	if force {
		os.Remove(dstAbs)
	} else if exists(dstAbs) {
		return nil
	}

	if err := os.Symlink(path, dstAbs); err != nil {
		return xerrors.Errorf("failed to create symlink. old: %s, new:%s. : %w", path, dstAbs, err)
	}

	fmt.Printf("created %s to %s\n", dstAbs, path)

	return nil
}
