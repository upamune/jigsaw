package drawer

import (
	"fmt"
	"strings"
)

const TypeMermaid = "mermaid"

type Mermaid struct{}

func (m *Mermaid) Draw(from, to, msg string) string {
	indent := strings.Repeat(" ", 4)
	return fmt.Sprintf("%s%s->>%s: %s", indent, from, to, msg)
}

func (m *Mermaid) Comment(s string) string {
	return fmt.Sprintf("%%%% %s", s)
}

func (m *Mermaid) Header() string {
	return "sequenceDiagram"
}

func (m *Mermaid) Footer() string {
	return ""
}
