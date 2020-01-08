package modules

type Model struct {
	Name  string
	Model interface{}
}

var models []Model

func RegisterModel(model Model) {
	models = append(models, model)
}

func GetModels() []Model {
	return models
}
