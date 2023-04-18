# Advertisement Gateway

        Microservice written in Go
        Fetching data from REST request
        Sending it further via event using NATS

## Requirements

        Installed Go, Docker and PowerShell

### Go1ang Version

        Written in version 1.19.2

### Nats 
 
    Run: docker pull nats to download nats image
    Nats port can be changed inside nats_configuration.go file

### Docker

    Run docker build -t gateway . to build docker image 
    Run docker run -p 8080:8080 gateway to run docker image
