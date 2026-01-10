package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/utils"

	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/gorilla/websocket"
)

type Config struct {
	Tunnel map[string]ConfigTunnel `json:"tunnel"`
}

type ConfigTunnel struct {
	Name   string `json:"name"`
	Port   int    `json:"port"`
	Proto  string `json:"proto"`
	Host   string `json:"host"`
	Secure bool   `json:"secure"`
}

type Tunnel struct {
	Conns   map[string]*websocket.Conn
	ConnIds map[string]string
	Config  *Config
}

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type TxRequestResponse struct {
	Type string          `json:"type"`
	Data *TxResponseData `json:"data"`
}
type TxResponseData struct {
	ID      string          `json:"id"`
	Body    []byte          `json:"body"`
	Status  int             `json:"status"`
	Headers json.RawMessage `json:"headers"`
}

type TxCloseData struct {
	ID string `json:"id"`
}

type TxClose struct {
	Type string      `json:"type"`
	Data TxCloseData `json:"data"`
}

type RxHandleCall struct {
	ID      string            `json:"id"`
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"header"`
	Host    string            `json:"host"`
	Body    json.RawMessage   `json:"body"`
}

func (t *Tunnel) Connect(ct *ConfigTunnel) bool {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s://%s/ws", utils.Ternary(ct.Secure, "wss", "ws"), ct.Host), nil)
	if err != nil {
		log.Fatal("Error dialing web socket: ", err.Error())
		return false
	}

	t.Conns[ct.Host] = conn
	log.Println("WebSocket Client Connected")

	conn = nil

	go t.Listen(ct)
	return true
}

func (t *Tunnel) Listen(ct *ConfigTunnel) {
	for {
		msg := &Message{}
		err := t.Conns[ct.Host].ReadJSON(msg)

		if err != nil {
			log.Fatal("Error reading message: ", err.Error())
		}

		switch msg.Type {
		case "request:handle":
			data := &RxHandleCall{}
			err := json.Unmarshal([]byte(msg.Data), data)
			if err != nil {
				log.Fatal("Error unmarshalling handle call data: ", err.Error())
			}

			go t.handleCalls(data, ct)

			log.Printf("Handle call ID received: %s", data.ID)
		default:
			log.Printf("Unknown message type: %s", msg.Type)
		}
	}
}

func (t *Tunnel) handleCalls(payload *RxHandleCall, ct *ConfigTunnel) {

	url := fmt.Sprintf("http://localhost:%s%s", strconv.Itoa(ct.Port), payload.Path)

	req, err := http.NewRequest(strings.ToUpper(payload.Method), url, bytes.NewReader([]byte(payload.Body)))
	for key, value := range payload.Headers {
		req.Header.Set(key, value)
	}
	if err != nil {
		log.Fatal("Error creating request: ", err.Error())
	}

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		fmt.Println("Error making request: " + err.Error())
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error reading response body: ", err.Error())
	}

	headers, err := json.Marshal(res.Header)
	if err != nil {
		log.Fatal("Error marshalling headers: ", err.Error())
	}

	parsedPayload, err := json.Marshal(&TxRequestResponse{Type: "request:response", Data: &TxResponseData{ID: payload.ID, Body: body, Status: res.StatusCode, Headers: headers}})

	if err != nil {
		log.Fatal("Error marshalling response: ", err.Error())
	}

	t.Conns[ct.Host].WriteMessage(websocket.TextMessage, parsedPayload)
}

func main() {
	tunnel := &Tunnel{}
	fmt.Printf("asda")
	InitializeConfig(tunnel)
	connectedConn := make(map[string]bool)
	tunnel.Conns = make(map[string]*websocket.Conn, 4)
	for k, v := range tunnel.Config.Tunnel {
		connectedConn[k] = tunnel.Connect(&v)
		defer tunnel.Conns[k].Close()
	}

	for k, v := range connectedConn {

		if v {
			fmt.Printf("Forwarding http://localhost:%s -> %s\n", strconv.Itoa(tunnel.Config.Tunnel[k].Port), tunnel.Config.Tunnel[k].Host)
		}
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	fmt.Println("\nClosing connection")
}
