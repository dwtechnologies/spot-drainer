# spot-drainer (deprecated)
**As per Sep 27, 2019 Amazon supports automated draining for spot instances running ECS services.
[announcement](https://aws.amazon.com/about-aws/whats-new/2019/09/amazon-ecs-supports-automated-draining-for-spot-instances-running-ecs-services/)

codebase will be left for educational purpose.**

Amazon EC2 Spot Instances offer spare compute capacity available in the AWS cloud at steep discounts compared to On-Demand instances.
Spot Instances can be interrupted by EC2 with two minutes of notification when EC2 needs the capacity back.
The two-minute warning for Spot instances is available via Amazon CloudWatch Events

This Lambda function reacts to such CloudWatch event annd set a spot ec2 instance to **DRAINING** state if the instance is used in an ECS cluster.


### CloudWatch event
```json
{
  "version": "0",
  "id": "12345678-1234-1234-1234-123456789012",
  "detail-type": "EC2 Spot Instance Interruption Warning",
  "source": "aws.ec2",
  "account": "123456789012",
  "time": "2018-07-18T09:40:00Z",
  "region": "eu-west-1",
  "resources": [
    "arn:aws:ec2:eu-west-1:123456789012:instance/i-1234567890abcdef0"
  ],
  "detail": {
    "instance-id": "i-055eef5d6b5fcdfc8",
    "instance-action": "action"
  }
}
```


### Deployment requirements
- docker
- make
- aws cli


### Deployment
Use the included `Makefile` to deploy the resources.

The `OWNER` env var is for tagging. So you can set this to what you want.
The `ENVIRONMENT` env var is also for naming + tagging, but will also be included in CloudWatch logs.
This so you can make out differences between dev, test and prod etc. if you're running them on the same AWS Account.

```bash
AWS_PROFILE=my-profile AWS_REGION=region OWNER=TeamName S3_BUCKET=my-artifact-bucket ECS_CLUSTER=target-ecs-cluster make deploy
```

Example
```bash
AWS_PROFILE=default AWS_REGION=eu-west-1 OWNER=cloudops S3_BUCKET=my-artifact-bucket ECS_CLUSTER=cluster-one-prod make deploy
```

