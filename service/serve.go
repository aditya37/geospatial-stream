package service

import (
	"fmt"
	"log"
	"net"
	"net/http"

	getenv "github.com/aditya37/get-env"
	"github.com/soheilhy/cmux"
)

func serve(httpHandler http.Handler) error {
	lis, err := net.Listen(
		"tcp",
		fmt.Sprintf(":%s", getenv.GetString("PORT", "8887")),
	)
	if err != nil {
		return err
	}

	m := cmux.New(lis)
	httpl := m.Match(cmux.HTTP1Fast())

	http := &http.Server{
		Handler: httpHandler,
	}
	log.Println("Geospatial stream service run on port", getenv.GetString("PORT", "8887"))
	// serve grpc
	//go grpc.Serve(grpcl)
	// serve http
	go http.Serve(httpl)
	return m.Serve()
}
