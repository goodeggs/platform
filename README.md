platform
========

The home of all things related to the Good Eggs [12 Factor App](http://12factor.net/) platform.  We heavily leverage [Convox](http://convox.com/), but also believe in the [Paasta Principles](http://paasta.readthedocs.io/en/latest/about/paasta_principles.html).

## Concepts

### Logging
Though Convox has per-app logging similar to Heroku's Logplex, we choose to enable logging at the platform level.  [logspout](https://github.com/gliderlabs/logspout) is a docker-oriented log shipper.  We maintain a [custom logspout module](https://github.com/goodeggs/logspout-http) that parses JSON application logs and ships to Sumo Logic.  Every docker container  will have its logs sent to Sumo Logic automatically.  You can opt out by adding `LOGSPOUT=ignore` to your container's docker environment.

## Projects

### [ami](./ami)
Contains the [packer.io](https://packer.io/) definition for our custom AMI.

### [logspout-goodeggs](./logspout-goodeggs)
Our custom build of [logspout-http](https://github.com/raychaser/logspout-http) that ships docker logs to a [Sumo Logic](https://www.sumologic.com/) collector.

### [ranch](./cmd/ranch)
The platform's CLI, which serves as the main interface for developers and CI/CD workers.

### [ranch-api](https://github.com/goodeggs/ranch-api)
The platform's API, to which the CLI communicates.

### [goodeggs-terraform](https://github.com/goodeggs/goodeggs-terraform)
Terraform for the entire infrastructure, including this platform.

