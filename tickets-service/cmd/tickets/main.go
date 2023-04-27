package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/umalmyha/gonats/tickets-service/internal/broker/pub"
	"github.com/velmie/broker/natsjs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/twitchtv/twirp"

	"github.com/umalmyha/gonats/tickets-service/internal/config"
	"github.com/umalmyha/gonats/tickets-service/internal/store"
	"github.com/umalmyha/gonats/tickets-service/internal/ticketserver"
	pb "github.com/umalmyha/gonats/tickets-service/rpc/ticket"
)

func main() {
	cfg, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	publisher, err := publisher(cfg.NATS)
	if err != nil {
		log.Fatal(err)
	}
	defer publisher.Close()

	ticketPublisher := pub.NewTicketEventPublisher(publisher)

	ticketStore := store.NewTicketStore(db)

	ticketServer := ticketserver.NewServer(ticketStore, ticketPublisher)

	ticketHandler := pb.NewTicketServiceServer(
		ticketServer,
		twirp.WithServerPathPrefix("/tickets"),
	)

	mux := http.NewServeMux()
	mux.Handle(ticketHandler.PathPrefix(), ticketHandler)

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}

	errCh := make(chan error, 1)
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt)

	go func() {
		log.Println("starting twirp server at port", cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

	select {
	case err = <-errCh:
		fmt.Println("unknown twirp server error:", err)
	case <-stopCh:
		fmt.Println("interrupt signal has been sent, stopping the server")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			fmt.Println("failed to stop server gracefully:", err)
		}
	}
}

func database(url string) (*sql.DB, error) {
	db, err := sql.Open("mysql", url)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func publisher(cfg config.NATS) (*natsjs.Publisher, error) {
	pub, err := natsjs.NewPublisher(
		cfg.StreamName,
		"tickets",
		natsjs.DefaultConnectionFactory(),
		natsjs.DefaultJetStreamFactory(),
		natsjs.PublisherConnURL(cfg.URL),
	)
	if err != nil {
		return nil, err
	}
	return pub, nil
}
