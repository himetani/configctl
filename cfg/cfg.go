package cfg

// Cfg is struct of operation configuration
type Cfg struct {
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
	Abs      string `json:"abs"`
	Username string `json:"username"`
}
