Convox AMI
----------

Our custom-built AMI that is used for all Convox instances.  The base AMI is whatever the current Convox base AMI is (currently the Amazon ECS Optimized image), which ccan be found in their [CloudFormation Definition](https://github.com/convox/rack/blob/master/api/dist/kernel.json#L13).

Our custom AMI runs a script on first boot that creates a [logspout ECS task](../../infra/logspout), which ensures every machine in the cluster has a logspout running on it.  This method is taken from Amazon's [Running an Amazon ECS Task on Every Instance](https://aws.amazon.com/blogs/compute/running-an-amazon-ecs-task-on-every-instance/) article.

## Development

We use [packer](https://packer.io/) to build the custom AMI.

   $ brew install packer
   $ packer build convox-ami.json

## Test

Once you have an AMI candidate, you should boot a new Convox rack and verify the instance associated to the ECS cluster correctly.

   $ convox install --ami <new ami> --instance-count 1 --stack-name test --key goodeggs-aws

## Release

??

## TODO

* document how to update the Convox rack once you have a new AMI.

