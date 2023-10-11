package flow

import "gobit-demo/model/pk"

type CreateFlowEvent struct {
	FlowID    pk.ID  `json:"flow_id,omitempty"`
	FlowName  string `json:"flow_name,omitempty"`
	OwnerID   pk.ID  `json:"owner_id,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

func (e CreateFlowEvent) Topic() string {
	return "flow.created"
}

type DeleteFlowEvent struct {
	FlowID    pk.ID  `json:"flow_id,omitempty"`
	FlowName  string `json:"flow_name,omitempty"`
	OwnerID   pk.ID  `json:"owner_id,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

func (e DeleteFlowEvent) Topic() string {
	return "flow.deleted"
}
