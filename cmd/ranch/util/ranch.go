package util

type RanchConfig struct {
	Name      string                `json:"name"`
	Version   int                   `json:"version"`
	Processes RanchConfigProcessMap `json:"processes"`
}

type RanchConfigProcess struct {
	Command   string `json:"command"`
	Instances int    `json:"instances"`
	Memory    int    `json:"memory"`
}

type RanchConfigProcessMap map[string]RanchConfigProcess
