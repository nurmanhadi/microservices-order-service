package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"order-service/internal/dto"
	"order-service/pkg/env"
)

func PaymentCreateTransaction(request dto.PaymentAddRequest) (*dto.PaymentResponse, error) {
	url := fmt.Sprintf("%s/api/payments", env.CONF.API.Payment)
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("unexpected api payment create transaction status code: %d", resp.StatusCode)
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	apiBody := new(dto.ApiResponse[*dto.PaymentResponse])
	err = json.Unmarshal(responseBody, apiBody)
	if err != nil {
		return nil, err
	}
	return apiBody.Data, nil
}
