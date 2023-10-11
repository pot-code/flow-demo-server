package flow

import "gobit-demo/model/pk"

type CreateFlowDto struct {
	Name        string `json:"name,omitempty" validate:"required,min=1,max=32"`
	Nodes       string `json:"nodes,omitempty"`
	Edges       string `json:"edges,omitempty"`
	Description string `json:"description,omitempty"`
}

type UpdateFlowDto struct {
	ID          pk.ID  `json:"id,omitempty" validate:"required"`
	Name        string `json:"name,omitempty" validate:"required,min=1,max=32"`
	Nodes       string `json:"nodes,omitempty"`
	Edges       string `json:"edges,omitempty"`
	Description string `json:"description,omitempty"`
}
