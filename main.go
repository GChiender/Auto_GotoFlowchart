package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"Auto_GotoFlowchart/parser"
	"Auto_GotoFlowchart/tools/writer"
)

func main() {
	inputPath := os.Getenv("GO_INPUT_FILE")
	if inputPath == "" {
		inputPath = "target/demo.go"
	}

	base := strings.TrimSuffix(filepath.Base(inputPath), ".go")
	dotPath := fmt.Sprintf("output/dot/%s.dot", base)
	drawioPath := fmt.Sprintf("output/drawio/%s.drawio", base)

	// Step 1: 构建流程图 Graph
	graph := parser.BuildFlowGraph(inputPath)

	// Step 2: 写入 .dot 文件
	os.MkdirAll(filepath.Dir(dotPath), os.ModePerm)
	dot := writer.WriteDOT(graph)
	err := os.WriteFile(dotPath, []byte(dot), 0644)
	if err != nil {
		fmt.Println("❌ 写入 .dot 文件失败:", err)
		return
	}
	fmt.Println("✅ 生成 DOT 文件:", dotPath)

	// Step 3: 调用 Python 脚本转换为 drawio
	cmd := exec.Command("python3", "tools/convert_dot_to_drawio.py", "--input", dotPath, "--output", drawioPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("❌ Python 脚本执行失败:", err)
		return
	}
	fmt.Println("✅ 已生成 drawio 文件:", drawioPath)

	// Step 4: 自动打开 draw.io 查看
	err = exec.Command("open", drawioPath).Start() // macOS: 使用 "open"
	// err = exec.Command("xdg-open", drawioPath).Start() // Linux
	// err = exec.Command("cmd", "/C", "start", drawioPath).Start() // Windows
	if err != nil {
		fmt.Println("⚠️ 无法自动打开 drawio 文件:", err)
		return
	}

	fmt.Println("✅ 已尝试打开 draw.io 查看图形文件")
}
