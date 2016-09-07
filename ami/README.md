Convox AMI
----------

Our custom-built AMI that is used for all Convox instances.  The base AMI is whatever the current Convox base AMI is (currently the Amazon ECS Optimized image), which ccan be found in their [CloudFormation Definition](https://github.com/convox/rack/blob/master/api/dist/kernel.json#L13).

Our custom AMI runs a script on first boot that creates a [logspout ECS task](../../infra/logspout), which ensures every machine in the cluster has a logspout running on it.  This method is taken from Amazon's [Running an Amazon ECS Task on Every Instance](https://aws.amazon.com/blogs/compute/running-an-amazon-ecs-task-on-every-instance/) article.

## Development

```
$ export CONVOX_VERSION=$( curl -s http://convox.s3.amazonaws.com/release/versions.json | jq -r 'map(select(.published)) | map(.version) | sort | last' )
$ echo source_ami=$( curl -s https://convox.s3.amazonaws.com/release/$CONVOX_VERSION/formation.json | jq -r '.Mappings.RegionConfig["us-east-1"].Ami' )
```

Update `variables.dev.json` with the new `source_ami` value from above.

We use [packer](https://packer.io/) to build the custom AMI.  This step should be done in the `dev` AWS account!

    $ brew install packer
    $ aws-vault exec dev -- packer build -var="env=dev" -var="version=$(git rev-parse --short HEAD)" -var-file="variables.dev.json" packer.json

You can also push a branch to Github and Travis will run this for you.

## Test

Once you have an AMI candidate, proceed to the [goodeggs/goodeggs-terraform](https://github.com/goodeggs/goodeggs-terraform) repo to apply your changes to dev.  Return here to manually test the development rack:

1. That `convox rack` still works and returns the correct information
2. The `hello-world` app is accessible via its ELB
3. The HTTP logs from step 1 made it into SumoLogic
4. The collectd metrics made it into Librato (try [here](https://metrics.librato.com/s/metrics/collectd.cpu.percent.user?q=collectd.cpu&source=aws.dev.%2a))

## Release

Switch to the `prod` AWS account and rebuild the AMI.  You should use the short git SHA as the version:

    $ aws-vault exec prod -- packer build -var="env=prod" -var="version=$(git rev-parse --short HEAD)" -var-file="variables.prod.json" packer.json

Same as before, proceed to [goodeggs/goodeggs-terraform](https://github.com/goodeggs/goodeggs-terraform) and apply your changes to prod.

