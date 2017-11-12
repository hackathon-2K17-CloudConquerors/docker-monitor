# docker-monitor
Docker events reporting microservice

----

## docker-monitor architecture

<img>

* docker-monitor: The main service that polls running containers (nginx, postgres and httpd in this case), pushes data to influx and listens for request from client

* influxdb: Time-series database used to store metrics of the containers

----

## Getting started with docker-monitor

1) Checkout the deployment files:
```
git clone https://github.com/hackathon-2K17-CloudConquerors/docker-monitor.git
cd docker-monitor/deployments/docker-compose
```

2) create the configuration file: (keeping everything by default should be fine)
```

```

3) Start the containers:
```
docker-compose up
```

## Prerequisites

* docker-monitor requires access to the Docker event API socket (`/var/run/docker.sock` by default)
* docker-monitor requires privileged access.
