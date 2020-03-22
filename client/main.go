package main

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

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
	wait := sync.WaitGroup{}
	user := createUser()

	client := createClient()
	stream, err := client.CreateStream(context.Background(), &user)
	if err != nil {
		panic(err)
	}

	wait.Add(1)
	go func() {
		defer wait.Done()

		for {
			message, err := stream.Recv()
			if err != nil {
				fmt.Printf("Error during receiving message: %v", err)
				return
			}
			fmt.Print(message)
		}
	}()

	wait.Add(1)
	go func() {
		defer wait.Done()

		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {

			message := proto.Message{
				Id:        user.GetId(),
				Content:   scanner.Text(),
				Timestamp: time.Now().String(),
			}

			_, err := client.BroadcastMessage(context.Background(), &message)
			if err != nil {
				fmt.Printf("Error during sending message: %v", err)
				return
			}

		}
	}()

	// message := proto.Message{
	// 	Id:        "1",
	// 	Content:   "Hello",
	// 	Timestamp: time.Now().String(),
	// }
	// if _, err := client.BroadcastMessage(context.Background(), &message); err != nil {
	// 	panic(err)
	// }

	wait.Wait()
}
