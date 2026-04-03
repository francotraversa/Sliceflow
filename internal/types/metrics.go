package types

type MachineMetric struct {
	MachineID   int     `json:"machine_id"`
	MachineName string  `json:"machine_name"`
	QueuedHours float64 `json:"queued_hours"`
}

type MaterialMetric struct {
	MaterialID   int     `json:"material_id"`
	MaterialName string  `json:"material_name"`
	MaterialType string  `json:"material_type"`
	QueuedKilos  float64 `json:"queued_kilos"`
}

type MetricsResponse struct {
	Machines  []MachineMetric  `json:"machines"`
	Materials []MaterialMetric `json:"materials"`
}
