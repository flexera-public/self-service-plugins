// Copyright (c) 2015 RightScale, Inc. - see LICENSE

package main

//"github.com/stripe/aws-go/aws"
//"github.com/stripe/aws-go/gen/endpoints"
import (
	"flag"
	"log"

	"github.com/zenazn/goji"
)

var metaDir = flag.String("metadir", "aws-metadata", "directory with metadata from aws-go/apis")

var services = Services{}

func main() {
	flag.Parse()

	// load the metadata for all the services we know
	err := services.Load(*metaDir, serviceFiles)
	if err != nil {
		log.Fatal(err)
	}

	serviceStats()
	defineHandlers()
	goji.Serve()
}

// serviceStats tallies what we got and prints some impressive stats
func serviceStats() {
	ops := 0
	shapes := 0
	res := 0
	act := 0
	for _, svc := range services {
		ops += len(svc.Operations)
		shapes += len(svc.Shapes)
		res += len(svc.Resources)
		act += len(svc.ServiceActions)
		for _, r := range svc.Resources {
			act += len(r.CrudActions) + len(r.CustomActions) + len(r.CollectionActions)
		}
	}
	log.Printf(
		"Done loading %d services with %d operations, %d shapes, %d resources and %d actions",
		len(services), ops, shapes, res, act)
}

func defineHandlers() {
	pref := "/:service/:region"
	goji.Get(pref+"/:resource", indexHandler)
	goji.Get(pref+"/:resource/:id", showHandler)
	goji.Put(pref+"/:resource/:id", updateHandler)
	goji.Post(pref+"/:resource", createHandler)
	goji.Delete(pref+"/:resource/:id", deleteHandler)
	goji.Post(pref+"/action/:action", serviceActionHandler)
	goji.Post(pref+"/:resource/action/:action", collectionActionHandler)
	goji.Post(pref+"/:resource/:id/action/:action", resourceActionHandler)
}
