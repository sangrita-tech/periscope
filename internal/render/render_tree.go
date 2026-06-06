package render

import (
	"bytes"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/sangrita-tech/periscope/internal/model"
)

type treeNode struct {
	Name     string
	IsFile   bool
	Children []*treeNode

	index map[string]*treeNode
}

func RenderTree(src model.Source, entries []model.Entry, generatedAt time.Time) string {
	var buffer bytes.Buffer

	root := buildTree(src.Name, entries)

	writeHeader(&buffer, src.Name, generatedAt)

	buffer.WriteString("```\n")
	buffer.WriteString(root.Name)

	if !root.IsFile {
		buffer.WriteString("/")
	}

	buffer.WriteString("\n")
	writeTree(&buffer, root, "")
	buffer.WriteString("```\n")

	return buffer.String()
}

func buildTree(rootPath string, entries []model.Entry) *treeNode {
	root := newTreeNode(path.Clean(rootPath), false)

	for _, entry := range entries {
		addTreeEntry(root, entry.RelPath)
	}

	sortTree(root)

	return root
}

func addTreeEntry(root *treeNode, relPath string) {
	relPath = path.Clean(relPath)

	if relPath == "." || relPath == "" {
		root.IsFile = true
		return
	}

	parts := strings.Split(relPath, "/")
	node := root

	for index, part := range parts {
		isFile := index == len(parts)-1
		node = node.child(part, isFile)
	}
}

func newTreeNode(name string, isFile bool) *treeNode {
	return &treeNode{
		Name:   name,
		IsFile: isFile,
		index:  make(map[string]*treeNode),
	}
}

func (n *treeNode) child(name string, isFile bool) *treeNode {
	if child, ok := n.index[name]; ok {
		if !isFile {
			child.IsFile = false
		}

		return child
	}

	child := newTreeNode(name, isFile)

	n.index[name] = child
	n.Children = append(n.Children, child)

	return child
}

func sortTree(node *treeNode) {
	sort.SliceStable(node.Children, func(i, j int) bool {
		left := node.Children[i]
		right := node.Children[j]

		if left.IsFile != right.IsFile {
			return !left.IsFile
		}

		return left.Name < right.Name
	})

	for _, child := range node.Children {
		sortTree(child)
	}
}

func writeTree(buffer *bytes.Buffer, node *treeNode, prefix string) {
	for index, child := range node.Children {
		isLast := index == len(node.Children)-1

		branch := "├── "
		nextPrefix := prefix + "│   "

		if isLast {
			branch = "└── "
			nextPrefix = prefix + "    "
		}

		buffer.WriteString(prefix)
		buffer.WriteString(branch)
		buffer.WriteString(child.Name)

		if !child.IsFile {
			buffer.WriteString("/")
		}

		buffer.WriteString("\n")

		writeTree(buffer, child, nextPrefix)
	}
}
