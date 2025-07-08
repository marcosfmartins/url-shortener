package main

import (
	"context"
	"errors"
	"fmt"

	ihttp "github.com/marcosfmartins/url_shortener/internal/http"
	"github.com/marcosfmartins/url_shortener/internal/service"
	mgo "github.com/marcosfmartins/url_shortener/internal/storage/mongo"
	"github.com/marcosfmartins/url_shortener/internal/storage/redis"
	"github.com/marcosfmartins/url_shortener/internal/stream/kafka"
	"github.com/marcosfmartins/url_shortener/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	envServerPort = "SERVER_PORT"
	envMongoURL   = "MONGO_URL"
	envRedisURL   = "REDIS_URL"
	envKafkaURL   = "KAFKA_URL"
)

func main() {
	log := logger.NewZerologAdapter()

	log.Debug("open connection with mongo")
	mCli, DB, err := mgo.NewClient(getEnv(envMongoURL))
	if err != nil {
		log.Err(err).Fatal("failed to connect to mongo")
	}
	defer func() {
		err := mCli.Disconnect(context.Background())
		if err != nil {
			log.Err(err).Error("failed to disconnect")
		}
	}()

	publisher := kafka.NewPublisher(getKafkaEnv(), "url_statistics")
	defer func() {
		err := publisher.Close()
		if err != nil {
			log.Err(err).Error("failed to close publisher")
		}
	}()
	consumer := kafka.NewConsumer(getKafkaEnv(), "url_statistics", "url_statistics_group")
	defer func() {
		err = consumer.Close()
		if err != nil {
			log.Err(err).Error("failed to close consumer")
		}
	}()

	urlStorage := mgo.NewURLStorage(DB)
	cacheStorage := redis.NewCache(getEnv(envRedisURL))
	defer func() {
		err := cacheStorage.Close()
		if err != nil {
			log.Err(err).Error("failed to close cache")
		}
	}()

	urlService := service.NewURLService(log, urlStorage, cacheStorage, publisher)
	statisticService := service.NewStatisticService(urlService, consumer)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go statisticService.Processor(ctx)

	handler := ihttp.NewHandler(log, urlService)

	srv := ihttp.NewServer(handler, getEnv(envServerPort))

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Err(err).Fatal("failed to start server")
		}
	}()

	log.Info(fmt.Sprintf("Server is up and running in port %s", getEnv(envServerPort)))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Err(err).Error("failed to stop server")
	}

	log.Info("Server stopped")
}

func getEnv(env string) string {
	return os.Getenv(env)
}

func getKafkaEnv() []string {
	data := getEnv(envKafkaURL)
	return strings.Split(data, ",")
}
