package db_test

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/derickit/go-rest-api/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestMongoDBCredentialFromSideCar(t *testing.T) {
	type mongoCredentialTestCase struct {
		Description string
		Input       string
		ExpectedOut *db.MongoDBCredentials
		ExoectedErr error
	}
	var testCases = []mongoCredentialTestCase{
		{
			Description: "Valid Credential",
			Input:       "../mockData/mongoDB_test_credentials.json",
			ExpectedOut: &db.MongoDBCredentials{
				HostName:   "localhost",
				Password:   "password",
				User:       "ecuser",
				ReplicaSet: "",
				Port:       "27017",
			},
		},
		{
			Description: "ensure MongoDBCredentialFromSideCar returns error when invalid file path is given",
			Input:       "../mockData/non-existent.json",
			ExoectedErr: db.ErrSideCarFileRead,
		},
		{
			Description: "expect nil when invalid json file is given",
			Input:       "../mockData/mongoDB_test_credentials_fail.json",
			ExpectedOut: nil,
			ExoectedErr: db.ErrSideCarFileFormat,
		},
	}

	for i, tc := range testCases {
		t.Setenv("MongoVaultSideCar", tc.Input)
		result, err := db.MongoDBCredentialFromSideCar(os.Getenv("MongoVaultSideCar"))
		if !assert.EqualValues(t, tc.ExpectedOut, result) || !errors.Is(err, tc.ExoectedErr) {
			t.Errorf("MongoDBCredentialFromSideCar Test Case:%d:%s failed:expected %v,%v;got %v,%v",
				i+1, tc.Description, tc.ExpectedOut, tc.ExoectedErr, result, err)
		}
	}
}

func TestMongoConnectionUrl(t *testing.T) {
	type connectionURLTestCase struct {
		Description string
		Input       *db.MongoDBCredentials
		ExpectedOut string
	}

	var testCases = []connectionURLTestCase{
		{
			Description: "ensure Connection Url is empty",
			Input:       &db.MongoDBCredentials{},
			ExpectedOut: "",
		},
		{
			Description: "atlas connection url should have options retryWrites and w",
			Input: &db.MongoDBCredentials{
				HostName: "mongodb.net",
				User:     "fuzzy-derick",
				Password: "hskdkldldl",
			},
		},
		{
			Description: "connection url should have username and password in url",
			Input: func() *db.MongoDBCredentials {
				var m db.MongoDBCredentials
				data, _ := os.ReadFile("../mockData/mongoDB_test_credentials.json")
				_ = json.Unmarshal(data, &m)
				return &m
			}(),
			ExpectedOut: "mongodb://ecuser:12345667@test:8899",
		},
		{
			Description: "should have replicaset",
			Input: &db.MongoDBCredentials{
				HostName:   "mongodb1.svc.com,mongodb2.svc.com",
				ReplicaSet: "mySet",
			},
			ExpectedOut: "mongodb://mongodb1.svc.com,mongodb2.svc.com/?replicaSet=mySet",
		},
		{
			Description: "do not include authentication if password is missing",
			Input: &db.MongoDBCredentials{
				HostName:   "mongodb1.svc.com,mongodb2.svc.com",
				ReplicaSet: "mySet",
				User:       "test",
			},
			ExpectedOut: "mongodb://mongodb1.svc.com,mongodb2.svc.com/?replicaSet=mySet",
		},
		{
			Description: "do not include authentication if username is missing",
			Input: &db.MongoDBCredentials{
				HostName:   "mongodb1.svc.com,mongodb2.svc.com",
				ReplicaSet: "mySet",
				Password:   "test",
			},
			ExpectedOut: "mongodb://mongodb1.svc.com,mongodb2.svc.com/?replicaSet=mySet",
		},
		{
			Description: "include port when its given",
			Input: &db.MongoDBCredentials{
				HostName:   "mongodb1.svc.com",
				ReplicaSet: "mySet",
				Password:   "test",
				Port:       "27107",
			},
			ExpectedOut: "mongodb://mongodb1.svc.com:27107/?replicaSet=mySet",
		},
		{
			Description: "discard port when its multiple hosts",
			Input: &db.MongoDBCredentials{
				HostName:   "mongodb1.svc.com,mongodb2.svc.com",
				ReplicaSet: "mySet",
				Password:   "test",
				Port:       "23456",
			},
			ExpectedOut: "mongodb://mongodb1.svc.com,mongodb2.svc.com/?replicaSet=mySet",
		},
	}
	for i, tc := range testCases {
		result := db.MongoConnectionURL(tc.Input)
		if result != tc.ExpectedOut {
			t.Errorf("TestMongoConnectionUrl test case %d:%s failed:expected %s;got %s",
				i, tc.Description, tc.ExpectedOut, result)
		}
	}
}

func TestMaskedMongoConnectionURL(t *testing.T) {
	type connectionURLTestCase struct {
		Description string
		Input       *db.MongoDBCredentials
		ExpectedOut string
	}
	var testCases = []connectionURLTestCase{
		{
			Description: "ensure masking doesnt fail for empty credentials",
			Input:       &db.MongoDBCredentials{},
			ExpectedOut: "",
		}, {
			Description: "ensure username and passwords are masked",
			Input: &db.MongoDBCredentials{
				HostName: "mongodb.net",
				User:     "deedede",
				Password: "hosnkdkkkdsd",
			},
			ExpectedOut: "mongodb://****:****@mongodb.net",
		},
	}

	for i, tc := range testCases {
		result := db.MaskedMongoConnectionURL(tc.Input)
		if result != tc.ExpectedOut {
			t.Errorf("TestMaskedMongoConnectionURL test case %d:%s failed:expected %s;got %s",
				i, tc.Description, tc.ExpectedOut, result)
		}
	}
}
