package flow

import "gobit-demo/model"

type CreateFlowEvent struct {
	FlowID    model.UUID `json:"flow_id,omitempty"`
	FlowName  string     `json:"flow_name,omitempty"`
	OwnerID   model.UUID `json:"owner_id,omitempty"`
	Timestamp int64      `json:"timestamp,omitempty"`
}

func (e CreateFlowEvent) Topic() string {
	return "flow.created"
}

type DeleteFlowEvent struct {
	FlowID    model.UUID `json:"flow_id,omitempty"`
	FlowName  string     `json:"flow_name,omitempty"`
	OwnerID   model.UUID `json:"owner_id,omitempty"`
	Timestamp int64      `json:"timestamp,omitempty"`
}

func (e DeleteFlowEvent) Topic() string {
	return "flow.deleted"
}
