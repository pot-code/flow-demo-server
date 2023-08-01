package flow

type CreateFlowRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
}

type ListFlowDto struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateFlowNodeRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
	FlowID      *uint   `json:"flow_id" validate:"required"`
	PrevID      *uint   `json:"prev_id"`
	NextID      *uint   `json:"next_id"`
}

type ListFlowNodeParams struct {
	FlowID *uint `json:"flow_id" query:"flow_id" validate:"required"`
}

type ListFlowNodeDto struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	FlowID      *uint   `json:"flow_id"`
	PrevID      *uint   `json:"prev_id"`
	NextID      *uint   `json:"next_id"`
}
