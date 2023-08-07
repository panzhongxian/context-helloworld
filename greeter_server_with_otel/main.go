/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"reflect"
	"unsafe"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc/example/config"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/metadata"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v\n", in.GetName())
	fmt.Printf("ctx:\n\t%v\n\n", ctx)
	if dl, ok := ctx.Deadline(); ok {
		fmt.Printf("ctx.Deadline():\n\t %v\n\n", dl)
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {

		fmt.Println("metadata:")
		for k, v := range md {
			fmt.Printf("\t%v: %+v\n", k, v)
		}
	}
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	tp, _ := config.Init()
	defer tp.Shutdown(context.Background())

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func dumpContextInternals(ctx interface{}) map[any]any {
	ret := make(map[any]any)

	contextKeys := reflect.TypeOf(ctx).Elem()
	contextValues := reflect.ValueOf(ctx).Elem()
	if contextKeys.Kind() != reflect.Struct {
		return ret
	}

	var key, val any
	found := false
	for i := 0; i < contextValues.NumField(); i++ {
		reflectValue := contextValues.Field(i)
		reflectValue = reflect.NewAt(reflectValue.Type(), unsafe.Pointer(reflectValue.UnsafeAddr())).Elem()

		reflectField := contextKeys.Field(i)

		if reflectField.Name == "Context" {
			tmpMap := dumpContextInternals(reflectValue.Interface())
			for k, v := range tmpMap {
				ret[k] = v
			}
		} else if reflectField.Name == "cancelCtx" {
			tmpMap := dumpContextInternals(reflectValue.FieldByName("Context").Interface())
			for k, v := range tmpMap {
				ret[k] = v
			}
		} else if reflectField.Name == "key" {
			found = true
			key = reflectValue.Interface()
		} else if reflectField.Name == "val" {
			val = reflectValue.Interface()
		}
	}
	if found {
		ret[key] = val
	}
	return ret
}
