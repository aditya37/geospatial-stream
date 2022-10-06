package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	sseDeliver "github.com/aditya37/geospatial-stream/delivery/sse"
	wsDeliver "github.com/aditya37/geospatial-stream/delivery/websocket"
	"github.com/aditya37/geospatial-stream/infra"
	channel "github.com/aditya37/geospatial-stream/repository/channel"
	gcppubsub "github.com/aditya37/geospatial-stream/repository/gcppubsub"
	deviceCase "github.com/aditya37/geospatial-stream/usecase/device"
	geofencing "github.com/aditya37/geospatial-stream/usecase/geofencing"
	geostrckProto "github.com/aditya37/geospatial-tracking/proto"
	getenv "github.com/aditya37/get-env"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	// deviceLogsChan...
	deviceLogsChan := channel.NewStreamDeviceLogs(1024)
	go deviceLogsChan.Run()

	// geofence Detect...
	geofenceDetectChan := channel.NewStreamGeofenceDetect()
	go geofenceDetectChan.Run()

	// grpcDialOpt..
	grpcDialtOpt := []grpc.DialOption{}
	grpcDialtOpt = append(grpcDialtOpt, grpc.WithTransportCredentials(insecure.NewCredentials()))
	svcTrackingConn, err := grpc.Dial(
		getenv.GetString("SERVICE_TRACKING_HOST", "127.0.0.1:1112"),
		grpcDialtOpt...,
	)
	if err != nil {
		return nil, err
	}
	// connect to geotracking (gRPC)
	geotrackingSvc := geostrckProto.NewGeotrackingClient(svcTrackingConn)

	// usecase
	geofencingCase := geofencing.NewGefencingUsecase(
		gcppubsubRepo,
		avgMobilityChan,
		geofenceDetectChan,
	)
	deviceManagerCase := deviceCase.NewDeviceManagerUsecase(
		gcppubsubRepo,
		geotrackingSvc,
		deviceLogsChan,
	)
	// async...
	go geofencingCase.SubscribeGeofencingDetect(
		ctx,
		getenv.GetString("GEOFENCING_TOPIC_DETECT", "topic-detect"),
		"geospatial-stream",
	)
	// streamDeviceLogs...
	go deviceManagerCase.SubscribeDeviceLogs(
		ctx,
		getenv.GetString("DEVICE_LOG_STREAM_TOPIC", "stream-device-logs"),
		"geospatial-stream",
	)
	// Socket Deliver...
	geofencingWsDeliv := wsDeliver.NewWebsocketGeofencing(geofencingCase, avgMobilityChan)
	deviceWsDeliv := wsDeliver.NewDeviceWebsocketDeliver(
		getenv.GetInt("INTERVAL_TICKER_TRIGGER_DEVICE_LOGS", 2),
		deviceManagerCase,
		deviceLogsChan,
	)

	// SSE Deliver...
	geofenceSSEDeliver := sseDeliver.NewSSEGeofencingDelivery(geofenceDetectChan)

	// http handler...
	httpHandler := NewHttpHandler(
		geofencingWsDeliv,
		deviceWsDeliv,
		geofenceSSEDeliver,
	)

	return &service{
		close: func() {
			log.Println("take break broh,,, everyrhing has been closed")
			gcppubsubRepo.Close()
			svcTrackingConn.Close()
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
