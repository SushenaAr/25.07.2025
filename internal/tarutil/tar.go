package tarutil

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

type Dependency struct {
	Name string `json:"name"`
	Ver  string `json:"ver,omitempty"`
}

// UnpackArchive распаковывает .tar.gz архив в директорию.
// Возвращает список зависимостей из файла dependencies.json (если он есть).
func UnpackArchive(srcPath, destDir string) ([]Dependency, error) {
	//1. Открываем архивный файл
	f, err := os.Open(srcPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	//2. Создаём целевую директорию, если она ещё не существует
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, err
	}
	//3. Создаём gzip ридер
	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}

	//4. Создаём tar ридер
	tr := tar.NewReader(gz)
	// Массив зависимостей из dependencies.json
	var deps []Dependency

	//5. Проходим по содержимому архива
	for {
		hdr, err := tr.Next()
		if err != nil {
			break // конец архива
		}
		// Полный путь, куда распаковывать файл или директорию
		path := filepath.Join(destDir, hdr.Name)
		switch hdr.Typeflag {
		case tar.TypeDir:
			// Если это директория — создаём
			//вроде сюда не входит, но я оставлю
			err := os.MkdirAll(path, 0755)
			if err != nil {
				return nil, err
			}

		case tar.TypeReg:
			// Если обычный файл — создаём и копируем содержимое
			//сразу создаю папку, если в первое не попало
			err := os.MkdirAll(filepath.Dir(path), 0755)
			if err != nil {
				return nil, err
			}
			//создаю файл
			out, err := os.Create(path)
			if err != nil {
				return nil, err
			}
			io.Copy(out, tr)
			out.Close()

			// Если это файл dependencies.json — парсим его
			if filepath.Base(hdr.Name) == "dependencies.json" {
				data, _ := os.ReadFile(path)
				json.Unmarshal(data, &deps)
			}
		}
	}
	return deps, nil
}
