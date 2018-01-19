package workspace

// Cfg is struct of operation configuration
type Cfg struct {
	Name       string `json:"name"`
	Hostname   string `json:"hostname"`
	Port       string `json:"port"`
	Abs        string `json:"abs"`
	Username   string `json:"username"`
	PrivateKey string `json:"private_key"`
}
