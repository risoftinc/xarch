package health

type (
	HealthMetric struct {
		Status map[string]interface{} `json:"status"`
		DB     interface{}            `json:"database"`
	}
)
