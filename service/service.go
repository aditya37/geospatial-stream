package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	wsDeliver "github.com/aditya37/geospatial-stream/delivery/websocket"
	"github.com/aditya37/geospatial-stream/infra"
	channel "github.com/aditya37/geospatial-stream/repository/channel"
	gcppubsub "github.com/aditya37/geospatial-stream/repository/gcppubsub"
	geofencing "github.com/aditya37/geospatial-stream/usecase/geofencing"
	getenv "github.com/aditya37/get-env"
)

type Service interface {
	Run()
}
type service struct {
	httpHandler *httpHandler
	close       func()
}

func NewService() (Service, error) {
	ctx := context.Background()

	// pubsub instance
	if err := infra.NewGCPPubsub(
		ctx,
		getenv.GetString("GCP_PROJECT_ID", ""),
	); err != nil {
		return nil, err
	}
	gcpPubsubClient := infra.GetGCPPubsubInstance()
	if gcpPubsubClient == nil {
		if err := infra.NewGCPPubsub(
			ctx,
			getenv.GetString("GCP_PROJECT_ID", ""),
		); err != nil {
			return nil, err
		}
		gcpPubsubClient = infra.GetGCPPubsubInstance()
	}

	// repo
	gcppubsubRepo, err := gcppubsub.NewGCPPubsubRepo(gcpPubsubClient)
	if err != nil {
		return nil, err
	}

	// channel repo
	avgMobilityChan := channel.NewStreamMobilityAvg()
	go avgMobilityChan.Run()

	// usecase
	geofencingCase := geofencing.NewGefencingUsecase(
		gcppubsubRepo,
		avgMobilityChan,
	)

	// async...
	go geofencingCase.SubscribeGeofencingDetect(
		ctx,
		getenv.GetString("GEOFENCING_TOPIC_DETECT", "topic-detect"),
		"geospatial-stream",
	)

	// deliver...
	geofencingWsDeliv := wsDeliver.NewWebsocketGeofencing(geofencingCase, avgMobilityChan)

	// http handler...
	httpHandler := NewHttpHandler(geofencingWsDeliv)

	return &service{
		close: func() {
			log.Println("take break broh,,, everyrhing has been closed")
			gcppubsubRepo.Close()
		},
		httpHandler: httpHandler,
	}, nil
}

func (s *service) Run() {
	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		errs <- fmt.Errorf("%s", <-c)
		defer s.close()
	}()
	go func() {
		errs <- serve(s.httpHandler.Handler())
	}()
	// blocking...
	log.Fatalf("Stop server with error detail: %v", <-errs)
}
