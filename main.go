package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"Auto_GotoFlowchart/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("❌ 用法: go run main.go <Go源文件路径>")
		os.Exit(1)
	}

	inputPath := os.Args[1]

	// 1. 校验文件是否存在
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		fmt.Printf("❌ 文件不存在: %s\n", inputPath)
		os.Exit(1)
	}

	// 2. 构建 AST 流程图
	graph := parser.BuildFlowGraph(inputPath)

	// 3. 确保 output 目录存在
	outputDir := "output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("❌ 创建输出目录失败: %v\n", err)
		os.Exit(1)
	}

	// 4. 生成 .dot 文件名（保留原文件名）
	baseName := filepath.Base(inputPath)
	baseNoExt := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	outputFile := filepath.Join(outputDir, baseNoExt+".dot")

	// 5. 写入 DOT 文件
	if err := parser.WriteDOT(graph, outputFile); err != nil {
		fmt.Printf("❌ 写入 DOT 文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ 成功生成 DOT 文件: %s\n", outputFile)
}
