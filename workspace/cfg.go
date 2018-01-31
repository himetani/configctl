package workspace

import "time"

// Job is struct of operation configuration
type Job struct {
	Name        string    `json:"name"`
	Hostname    string    `json:"hostname"`
	Port        string    `json:"port"`
	Abs         string    `json:"abs"`
	Username    string    `json:"username"`
	PrivateKey  string    `json:"private_key"`
	LastUpdated time.Time `json:"last_updated"`
	LatestIdx   int       `json:"latest_idx"`
}
