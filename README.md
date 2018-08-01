# spot-drainer
Amazon EC2 Spot Instances offer spare compute capacity available in the AWS cloud at steep discounts compared to On-Demand instances.
Spot Instances can be interrupted by EC2 with two minutes of notification when EC2 needs the capacity back.
The two-minute warning for Spot instances is available via Amazon CloudWatch Events

This Lambda function reacts to such CloudWatch event annd set a spot ec2 instance to **DRAINING** state if the instance is used in an ECS cluster.

CloudWatchEvent

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



Build

```sh
cd source; GOOS=linux go build -o main handler.go && zip deployment.zip main
aws cloudformation package \
	--template-file sam.yaml \
	--output-template-file output_sam.yaml \
	--s3-bucket <some-bucket>
```



Deploy

```sh
aws cloudformation deploy \
	--template-file output_sam.yaml \
	--capabilities CAPABILITY_IAM \
	--stack-name <some-name> \
	--no-fail-on-empty-changeset
```


Test (SAM CLI)
```sh
AWS_PROFILE=<someprofile> sam local invoke spotDrainer --event sample_event.json --template sam.yaml
```
