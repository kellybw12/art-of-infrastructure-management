package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// Script to clear all of the exiting s3 buckets

func cleanup() error {
	listBucketsCommand := exec.Command("aws", "s3", "ls", "--endpoint=http://localhost:4566")
	listOutputS3Buckets, err := listBucketsCommand.Output()
	if err != nil {
		fmt.Println("error listing s3 buckets: ", err)
		return err
	}

	bucketLines := strings.Split(string(listOutputS3Buckets), "\n")
	for _, bucketLine := range bucketLines {
		if bucketLine != "" {
			bucketName := strings.Fields(bucketLine)[2]
			fmt.Printf("deleting S3 bucket %s \n", bucketName)
			deleteBucketCommand := exec.Command("aws", "s3", "rb", "s3://"+bucketName, "--endpoint=http://localhost:4566")

			err := deleteBucketCommand.Run()
			if err != nil {
				fmt.Printf("error deleting bucket %s: %v\n", bucketName, err)
			} else {
				fmt.Printf("\u2714"+" successfully deleted bucket %s\n", bucketName)
			}

		}
	}
	fmt.Println("All S3 buckets deleted")
	return nil
}

func main() {
	err := cleanup()
	if err != nil {
		fmt.Println("error deleting existing aws s3 buckets")
	}
}
