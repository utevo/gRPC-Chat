package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"math/rand"
	"strconv"

	proto "github.com/utevo/gRPC-Chat/proto"
	"google.golang.org/grpc"
)

func createClient() proto.BroadcastClient {
	conn, err := grpc.Dial("localhost:5050", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := proto.NewBroadcastClient(conn)
	return client
}

func createUser() proto.User {
	userName := flag.String("name", "Anonymus", "The name of user")

	randInt := rand.Int31()
	randIntAsString := strconv.Itoa(int(randInt))
	idAsBytes := sha256.Sum256([]byte(*userName + randIntAsString))
	id := hex.EncodeToString(idAsBytes[:])

	user := proto.User{
		Id:   id,
		Name: *userName,
	}
	return user
}

func main() {
	user := createUser()

	client := createClient()
	stream, err := client.CreateStream(context.Background(), &user)
	if err != nil {
		panic(err)
	}
	_ = stream
}
