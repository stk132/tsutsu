package tsutsu

import (
	"errors"
	"fmt"
	"github.com/fireworq/fireworq/model"
	"io"
	"io/ioutil"
	"net/http"
)

type Tsutsu struct {
	baseURL string
}

func NewTsutsu(baseURL string) *Tsutsu {
	return &Tsutsu{baseURL}
}

func get(url string) (*httpBodyDecoder, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
			return nil, err
		}
		defer res.Body.Close()
		return nil, errors.New(fmt.Sprintf("status_code: %d", res.StatusCode))
	}

	return newHttpBodyDecoder(res.Body), nil
}

func (t *Tsutsu) Queues() ([]model.Queue, error) {
	decoder, err := get(t.baseURL + "/queues")
	if err != nil {
		return nil, err
	}

	defer decoder.Close()

	var queues []model.Queue
	if err := decoder.Decode(&queues); err != nil {
		return nil, err
	}

	return queues, nil
}

func (t *Tsutsu) Routings() ([]model.Routing, error) {
	decoder, err := get(t.baseURL + "/routings")
	if err != nil {
		return nil, err
	}

	defer decoder.Close()

	var routings []model.Routing
	if err := decoder.Decode(&routings); err != nil {
		return nil, err
	}

	return routings, err

}