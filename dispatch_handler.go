package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bigroom/vision/models"
	"github.com/bigroom/zombies"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
)

func messageLoop() {
	/*
		for {
			log.Println("Waiting on message...")
			m := <-messages
			log.Printf("Dispatching message '%v' to channel with key: '%v'", m.Content, m.Key())
			for _, u := range clients[m.Key()] {
				fmt.Println("Writing message")

					err := u.WSConn.WriteJSON(m)
					if err != nil {
						fmt.Println("error (sending message):", err)
					}
			}
		}

	*/
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func dispatchHandler(w http.ResponseWriter, r *http.Request, t *jwt.Token) {
	fmt.Println(t.Claims["id"])

	u, err := models.FetchUser("id", t.Claims["id"])
	if err != nil {
		log.Println("COuldnt get user")
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("couldn't upgrade:", err)
		return
	}

	defer c.Close()

	server := r.FormValue("server")
	if server == "" {
		server = "chat.freenode.net:6667"
	}

	if bath.Exists(u.ID) {
		/*
			log.Println("Reviving zombie")
			user, err = bath.Revive(u.ID, *c)
		*/
		fmt.Println("failing because not implemented")
		return
	} else {
		log.Println("Creating zombie")
		/*user, err = bath.New(u.ID,
			server,
			u.Username,
			c,
		)*/

		add := zombies.Add{
			ID:     u.ID,
			Nick:   u.Username,
			Server: server,
		}

		resp, err := pool.Tell("add", add)
		if err != nil {
			panic(err)
		}

		log.Println(resp)
	}

	if err != nil {
		log.Println("couldnt create connection", err)
		return
	}

	for {
		var a action
		err := c.ReadJSON(&a)
		if err != nil {
			log.Println("Closing connection. Error reading:", err)
			return
		}

		if a.Name == "SET" {
			log.Println("Adding user to chanel", a.Message)

			_, err := pool.Tell("join", zombies.Join{
				ID:      u.ID,
				Channel: a.Message,
			})

			if err != nil {
				log.Println("Closing connection. Error joining chanel:", err)
				return
			}
		} else if a.Name == "SEND" {
			log.Printf("Sending message '%v' to channel '%v'", a.Message, a.Channel)
			_, err := pool.Tell("send", zombies.Send{
				ID:      u.ID,
				Channel: a.Channel,
				Message: a.Message,
			})

			if err != nil {
				log.Println("Closing connection. Error sending message:", err)
				return
			}
		}
	}
	/*
		for {
			var a action
			err := user.WSConn.ReadJSON(&a)
			if err != nil {
				log.Println("error reading:", err)
				return
			}

			if a.Name == "SET" {
				log.Println("User joined channel", a.Message)

				// Prevent duplicate users
				add := true
				for _, client := range clients[a.Message] {
					if client == user {
						add = false
					}
				}

				if add {
					clients[a.Message] = append(clients[a.Message], user)
				}
			} else if a.Name == "SEND" {
				log.Printf("Sending message '%s' to channel '%s'", a.Message, a.Channel)
				user.Messages <- a.Message
			} else if a.Name == "NICK" {
				log.Printf("Changing nick to '%v'", a.Message)
				user.SetNick(a.Message)
			}
		}*/
}

type action struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Channel string `json:"channel"`
}
