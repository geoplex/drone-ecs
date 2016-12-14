Use this plugin for deploying a docker container application to AWS EC2 Container Service (ECS).

### Settings

* `access_key` - AWS access key ID, MUST be an IAM user with the AmazonEC2ContainerServiceFullAccess policy attached
* `secret_key` - AWS secret access key
* `region` - AWS availability zone
* `service` - Name of the service in the cluster, **MUST** be created already in ECS
* `cluster` - Name of the cluster. Optional. Default cluster is used if not specified
* `family` - Family name of the task definition to create or update with a new revision
* `deployment_configuration` - Deployment configuration, format is `minimumHealthyPercent maximumPercent`
* `desired_count`, Desired number of instances to run
* `network_mode` - Container network mode: bridge, host, none. Defaults to bridge`
* `container_definitions` - An array of container definitions
  * `container_name` - Container name to use
  * `image_name`, Container image to use, do not include the tag here
  * `image_tag` - Tag of the image to use, defaults to latest
  * `port_mappings` - Port mappings from host to container, format is `hostPort containerPort`, protocol is automatically set to TransportProtocol
  * `memory`, Amount of memory to assign to the container, defaults to 128
  * `memoryReservation`, Amount of memoryReservation to assign to the container, defaults to 128
  * `cpu`, Amount of CPU shares to assign to the container, defaults to 1024
  * `environment_variables` - List of Environment Variables to be passed to the container, format is `NAME=VALUE`
  * `docker_labels` - Optional docker labels`
  * `links` - Optional links`
  * `log_driver` - Log driver`
  * `log_driver_options` - Log driver options`

## Example

```yaml
deploy:
  ecs:
    image: geoplex/drone-ecs

    region: eu-west-1
    access_key: $$ACCESS_KEY_ID
    secret_key: $$SECRET_ACCESS_KEY
    family: my-ecs-task
    service: my-ecs-service
    desired_count: 1
    deployment_configuration: 50 200

    container_definitions:
      - container_name: flask
        image_name: namespace/repo
        image_tag: latest
        environment_variables:
          - DATABASE_URI=$$MY_DATABASE_URI
        memory: 256
        cpu: 512
      - container_name: nginx
        image_name: namespace/repo
        image_tag: latest
        port_mappings:
          - 80 80
        memory: 235
        cpu: 512
        links:
          - "flask"
```
