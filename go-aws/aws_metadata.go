package main

import (
	"bytes"
	"fmt"
)

type Service struct {
	Metadata   *Metadata            `json:"metadata"`
	Version    string               `json:"apiVersion"`
	Operations map[string]Operation `json:"operations"`
	Shapes     map[string]Shape     `json:"shapes"`
}

func (s *Service) String(name string) string {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "===== Service %s =====\n", name)
	fmt.Fprintf(buf, "Version: %s\n", s.Version)
	fmt.Fprint(buf, s.Metadata.String())
	fmt.Fprintf(buf, "Operations: %d\n", len(s.Operations))
	for _, op := range s.Operations {
		buf.WriteString(op.String())
	}
	fmt.Fprintf(buf, "Shapes: %d\n", len(s.Shapes))
	for n, sh := range s.Shapes {
		buf.WriteString(sh.String(n))
	}
	return buf.String()
}

type Metadata struct {
	ApiVersion       string `json:"apiVersion"`
	EndpointPrefix   string `json:"endpointPrefix"`
	Protocol         string `json:"protocol"`
	ServiceFullName  string `json:"serviceFullName"`
	SignatureVersion string `json:"signatureVersion"`
	XmlNamespace     string `json:"xmlNamespace"`
}

func (m *Metadata) String() string {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "Metadata:\n")
	fmt.Fprintf(buf, "  ApiVersion      : %s\n", m.ApiVersion)
	fmt.Fprintf(buf, "  EndpointPrefix  : %s\n", m.EndpointPrefix)
	fmt.Fprintf(buf, "  Protocol        : %s\n", m.Protocol)
	fmt.Fprintf(buf, "  ServiceFullName : %s\n", m.ServiceFullName)
	fmt.Fprintf(buf, "  SignatureVersion: %s\n", m.SignatureVersion)
	fmt.Fprintf(buf, "  XmlNamespace    : %s\n", m.XmlNamespace)
	return buf.String()
}

type Operation struct {
	Name string `json:"name"`
	Http struct {
		Method     string `json:"method"`
		RequestUri string `json:"requestUri"`
	} `json:"http"`
	Input struct {
		Shape string `json:"shape"`
	} `json:"input"`
	Output struct {
		ResultWrapper string `json:"resultWrapper"`
		Shape         string `json:"shape"`
	} `json:"output"`
	Errors []struct {
		Error struct {
			Code           string  `json:"code"`
			HttpStatusCode float64 `json:"httpStatusCode"`
			SenderFault    bool    `json:"senderFault"`
		} `json:"error"`
		Exception bool   `json:"exception"`
		Shape     string `json:"shape"`
	} `json:"errors"`
}

func (o *Operation) String() string {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "  %-20s: %s %s\n", o.Name, o.Http.Method, o.Http.RequestUri)
	fmt.Fprintf(buf, "    In: %s Out: %s Err: %d\n",
		o.Input.Shape, o.Output.Shape, len(o.Errors))
	return buf.String()
}

type Shape struct {
	Type     string   `json:"type"`
	Required []string `json:"required"`
	Members  map[string]struct {
		Shape        string `json:"shape"`
		Location     string `json:"location"`
		LocationName string `json:"locationName"`
	} `json:"members"`
	Pattern string  `json:"pattern"`
	Max     float64 `json:"max"`
	Min     float64 `json:"min"`
}

func (s *Shape) String(name string) string {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "  Type: %s Members:", s.Type)
	for n, m := range s.Members {
		fmt.Fprintf(buf, " %s (%s)", n, m.Shape)
	}
	fmt.Fprint(buf, "\n")
	return buf.String()
}
