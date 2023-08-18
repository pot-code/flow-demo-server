package flow

import (
	"gobit-demo/model"
)

type CreateFlowRequest struct {
	Name        string `json:"name,omitempty" validate:"required,min=1,max=32"`
	Nodes       string `json:"nodes,omitempty"`
	Edges       string `json:"edges,omitempty"`
	Description string `json:"description,omitempty"`
}

type UpdateFlowRequest struct {
	ID          model.UUID `json:"id,omitempty" validate:"required"`
	Name        string     `json:"name,omitempty" validate:"required,min=1,max=32"`
	Nodes       string     `json:"nodes,omitempty"`
	Edges       string     `json:"edges,omitempty"`
	Description string     `json:"description,omitempty"`
}
