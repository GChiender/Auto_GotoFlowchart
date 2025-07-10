package parser

import (
	"fmt"
	"os"
	"strings"
)

func WriteDOT(graph Graph, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, _ = file.WriteString("digraph G {\n")

	// 写节点
	for id, label := range graph.Nodes {
		safeLabel := escapeLabel(label)
		_, _ = file.WriteString(fmt.Sprintf("    %s [label=\"%s\"];\n", id, safeLabel))
	}

	// 写边
	for _, edge := range graph.Edges {
		_, _ = file.WriteString(fmt.Sprintf("    %s -> %s;\n", edge[0], edge[1]))
	}

	_, _ = file.WriteString("}\n")
	return nil
}

func escapeLabel(label string) string {
	label = strings.ReplaceAll(label, "\"", "\\\"")
	label = strings.ReplaceAll(label, "\n", "\\n")
	return label
}
