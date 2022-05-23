package main

import (
	"collector/pkg/athena"
	"collector/pkg/credentials"
	"collector/pkg/models"
	"collector/pkg/util"
	"context"
	"fmt"
	"github.com/NYTimes/gizmo/config"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/go-github/v44/github"
	"log"
	"strings"
	"time"
)

type GithubConf struct {
	Owner string
	Repo  string
	ReqPerPage int
}

type Config struct {
	Region     string
	Credential credentials.Credentials
	Athena     athena.Config
	Github     GithubConf
}

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
	return fmt.Sprintf("Hello %s!", name.Name ), nil
}

func main() {
	var cfg *Config
	config.LoadJSONFile("./config/config.json", &cfg)

	_, err := athena.NewClient(cfg.Region, cfg.Credential.AccessKey, cfg.Credential.SecretKey, cfg.Athena)
	if err != nil {

	}
	//results, err := athenaClient.Execute("select max(CreatedAt) from records;")
	latestDate, _ := time.Parse(util.TimeFormat, "2008-09-15 03:04:05.324")

	var dataList []models.Data

	//ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	ctx := context.Background()
	//tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(nil)

	iOp := &github.IssueListByRepoOptions{
		State: "all",
		Since: latestDate,
		ListOptions: github.ListOptions{
			Page: 1,
			PerPage: cfg.Github.ReqPerPage,
		},
	}
	next := true
	for next {
		issues, resp, err := client.Issues.ListByRepo(ctx, cfg.Github.Owner, cfg.Github.Repo, iOp)
		if err != nil {
			log.Printf("List Issues error: %s. Since[%v]. Page[%d]. PerPage[%d]", err.Error(), latestDate, iOp.Page, iOp.PerPage)
			// TODO
		}
		iOp.Page += 1
		if iOp.Page > resp.LastPage {
			next = false
		}

		for _, issue := range issues {
			company  := issue.GetUser().GetCompany()
			userType := models.UserTypeContributor

			if company != "" && isPingCaper(company) {
				userType = models.UserTypePingCaper
			}

			data := models.Data{
				ID:        issue.GetID(),
				Number:    issue.GetNumber(),
				Type:      models.DataTypeIssue,
				UserID:    issue.GetUser().GetID(),
				UserType:  models.UserType(userType),
				UserLogin: issue.GetUser().GetLogin(),
				CreatedAt: issue.GetCreatedAt(),
			}

			dataList = append(dataList, data)
		}
	}

	prOp := &github.PullRequestListOptions{
		State: "all",
		Sort: "created",
		ListOptions: github.ListOptions{
			Page: 1,
			PerPage: cfg.Github.ReqPerPage,
		},
	}

	next = true
	for next {
		pullRequests, resp, err := client.PullRequests.List(ctx, cfg.Github.Owner, cfg.Github.Repo, prOp)
		if err != nil {
			log.Printf("List prs error: %s. Page[%d]. PerPage[%d]", err.Error(), iOp.Page, iOp.PerPage)
			// TODO
		}
		iOp.Page += 1
		if iOp.Page > resp.LastPage {
			next = false
		}

		for _, pr := range pullRequests {
			if pr.GetCreatedAt().Before(latestDate) {
				next = false
				break
			}

			company  := pr.GetUser().GetCompany()
			userType := models.UserTypeContributor

			if company != "" && isPingCaper(company) {
				userType = models.UserTypePingCaper
			}

			data := models.Data{
				ID:        pr.GetID(),
				Number:    pr.GetNumber(),
				Type:      models.DataTypeIssue,
				UserID:    pr.GetUser().GetID(),
				UserType:  models.UserType(userType),
				UserLogin: pr.GetUser().GetLogin(),
				CreatedAt: pr.GetCreatedAt(),
			}

			dataList = append(dataList, data)
		}
	}

	lambda.Start(HandleRequest)
}

func isPingCaper(company string) bool {
	userCompany := strings.ToLower(company)
	return strings.Contains(userCompany, "pingcap")
}

