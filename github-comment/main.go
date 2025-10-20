package main

import (
	"context"
	"strings"
	"time"

	"github.com/google/go-github/v73/github"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load environment variables
	env, err := loadEnv()
	if err != nil {
		logrus.Fatalf("Error loading environment variables: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if env.Plugin.GithubToken == "" && env.Plugin.GithubTokenPath == "" {
		logrus.Fatal("Either settings.github_token or settings.github_token_path must be set")
	}

	if env.Plugin.Comment == "" && env.Plugin.CommentPath == "" {
		logrus.Fatal("Either settings.comment or settings.comment_path must be set")
	}

	if env.Plugin.UpdateInPlace && env.Plugin.MessageID == "" {
		logrus.Fatal("settings.message_id must be set when update_in_place is true")
	}

	githubToken := coalesce(env.Plugin.GithubToken, env.Plugin.GithubTokenPath)
	comment := coalesce(env.Plugin.Comment, env.Plugin.CommentPath)
	repo := coalesce(env.Plugin.Repo, env.CI.Repo)
	pr := coalesce(env.Plugin.PullRequest, env.CI.Commit.PullRequest)

	owner, repoName, found := strings.Cut(repo, "/")
	if !found {
		logrus.Fatalf("Invalid repo format: %s, expected owner/repo", repo)
	}

	githubClient := github.NewClient(nil).WithAuthToken(githubToken)
	logrus.Infof("GitHub client initialized for: %s", env.CI.Repo)

	// Send or update comment
	if env.Plugin.UpdateInPlace {
		logrus.Infof("Updating comment in place for PR #%d in repo %s", pr, repo)
		err = updateComment(ctx, githubClient, owner, repoName, pr, comment, env.Plugin.MessageID)
		if err != nil {
			logrus.Fatalf("Error updating comment: %v", err)
		}
	} else {
		logrus.Infof("Sending new comment for PR #%d in repo %s", pr, repo)
		err = sendComment(ctx, githubClient, owner, repoName, pr, comment, env.Plugin.MessageID)
		if err != nil {
			logrus.Fatalf("Error sending comment: %v", err)
		}
	}

	logrus.Infof("Comment operation completed successfully")
}

func sendComment(ctx context.Context, client *github.Client, owner string, repo string, pr int, comment string, messageID string) error {
	issueComment := &github.IssueComment{Body: makeBody(comment, messageID)}
	_, _, err := client.Issues.CreateComment(ctx, owner, repo, pr, issueComment)
	return err
}

func updateComment(ctx context.Context, client *github.Client, owner string, repo string, pr int, comment string, messageID string) error {
	// Find existing comment by messageID and update it
	comments, _, err := client.Issues.ListComments(ctx, owner, repo, pr, nil)
	if err != nil {
		return err
	}

	for _, c := range comments {
		body := c.GetBody()
		if body != "" && messageID != "" && strings.Contains(body, messageID) {
			// Update the comment
			updatedComment := &github.IssueComment{Body: makeBody(comment, messageID)}
			_, _, err := client.Issues.EditComment(ctx, owner, repo, c.GetID(), updatedComment)
			return err
		}
	}

	return sendComment(ctx, client, owner, repo, pr, comment, messageID)
}

func makeBody(comment string, messageID string) *string {
	var body = comment
	body += "\n\n"
	body += messageFooter(messageID)
	return &body
}

func messageFooter(id string) string {
	messageID := "MessageID: #" + id
	if id == "" {
		messageID = ""
	}

	return `<details><summary>Details</summary>
<p>
<sub>
` + messageID + `  <br>
Written at: ` + time.Now().Format(time.RFC3339) + `  <br>
Written by [woodpecker-plugins/github-comment](https://github.com/yyewolf/woodpecker-plugins/github-comment)  <br>
</p>
</details>`
}
