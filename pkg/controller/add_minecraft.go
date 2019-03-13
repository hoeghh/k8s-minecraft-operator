package controller

import (
	"github.com/softica/minecraft-operator/pkg/controller/minecraft"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, minecraft.Add)
}
