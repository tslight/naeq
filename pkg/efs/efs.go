package efs

import (
	"embed"
	"io/fs"
	"path/filepath"
	"strings"
)

func GetPaths(efs *embed.FS) (files []string, err error) {
	if err := fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	}); err != nil {
		return nil, err
	}
	return files, nil
}

func GetBaseNamesSansExt(efs *embed.FS) ([]string, error) {
	var baseNamesSansExt []string
	paths, err := GetPaths(efs)
	if err != nil {
		return nil, err
	}
	for _, v := range paths {
		name := strings.TrimSuffix(filepath.Base(v), filepath.Ext(v))
		baseNamesSansExt = append(baseNamesSansExt, name)
	}
	return baseNamesSansExt, nil
}
