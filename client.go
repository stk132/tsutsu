package tsutsu

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fireworq/fireworq/model"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type Tsutsu struct {
	baseURL string
}

func NewTsutsu(baseURL string) *Tsutsu {
	return &Tsutsu{baseURL}
}

func get(url string) (*httpBodyDecoder, error) {
	return getWithContext(context.Background(), url)
}

func getWithContext(ctx context.Context, url string) (*httpBodyDecoder, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return do(req)
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

func put(uri string, r io.Reader) (*httpBodyDecoder, error) {
	return putWithContext(context.Background(), uri, r)
}

func putWithContext(ctx context.Context, uri string, r io.Reader) (*httpBodyDecoder, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, r)
	if err != nil {
		return nil, err
	}

	return do(req)
}

func httpDelete(uri string) (*httpBodyDecoder, error) {
	return httpDeleteWithContext(context.Background(), uri)
}

func httpDeleteWithContext(ctx context.Context, uri string) (*httpBodyDecoder, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, err
	}

	return do(req)
}

func (t *Tsutsu) Queues() ([]model.Queue, error) {
	return t.QueuesWithContext(context.Background())
}

func (t *Tsutsu) QueuesWithContext(ctx context.Context) ([]model.Queue, error) {
	decoder, err := getWithContext(ctx, t.baseURL+"/queues")
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
	return t.QueueWithContext(context.Background(), name)
}

func (t *Tsutsu) QueueWithContext(ctx context.Context, name string) (model.Queue, error) {
	decoder, err := getWithContext(ctx, fmt.Sprintf("%s/queue/%s", t.baseURL, name))
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
	return t.StatsWithContext(context.Background(), queueName)
}

func (t *Tsutsu) StatsWithContext(ctx context.Context, queueName string) (QueueStats, error) {
	uri := fmt.Sprintf("%s/queue/%s/stats", t.baseURL, queueName)
	decoder, err := getWithContext(ctx, uri)
	if err != nil {
		return QueueStats{}, err
	}

	defer decoder.Close()

	var stats QueueStats
	if err := decoder.Decode(&stats); err != nil {
		return QueueStats{}, err
	}

	return stats, nil
}

func (t *Tsutsu) Node(queueName string) (NodeInfo, error) {
	return t.NodeWithContext(context.Background(), queueName)
}

func (t *Tsutsu) NodeWithContext(ctx context.Context, queueName string) (NodeInfo, error) {
	uri := fmt.Sprintf("%s/queue/%s/node", t.baseURL, queueName)
	decoder, err := getWithContext(ctx, uri)
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
	return t.CreateQueueWithContext(context.Background(), name, pollingInterval, maxWorkers)
}

func (t *Tsutsu) CreateQueueWithContext(ctx context.Context, name string, pollingInterval, maxWorkers uint) (model.Queue, error) {
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
	uri := fmt.Sprintf("%s/queue/%s", t.baseURL, name)
	decoder, err := putWithContext(ctx, uri, r)
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

func (t *Tsutsu) DeleteQueue(name string) (model.Queue, error) {
	return t.DeleteQueueWithContext(context.Background(), name)
}

func (t *Tsutsu) DeleteQueueWithContext(ctx context.Context, name string) (model.Queue, error) {
	uri := fmt.Sprintf("%s/queue/%s", t.baseURL, name)
	decoder, err := httpDeleteWithContext(ctx, uri)
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

func (t *Tsutsu) Routings() ([]model.Routing, error) {
	return t.RoutingsWithContext(context.Background())
}

func (t *Tsutsu) RoutingsWithContext(ctx context.Context) ([]model.Routing, error) {
	decoder, err := getWithContext(ctx, t.baseURL+"/routings")
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
	return t.RoutingWithContext(context.Background(), jobCategory)
}

func (t *Tsutsu) RoutingWithContext(ctx context.Context, jobCategory string) (model.Routing, error) {
	decoder, err := getWithContext(ctx, fmt.Sprintf("%s/routing/%s", t.baseURL, jobCategory))
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
	return t.CreateRoutingWithContext(context.Background(), jobCategory, queueName)
}

func (t *Tsutsu) CreateRoutingWithContext(ctx context.Context, jobCategory, queueName string) (model.Routing, error) {
	rt := model.Routing{
		QueueName:   queueName,
		JobCategory: jobCategory,
	}

	buf, err := json.Marshal(&rt)
	if err != nil {
		return model.Routing{}, err
	}

	r := bytes.NewReader(buf)
	uri := fmt.Sprintf("%s/routing/%s", t.baseURL, jobCategory)
	decoder, err := putWithContext(ctx, uri, r)
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

func (t *Tsutsu) DeleteRouting(jobCategory string) (model.Routing, error) {
	return t.DeleteRoutingWithContext(context.Background(), jobCategory)
}

func (t *Tsutsu) DeleteRoutingWithContext(ctx context.Context, jobCategory string) (model.Routing, error) {
	uri := fmt.Sprintf("%s/routing/%s", t.baseURL, jobCategory)
	decoder, err := httpDeleteWithContext(ctx, uri)
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

func (t *Tsutsu) Job() *JobInspector {
	return newJobInspector(t)
}

type JobInspector struct {
	client *Tsutsu
	limit  uint
	cursor string
	order  string
}

func newJobInspector(client *Tsutsu) *JobInspector {
	return &JobInspector{
		client: client,
		limit:  100,
		cursor: "",
		order:  "desc",
	}
}

func (j *JobInspector) Limit(limit uint) *JobInspector {
	j.limit = limit
	return j
}

func (j *JobInspector) Asc() *JobInspector {
	j.order = "asc"
	return j
}

func (j *JobInspector) Desc() *JobInspector {
	j.order = "desc"
	return j
}

func (j *JobInspector) Cursor(cursor string) *JobInspector {
	j.cursor = cursor
	return j
}

func (j *JobInspector) queryString() string {
	query := url.Values{}
	query.Set("limit", strconv.Itoa(int(j.limit)))
	query.Set("order", j.order)
	if j.cursor != "" {
		query.Set("cursor", j.cursor)
	}
	return query.Encode()
}

func (j *JobInspector) do(ctx context.Context, uri string) (JobsInfo, error) {
	decoder, err := getWithContext(ctx, uri)
	if err != nil {
		return JobsInfo{}, err
	}

	defer decoder.Close()
	var jobsInfo JobsInfo

	if err := decoder.Decode(&jobsInfo); err != nil {
		return JobsInfo{}, err
	}

	return jobsInfo, nil
}

func (j *JobInspector) Grabbed(queueName string) (JobsInfo, error) {
	return j.GrabbedWithContext(context.Background(), queueName)
}

func (j *JobInspector) GrabbedWithContext(ctx context.Context, queueName string) (JobsInfo, error) {
	uri := fmt.Sprintf("%s/queue/%s/grabbed?%s", j.client.baseURL, queueName, j.queryString())
	return j.do(ctx, uri)
}

func (j *JobInspector) Waiting(queueName string) (JobsInfo, error) {
	return j.WaitingWithContext(context.Background(), queueName)
}

func (j *JobInspector) WaitingWithContext(ctx context.Context, queueName string) (JobsInfo, error) {
	uri := fmt.Sprintf("%s/queue/%s/waiting?%s", j.client.baseURL, queueName, j.queryString())
	return j.do(ctx, uri)
}

func (j *JobInspector) DeferredWithContext(ctx context.Context, queueName string) (JobsInfo, error) {
	uri := fmt.Sprintf("%s/queue/%s/deferred?%s", j.client.baseURL, queueName, j.queryString())
	return j.do(ctx, uri)
}

func (j *JobInspector) Deferred(queueName string) (JobsInfo, error) {
	return j.DeferredWithContext(context.Background(), queueName)
}

func (j *JobInspector) FailedWithContext(ctx context.Context, queueName string) (FailedJobsInfo, error) {
	uri := fmt.Sprintf("%s/queue/%s/failed?%s", j.client.baseURL, queueName, j.queryString())
	decoder, err := getWithContext(ctx, uri)
	if err != nil {
		return FailedJobsInfo{}, err
	}

	defer decoder.Close()
	var failedJobsInfo FailedJobsInfo

	if err := decoder.Decode(&failedJobsInfo); err != nil {
		return FailedJobsInfo{}, err
	}

	return failedJobsInfo, nil
}

func (j *JobInspector) Failed(queueName string) (FailedJobsInfo, error) {
	return j.FailedWithContext(context.Background(), queueName)
}
