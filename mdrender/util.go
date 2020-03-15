package md

import (
	"io"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type stopWalkError struct{}

func (stopWalkError) Error() string {
	return "StapWalk"
}

func extractTextFromNode(n ast.Node, reader text.Reader, w io.Writer) error {
	return ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && n.Kind() == ast.KindText {
			textNode := n.(*ast.Text)
			cnt, err := w.Write(reader.Value(textNode.Segment))
			if textNode.HardLineBreak() {
				w.Write([]byte("\n"))
			}
			if textNode.SoftLineBreak() {
				w.Write([]byte(" "))
			}
			if cnt == 0 || err != nil {
				if err == nil {
					err = stopWalkError{}
				}
				return ast.WalkStop, err
			}
		}
		return ast.WalkContinue, nil
	})
}
