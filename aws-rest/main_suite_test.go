// Copyright (c) 2015 RightScale, Inc. - see LICENSE

package main

import (
	"log"
	"net"
	"net/http"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/bind"
)

const testMetaDir = "/home/src/aws-sdk-core-ruby/aws-sdk-core/apis/"

func TestAwsRest(t *testing.T) {

	// load the metadata for all the services we know
	err := services.Load(testMetaDir, serviceFiles)
	if err != nil {
		log.Fatal(err)
	}

	serviceStats()
	defineHandlers()

	format.UseStringerRepresentation = true
	RegisterFailHandler(Fail)
	RunSpecs(t, "AWS Rest Suite")
}

var SRV = "http://localhost:8765"
var listener net.Listener

var _ = BeforeSuite(func() {
	goji.DefaultMux.Compile()
	//http.Handle("/", DefaultMux)

	listener = bind.Socket(":8765")
	go http.Serve(listener, goji.DefaultMux)
	time.Sleep(10 * time.Millisecond)
	log.Printf("Serving on port 8765")
})

var _ = AfterSuite(func() {
	listener.Close()
})
