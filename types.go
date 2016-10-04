package main

import (
	"github.com/drone/drone-go/drone"
)

type Params struct {
	AccessKey               string            `json:"access_key"`
	SecretKey               string            `json:"secret_key"`
	Region                  string            `json:"region"`
	Family                  string            `json:"family"`
	Image                   string            `json:"image_name"`
	Tag                     string            `json:"image_tag"`
	Service                 string            `json:"service"`
	Cluster                 string            `json:"cluster"`
	ContainerName           string            `json:"container_name"`
	DeploymentConfiguration string            `json:"deployment_configuration"`
	CPU                     int64             `json:"cpu"`
	DesiredCount            int64             `json:"desired_count"`
	Memory                  int64             `json:"memory"`
	MemoryReservation       int64             `json:"memoryReservation"`
	Environment             drone.StringSlice `json:"environment_variables"`
	PortMappings            drone.StringSlice `json:"port_mappings"`
}
