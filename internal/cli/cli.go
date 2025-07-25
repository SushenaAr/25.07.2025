package cli

import (
	"awesomeProject1/internal/sshutil"
	"awesomeProject1/internal/tarutil"
	"awesomeProject1/internal/version"
	"fmt"
	"path/filepath"
)

type PackageSpec struct {
	Name string `json:"name"`
	Ver  string `json:"ver,omitempty"`
}

type UpdateFile struct {
	Packages []PackageSpec `json:"packages"`
}

func Process(name, ver string, files []string, visited map[string]bool, cfg *sshutil.Config, remotePath string) error {
	//1. Помечаем файл пройденным
	key := name + ":" + ver
	if visited[key] {
		//пропускаем пройденные пакеты
		return nil
	}
	visited[key] = true

	//2. Ищем лучший подходящий архив по имени и constraint (версия)
	if ver == "" {
		ver = ">=0"
	}
	match, dir, err := version.FindBestMatch(files, name, ver)
	if err != nil {
		return fmt.Errorf("%s %s: %w", name, ver, err)
	}

	// Пути для скачивания и распаковки
	downloadPath := filepath.Join("downloads", match)
	unpackPath := filepath.Join("unpacked", dir)
	fmt.Println("Downloading", match)
	if err := sshutil.DownloadFile(*cfg, remotePath+"/"+match, downloadPath); err != nil {
		return err
	}

	fmt.Println("Unpacking to", unpackPath)
	deps, err := tarutil.UnpackArchive(downloadPath, unpackPath)
	if err != nil {
		return err
	}

	//3. Обрабатываем зависимости из dependencies.json
	for _, d := range deps {
		if err := Process(d.Name, d.Ver, files, visited, cfg, remotePath); err != nil {
			return err
		}
	}
	return nil
}
