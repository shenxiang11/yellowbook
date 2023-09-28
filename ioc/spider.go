package ioc

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"yellowbook/internal/domain"
	"yellowbook/internal/service"
)

type Spider struct {
	srv service.IArticleService
}

func NewSpider(srv service.IArticleService) *Spider {
	return &Spider{srv: srv}
}

func (s *Spider) Run() Spider {
	type Message struct {
		Title     string   `json:"title"`
		Content   string   `json:"content"`
		ImageList []string `json:"imageList"`
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{
			"localhost:9092",
			"localhost:9093",
			"localhost:9094",
			"localhost:9095",
		},
		Topic:       "weimi_luyao",
		StartOffset: kafka.SeekStart,
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			continue
		}

		var message Message
		err = json.Unmarshal(m.Value, &message)
		if err != nil {
			log.Println("unmarshal failed")
			continue
		}

		_, err = s.srv.Save(context.Background(), domain.Article{
			Title:     message.Title,
			Content:   message.Content,
			ImageList: message.ImageList,
			Author: domain.Author{
				Id: 1,
			},
		})
		if err != nil {
			log.Println("unmarshal failed")
			continue
		}

		//fmt.Println("", string(m.Key), string(m.Value))
	}
}
