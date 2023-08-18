package flow

import (
	"gobit-demo/model"
)

type Position struct {
	X float32 `json:"x" validate:"required"`
	Y float32 `json:"y" validate:"required"`
}

type Node struct {
	ID       *string                `json:"id" validate:"required"`
	Type     string                 `json:"type" validate:"required"`
	Data     map[string]interface{} `json:"data"`
	Position *Position              `json:"position" validate:"required"`
}

type Edge struct {
	ID           *string `json:"id" validate:"required"`
	Source       *string `json:"source" validate:"required"`
	Target       *string `json:"target" validate:"required"`
	SourceHandle string  `json:"sourceHandle"`
	TargetHandle string  `json:"targetHandle"`
}

type CreateFlowRequest struct {
	Name        string  `json:"name,omitempty" validate:"required,min=1,max=32"`
	Nodes       []*Node `json:"nodes,omitempty" validate:"required"`
	Edges       []*Edge `json:"edges,omitempty" validate:"required"`
	Description string  `json:"description,omitempty"`
}

type UpdateFlowRequest struct {
	ID          model.UUID `json:"id,omitempty" validate:"required"`
	Name        string     `json:"name,omitempty" validate:"required,min=1,max=32"`
	Nodes       []*Node    `json:"nodes,omitempty" validate:"required"`
	Edges       []*Edge    `json:"edges,omitempty" validate:"required"`
	Description string     `json:"description,omitempty"`
}

type FlowObjectResponse struct {
	ID          model.UUID `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Nodes       []*Node    `json:"nodes,omitempty"`
	Edges       []*Edge    `json:"edges,omitempty"`
	Description string     `json:"description,omitempty"`
}

type ListFlowResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
