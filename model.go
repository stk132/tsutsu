package tsutsu

type QueueStats struct {
	TotalPushes            int64 `json:"total_pushes"`
	TotalPops              int64 `json:"total_pops"`
	TotalSuccesses         int64 `json:"total_successes"`
	TotalFailures          int64 `json:"total_failures"`
	TotalPermanentFailures int64 `json:"total_permanent_failures"`
	TotalCompletes         int64 `json:"total_completes"`
	TotalElapsed           int64 `json:"total_elapsed"`
	PushesPerSecond        int64 `json:"pushes_per_second"`
	PopPerSecond           int64 `json:"pop_per_second"`
	TotalWorkers           int64 `json:"total_workers"`
	IdleWorkers            int64 `json:"idle_workers"`
	ActiveNodes            int64 `json:"active_nodes"`
}
