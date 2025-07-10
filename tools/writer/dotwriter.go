package writer

import (
	"Auto_GotoFlowchart/parser" // 替换成你实际的 module 路径
	"fmt"
	"strings"
)

func WriteDOT(graph parser.Graph) string {
	var sb strings.Builder

	sb.WriteString("digraph G {\n")
	sb.WriteString("  rankdir=TB;\n") // 从上到下
	sb.WriteString("  node [fontname=\"Helvetica\"];\n\n")

	// 写节点
	for id, node := range graph.Nodes {
		label := escapeLabel(node.Label)
		shape := node.Shape
		sb.WriteString(fmt.Sprintf("  %s [label=\"%s\", shape=%s];\n", id, label, shape))
	}

	sb.WriteString("\n")

	// 写边
	for _, edge := range graph.Edges {
		if edge.Label != "" {
			sb.WriteString(fmt.Sprintf("  %s -> %s [label=\"%s\"];\n", edge.From, edge.To, edge.Label))
		} else {
			sb.WriteString(fmt.Sprintf("  %s -> %s;\n", edge.From, edge.To))
		}
	}

	sb.WriteString("}\n")
	return sb.String()
}

// 替换双引号/换行，防止 DOT label 报错
func escapeLabel(s string) string {
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	return s
}
