Use this plugin for deploying a docker container application to AWS EC2 Container Service (ECS).

### Settings

* `access_key` - AWS access key ID, MUST be an IAM user with the AmazonEC2ContainerServiceFullAccess policy attached
* `secret_key` - AWS secret access key
* `region` - AWS availability zone
* `service` - Name of the service in the cluster, **MUST** be created already in ECS
* `cluster` - Name of the cluster. Optional. Default cluster is used if not specified
* `family` - Family name of the task definition to create or update with a new revision
* `image_name`, Container image to use, do not include the tag here
* `image_tag` - Tag of the image to use, defaults to latest
* `container_name` - Container name to use
* `port_mappings` - Port mappings from host to container, format is `hostPort containerPort`, protocol is automatically set to TransportProtocol
* `deployment_configuration` - Deployment configuration, format is `minimumHealthyPercent maximumPercent`
* `memory`, Amount of memory to assign to the container, defaults to 128
* `memoryReservation`, Amount of memoryReservation to assign to the container, defaults to 128
* `cpu`, Amount of CPU shares to assign to the container, defaults to 1024
* `desired_count`, Desired number of instances to run
* `environment_variables` - List of Environment Variables to be passed to the container, format is `NAME=VALUE`

## Example

```yaml
deploy:
  ecs:
    image: clicktripz/drone-ecs

    region: eu-west-1
    access_key: $$ACCESS_KEY_ID
    secret_key: $$SECRET_ACCESS_KEY
    family: my-ecs-task
    image_name: namespace/repo
    image_tag: latest
    service: my-ecs-service
    container_name: my-container-name
    environment_variables:
      - DATABASE_URI=$$MY_DATABASE_URI
    port_mappings:
      - 80 9000
    memory: 128
    memoryReservation: 128
    cpu: 1024
    desired_count: 1
    deployment_configuration: 50 200
```
