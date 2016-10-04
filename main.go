package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin"
)

var (
	buildCommit string
)

func main() {
	fmt.Printf("Drone AWS ECS Plugin built from %s\n", buildCommit)

	workspace := drone.Workspace{}
	repo := drone.Repo{}
	build := drone.Build{}
	vargs := Params{}

	plugin.Param("workspace", &workspace)
	plugin.Param("repo", &repo)
	plugin.Param("build", &build)
	plugin.Param("vargs", &vargs)
	plugin.MustParse()

	if len(vargs.AccessKey) == 0 {
		fmt.Println("Please provide an access key")

		os.Exit(1)
		return
	}

	if len(vargs.SecretKey) == 0 {
		fmt.Println("Please provide a secret key")

		os.Exit(1)
		return
	}

	if len(vargs.Region) == 0 {
		fmt.Println("Please provide a region")

		os.Exit(1)
		return
	}

	if len(vargs.Family) == 0 {
		fmt.Println("Please provide a task definition family name")

		os.Exit(1)
		return
	}

	if len(vargs.Image) == 0 {
		fmt.Println("Please provide an image name")

		os.Exit(1)
		return
	}

	if len(vargs.Tag) == 0 {
		vargs.Tag = "latest"
	}

	if len(vargs.Service) == 0 {
		fmt.Println("Please provide a service name")

		os.Exit(1)
		return
	}

	if len(vargs.Cluster) == 0 {
		fmt.Println("Cluster: default")
	} else {
		fmt.Printf("Cluster: %s\n", vargs.Cluster)
	}

	if len(vargs.ContainerName) == 0 {
		vargs.ContainerName = vargs.Family + "-container"
	}

	svc := ecs.New(
		session.New(&aws.Config{
			Region:      aws.String(vargs.Region),
			Credentials: credentials.NewStaticCredentials(vargs.AccessKey, vargs.SecretKey, ""),
		}))

	Image := vargs.Image + ":" + vargs.Tag

	definition := ecs.ContainerDefinition{
		Command: []*string{},

		DnsSearchDomains:      []*string{},
		DnsServers:            []*string{},
		DockerLabels:          map[string]*string{},
		DockerSecurityOptions: []*string{},
		EntryPoint:            []*string{},
		Environment:           []*ecs.KeyValuePair{},
		Essential:             aws.Bool(true),
		ExtraHosts:            []*ecs.HostEntry{},

		Image:        aws.String(Image),
		Links:        []*string{},
		MountPoints:  []*ecs.MountPoint{},
		Name:         aws.String(vargs.ContainerName),
		PortMappings: []*ecs.PortMapping{},

		Ulimits: []*ecs.Ulimit{},
		//User: aws.String("String"),
		VolumesFrom: []*ecs.VolumeFrom{},
		//WorkingDirectory: aws.String("String"),
	}

	if vargs.CPU != 0 {
		definition.Cpu = aws.Int64(vargs.CPU)
	}

	if vargs.Memory == 0 && vargs.MemoryReservation == 0 {
		definition.MemoryReservation = aws.Int64(128)
	} else {
		if vargs.Memory != 0 {
			definition.Memory = aws.Int64(vargs.Memory)
		}
		if vargs.MemoryReservation != 0 {
			definition.MemoryReservation = aws.Int64(vargs.MemoryReservation)
		}
	}

	// Port mappings
	for _, portMapping := range vargs.PortMappings.Slice() {
		cleanedPortMapping := strings.Trim(portMapping, " ")
		parts := strings.SplitN(cleanedPortMapping, " ", 2)
		hostPort, hostPortErr := strconv.ParseInt(parts[0], 10, 64)
		if hostPortErr != nil {
			fmt.Println(hostPortErr.Error())
			os.Exit(1)
			return
		}
		containerPort, containerPortError := strconv.ParseInt(parts[1], 10, 64)
		if containerPortError != nil {
			fmt.Println(containerPortError.Error())
			os.Exit(1)
			return
		}

		pair := ecs.PortMapping{
			ContainerPort: aws.Int64(containerPort),
			HostPort:      aws.Int64(hostPort),
			Protocol:      aws.String("TransportProtocol"),
		}

		definition.PortMappings = append(definition.PortMappings, &pair)
	}

	// Environment variables
	for _, envVar := range vargs.Environment.Slice() {
		parts := strings.SplitN(envVar, "=", 2)
		pair := ecs.KeyValuePair{
			Name:  aws.String(strings.Trim(parts[0], " ")),
			Value: aws.String(strings.Trim(parts[1], " ")),
		}
		definition.Environment = append(definition.Environment, &pair)
	}
	params := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: []*ecs.ContainerDefinition{
			&definition,
		},
		Family:  aws.String(vargs.Family),
		Volumes: []*ecs.Volume{},
	}
	resp, err := svc.RegisterTaskDefinition(params)

	if err != nil {
		fmt.Println(err.Error())

		os.Exit(1)
		return
	}

	val := *(resp.TaskDefinition.TaskDefinitionArn)
	sparams := &ecs.UpdateServiceInput{
		Cluster:        aws.String(vargs.Cluster),
		Service:        aws.String(vargs.Service),
		TaskDefinition: aws.String(val),
	}

	if vargs.DesiredCount != 0 {
		sparams.DesiredCount = aws.Int64(vargs.DesiredCount)
	}

	cleanedDeploymentConfiguration := strings.Trim(vargs.DeploymentConfiguration, " ")
	parts := strings.SplitN(cleanedDeploymentConfiguration, " ", 2)
	minimumHealthyPercent, minimumHealthyPercentError := strconv.ParseInt(parts[0], 10, 64)
	if minimumHealthyPercentError != nil {
		fmt.Println(minimumHealthyPercentError.Error())
		os.Exit(1)
		return
	}
	maximumPercent, maximumPercentErr := strconv.ParseInt(parts[1], 10, 64)
	if maximumPercentErr != nil {
		fmt.Println(maximumPercentErr.Error())
		os.Exit(1)
		return
	}

	sparams.DeploymentConfiguration = &ecs.DeploymentConfiguration{
		MaximumPercent:        aws.Int64(maximumPercent),
		MinimumHealthyPercent: aws.Int64(minimumHealthyPercent),
	}

	sresp, serr := svc.UpdateService(sparams)

	if serr != nil {
		fmt.Println(serr.Error())
		os.Exit(1)
		return
	}

	fmt.Println(sresp)

	fmt.Println(resp)
}
