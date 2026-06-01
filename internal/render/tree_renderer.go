package render

import (
	"bytes"
	"slices"
	"strings"
	"time"

	"github.com/sangrita-tech/periscope/internal/domain"
)

type treeNode struct {
	name     string
	isFile   bool
	children map[string]*treeNode
}

type TreeRenderer struct {
}

func NewTreeRenderer() *TreeRenderer {
	return &TreeRenderer{}
}

func (r *TreeRenderer) Render(source domain.Source, entries []domain.Entry) string {
	root := newTreeNode(source.Name)

	for _, entry := range entries {
		relPath := entry.RelPath
		if relPath == "" || relPath == "." {
			root.isFile = true
			continue
		}

		parts := strings.Split(relPath, "/")
		current := root

		for index, part := range parts {
			if part == "" || part == "." {
				continue
			}

			child, ok := current.children[part]
			if !ok {
				child = newTreeNode(part)
				current.children[part] = child
			}

			if index == len(parts)-1 {
				child.isFile = true
			}

			current = child
		}
	}

	var buffer bytes.Buffer

	buffer.WriteString("# Periscoped project " + root.name + " " + time.Now().Format("2006-01-02 15:04:05") + "\n\n")

	buffer.WriteString("```\n")
	buffer.WriteString(root.name)

	if !root.isFile {
		buffer.WriteString("/")
	}

	buffer.WriteString("\n")
	writeTree(&buffer, root, "")
	buffer.WriteString("```\n\n")

	return buffer.String()
}

func newTreeNode(name string) *treeNode {
	return &treeNode{
		name:     name,
		children: make(map[string]*treeNode),
	}
}

func writeTree(buffer *bytes.Buffer, node *treeNode, prefix string) {
	children := sortedTreeChildren(node)

	for index, child := range children {
		isLast := index == len(children)-1

		branch := "├── "
		nextPrefix := prefix + "│   "

		if isLast {
			branch = "└── "
			nextPrefix = prefix + "    "
		}

		buffer.WriteString(prefix)
		buffer.WriteString(branch)
		buffer.WriteString(child.name)

		if !child.isFile {
			buffer.WriteString("/")
		}

		buffer.WriteString("\n")

		writeTree(buffer, child, nextPrefix)
	}
}

func sortedTreeChildren(node *treeNode) []*treeNode {
	children := make([]*treeNode, 0, len(node.children))

	for _, child := range node.children {
		children = append(children, child)
	}

	slices.SortFunc(children, func(a, b *treeNode) int {
		if a.isFile != b.isFile {
			if a.isFile {
				return 1
			}

			return -1
		}

		if a.name < b.name {
			return -1
		}

		if a.name > b.name {
			return 1
		}

		return 0
	})

	return children
}
