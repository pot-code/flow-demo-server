package flow

type CreateFlowRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type ListFlowResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SaveFlowNodeRequest struct {
	ID          *uint  `json:"id"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	PrevID      *uint  `json:"prev_id"`
	NextID      *uint  `json:"next_id"`
}

type ListFlowNodeQueryParams struct {
	FlowID *uint `query:"flow_id" validate:"required"`
}

type ListFlowNodeResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	FlowID      *uint  `json:"flow_id"`
	PrevID      *uint  `json:"prev_id"`
	NextID      *uint  `json:"next_id"`
}
