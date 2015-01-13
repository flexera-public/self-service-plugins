// Copyright (c) 2015 RightScale, Inc. - see LICENSE

package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"sync"

	"bitbucket.org/pkg/inflect"
	"github.com/rightscale/self-service-plugins/aws-rest/Godeps/_workspace/src/gopkg.in/yaml.v2"
	"github.com/stripe/aws-go/model"
)

// AWS service, first part copied from https://github.com/stripe/aws-go/blob/master/model/model.go
// second part added from Raphael's aws_analyzer
type Service struct {
	// From aws-go
	Name          string
	Short         string
	Metadata      model.Metadata
	Documentation string
	Operations    map[string]model.Operation
	Shapes        map[string]*model.Shape
	// From aws_analyzer
	Resources      map[string]*Resource
	ServiceActions map[string]*Action `yaml:"actions"`
}

type Resource struct {
	Name              string
	ShapeName         string             `yaml:"shape"`
	PrimaryId         string             `yaml:"primary_id"`
	SecondaryIds      []string           `yaml:"secondary_ids"`
	CrudActions       map[string]*Action `yaml:"actions"`
	CustomActions     map[string]*Action `yaml:"custom_actions"`
	CollectionActions map[string]*Action `yaml:"collection_actions"`
	Links             map[string]string
}

type Action struct {
	Name          string
	Verb          string
	Path          string
	PayloadShape  string   `yaml:"payload"`
	ParamShapes   []string `yaml:"params"`
	ResponseShape string   `yaml:"response"`
}

type Services map[string]Service

// Load the metadata for a service both from the AWS API docs via aws-go and from the
// Ruby aws_analyzer via yaml
func (sm Services) Load(dir string, services map[string]string) error {
	var mapMutex sync.Mutex
	var wg sync.WaitGroup
	for s, f := range services {
		wg.Add(1)
		go func(s, f string) {
			defer wg.Done()
			// Open the AWS metadata file
			in, err := os.Open(filepath.Join(dir, f))
			if err != nil {
				log.Printf("Cannot load %s from %s: %s", s, f, err.Error())
				return
			}
			defer in.Close()

			// Load the metadata from JSON
			service := Service{}
			if err := json.NewDecoder(in).Decode(&service); err != nil {
				log.Printf("Cannot load %s from %s: %s", s, f, err.Error())
				return
			}

			log.Printf("Loaded AWS metadata for %s: %d operations, %d shapes", s,
				len(service.Operations), len(service.Shapes))

			// Run the ruby aws_analyzer on this thing and load its yaml output
			yamlData, err := exec.Command("../aws_analyzer/bin/analyze",
				"--service", s, "--path", dir, "--resource-only",
				"--force").Output()
			if _, ok := err.(*exec.ExitError); ok {
				log.Printf("Analyze for %s failed: %s", s, err.Error())
			} else if err != nil {
				log.Printf("Running analyze for %s failed: %s", s, err.Error())
				return
			} else if len(yamlData) > 100 {
				err = yaml.Unmarshal(yamlData, &service)
				if err != nil {
					log.Printf("Parsing yaml for %s failed: %s", s, err.Error())
					log.Printf("Ooops: %s", yamlData[0:300])
					return
				}
				log.Printf("Loaded resource metadata for %s: %d resources, %d actions",
					s, len(service.Resources), len(service.ServiceActions))
			} else {
				log.Printf("Analyze returned too little to parse: %s", yamlData)
			}

			// Fix up the top-level Service struct a bit
			for n, shape := range service.Shapes {
				shape.Name = n
			}
			service.Short = service.Metadata.EndpointPrefix
			service.Name = service.Metadata.ServiceAbbreviation

			mapMutex.Lock()
			defer mapMutex.Unlock()
			sm[service.Short] = service
		}(s, f)
	}
	wg.Wait()
	return nil
}

func (s *Service) ResourceNames() []string {
	l := make([]string, len(s.Resources))
	i := 0
	for _, r := range s.Resources {
		l[i] = inflect.Pluralize(r.Name)
		i++
	}
	sort.Strings(l)
	return l
}

func ActionNames(am map[string]*Action) []string {
	l := make([]string, len(am))
	i := 0
	for _, a := range am {
		l[i] = a.Name
		i++
	}
	sort.Strings(l)
	return l
}

// Metadata files from the Ruby AWS SDK
var serviceFiles = map[string]string{
	//"AutoScaling":          "AutoScaling.api.json",
	"CloudFormation": "CloudFormation.api.json",
	//"CloudFront":           "CloudFront.api.json",
	//"CloudSearch":          "CloudSearch.api.json",
	//"CloudSearchDomain":    "CloudSearchDomain.api.json",
	//"CloudTrail":           "CloudTrail.api.json",
	//"CloudWatch":           "CloudWatch.api.json",
	//"CloudWatchLogs":       "CloudWatchLogs.api.json",
	//"CodeDeploy":           "CodeDeploy.api.json",
	//"CognitoIdentity":      "CognitoIdentity.api.json",
	//"CognitoSync":          "CognitoSync.api.json",
	//"ConfigService":        "ConfigService.api.json",
	//"DataPipeline":         "DataPipeline.api.json",
	//"DirectConnect":        "DirectConnect.api.json",
	//"DynamoDB":             "DynamoDB.api.json",
	"EC2": "EC2.api.json",
	//"ElastiCache":          "ElastiCache.api.json",
	//"ElasticBeanstalk":     "ElasticBeanstalk.api.json",
	//"ElasticLoadBalancing": "ElasticLoadBalancing.api.json",
	//"ElasticTranscoder":    "ElasticTranscoder.api.json",
	//"EMR":                  "EMR.api.json",
	"Glacier": "Glacier.api.json",
	"IAM":     "IAM.api.json",
	//"ImportExport":         "ImportExport.api.json",
	//"Kinesis":              "Kinesis.api.json",
	//"KMS":                  "KMS.api.json",
	//"Lambda":               "Lambda.api.json",
	//"OpsWorks":             "OpsWorks.api.json",
	//"RDS":                  "RDS.api.json",
	//"Redshift":             "Redshift.api.json",
	//"Route53":              "Route53.api.json",
	//"Route53Domains":       "Route53Domains.api.json",
	"S3": "S3.api.json",
	//"SES":                  "SES.api.json",
	//"SimpleDB":             "SimpleDB.api.json",
	"SNS": "SNS.api.json",
	"SQS": "SQS.api.json",
	//"StorageGateway":       "StorageGateway.api.json",
	//"STS":                  "STS.api.json",
	//"Support":              "Support.api.json",
	//"SWF":                  "SWF.api.json",
}

/*
// Metadata files from the aws-go repository, which come from Boto
var serviceFiles = map[string]string{
	"AutoScaling":       "autoscaling/2011-01-01.normal.json",
	"CloudFormation":    "cloudformation/2010-05-15.normal.json",
	"CloudFront":        "cloudfront/2014-10-21.normal.json",
	"CloudTrail":        "cloudtrail/2013-11-01.normal.json",
	"CloudSearch":       "cloudsearch/2013-01-01.normal.json",
	"CloudSearchDomain": "cloudsearchdomain/2013-01-01.normal.json",
	"CloudWatch":        "cloudwatch/2010-08-01.normal.json",
	"CognitoIdentity":   "cognito-identity/2014-06-30.normal.json",
	"CognitoSync":       "cognito-sync/2014-06-30.normal.json",
	"CodeDeploy":        "codedeploy/2014-10-06.normal.json",
	"Config":            "config/2014-11-12.normal.json",
	"DataPipeline":      "datapipeline/2012-10-29.normal.json",
	"DirectConnect":     "directconnect/2012-10-25.normal.json",
	"DynamoDB":          "dynamodb/2012-08-10.normal.json",
	"EC2":               "ec2/2014-10-01.normal.json",
	"ElasticCache":      "elasticache/2014-09-30.normal.json",
	"ElasticBeanstalk":  "elasticbeanstalk/2010-12-01.normal.json",
	"ElasticTranscoder": "elastictranscoder/2012-09-25.normal.json",
	"ELB":               "elb/2012-06-01.normal.json",
	"EMR":               "emr/2009-03-31.normal.json",
	"IAM":               "iam/2010-05-08.normal.json",
	"ImportExport":      "importexport/2010-06-01.normal.json",
	"Kinesis":           "kinesis/2013-12-02.normal.json",
	"KMS":               "kms/2014-11-01.normal.json",
	"Lambda":            "lambda/2014-11-11.normal.json",
	"Logs":              "logs/2014-03-28.normal.json",
	"OpsWorks":          "opsworks/2013-02-18.normal.json",
	"RDS":               "rds/2014-09-01.normal.json",
	"RedShift":          "redshift/2012-12-01.normal.json",
	"Route53":           "route53/2013-04-01.normal.json",
	"Route53Domains":    "route53domains/2014-05-15.normal.json",
	"S3":                "s3/2006-03-01.normal.json",
	"SDB":               "sdb/2009-04-15.normal.json",
	"SES":               "ses/2010-12-01.normal.json",
	"SNS":               "sns/2010-03-31.normal.json",
	"SQS":               "sqs/2012-11-05.normal.json",
	"StorageGateway":    "storagegateway/2013-06-30.normal.json",
	"STS":               "sts/2011-06-15.normal.json",
	"Support":           "support/2013-04-15.normal.json",
	"SWF":               "swf/2012-01-25.normal.json",
}
*/
