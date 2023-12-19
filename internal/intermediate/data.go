package intermediate

import "github.com/tflexsoom/duffle/internal/container"

type DataValue struct {
	Type      TypeId
	TextValue string
}

type DataConfig struct {
	TypeHash   string
	FirstName  string
	SecondName string
	Values     container.Tree[DataValue]
}
