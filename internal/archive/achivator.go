package archive

import (
	"archive/tar"
	"awesomeProject1/internal/model"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CreateTarGzWithDeps создает .tar.gz архив из файлов + dependencies.json
func CreateTarGzWithDeps(outputPath string, files []string, dependencies []model.Dependency) error {
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("ошибка создания архива: %w", err)
	}
	defer outFile.Close()

	gzipWriter := gzip.NewWriter(outFile)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Добавляем обычные файлы
	for _, file := range files {
		err := addFileToTarWriter(file, tarWriter)
		if err != nil {
			return fmt.Errorf("ошибка добавления файла %s: %w", file, err)
		}
	}

	// Добавляем dependencies.json
	err = addDependenciesFile(tarWriter, dependencies)
	if err != nil {
		return fmt.Errorf("ошибка добавления dependencies.json: %w", err)
	}

	return nil
}

func addFileToTarWriter(filePath string, tarWriter *tar.Writer) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	header.Name = filepath.ToSlash(filePath)
	header.Name = header.Name[strings.Index(filePath, "\\")+1:]
	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(tarWriter, file)
	return err
}

func addDependenciesFile(tarWriter *tar.Writer, dependencies []model.Dependency) error {
	data, err := json.MarshalIndent(dependencies, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации зависимостей: %w", err)
	}

	buf := bytes.NewBuffer(data)
	header := &tar.Header{
		Name: "dependencies.json",
		Mode: 0644,
		Size: int64(buf.Len()),
	}

	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(tarWriter, buf)
	return err
}
