package runFunc

import (
	"awesomeProject1/internal/archive"
	"awesomeProject1/internal/filepicker"
	"awesomeProject1/internal/model"
	"awesomeProject1/internal/sshutil"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

const prefixPath = "packages/"

func Create(cmd *cobra.Command, args []string) {
	// 1. Читаем файл
	data, err := os.ReadFile(args[0])
	if err != nil {
		log.Fatalf("ошибка чтения файла: %v", err)
	}

	// 2. Парсим JSON в Packet
	var p *model.Packet
	if err := json.Unmarshal(data, &p); err != nil {
		log.Fatalf("ошибка парсинга JSON: %v", err)
	}
	fileName := p.Name + "-" + p.Version + ".tar.gz"

	// 3. Выводим результат
	printPacket(p)

	//4. Получаем набор путей с нужными файлами
	paths := getPathsFiles(p)

	//5. Упаковываем
	err = archive.CreateTarGzWithDeps(fileName, paths, p.Packets)
	if err != nil {
		return
	}

	//6. Загружаем на сервер
	err = sshutil.UploadFile(*sshutil.SSHCfg, fileName, sshutil.RemotePath+"/"+fileName)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("ok")
}

func getPathsFiles(p *model.Packet) []string {
	paths := make([]string, 0)
	for _, target := range p.Targets {
		target.Path = prefixPath + target.Path[2:]
		files, err := filepicker.CollectFiles(target.Path, target.Exclude)
		if err != nil {
			log.Fatal(err.Error())
		}
		paths = append(paths, files...)
	}
	return paths
}

func printPacket(p *model.Packet) {
	fmt.Printf("Пакет: %s (версия %s)\n", p.Name, p.Version)
	fmt.Println("Цели архивации:")
	for _, t := range p.Targets {
		fmt.Printf("- Путь: %s", t.Path)
		if t.Exclude != "" {
			fmt.Printf(" (исключить: %s)", t.Exclude)
		}
		fmt.Println()
	}
	for _, dep := range p.Packets {
		fmt.Printf("Зависимость: %s", dep.Name)
		if dep.Version != "" {
			fmt.Printf(" (версия %s)", dep.Version)
		}
		fmt.Println()
	}
}
