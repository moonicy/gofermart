package accrual

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/moonicy/gofermart/internal/models"
	"io"
	"log"
	"net/http"
	"time"
)

var ErrNotFound = errors.New("not found")
var ErrTooManyRequests = errors.New("rate limited")

type Client struct {
	cl   http.Client
	host string
}

func NewClient(host string) *Client {
	return &Client{cl: http.Client{Timeout: 5 * time.Second}, host: host}
}

func (cl *Client) GetOrderInfo(number string) (models.Order, error) {
	uri := fmt.Sprintf("%s/api/orders/%s", cl.host, number)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Print(err)
		return models.Order{}, err
	}
	resp, err := cl.cl.Do(req)
	if err != nil {
		log.Print(err)
		return models.Order{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNoContent {
		return models.Order{}, ErrNotFound
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return models.Order{}, ErrTooManyRequests
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Order{}, err
	}
	var order models.Order
	if err = json.Unmarshal(body, &order); err != nil {
		return models.Order{}, err
	}
	return order, nil
}
