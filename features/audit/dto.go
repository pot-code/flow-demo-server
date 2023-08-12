package audit

import "gobit-demo/model"

type auditLog struct {
	model.AuditLog
	rawPayload any
}
