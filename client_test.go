package tsutsu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fireworq/fireworq/model"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"testing"
)

var FIREWORQ_PORT = os.Getenv("TEST_FIREWORQ_PORT")
var FIREWORQ_URL = fmt.Sprintf("http://localhost:%s", FIREWORQ_PORT)

func initRouting() (model.Routing, error) {
	routing := model.Routing{
		QueueName: "default",
		JobCategory: "test_category",
	}
	buf, err := json.Marshal(&routing)
	if err != nil {
		return model.Routing{}, err
	}
	r := bytes.NewReader(buf)
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(FIREWORQ_URL + "/routing/" + routing.JobCategory), r)
	if err != nil {
		return model.Routing{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return model.Routing{}, err
	}

	defer res.Body.Close()
	if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
		return model.Routing{}, err
	}

	return routing, nil
}

func TestTsutsu_Queues(t1 *testing.T) {
	type fields struct {
		baseURL string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []model.Queue
		wantErr bool
	}{
		{
			name:   "can parse json",
			fields: fields{baseURL: FIREWORQ_URL},
			want: []model.Queue{
				{
					Name:            "default",
					PollingInterval: 200,
					MaxWorkers:      20,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Tsutsu{
				baseURL: tt.fields.baseURL,
			}
			got, err := t.Queues()
			if (err != nil) != tt.wantErr {
				t1.Errorf("Queues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Queues() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTsutsu_CreateQueue(t1 *testing.T) {
	wantQueue := model.Queue{
		Name: "test_queue",
		PollingInterval: 100,
		MaxWorkers: 1,
	}
	defer func(){
		req, err := http.NewRequest(http.MethodDelete, FIREWORQ_URL + "/queue/" + wantQueue.Name, nil)
		if err != nil {
			t1.Error(err)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t1.Error(err)
		}
		defer res.Body.Close()
		if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
			t1.Error(err)
		}
	}()
	type fields struct {
		baseURL string
	}
	type args struct {
		name            string
		pollingInterval uint
		maxWorkers      uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Queue
		wantErr bool
	}{
		{
			name:   "should be created",
			fields: fields{baseURL: FIREWORQ_URL},
			args: args{
				name: wantQueue.Name,
				pollingInterval: wantQueue.PollingInterval,
				maxWorkers:      wantQueue.MaxWorkers},
			want:    wantQueue,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Tsutsu{
				baseURL: tt.fields.baseURL,
			}
			got, err := t.CreateQueue(tt.args.name, tt.args.pollingInterval, tt.args.maxWorkers)
			if (err != nil) != tt.wantErr {
				t1.Errorf("CreateQueue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("CreateQueue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTsutsu_DeleteQueue(t1 *testing.T) {
	queue := model.Queue{
		Name:            "be_delete",
		PollingInterval: 100,
		MaxWorkers:      20,
	}
	
	buf, err := json.Marshal(&queue)
	if err != nil {
		t1.Error(err)
	}
	
	r := bytes.NewReader(buf)
	req, err := http.NewRequest(http.MethodPut, FIREWORQ_URL + "/queue/" + queue.Name, r)
	if err != nil {
		t1.Error(err)
	}
	
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t1.Error(err)
	}
	
	defer res.Body.Close()
	if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
		t1.Error(err)
	}
	
	type fields struct {
		baseURL string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Queue
		wantErr bool
	}{
		{
			name:    "should be delete",
			fields:  fields{baseURL: FIREWORQ_URL},
			args:    args{name: queue.Name	},
			want:    model.Queue{
				Name: queue.Name,
				PollingInterval: queue.PollingInterval,
				MaxWorkers: queue.MaxWorkers,
			},
			wantErr: false,
		},
		{
			name:    "should be error",
			fields:  fields{baseURL: FIREWORQ_URL},
			args:    args{name: "not_found"},
			want:    model.Queue{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Tsutsu{
				baseURL: tt.fields.baseURL,
			}
			got, err := t.DeleteQueue(tt.args.name)
			if (err != nil) != tt.wantErr {
				t1.Errorf("DeleteQueue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("DeleteQueue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTsutsu_Routings(t1 *testing.T) {
	routing, err := initRouting()
	if err != nil {
		t1.Error(err)
	}

	defer func() {
		req, err := http.NewRequest(http.MethodDelete, FIREWORQ_URL + "/queue/" + routing.JobCategory, nil)
		if err != nil {
			t1.Error(err)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t1.Error(err)
		}

		defer res.Body.Close()
		if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
			t1.Error(err)
		}
	}()

	type fields struct {
		baseURL string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []model.Routing
		wantErr bool
	}{
		{
			name:    "should be add",
			fields:  fields{baseURL: FIREWORQ_URL},
			want:    []model.Routing{
				routing,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Tsutsu{
				baseURL: tt.fields.baseURL,
			}
			got, err := t.Routings()
			if (err != nil) != tt.wantErr {
				t1.Errorf("Routings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Routings() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTsutsu_Routing(t1 *testing.T) {
	routing, err := initRouting()
	if err != nil {
		t1.Error(err)
	}

	type fields struct {
		baseURL string
	}
	type args struct {
		jobCategory string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Routing
		wantErr bool
	}{
		{
			name:    "should be return routing",
			fields:  fields{baseURL: FIREWORQ_URL},
			args:    args{jobCategory: routing.JobCategory},
			want:    routing,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Tsutsu{
				baseURL: tt.fields.baseURL,
			}
			got, err := t.Routing(tt.args.jobCategory)
			if (err != nil) != tt.wantErr {
				t1.Errorf("Routing() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Routing() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTsutsu_CreateRouting(t1 *testing.T) {
	routing := model.Routing{
		QueueName: "default",
		JobCategory: "test_category",
	}

	type fields struct {
		baseURL string
	}
	type args struct {
		jobCategory string
		queueName   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Routing
		wantErr bool
	}{
		{
			name:    "should be created",
			fields:  fields{baseURL: FIREWORQ_URL},
			args:    args{
				jobCategory: routing.JobCategory,
				queueName: routing.QueueName,
			},
			want:    routing,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Tsutsu{
				baseURL: tt.fields.baseURL,
			}
			got, err := t.CreateRouting(tt.args.jobCategory, tt.args.queueName)
			if (err != nil) != tt.wantErr {
				t1.Errorf("CreateRouting() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("CreateRouting() got = %v, want %v", got, tt.want)
			}
		})
	}

	req, err := http.NewRequest(http.MethodDelete, FIREWORQ_URL + "/routing/" + routing.JobCategory, nil)
	if err != nil {
		t1.Error(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t1.Error(err)
	}
	defer res.Body.Close()

	if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
		t1.Error(err)
	}
}

func TestTsutsu_DeleteRouting(t1 *testing.T) {
	routing, err := initRouting()
	if err != nil {
		t1.Error(err)
	}

	type fields struct {
		baseURL string
	}
	type args struct {
		jobCategory string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Routing
		wantErr bool
	}{
		{
			name:    "should be delete",
			fields:  fields{baseURL: FIREWORQ_URL},
			args:    args{
				jobCategory: routing.JobCategory,
			},
			want:    routing,
			wantErr: false,
		},
		{
			name:    "should be error",
			fields:  fields{baseURL: FIREWORQ_URL},
			args:    args{jobCategory: "nothing"},
			want:    model.Routing{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Tsutsu{
				baseURL: tt.fields.baseURL,
			}
			got, err := t.DeleteRouting(tt.args.jobCategory)
			if (err != nil) != tt.wantErr {
				t1.Errorf("DeleteRouting() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("DeleteRouting() got = %v, want %v", got, tt.want)
			}
		})
	}
}