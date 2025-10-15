package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"order-service/internal/dto"
	"order-service/pkg/env"
	"strconv"
)

func GetProductByID(productID int64) (*dto.ApiProductResponse, error) {
	url := fmt.Sprintf("%s/api/products/%d", env.CONF.API.Product, productID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("unexpected api get product by id status: " + strconv.Itoa(resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	product := new(dto.ApiResponse[*dto.ApiProductResponse])
	if err := json.Unmarshal(body, &product); err != nil {
		return nil, err
	}
	return product.Data, nil
}
