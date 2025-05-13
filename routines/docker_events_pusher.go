package routines

//
//import (
//	"bytes"
//	"context"
//	"encoding/json"
//	"fmt"
//	"time"
//
//	"github.com/cego/dopetes/model"
//	"github.com/cego/go-lib"
//	"github.com/elastic/go-elasticsearch/v9"
//	"github.com/elastic/go-elasticsearch/v9/esutil"
//)
//
//func PushDockerEventsToElastic(ctx context.Context, m *model.Model, logger cego.Logger) {
//	elasticsearchConfig := m.GetElasticsearchConfig()
//	es, err := elasticsearch.NewClient(elasticsearch.Config{
//		Addresses: elasticsearchConfig.Hosts,
//		APIKey:    elasticsearchConfig.ApiKey,
//		Username:  elasticsearchConfig.Username,
//		Password:  elasticsearchConfig.Password,
//	})
//	if err != nil {
//		logger.Error(err.Error())
//		return
//	}
//
//	var bulkErrors []error
//
//	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
//		Index:  "filebeat-prod-8.18.1",
//		Client: es,
//
//		OnError: func(ctx context.Context, err error) {
//			bulkErrors = append(bulkErrors, err)
//		},
//	})
//	if err != nil {
//		logger.Error(err.Error())
//		return
//	}
//
//	for _, dockerPullEvent := range m.GetDockerPullEvents() {
//		var b []byte
//		b, err = json.Marshal(&dockerPullEvent)
//		if err != nil {
//			logger.Error(err.Error())
//			return
//		}
//
//		err = bi.Add(ctx, esutil.BulkIndexerItem{
//			Action: "create",
//			Body:   bytes.NewReader(b),
//			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
//				if err != nil {
//					logger.Error(err.Error())
//				} else {
//					logger.Error(fmt.Sprintf("%s %s", res.Error.Type, res.Error.Reason))
//				}
//			},
//		})
//		if err != nil {
//			logger.Error(err.Error())
//			return
//		}
//	}
//
//	err = bi.Close(ctx)
//	if err != nil {
//		logger.Error(err.Error())
//		return
//	}
//
//	if len(bulkErrors) > 0 {
//		errorText := ""
//		for bulkError := range bulkErrors {
//			errorText += bulkErrors[bulkError].Error()
//		}
//		logger.Error(errorText)
//		return
//	}
//
//	biStats := bi.Stats()
//	logger.Debug(fmt.Sprintf("Pushed %d to %s", biStats.NumCreated, elasticsearchConfig.Hosts))
//
//	m.ClearDockerPullEvents()
//}
//
//func StartDockerEventsPusher(ctx context.Context, m *model.Model, logger cego.Logger) {
//	go func() {
//		for {
//			PushDockerEventsToElastic(ctx, m, logger)
//			time.Sleep(2000 * time.Millisecond)
//		}
//	}()
//}
