package util

import "time"

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

type Process struct {
	Id      string    `json:"id"`
	App     string    `json:"app"`
	Command string    `json:"command"`
	Host    string    `json:"host"`
	Image   string    `json:"image"`
	Name    string    `json:"name"`
	Ports   []string  `json:"ports"`
	Release string    `json:"release"`
	Cpu     float64   `json:"cpu"`
	Memory  float64   `json:"memory"`
	Started time.Time `json:"started"`
}

type Processes []Process

type Release struct {
	Id      string    `json:"id"`
	App     string    `json:"app"`
	Created time.Time `json:"created"`
	Status  string    `json:"status"`
}

type Releases []Release
