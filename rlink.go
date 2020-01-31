package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type Main struct {
	relative bool
	force bool
	dryrun bool
	ignore *regexp.Regexp
}

func New(relative, force, dryrun bool, ignore *regexp.Regexp) Main {
	return Main{
		relative: relative,
		force:force,
		dryrun:dryrun,
		ignore: ignore,
	}
}

func (m Main) MkLinks(srcRoot, dstRoot string) error {
	return filepath.Walk(srcRoot, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || m.ignore.MatchString(path) {
			return nil
		}

		rel, err := filepath.Rel(srcRoot, path)
		if err != nil {
			return fmt.Errorf("failed to filepath.Rel of %s, %s: %w", srcRoot, path, err)
		}

		src, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to filepath.Abs of %s: %w", path, err)
		}

		dst, err := filepath.Abs(filepath.Join(dstRoot, rel))
		if err != nil {
			return fmt.Errorf("failed to filepath.Abs of %s: %w", filepath.Join(dstRoot, rel), err)
		}

		if err := m.mklink(src, dst); err != nil {
			return fmt.Errorf("failed to mklink: %w", err)
		}

		return nil
	})
}

func (m Main) mklink(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("failed to MkdirAll %s: %w", filepath.Dir(dst), err)
	}

	var path string
	if m.relative {
		var err error
		path, err = filepath.Rel(filepath.Dir(dst), src)
		if err != nil {
			return fmt.Errorf("failed to filepath.Rel of %s, %s: %w", dst, src, err)
		}
	} else {
		path = src
	}

	if m.force {
		if !m.dryrun {
			os.Remove(dst)
		}
	} else if exists(dst) {
		return nil
	}

	if !m.dryrun {
		if err := os.Symlink(path, dst); err != nil {
			return fmt.Errorf("failed to create symlink. old: %s, new:%s. : %w", path, dst, err)
		}
	}

	fmt.Printf("created %s to %s\n", dst, path)

	return nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
