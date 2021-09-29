package drawer

type Drawer interface {
	Draw(from, to, msg string) string
	Comment(s string) string
	Header() string
	Footer() string
}
