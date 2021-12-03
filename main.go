package main

import (
	"log"
	"net"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/jalexanderII/urban-lamp/sensor"
	"github.com/jalexanderII/urban-lamp/server/sensorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const port = ":9092"

// server is a gRPC server it implements the methods defined by the server interface
type server struct {
	log   hclog.Logger
	Sensor *sensor.Sensor
}

// NewServer creates a new server
func NewServer(l hclog.Logger, s *sensor.Sensor) *server {
	return &server{l, s}
}

func (s *server) TempSensor(req *sensorpb.SensorRequest, stream sensorpb.Sensor_TempSensorServer) error {
	for {
		time.Sleep(time.Second * 5)

		temp := s.Sensor.GetTempSensor()
		err := stream.Send(&sensorpb.SensorResponse{Value: temp})
		if err != nil {
			log.Println("Error sending metric message ", err)
		}
	}
	return nil
}

func (s *server) HumiditySensor(req *sensorpb.SensorRequest, stream sensorpb.Sensor_HumiditySensorServer) error {
	for {
		time.Sleep(time.Second * 2)

		humd := s.Sensor.GetHumiditySensor()

		err := stream.Send(&sensorpb.SensorResponse{Value: humd})
		if err != nil {
			log.Println("Error sending metric message ", err)
		}
	}
	return nil
}

func main() {
	hlog := hclog.Default()

	// create a TCP socket for inbound server connections
	lstnr, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to start server:", err)
	}

	// setup and register currency service
	// create a new gRPC server, use WithInsecure to allow http connections
	grpcServer := grpc.NewServer()
	// create an instance of the sensor
	sns := sensor.NewSensor()
	sns.StartMonitoring()

	sensorService := NewServer(hlog, sns)
	sensorpb.RegisterSensorServer(grpcServer, sensorService)

	// register the reflection service which allows clients to determine the methods
	// for this gRPC service
	reflection.Register(grpcServer)

	// start service's server
	log.Println("starting rpc server on", port)
	if err := grpcServer.Serve(lstnr); err != nil {
		log.Fatal(err)
	}
}
