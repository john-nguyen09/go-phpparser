package lexer

type ChangeEvent struct {
	Start int
	End   int
	Text  []rune
}
