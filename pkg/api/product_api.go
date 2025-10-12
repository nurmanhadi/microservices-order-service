package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"order-service/internal/dto"
	"strconv"
)

type ProductAPI interface {
	GetProductByID(productID int64) (*dto.ProductResponse, error)
}
type productAPI struct {
	baseURL string
}

func NewProductAPI(baseURL string) ProductAPI {
	return &productAPI{baseURL: baseURL}
}
func (a *productAPI) GetProductByID(productID int64) (*dto.ProductResponse, error) {
	url := fmt.Sprintf("%s/api/products/%d", a.baseURL, productID)
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
	product := new(dto.ProductResponse)
	if err := json.Unmarshal(body, &product); err != nil {
		return nil, err
	}
	return product, nil
}
