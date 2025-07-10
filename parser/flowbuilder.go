package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
)

type Node struct {
	Label string
	Shape string // "box", "diamond", "ellipse"
}

type Edge struct {
	From  string
	To    string
	Label string // 可用于 Yes/No/back
}

type Graph struct {
	Nodes map[string]Node
	Edges []Edge
}

var nodeCounter int

// BuildFlowGraph 构建完整流程图
func BuildFlowGraph(path string) Graph {
	nodeCounter = 0 // 每次构建时重置计数器
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	graph := Graph{
		Nodes: make(map[string]Node),
		Edges: []Edge{},
	}

	startID := addNode(&graph, "Start", "ellipse", "", false)
	endID := addNode(&graph, "End", "ellipse", "", false)

	for _, decl := range astFile.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			funcID := addNode(&graph, "Function: "+fn.Name.Name, "box", startID, false)
			preID := funcID
			var outputID string

			for _, stmt := range fn.Body.List {
				outputID = visitNode(&graph, stmt, preID, endID)
				preID = outputID
			}

			graph.Edges = append(graph.Edges, Edge{From: preID, To: endID})
		}
	}

	return graph
}

// visitNode 递归构建流程结构
func visitNode(graph *Graph, node ast.Node, preID string, endID string) string {
	if node == nil {
		return preID
	}

	switch n := node.(type) {

	case *ast.AssignStmt:
		label := prettyPrint(n)
		return addNode(graph, label, "box", preID, false)

	case *ast.ExprStmt:
		label := prettyPrint(n.X)
		return addNode(graph, label, "box", preID, false)

	case *ast.ReturnStmt:
		label := "return"
		if len(n.Results) > 0 {
			label += " " + prettyPrint(n.Results[0])
		}
		return addNode(graph, label, "box", preID, false)

	case *ast.IfStmt:
		condLabel := "if " + prettyPrint(n.Cond)
		condID := addNode(graph, condLabel, "diamond", preID, false)

		// YES 分支
		preYESID := condID
		for _, stmt := range n.Body.List {
			preYESID = visitNodeSkippingEdge(graph, stmt, preYESID, endID)
		}
		graph.Edges = append(graph.Edges, Edge{From: condID, To: preYESID, Label: "Yes"})

		// NO 分支
		var preNOID string
		if blk, ok := n.Else.(*ast.BlockStmt); ok {
			preNOID = condID
			for _, stmt := range blk.List {
				preNOID = visitNodeSkippingEdge(graph, stmt, preNOID, endID)
			}
		} else {
			preNOID = addNode(graph, "empty else", "box", condID, true)
		}
		graph.Edges = append(graph.Edges, Edge{From: condID, To: preNOID, Label: "No"})

		// 合并出口
		mergeID := addNode(graph, "merge", "ellipse", "", false)
		graph.Edges = append(graph.Edges, Edge{From: preYESID, To: mergeID})
		graph.Edges = append(graph.Edges, Edge{From: preNOID, To: mergeID})

		return mergeID

	case *ast.ForStmt:
		label := "for"
		if n.Init != nil {
			label += "\ninit: " + prettyPrint(n.Init)
		}
		if n.Cond != nil {
			label += "\ncond: " + prettyPrint(n.Cond)
		}
		if n.Post != nil {
			label += "\npost: " + prettyPrint(n.Post)
		}
		loopID := addNode(graph, label, "diamond", preID, false)

		bodyID := loopID
		for _, stmt := range n.Body.List {
			bodyID = visitNode(graph, stmt, bodyID, endID)
		}

		// 循环回跳
		graph.Edges = append(graph.Edges, Edge{From: bodyID, To: loopID, Label: "back"})
		return loopID

	default:
		// 其他语法结构略过
		return preID
	}
}

func visitNodeSkippingEdge(graph *Graph, node ast.Node, preID string, endID string) string {
	switch n := node.(type) {
	case *ast.AssignStmt:
		label := prettyPrint(n)
		return addNode(graph, label, "box", preID, true)
	case *ast.ExprStmt:
		label := prettyPrint(n.X)
		return addNode(graph, label, "box", preID, true)
	case *ast.ReturnStmt:
		label := "return"
		if len(n.Results) > 0 {
			label += " " + prettyPrint(n.Results[0])
		}
		return addNode(graph, label, "box", preID, true)
	default:
		return visitNode(graph, node, preID, endID) // fall back
	}
}

// 添加节点并自动连接到父节点
func addNode(graph *Graph, label, shape, parentID string, noEdge bool) string {
	nodeCounter++
	nodeID := fmt.Sprintf("n%d", nodeCounter)
	graph.Nodes[nodeID] = Node{Label: label, Shape: shape}
	if parentID != "" && !noEdge {
		graph.Edges = append(graph.Edges, Edge{From: parentID, To: nodeID})
	}
	return nodeID
}

// 美化打印 AST 节点为 Go 源码字符串
func prettyPrint(node ast.Node) string {
	var buf bytes.Buffer
	if err := format.Node(&buf, token.NewFileSet(), node); err != nil {
		return ""
	}
	return buf.String()
}

// 可选：检查边是否存在（调试用）
func checkEdge(from, to string, graph *Graph) bool {
	for _, edge := range graph.Edges {
		if edge.From == from && edge.To == to {
			return true
		}
	}
	return false
}
