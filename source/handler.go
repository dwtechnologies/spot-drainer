package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

type cloudwatchEventDetail struct {
	InstanceID string `json:"instance-id"`
}

func main() {
	lambda.Start(Handler)
}

func Handler(event events.CloudWatchEvent) error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	eventDetail := &cloudwatchEventDetail{}

	err := json.Unmarshal(event.Detail, eventDetail)
	if err != nil {
		log.Println("error unmarshaling event")
		return err
	}
	log.Printf("'EC2 Spot Instance Interruption Warning' event received from instance %s", eventDetail.InstanceID)

	client, err := newECSClient()
	if err != nil {
		log.Printf("can't get ECS client %s", err)
		return err
	}

	cluster, containerInstanceArn, err := getClusterAndContainerInstance(client, eventDetail.InstanceID)
	if err != nil {
		log.Printf("can't get container instance arn %s", err)
		return err
	}

	if containerInstanceArn == "" {
		log.Printf("instance %s is not part of any cluster, ignoring", eventDetail.InstanceID)
		return err
	}

	err = drainContainerInstance(client, cluster, containerInstanceArn)
	if err != nil {
		return err
	}

	log.Printf("instance %s in cluster %s set to DRAINING state", containerInstanceArn, cluster)
	return err
}

func newECSClient() (*ecs.ECS, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}

	return ecs.New(cfg), err
}

func getClusterAndContainerInstance(client *ecs.ECS, instanceID string) (cluster, containerInstanceArn string, err error) {
	// list ecs clusters
	r, err := client.ListClustersRequest(&ecs.ListClustersInput{}).Send()
	if err != nil {
		return "", "", err
	}

	for _, cluster := range r.ClusterArns {
		// list container instances
		r, err := client.ListContainerInstancesRequest(&ecs.ListContainerInstancesInput{
			Cluster: aws.String(cluster),
			Status:  "ACTIVE",
		}).Send()
		if err != nil {
			return "", "", err
		}

		// describe container instances
		if len(r.ContainerInstanceArns) > 0 {
			resci, err := client.DescribeContainerInstancesRequest(&ecs.DescribeContainerInstancesInput{
				Cluster:            aws.String(cluster),
				ContainerInstances: r.ContainerInstanceArns,
			}).Send()
			if err != nil {
				return "", "", err
			}

			for _, ci := range resci.ContainerInstances {
				if *ci.Ec2InstanceId == instanceID {
					return cluster, *ci.ContainerInstanceArn, nil
				}
			}

		}
	}

	return "", "", err
}

func drainContainerInstance(client *ecs.ECS, cluster, containerInstanceArn string) error {
	_, err := client.UpdateContainerInstancesStateRequest(&ecs.UpdateContainerInstancesStateInput{
		Cluster:            aws.String(cluster),
		ContainerInstances: []string{containerInstanceArn},
		Status:             "DRAINING",
	}).Send()

	return err
}
