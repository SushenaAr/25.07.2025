package main

import (
	"awesomeProject1/internal/runFunc"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pm",
	Short: "Packet Manager CLI",
}

var createCmd = &cobra.Command{
	Use:   "create <file>",
	Short: "Создать архив из описания",
	Args:  cobra.ExactArgs(1),
	Run:   runFunc.Create,
}

var updateCmd = &cobra.Command{
	Use:   "update <file>",
	Short: "Загрузить и распаковать пакеты",
	Args:  cobra.ExactArgs(1),
	Run:   runFunc.Update,
}

// Мне стало интересно написать пакетный менелджер с полного 0
// написать комменты, зависимости единая точка входа прописать
func main() {
	rootCmd.AddCommand(createCmd, updateCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Ошибка:", err)
		os.Exit(1)
	}
}
