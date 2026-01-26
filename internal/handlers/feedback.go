package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/care-giver-app/care-giver-api/internal/response"
	"go.uber.org/zap"
)

const (
	submitFeedback = "submit feedback"
)

type FeedbackRequest struct {
	Message string `json:"message" validate:"required"`
}

type FeedbackResponse struct {
	Status string `json:"status"`
}

type Notification struct {
	NotificationType string   `json:"notification_type"`
	Channel          []string `json:"channel"`
	ExecutionData    any      `json:"execution_data"`
}

type FeedbackNotification struct {
	Email   string `json:"email"`
	Message string `json:"message"`
}

func HandleFeedbackRequest(ctx context.Context, params HandlerParams) (events.APIGatewayProxyResponse, error) {
	params.AppCfg.Logger.Sugar().Infof(handlerStart, submitFeedback)

	var feedbackRequest FeedbackRequest
	err := readRequestBody(params.Request.Body, &feedbackRequest)
	if err != nil {
		params.AppCfg.Logger.Error(requestBodyError, zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	if params.AppCfg.FeedbackQueueURL == "" {
		params.AppCfg.Logger.Error("feedback queue URL not configured")
		return response.CreateInternalServerErrorResponse(), nil
	}

	sqsClient := sqs.NewFromConfig(params.AppCfg.AWSConfig)

	sqsMessage := Notification{
		NotificationType: "feedback",
		Channel:          []string{"email"},
		ExecutionData: FeedbackNotification{
			Email:   "twilliams0095@gmail.com",
			Message: feedbackRequest.Message,
		},
	}

	messageBody, err := json.Marshal(sqsMessage)
	if err != nil {
		params.AppCfg.Logger.Error("error marshaling feedback message", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(params.AppCfg.FeedbackQueueURL),
		MessageBody: aws.String(string(messageBody)),
	}

	_, err = sqsClient.SendMessage(ctx, input)
	if err != nil {
		params.AppCfg.Logger.Error("error sending message to SQS", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	params.AppCfg.Logger.Sugar().Infof(handlerSuccessful, submitFeedback)

	resp := FeedbackResponse{
		Status: response.Success,
	}

	return response.FormatResponse(resp, http.StatusOK), nil
}
