package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
)

type Graph struct {
	Nodes map[string]string // 节点 ID -> 标签
	Edges [][2]string       // 有向边 [from, to]
}

var nodeCounter int

// 入口函数：构建流程图
func BuildFlowGraph(path string) Graph {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	graph := Graph{
		Nodes: make(map[string]string),
		Edges: make([][2]string, 0),
	}

	for _, decl := range node.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			visitNode(&graph, fn, "")
		}
	}

	return graph
}

// 遍历语义节点
func visitNode(graph *Graph, node ast.Node, parentID string) {
	if node == nil {
		return
	}

	var label string
	var children []ast.Stmt

	switch n := node.(type) {
	case *ast.FuncDecl:
		label = "Function: " + n.Name.Name
		children = n.Body.List
	case *ast.AssignStmt:
		label = prettyPrint(n)
	case *ast.ReturnStmt:
		label = "return " + prettyPrint(n.Results[0])
	case *ast.ExprStmt:
		label = prettyPrint(n.X)
	case *ast.IfStmt:
		label = "if " + prettyPrint(n.Cond)
	case *ast.ForStmt:
		label = "for " + prettyPrint(n.Cond)
	default:
		// 忽略非语义结构
		return
	}

	// 添加节点和边
	nodeID := addNode(graph, label, parentID)

	// 递归处理语句体
	switch n := node.(type) {
	case *ast.FuncDecl:
		for _, stmt := range children {
			visitNode(graph, stmt, nodeID)
		}
	case *ast.IfStmt:
		for _, stmt := range n.Body.List {
			visitNode(graph, stmt, nodeID)
		}
		if n.Else != nil {
			if elseBlock, ok := n.Else.(*ast.BlockStmt); ok {
				for _, stmt := range elseBlock.List {
					visitNode(graph, stmt, nodeID)
				}
			}
		}
	case *ast.ForStmt:
		for _, stmt := range n.Body.List {
			visitNode(graph, stmt, nodeID)
		}
	}
}

func addNode(graph *Graph, label string, parentID string) string {
	nodeID := genNodeID()
	graph.Nodes[nodeID] = label
	if parentID != "" {
		graph.Edges = append(graph.Edges, [2]string{parentID, nodeID})
	}
	return nodeID
}

func genNodeID() string {
	nodeCounter++
	return fmt.Sprintf("n%d", nodeCounter)
}

func prettyPrint(node ast.Node) string {
	var buf bytes.Buffer
	_ = format.Node(&buf, token.NewFileSet(), node)
	return buf.String()
}
