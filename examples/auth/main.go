package main

import (
	"context"
	"log"
	"os"
	"time"

	"flag"

	"github.com/10gen/mongo-go-driver/bson"
	"github.com/10gen/mongo-go-driver/cluster"
	"github.com/10gen/mongo-go-driver/connstring"
	"github.com/10gen/mongo-go-driver/ops"
	"github.com/10gen/mongo-go-driver/readpref"
)

var col = flag.String("c", "test", "the collection name to use")

func main() {

	flag.Parse()

	mongodbURI := os.Getenv("MONGODB_URI")
	if mongodbURI == "" {
		log.Fatalf("MONGODB_URI was not set")
	}

	cs, err := connstring.Parse(mongodbURI)
	if err != nil {
		log.Fatal(err)
	}

	c, err := cluster.New(
		cluster.WithConnString(cs),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	selectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	s, err := c.SelectServer(selectCtx, cluster.WriteSelector())
	if err != nil {
		log.Fatalf("%v: %v", err, c.Model().Servers[0].LastError)
	}

	dbname := cs.Database
	if dbname == "" {
		dbname = "test"
	}

	var result bson.D
	err = ops.Run(
		ctx,
		&ops.SelectedServer{
			Server:   s,
			ReadPref: readpref.Primary(),
		},
		dbname,
		bson.D{{"count", *col}},
		&result)
	if err != nil {
		log.Fatalf("failed executing count command on %s.%s: %v", dbname, *col, err)
	}

	log.Println(result)
}
