package athena

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
	"time"
)

type Config struct {
	DataBase    string
	OutLocation string
}

type Client struct {
	Config  Config
	client  *athena.Athena
}

func NewClient(region, id, secret string, conf Config) (*Client, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(id, secret, ""),
	})
	if err != nil {
		return nil, err
	}

	c := athena.New(sess, aws.NewConfig().WithLogLevel(aws.LogDebugWithHTTPBody))

	return &Client{
		Config: conf,
		client: c,
	}, nil
}

func (c *Client) Execute(stmt string) (string, error) {
	input := &athena.StartQueryExecutionInput{
		QueryExecutionContext: &athena.QueryExecutionContext{
			Database: aws.String(c.Config.DataBase),
		},
		QueryString: aws.String(stmt),
		ResultConfiguration: &athena.ResultConfiguration{
			OutputLocation: aws.String(c.Config.OutLocation),
		},
	}
	out, err := c.client.StartQueryExecution(input)
	if err != nil {
		return "", err
	}
	queryID := aws.StringValue(out.QueryExecutionId)
	return queryID, nil
}

func (c *Client) GetResult(queryID string) (*athena.ResultSet, error) {
	var result *athena.GetQueryResultsOutput
	var err error

	query := true
	for query {
		if result, err = c.client.GetQueryResults(&athena.GetQueryResultsInput{
			QueryExecutionId: aws.String(queryID),
		}); err == nil {
			break
		} else {
			e, ok := err.(*athena.InvalidRequestException)
			if !ok {
				return nil, err
			}
			errCode := aws.StringValue(e.AthenaErrorCode)
			if errCode != "INVALID_QUERY_EXECUTION_STATE" {
				return nil, err
			}
		}

		time.Sleep(1 * time.Second)
	}

	return result.ResultSet, nil
}
