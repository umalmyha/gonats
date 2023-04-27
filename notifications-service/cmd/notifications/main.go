package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/velmie/broker/natsjs"

	"github.com/umalmyha/gonats/notifications-service/internal/broker/sub"
	"github.com/umalmyha/gonats/notifications-service/internal/config"
)

func main() {
	cfg, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}

	subscriber, err := subscriber(cfg.NATS)
	if err != nil {
		log.Fatal(err)
	}
	defer subscriber.Close()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt)

	ticketsSub := sub.NewTicketEventsSubscriber(subscriber)

	ctx, cancel := context.WithCancel(context.Background())
	go ticketsSub.Listen(ctx, "tickets.tickets.created")

	<-stopCh
	cancel()
}

func subscriber(cfg config.NATS) (*natsjs.Subscriber, error) {
	js, err := natsjs.NewSubscriber(
		cfg.StreamName,
		"notifications",
		natsjs.DefaultConnectionFactory(),
		natsjs.DefaultJetStreamFactory(),
		natsjs.DefaultSubscriptionFactory(),
		natsjs.DefaultConsumerFactory(nil),
		natsjs.SubscriberConnURL(cfg.URL),
	)
	if err != nil {
		return nil, err
	}
	return js, nil
}
