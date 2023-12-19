package intermediate

import "github.com/tflexsoom/duffle/internal/container"

type DataConfig struct {
	FirstName  string
	SecondName string
	Values     container.Tree[string]
}
