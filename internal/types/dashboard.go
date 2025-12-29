package types

type ProductionDashboardResponse struct {
	ActiveJobs      int64   `json:"active_jobs"`      // Cantidad de órdenes en curso
	UtilizationRate float64 `json:"utilization_rate"` // % de máquinas ocupadas

	TotalRevenueFDM float64 `json:"revenue_fdm,omitempty"` // Opcional por ahora
	TotalRevenueSLS float64 `json:"revenue_sls,omitempty"` // Opcional por ahora

	Machines []Machine         `json:"machines"`      // Lista para mostrar estado de c/u
	Orders   []ProductionOrder `json:"active_orders"` // Lista de lo que se está haciendo ya
}
