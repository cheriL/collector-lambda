package athena

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
)

type Config struct {
	DataBase    string
	OutLocation string
}

type Client struct {
	Config  Config
	session *session.Session
}

func NewClient(region, id, secret string, conf Config) (*Client, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(id, secret, ""),
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		Config: conf,
		session: sess,
	}, nil
}

func (c *Client) Execute(stmt string) (*athena.ResultSet, error) {
	_athena := athena.New(c.session, aws.NewConfig().WithLogLevel(aws.LogDebugWithHTTPBody))
	input := &athena.StartQueryExecutionInput{
		QueryExecutionContext: &athena.QueryExecutionContext{
			Database: aws.String(c.Config.DataBase),
		},
		QueryString: aws.String(stmt),
		ResultConfiguration: &athena.ResultConfiguration{
			OutputLocation: aws.String(c.Config.OutLocation),
		},
	}
	out, err := _athena.StartQueryExecution(input)
	if err != nil {
		return nil, err
	}
	queryID := out.QueryExecutionId
	result, err := _athena.GetQueryResults(&athena.GetQueryResultsInput{
		QueryExecutionId: queryID,
	})
	if err != nil {
		return nil, err
	}
	return result.ResultSet, nil
}
