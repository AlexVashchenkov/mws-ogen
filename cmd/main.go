package main

import (
	"context"
	"fmt"
	"log"
	api "mws-ogen/gen"
	serv "mws-ogen/server"
	"net/http"
	"time"
)

func main() {
	h := serv.NewUserStorage()
	srv, err := api.NewServer(h)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Println("listening :8080")
		log.Fatal(http.ListenAndServe(":8080", srv))
	}()

	time.Sleep(200 * time.Millisecond)

	cl, err := api.NewClient("http://localhost:8080", api.WithClient(http.DefaultClient))
	if err != nil {
		log.Fatal(err)
	}

	u := api.User{Name: "Alexey", Email: api.NewOptString("litmo@litmo.ru")}
	if _, err := cl.AddUser(context.Background(), &u); err != nil {
		log.Fatal(err)
	}

	lst, err := cl.ListUsers(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	switch v := lst.(type) {
	case *api.ListUsersOKApplicationJSON:
		fmt.Println("users:", []api.User(*v))
	default:
		fmt.Printf("unexpected list response: %#v\n", v)
	}
}
