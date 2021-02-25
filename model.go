package tsutsu

import (
	"encoding/json"
	"time"
)

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

type JobInfo struct {
	ID         uint64          `json:"id"`
	Category   string          `json:"category"`
	URL        string          `json:"url"`
	Payload    json.RawMessage `json:"payload,omitempty"`
	Status     string          `json:"status"`
	CreatedAt  time.Time       `json:"created_at"`
	NextTry    time.Time       `json:"next_try"`
	Timeout    uint            `json:"timeout"`
	FailCount  uint            `json:"fail_count"`
	MaxRetries uint            `json:"max_retries"`
	RetryDelay uint            `json:"retry_delay"`
}

type JobsInfo struct {
	Jobs       []JobInfo `json:"jobs"`
	NextCursor string    `json:"next_cursor"`
}

type FailedJobsInfo struct {
	FailedJobs []JobInfo `json:"failed_jobs"`
	NextCursor string    `json:"next_cursor"`
}

type NodeInfo struct {
	ID   string `json:"id"`
	Host string `json:"host"`
}
