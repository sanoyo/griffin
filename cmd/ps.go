/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

// psCmd represents the ps command
var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// tview.TextViewウィジェットを作成
		textView := tview.NewTextView().
			SetDynamicColors(true).
			SetRegions(true).
			SetWrap(false).
			SetWordWrap(false)

		// tview.Applicationを作成
		tviewApp := tview.NewApplication()

		// TextViewにキーハンドラを設定
		textView.SetChangedFunc(func() {
			tviewApp.Draw()
		})

		// psコマンドの出力を定期的に更新するゴルーチンを起動
		go func() {
			for {
				psOutput, err := getPsOutput()
				if err != nil {
					fmt.Printf("Error executing ps command: %v\n", err)
					return
				}
				colorizedOutput := colorizePsOutput(psOutput)

				// TextViewの内容を更新
				tviewApp.QueueUpdateDraw(func() {
					textView.SetText(colorizedOutput)
				})

				// 1秒間隔で更新
				time.Sleep(1 * time.Second)
			}
		}()

		// Flexレイアウトを作成し、TextViewを追加
		flex := tview.NewFlex().
			AddItem(textView, 0, 1, false)

		// Applicationを起動
		if err := tviewApp.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(psCmd)
}

// psコマンドを実行して出力を取得する関数
func getPsOutput() (string, error) {
	cmd := exec.Command("ps")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// psコマンドの出力に色を付ける関数
func colorizePsOutput(psOutput string) string {
	lines := strings.Split(psOutput, "\n")
	if len(lines) == 0 {
		return psOutput
	}

	header := lines[0]
	body := lines[1:]

	// ヘッダーに色を付ける
	colorizedHeader := fmt.Sprintf("[yellow]%s[white]", header)

	// ボディに色を付ける
	var colorizedBody []string
	for _, line := range body {
		if strings.TrimSpace(line) == "" {
			continue
		}
		colorizedBody = append(colorizedBody, fmt.Sprintf("[green]%s[white]", line))
	}

	return fmt.Sprintf("%s\n%s", colorizedHeader, strings.Join(colorizedBody, "\n"))
}
