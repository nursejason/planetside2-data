package main

import (
	"encoding/json"
	"flag"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"

	log "github.com/Sirupsen/logrus"
)

const wssHost = "push.planetside2.com"
const wssPath = "/streaming"
const wssRawQuery = "environment=ps2&service-id=s:example"
const wssParams = `{"service":"event","action":"subscribe","worlds":["all"],"eventNames":["PlayerLogin","PlayerLogout", "MetagameEvent"]}`

var (
	logLevel = flag.String("logLevel", "info",
		"Options include 'panic', 'fatal', 'error', 'warn', 'warning', 'info', 'debug'")
	elasticURL  = flag.String("elasticURL", "elastic", "URL for ElasticSearch")
	elasticPort = flag.String("elasticPort", "9300", "Port for ElasticSearch")
)

// ParentEvent contains a single PlanetsideEvent
type ParentEvent struct {
	Payload PlanetsideEvent `json:"payload"`
	Service string          `json:"service"`
	Type    string          `json:"type"`
}

// PlanetsideEvent is the struct representation of Planetside JSON Events
type PlanetsideEvent struct {
	// PlayerEvents
	CharacterID string `json:"character_id"`
	EventName   string `json:"event_name"`
	Timestamp   string `json:"timestamp"`
	WorldID     string `json:"world_id"`

	// MetagameEvents
	ExperienceBonus    string `json:"experience_bonus"`
	FactionNC          string `json:"faction_nc"`
	FactionTR          string `json:"faction_tr"`
	FactionVS          string `json:"faction_vs"`
	MetagameEventID    string `json:"metagame_event_id"`
	MetagameEventState string `json:"metagame_event_state"`
	ZoneID             string `json:"zone_id"`
}

func processMessage(message []byte) error {
	var event ParentEvent
	if err := json.Unmarshal(message, &event); err != nil {
		log.WithFields(log.Fields{
			"err":     err,
			"message": string(message),
		}).Error("Error unmarshaling message to PlanetsideEvent JSON")
		return err
	}

	// Sanity check requests come in every so often
	if event.Payload == (PlanetsideEvent{}) {
		log.Info("Non PlanetsideEvent received, skipping message.")
		return nil
	}

	log.WithFields(log.Fields{
		"event": event,
	}).Debug("Successfully read message.")

	// If CharacterID -- Get Character info, then save to Elastic
	// Elif MetagameEventID -- Save to Elastic
	return nil
}

func main() {
	flag.Parse()

	if level, err := log.ParseLevel(*logLevel); err != nil {
		log.Fatal("Error parsing log level. Please validate entered value.")
	} else {
		log.SetLevel(level)
	}

	log.WithFields(log.Fields{
		"logLevel": *logLevel,
	}).Info("Beginning execution of application.")

	// TODO What's this doing?
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: wssHost, Path: wssPath, RawQuery: wssRawQuery}
	log.WithFields(log.Fields{
		"url": u.String(),
	}).Info("Connecting to client.")

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Error connecting to dialer")
	}

	defer func() {
		if err := c.Close(); err != nil {
			log.Error("Error closing websocket.")
		}
	}()

	done := make(chan struct{})
	go func() {
		defer func() {
			if err := c.Close(); err != nil {
				log.Error("Error closing websocket.")
			}
		}()
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("Error reading message")
				return
			}

			if err = processMessage(message); err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("Error processing message")
				return
			}
		}
	}()

	err = c.WriteMessage(websocket.TextMessage, []byte(wssParams))
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Error writing message")
		return
	}

	for {
		select {
		case <-interrupt:
			log.Info("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("Error reading message")
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			if err := c.Close(); err != nil {
				log.Error("Error closing websocket.")
			}
			return
		}
	}
}
