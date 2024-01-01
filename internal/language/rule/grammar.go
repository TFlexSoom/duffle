package rule

type Token struct {
	Name    string
	Pattern string
}

type Generator[T any] func[T any](parseFn)

type Rule struct {
	DependsOn []Token
	Generator Generator
}
