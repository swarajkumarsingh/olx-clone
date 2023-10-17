# Go boilerplate code gin

[[_TOC_]]

## Pre-requisites

1. Install golang on your machine. You can get it from [here](https://go.dev/doc/install)
2. Now install **Docker for Mac** and set it up by simply opening it, and entering your mac password
3. Create the folder to hold the logs by running the following command:
   ```sh
   mkdir -p ~/logs/go-docker
   ```
## Local Setup

> :warning: **NOTE:** DO NOT use docker compose for any non local setups (Production, Dev or UAT Environments)

You can use the following command to start all the containers on your localhost:

1.  Clone and open the directory
    ```sh
    git clone https://github.com/swarajkumarsingh/go-biolderplate-code-gin
    cd go-biolderplate-code-gin
    ```
2.  Run the project
    ```sh
        make run
    ```
3. And Voilà It's started.

## Important Instructions

### VS Code Checks
In case you use **VS Code**, Make sure [extension](https://marketplace.visualstudio.com/items?itemName=golang.go) is installed, along with `gopls`. While working on new files, keep checking `Problems` section in bottom tab

### Logging

Follow the below listed methods for logging. Taking into care correct log levels are used.

> NOTE: Production skips debug logs, while UAT and dev have debug level logs enabled as well.

#### Logging with request details
Below will automatically log `userID` as well in case of SDK APIs (using `AuthFilter`)

```go
logger.WithRequest(r).Debug(<log>)
logger.WithRequest(r).Info(<log>)
logger.WithRequest(r).Warn(<log>)
logger.WithRequest(r).Error(<log>)
```

#### Logging with `userID`

```go
logger.WithUser(userID).Debug(<log>)
logger.WithUser(userID).Info(<log>)
logger.WithUser(userID).Warn(<log>)
logger.WithUser(userID).Error(<log>)
```

#### Logging with `loanApplicationID`

```go
logger.WithLoanApplication(loanApplicationID).Debug(<log>)
logger.WithLoanApplication(loanApplicationID).Info(<log>)
logger.WithLoanApplication(loanApplicationID).Warn(<log>)
logger.WithLoanApplication(loanApplicationID).Error(<log>)
```

### Tracing External Requests

To trace external HTTP requests in DataDog, use GetTraceableHTTPClient() defined in tracer.go to initialize an http client.

The method accepts two arguments: a pointer to time.Duration and resource name as string.

## Commands

### Running Test Cases
```sh
go test -v ./...
```

## Adding a new migration script

To add an migration script for a given table, you have to create a blank SQL file first. You can club multiple SQL files in a single file by making sure each of them ends with a semicolon (;). To generate a blank sql file for a table `tablename` use following command after changing directory in terminal to inside your repo folder:

```sh
go run commands/generatesql.go tablename
```

#### Starting the datadog agent

When starting the `datadog agent` try to run it in the same network as the `deployable project` docker image. If we do this we'll simply be able to pass the docker image name of the `deployable project` docker image as `DD_AGENT_HOST` and we'll easily be connected to the `datadog agent`.

```sh
docker run -d --name datadog-agent \
              --network kycnet \
              --cgroupns host \
              --security-opt apparmor:unconfined \
              --cap-add=SYS_ADMIN \
              --cap-add=SYS_RESOURCE \
              --cap-add=SYS_PTRACE \
              --cap-add=NET_ADMIN \
              --cap-add=NET_BROADCAST \
              --cap-add=NET_RAW \
              --cap-add=IPC_LOCK \
              --cap-add=CHOWN \
              -v /var/run/docker.sock:/var/run/docker.sock:ro \
              -v /proc/:/host/proc/:ro \
              -v /sys/fs/cgroup/:/host/sys/fs/cgroup:ro \
              -v /etc/passwd:/etc/passwd:ro \
              -v /sys/kernel/debug:/sys/kernel/debug \
              -e DD_API_KEY="<<DD_API_KEY_HERE>>" \
              -e DD_APM_ENABLED=true \
              -e DD_SITE="datadoghq.eu" \
              -e DD_APM_NON_LOCAL_TRAFFIC=true \
              -e DD_HOSTNAME="kyc-dev" \
              -e DD_PROCESS_AGENT_ENABLED=true \
              -e DD_CONTAINER_EXCLUDE="name:datadog-agent" \
              -e DD_SYSTEM_PROBE_ENABLED=true \
              -e DD_DOGSTATSD_NON_LOCAL_TRAFFIC=true \
              gcr.io/datadoghq/agent:latest
```

In case there are multiple hosts, change the DD_HOSTNAME parameter above to include the machine number

---
