package filepicker

import (
	"os"
	"path/filepath"
)

func CollectFiles(pattern string, exclude string) ([]string, error) {
	files, err := filepath.Glob(pattern) //picking files
	if err != nil {
		return nil, err
	}

	var result []string
	for _, f := range files { //excluding files
		if exclude != "" && match(exclude, f) { //if excluding turn on and match
			continue
		}

		info, err := os.Stat(f)
		if err == nil && !info.IsDir() {
			result = append(result, f)
		}
	}
	return result, nil
}

func match(pattern, path string) bool {
	base := filepath.Base(path)
	ok, _ := filepath.Match(pattern, base)
	return ok
}
