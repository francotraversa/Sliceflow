package types

type ProductionDashboardResponse struct {
	ActiveJobs      int64   `json:"active_jobs"`      // Number of active orders
	UtilizationRate float64 `json:"utilization_rate"` // % of machines in use

	TotalRevenueFDM float64 `json:"revenue_fdm,omitempty"` // Optional for now
	TotalRevenueSLS float64 `json:"revenue_sls,omitempty"` // Optional for now

	Machines []Machine         `json:"machines"`      // List of machines with their status
	Orders   []ProductionOrder `json:"active_orders"` // Currently active orders
}
