package runFunc

import (
	"awesomeProject1/internal/cli"
	"awesomeProject1/internal/sshutil"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func Update(cmd *cobra.Command, args []string) {
	fmt.Println("Обновляем по:", args[0])
	//1. Читаю файл packages.json
	content, err := os.ReadFile(args[0])
	if err != nil {
		log.Fatal(err)
	}

	//2. Паршу файл package.json
	var spec cli.UpdateFile
	if err := json.Unmarshal(content, &spec); err != nil {
		log.Fatal(err)
	}

	//3. Список всех архивов на сервере
	remoteFiles, err := sshutil.ListRemoteFiles(*sshutil.SSHCfg, sshutil.RemotePath)
	if err != nil {
		log.Fatal(err)
	}
	//4. Скачиваю пакеты.
	//packet-3:<=2.0" - ключ.
	visited := map[string]bool{} //чтобы не зацикливаться, посещенные пакеты
	for _, pkg := range spec.Packages {
		if err := cli.Process(pkg.Name, pkg.Ver, remoteFiles, visited, sshutil.SSHCfg, sshutil.RemotePath); err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("ok")
}
