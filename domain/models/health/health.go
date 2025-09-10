package health

import "database/sql"

type (
	HealthMetric struct {
		Status map[string]interface{} `json:"status"`
		DB     sql.DBStats            `json:"database"`
	}
)
