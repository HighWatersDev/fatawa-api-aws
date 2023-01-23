package main

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
	"github.com/google/uuid"
)

const TableName = "Fatawa"

//var fatwa routes.CreateFatwa

var db dynamodb.Client

func init() {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	db = *dynamodb.NewFromConfig(sdkConfig)
}

type Fatwa struct {
	Id       string `json:"id" dynamodbav:"id"`
	Title    string `json:"title" dynamodbav:"title"`
	Question string `json:"question" dynamodbav:"question"`
	Answer   string `json:"answer" dynamodbav:"answer"`
	Link     string `json:"link" dynamodbav:"link"`
	Author   string `json:"author" dynamodbav:"author"`
	Topic    string `json:"topic" dynamodbav:"topic"`
	Lang     string `json:"lang" dynamodbav:"lang"`
	//	CreatedAt time.Time `json:"createdat" dynamodbav:"createdat"`
	//	UpdatedAt time.Time `json:"updatedat" dynamodbav:"updatedat"`
}

func GetItem(ctx context.Context, id string) (*Fatwa, error) {
	key, err := attributevalue.Marshal(id)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]types.AttributeValue{
			"id": key,
		},
	}

	log.Printf("Calling Dynamodb with input: %v", input)
	result, err := db.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}
	log.Printf("Executed GetItem DynamoDb successfully. Result: %#v", result)

	if result.Item == nil {
		return nil, nil
	}

	fatwa := new(Fatwa)
	err = attributevalue.UnmarshalMap(result.Item, fatwa)
	if err != nil {
		return nil, err
	}

	return fatwa, nil
}

func ListItems(ctx context.Context) ([]Fatwa, error) {
	fatawa := make([]Fatwa, 0)
	var token map[string]types.AttributeValue

	for {
		input := &dynamodb.ScanInput{
			TableName:         aws.String(TableName),
			ExclusiveStartKey: token,
		}

		result, err := db.Scan(ctx, input)
		if err != nil {
			return nil, err
		}

		var fetchedFatawa []Fatwa
		err = attributevalue.UnmarshalListOfMaps(result.Items, &fetchedFatawa)
		if err != nil {
			return nil, err
		}

		fatawa = append(fatawa, fetchedFatawa...)
		token = result.LastEvaluatedKey
		if token == nil {
			break
		}
	}

	return fatawa, nil
}

func InsertItem(ctx context.Context, createFatwa CreateFatwa) (*Fatwa, error) {
	fatwa := Fatwa{
		Id:       uuid.NewString(),
		Title:    createFatwa.Title,
		Question: createFatwa.Question,
		Answer:   createFatwa.Answer,
		Link:     createFatwa.Link,
		Author:   createFatwa.Author,
		Topic:    createFatwa.Topic,
		Lang:     createFatwa.Lang,
	}

	item, err := attributevalue.MarshalMap(fatwa)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item:      item,
	}

	res, err := db.PutItem(ctx, input)
	if err != nil {
		return nil, err
	}

	err = attributevalue.UnmarshalMap(res.Attributes, &fatwa)
	if err != nil {
		return nil, err
	}

	return &fatwa, nil
}

func DeleteItem(ctx context.Context, id string) (*Fatwa, error) {
	key, err := attributevalue.Marshal(id)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(TableName),
		Key: map[string]types.AttributeValue{
			"id": key,
		},
		ReturnValues: types.ReturnValue(*aws.String("ALL_OLD")),
	}

	res, err := db.DeleteItem(ctx, input)
	if err != nil {
		return nil, err
	}

	if res.Attributes == nil {
		return nil, nil
	}

	fatwa := new(Fatwa)
	err = attributevalue.UnmarshalMap(res.Attributes, fatwa)
	if err != nil {
		return nil, err
	}

	return fatwa, nil
}

func UpdateItem(ctx context.Context, id string, updateFatwa UpdateFatwa) (*Fatwa, error) {
	key, err := attributevalue.Marshal(id)
	if err != nil {
		return nil, err
	}

	expr, err := expression.NewBuilder().WithUpdate(
		expression.Set(
			expression.Name("title"),
			expression.Value(updateFatwa.Title),
		).Set(
			expression.Name("question"),
			expression.Value(updateFatwa.Question),
		).Set(
			expression.Name("answer"),
			expression.Value(updateFatwa.Answer),
		).Set(
			expression.Name("link"),
			expression.Value(updateFatwa.Link),
		).Set(
			expression.Name("author"),
			expression.Value(updateFatwa.Author),
		).Set(
			expression.Name("topic"),
			expression.Value(updateFatwa.Topic),
		).Set(
			expression.Name("lang"),
			expression.Value(updateFatwa.Lang),
		),
	).WithCondition(
		expression.Equal(
			expression.Name("id"),
			expression.Value(id),
		),
	).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"id": key,
		},
		TableName:                 aws.String(TableName),
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       expr.Condition(),
		ReturnValues:              types.ReturnValue(*aws.String("ALL_NEW")),
	}

	res, err := db.UpdateItem(ctx, input)
	if err != nil {
		var smErr *smithy.OperationError
		if errors.As(err, &smErr) {
			var condCheckFailed *types.ConditionalCheckFailedException
			if errors.As(err, &condCheckFailed) {
				return nil, nil
			}
		}

		return nil, err
	}

	if res.Attributes == nil {
		return nil, nil
	}

	fatwa := new(Fatwa)
	err = attributevalue.UnmarshalMap(res.Attributes, fatwa)
	if err != nil {
		return nil, err
	}

	return fatwa, nil
}
