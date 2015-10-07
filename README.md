platform
========

The home of all things related to the Good Eggs [12 Factor App](http://12factor.net/) platform, currently based on [Convox](http://convox.com/).

## Projects

### [ami](./ami)
Contains the [packer.io]() definition for our custom AMI.  It includes extra init scripts to start a logspout-goodeggs on each instance using ECS.

### [logspout-goodeggs](./logspout-goodeggs)
Our custom build of [logspout-http](https://github.com/raychaser/logspout-http) that ships docker logs to a [Sumo Logic](https://www.sumologic.com/) collector.

### [cmd/convox-install](./cmd/convox-install)
A command-line wrapper for `convox install`, which handles patching the base [Convox CloudFormation Definition](https://github.com/convox/rack/blob/master/api/dist/kernel.json) with our changes.

