package tsutsu

import (
	"bytes"
	"encoding/json"
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

func do(req *http.Request) (*httpBodyDecoder, error) {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		return nil, errors.New(fmt.Sprintf("status_code: %d", res.StatusCode))
	}

	return newHttpBodyDecoder(res.Body), nil
}

func put(url string, r io.Reader) (*httpBodyDecoder, error) {
	req, err := http.NewRequest(http.MethodPut, url, r)
	if err != nil {
		return nil, err
	}

	return do(req)
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

func (t *Tsutsu) Queue(name string) (model.Queue, error) {
	decoder, err := get(fmt.Sprintf("%s/queue/%s", t.baseURL, name))
	if err != nil {
		return model.Queue{}, err
	}

	defer decoder.Close()
	var queue model.Queue

	if err := decoder.Decode(&queue); err != nil {
		return model.Queue{}, err
	}

	return queue, nil
}

func (t *Tsutsu) Stats(queueName string) (QueueStats, error) {
	url := fmt.Sprintf("%s/queue/%s/stats", t.baseURL, queueName)
	decoder, err := get(url)
	if err != nil {
		return QueueStats{}, err
	}

	decoder.Close()

	var stats QueueStats
	if err := decoder.Decode(&stats); err != nil {
		return QueueStats{}, err
	}

	return stats, nil
}

func (t *Tsutsu) Node(queueName string) (NodeInfo, error) {
	url := fmt.Sprintf("%s/queue/%s/node", t.baseURL, queueName)
	decoder, err := get(url)
	if err != nil {
		return NodeInfo{}, err
	}

	defer decoder.Close()

	var node NodeInfo
	if err := decoder.Decode(&node); err != nil {
		return NodeInfo{}, err
	}
	return node, nil
}

func (t *Tsutsu) CreateQueue(name string, pollingInterval, maxWorkers uint) (model.Queue, error) {
	m := model.Queue{
		Name:            name,
		PollingInterval: pollingInterval,
		MaxWorkers:      maxWorkers,
	}

	buf, err := json.Marshal(&m)
	if err != nil {
		return model.Queue{}, err
	}

	r := bytes.NewReader(buf)
	url := fmt.Sprintf("%s/queue/%s", t.baseURL, name)
	req, err := http.NewRequest(http.MethodPut, url, r)
	if err != nil {
		return model.Queue{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return model.Queue{}, err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		decoder := json.NewDecoder(res.Body)
		var queue model.Queue
		if err := decoder.Decode(&queue); err != nil {
			return model.Queue{}, err
		}
		return queue, err
	} else {
		return model.Queue{}, errors.New(fmt.Sprintf("create queue error. status_code: %d", res.StatusCode))
	}
}

func (t *Tsutsu) DeleteQueue(name string) (model.Queue, error) {
	url := fmt.Sprintf("%s/queue/%s", t.baseURL, name)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return model.Queue{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return model.Queue{}, err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		decoder := json.NewDecoder(res.Body)
		var queue model.Queue
		if err := decoder.Decode(&queue); err != nil {
			return model.Queue{}, err
		}
		return queue, nil
	} else {
		return model.Queue{}, errors.New(fmt.Sprintf("delete error. status_code: %d", res.StatusCode))
	}
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

func (t *Tsutsu) Routing(jobCategory string) (model.Routing, error) {
	decoder, err := get(fmt.Sprintf("%s/routing/%s", t.baseURL, jobCategory))
	if err != nil {
		return model.Routing{}, err
	}

	defer decoder.Close()

	var routing model.Routing
	if err := decoder.Decode(&routing); err != nil {
		return model.Routing{}, err
	}

	return routing, nil
}

func (t *Tsutsu) CreateRouting(jobCategory, queueName string) (model.Routing, error) {
	rt := model.Routing{
		QueueName:   queueName,
		JobCategory: jobCategory,
	}

	buf, err := json.Marshal(&rt)
	if err != nil {
		return model.Routing{}, err
	}

	r := bytes.NewReader(buf)
	url := fmt.Sprintf("%s/routing/%s", t.baseURL, jobCategory)
	req, err := http.NewRequest(http.MethodPut, url, r)
	if err != nil {
		return model.Routing{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return model.Routing{}, err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		decoder := json.NewDecoder(res.Body)
		var routing model.Routing
		if err := decoder.Decode(&routing); err != nil {
			return model.Routing{}, err
		}
		return routing, nil
	} else {
		return model.Routing{}, errors.New(fmt.Sprintf("create routing error. status_code: %d", res.StatusCode))
	}
}

func (t *Tsutsu) DeleteRouting(jobCategory string) (model.Routing, error) {
	url := fmt.Sprintf("%s/routing/%s", t.baseURL, jobCategory)
	decoder, err := put(url, nil)
	if err != nil {
		return model.Routing{}, err
	}

	defer decoder.Close()

	var routing model.Routing
	if err := decoder.Decode(&decoder); err != nil {
		return model.Routing{}, err
	}
	return routing, nil
}
