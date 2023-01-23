package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

type UpdateFatwa struct {
	Title    string `json:"title"`
	Question string `json:"question" validate:"required"`
	Answer   string `json:"answer" validate:"required"`
	Link     string `json:"link"`
	Author   string `json:"author" validate:"required"`
	Topic    string `json:"topic"`
	Lang     string `json:"lang" validate:"required"`
}

type CreateFatwa struct {
	//Id       string `json:"id" validate:"required"`
	Title    string `json:"title"`
	Question string `json:"question" validate:"required"`
	Answer   string `json:"answer" validate:"required"`
	Link     string `json:"link"`
	Author   string `json:"author" validate:"required"`
	Topic    string `json:"topic"`
	Lang     string `json:"lang" validate:"required"`
	//	CreatedAt time.Time `json:"createdat" dynamodbav:"createdat"`
	//	UpdatedAt time.Time `json:"updatedat" dynamodbav:"updatedat"`
}

var validate *validator.Validate = validator.New()

func Router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received req %#v", req)

	switch req.HTTPMethod {
	case "GET":
		return processGet(ctx, req)
	case "POST":
		return processPost(ctx, req)
	case "DELETE":
		return processDelete(ctx, req)
	case "PUT":
		return processPut(ctx, req)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

func processGet(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := req.PathParameters["id"]
	if !ok {
		return processGetFatawa(ctx)
	} else {
		return processGetFatwa(ctx, id)
	}
}

func processGetFatwa(ctx context.Context, id string) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received GET fatwa request with id = %s", id)

	fatwa, err := GetItem(ctx, id)
	if err != nil {
		return serverError(err)
	}

	if fatwa == nil {
		return clientError(http.StatusNotFound)
	}

	json, err := json.Marshal(fatwa)
	if err != nil {
		return serverError(err)
	}
	log.Printf("Successfully fetched fatwa item %s", json)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(json),
	}, nil
}

func processGetFatawa(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	log.Print("Received GET fatawa request")

	fatawa, err := ListItems(ctx)
	if err != nil {
		return serverError(err)
	}

	json, err := json.Marshal(fatawa)
	if err != nil {
		return serverError(err)
	}
	log.Printf("Successfully fetched fatawa: %s", json)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(json),
	}, nil
}

func processPost(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var createFatwa CreateFatwa
	err := json.Unmarshal([]byte(req.Body), &createFatwa)
	if err != nil {
		log.Printf("Can't unmarshal body: %v", err)
		return clientError(http.StatusUnprocessableEntity)
	}

	err = validate.Struct(&createFatwa)
	if err != nil {
		log.Printf("Invalid body: %v", err)
		return clientError(http.StatusBadRequest)
	}
	log.Printf("Received POST request with item: %+v", createFatwa)

	res, err := InsertItem(ctx, createFatwa)
	if err != nil {
		return serverError(err)
	}
	log.Printf("Inserted new fatwa: %+v", res)

	json, err := json.Marshal(res)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(json),
		Headers: map[string]string{
			"Location": fmt.Sprintf("/fatwa/%s", res.Id),
		},
	}, nil
}

func processDelete(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := req.PathParameters["id"]
	if !ok {
		return clientError(http.StatusBadRequest)
	}
	log.Printf("Received DELETE request with id = %s", id)

	fatwa, err := DeleteItem(ctx, id)
	if err != nil {
		return serverError(err)
	}

	if fatwa == nil {
		return clientError(http.StatusNotFound)
	}

	json, err := json.Marshal(fatwa)
	if err != nil {
		return serverError(err)
	}
	log.Printf("Successfully deleted fatwa item %+v", fatwa)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(json),
	}, nil
}

func processPut(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := req.PathParameters["id"]
	if !ok {
		return clientError(http.StatusBadRequest)
	}

	var updateFatwa UpdateFatwa
	err := json.Unmarshal([]byte(req.Body), &updateFatwa)
	if err != nil {
		log.Printf("Can't unmarshal body: %v", err)
		return clientError(http.StatusUnprocessableEntity)
	}

	err = validate.Struct(&updateFatwa)
	if err != nil {
		log.Printf("Invalid body: %v", err)
		return clientError(http.StatusBadRequest)
	}
	log.Printf("Received PUT request with item: %+v", updateFatwa)

	res, err := UpdateItem(ctx, id, updateFatwa)
	if err != nil {
		return serverError(err)
	}

	if res == nil {
		return clientError(http.StatusNotFound)
	}

	log.Printf("Updated fatwa: %+v", res)

	json, err := json.Marshal(res)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(json),
		Headers: map[string]string{
			"Location": fmt.Sprintf("/fatwa/%s", res.Id),
		},
	}, nil
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {

	return events.APIGatewayProxyResponse{
		Body:       http.StatusText(status),
		StatusCode: status,
	}, nil
}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	log.Println(err.Error())

	return events.APIGatewayProxyResponse{
		Body:       http.StatusText(http.StatusInternalServerError),
		StatusCode: http.StatusInternalServerError,
	}, nil
}
