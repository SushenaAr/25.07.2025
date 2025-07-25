package sshutil

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
)

func UploadFile(cfg Config, localPath, remotePath string) error {
	// 1. Подключение по SSH
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	sshConfig := &ssh.ClientConfig{
		User:            cfg.User,
		Auth:            []ssh.AuthMethod{ssh.Password(cfg.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // для тестов
	}
	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 2. Создаем SFTP
	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	// 3. Открываем локальный файл
	srcFile, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 4. Создаем файл на сервере
	dstFile, err := sftpClient.Create(remotePath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 5. Копируем
	_, err = io.Copy(dstFile, srcFile)
	return err
}
