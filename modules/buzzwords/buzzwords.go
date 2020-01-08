package buzzwords

import "botBoilerplate/modules"

func init() {
	modules.RegisterBuzzword(modules.Buzzword{
		Name:        "Hello",
		Class:       modules.FUN,
		Trigger:     []string{"Hello"},
		Reaction:    []string{"Moin!"},
		Probability: 0,
	})
}