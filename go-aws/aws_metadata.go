package main

type Service struct {
	Metadata   *Metadata            `json:"metadata"`
	Version    string               `json:"apiVersion"`
	Operations map[string]Operation `json:"operations"`
	Shapes     map[string]Shape     `json:"shapes"`
}

type Metadata struct {
	ApiVersion       string `json:"apiVersion"`
	EndpointPrefix   string `json:"endpointPrefix"`
	Protocol         string `json:"protocol"`
	ServiceFullName  string `json:"serviceFullName"`
	SignatureVersion string `json:"signatureVersion"`
	XmlNamespace     string `json:"xmlNamespace"`
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
