package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

// Item is a struct that represents an item in the DynamoDB table.
//type Item struct {
//	ID   string `json:"id"`
//	Name string `json:"name"`
//}

func main() {
	lambda.Start(Router)
}

//
//func handleRequest(ctx context.Context, item Item) (Item, error) {
//	// Create a new AWS session and client for DynamoDB.
//	sess, err := session.NewSession()
//	if err != nil {
//		return Item{}, fmt.Errorf("Error creating new AWS session: %v", err)
//	}
//	client := dynamodb.New(sess)
//
//	// Write the item to the DynamoDB table.
//	_, err = client.PutItem(&dynamodb.PutItemInput{
//		TableName: aws.String("Fatawa"),
//		Item: map[string]*dynamodb.AttributeValue{
//			"id":   {S: aws.String(item.ID)},
//			"name": {S: aws.String(item.Name)},
//		},
//	})
//	if err != nil {
//		return Item{}, fmt.Errorf("Error writing item to DynamoDB: %v", err)
//	}
//
//	// Read the item from the DynamoDB table.
//	result, err := client.GetItem(&dynamodb.GetItemInput{
//		TableName: aws.String("Fatawa"),
//		Key: map[string]*dynamodb.AttributeValue{
//			"id": {S: aws.String(item.ID)},
//		},
//	})
//	if err != nil {
//		return Item{}, fmt.Errorf("Error reading item from DynamoDB: %v", err)
//	}
//
//	// Extract the item from the result.
//	item.Name = *result.Item["name"].S
//
//	return item, nil
//}
