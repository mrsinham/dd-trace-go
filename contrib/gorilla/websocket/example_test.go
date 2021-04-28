// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2021 Datadog, Inc.

package websocket_test

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	muxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
	websocketTrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/websocket"
)

// This example illustrates the usage of WrapConn to trace the websockets
// read and write operations.
func ExampleWrapConn() {
	mux := muxtrace.NewRouter()

	var upgrader websocket.Upgrader
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		oconn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer oconn.Close()
		conn := websocketTrace.WrapConn(r.Context(), oconn)
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("Received message: %s\n", message)

			err = conn.WriteMessage(mt, message)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	})

	log.Fatal(http.ListenAndServe(":8080", mux))
}
