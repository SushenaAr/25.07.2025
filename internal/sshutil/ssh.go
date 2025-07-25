package sshutil

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
}

var SSHCfg = &Config{
	Host:     "localhost",
	Port:     22,
	User:     "devuser",
	Password: "MyPass123",
}
var RemotePath = "C:/Users/devuser/packages"

// newSFTPClient устанавливает SSH-соединение и создаёт SFTP-клиент для работы с файлами на сервере.
func newSFTPClient(cfg Config) (*sftp.Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	sshConf := &ssh.ClientConfig{
		User:            cfg.User,
		Auth:            []ssh.AuthMethod{ssh.Password(cfg.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	//1. Устанавливаем SSH-соединение
	conn, err := ssh.Dial("tcp", addr, sshConf)
	if err != nil {
		return nil, err
	}

	//2. Создаём SFTP-клиент
	return sftp.NewClient(conn)
}

// DownloadFile скачивает файл с удалённого сервера (remotePath) на локальную машину (localPath).
func DownloadFile(cfg Config, remotePath, localPath string) error {
	client, err := newSFTPClient(cfg)
	if err != nil {
		return err
	}
	defer client.Close()

	//1. Открываем удалённый файл для чтения
	rf, err := client.Open(remotePath)
	if err != nil {
		return err
	}

	//2. Создаём локальный файл для записи
	lf, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer lf.Close()

	//3. Копируем содержимое файла
	_, err = io.Copy(lf, rf)
	return err
}

// ListRemoteFiles возвращает список файлов в удалённой директории (без директорий).
func ListRemoteFiles(cfg Config, dir string) ([]string, error) {
	client, err := newSFTPClient(cfg)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	//1. Читаем содержимое директории (без директорий)
	entries, err := client.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}
