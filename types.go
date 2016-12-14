package main

import (
	"github.com/drone/drone-go/drone"
)

type ContainerDefinition struct {
	ContainerName     string            `json:"container_name"`
	Image             string            `json:"image_name"`
	Tag               string            `json:"image_tag"`
	CPU               int64             `json:"cpu"`
	Memory            int64             `json:"memory"`
	MemoryReservation int64             `json:"memoryReservation"`
	Environment       drone.StringSlice `json:"environment_variables"`
	PortMappings      drone.StringSlice `json:"port_mappings"`
	LogDriver         string            `json:"log_driver"`
	LogDriverOptions  drone.StringSlice `json:"log_driver_options"`
	DockerLabels      drone.StringSlice `json:"docker_labels"`
	Links             drone.StringSlice `json:"links"`
}

type Params struct {
	AccessKey               string                 `json:"access_key"`
	SecretKey               string                 `json:"secret_key"`
	Region                  string                 `json:"region"`
	Family                  string                 `json:"family"`
	Service                 string                 `json:"service"`
	Cluster                 string                 `json:"cluster"`
	NetworkMode             string                 `json:"network_mode"`
	DeploymentConfiguration string                 `json:"deployment_configuration"`
	ContainerDefinitions    []*ContainerDefinition `json:"container_definitions"`
	DesiredCount            int64                  `json:"desired_count"`
}
