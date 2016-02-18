Convox AMI
----------

Our custom-built AMI that is used for all Convox instances.  The base AMI is whatever the current Convox base AMI is (currently the Amazon ECS Optimized image), which ccan be found in their [CloudFormation Definition](https://github.com/convox/rack/blob/master/api/dist/kernel.json#L13).

Our custom AMI runs a script on first boot that creates a [logspout ECS task](../../infra/logspout), which ensures every machine in the cluster has a logspout running on it.  This method is taken from Amazon's [Running an Amazon ECS Task on Every Instance](https://aws.amazon.com/blogs/compute/running-an-amazon-ecs-task-on-every-instance/) article.

## Development

```
$ CONVOX_VERSION=$( curl http://convox.s3.amazonaws.com/release/versions.json | jq -r 'map(select(.published)) | map(.version) | sort | last') )
$ curl https://convox.s3.amazonaws.com/release/$CONVOX_VERSION/formation.json > convox-formation.json
$ SOURCE_AMI=$( cat convox-formation.json | jq -r '.Mappings.RegionConfig["us-east-1"].Ami' )
```

Update `packer.json` with the new `source_ami` value from above.

We use [packer](https://packer.io/) to build the custom AMI.  This step should be done in the `dev` AWS account!

    $ brew install packer
    $ packer build \
      -var 'env=dev' \
      -var 'librato_email=...' \
      -var 'librato_token=...' \
      -var 'logspout_token=...' \
      packer.json

## Test

Once you have an AMI candidate, you should upload the `convox-formation.json` file, update the `Ami` and `Version` CloudFormation parameters in the dev cluster and verify:

1. That `convox rack` still works and returns the correct information
2. The `hello-world` app is accessible via its ELB
3. The HTTP logs from step 1 made it into SumoLogic
4. The collectd metrics made it into Librato (try [here](https://metrics.librato.com/s/metrics/collectd.cpu.percent.user?q=collectd.cpu&source=aws.dev.%2a))

## Release

Switch to the `prod` AWS account and rebuild the AMI.  You should use the short git SHA as the version:

    $ packer build \
      -var 'env=prod' \
      -var 'version=abcdef2' \
      -var 'librato_email=...' \
      -var 'librato_token=...' \
      -var 'logspout_token=...' \
      packer.json

Now you can upload the `convox-formation.json`, update the `Ami` and `Version` CloudFormation parameters, and verify as before.

