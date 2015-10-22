// Copyright (c) 2015 RightScale, Inc. - see LICENSE

package main

//"github.com/stripe/aws-go/aws"
//"github.com/stripe/aws-go/gen/endpoints"
import (
	"fmt"
	"log"
	"net/http"

	"bitbucket.org/pkg/inflect"
	"github.com/zenazn/goji/web"
)

// httpError writes an error to the HTTP response and returns an error object with that same
// error for logging
func httpError(w http.ResponseWriter, httpCode int, format string, args ...interface{}) error {
	w.WriteHeader(httpCode)
	fmt.Fprintf(w, format, args...)
	return fmt.Errorf(format, args...)
}

func serviceMeta(svcName, region string, w http.ResponseWriter) (*Service, error) {
	svcMeta, ok := services[svcName]
	if !ok {
		return nil, httpError(w, 404,
			"Sorry, service %s does not exist or is not supported", svcName)
	}
	return &svcMeta, nil
}

func serviceResource(svc *Service, resource string, w http.ResponseWriter) (*Resource, error) {
	res, ok := svc.Resources[inflect.Singularize(resource)]
	if !ok {
		return nil, httpError(w, 404,
			"Sorry, %s does not have resource %s, available resources: %+v",
			svc.Name, resource, svc.ResourceNames())
	}
	return res, nil
}

func findAction(actionMap map[string]*Action, action, actionType, container string,
	w http.ResponseWriter) (*Action, error) {
	act, ok := actionMap[(action)]
	if !ok {
		return nil, httpError(w, 404,
			"Sorry, %s does not have %s action %s, available aactions: %+v",
			container, actionType, action, ActionNames(actionMap))
	}
	return act, nil
}

// GET /:service/:region/:resource
func indexHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	log.Printf("Index params: %+v", c.URLParams)
	svc, err := serviceMeta(c.URLParams["service"], c.URLParams["region"], w)
	if err != nil {
		log.Print(err.Error())
		return
	}

	res, err := serviceResource(svc, c.URLParams["resource"], w)
	if err != nil {
		log.Print(err.Error())
		return
	}

	act, err := findAction(res.CrudActions, "index", "collection", res.Name, w)
	if err != nil {
		log.Print(err.Error())
		return
	}

	fmt.Fprintf(w, "-> %s %s%s", act.Verb, act.Path, act.Name)
}

func showHandler(c web.C, w http.ResponseWriter, r *http.Request) {
}

func updateHandler(c web.C, w http.ResponseWriter, r *http.Request) {
}

func createHandler(c web.C, w http.ResponseWriter, r *http.Request) {
}

func deleteHandler(c web.C, w http.ResponseWriter, r *http.Request) {
}

func serviceActionHandler(c web.C, w http.ResponseWriter, r *http.Request) {
}

func collectionActionHandler(c web.C, w http.ResponseWriter, r *http.Request) {
}

func resourceActionHandler(c web.C, w http.ResponseWriter, r *http.Request) {
}
