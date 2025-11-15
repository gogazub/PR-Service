package pullrequest_usecase

import (
	"PRService/internal/domain/pullrequest"
	"PRService/internal/domain/user"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, cmd CreatePRCommand) (*pullrequest.PullRequest, error)
	GetByID(ctx context.Context, id pullrequest.ID) (*pullrequest.PullRequest, error)
	UpdateStatus(ctx context.Context, cmd UpdateStatusCommand) error
	AssignReviewers(ctx context.Context, cmd AssignReviewersCommand) error
	ReassignReviewers(ctx context.Context, cmd ReassignReviewerCommand) error
	ListByUserID(ctx context.Context, id user.ID) ([]*pullrequest.PullRequest, error)
}

type service struct {
	prRepo pullrequest.Repo
}

func New(repo pullrequest.Repo) Service {
	return &service{prRepo: repo}
}

func (svc *service) Create(ctx context.Context, cmd CreatePRCommand) (*pullrequest.PullRequest, error) {

	id := uuid.New().String()
	pr := pullrequest.NewPullRequest(
		id,
		cmd.Name,
		cmd.Author,
		pullrequest.OPEN,
		cmd.Reviewers,
	)

	if err := svc.prRepo.Save(ctx, pr); err != nil {
		return nil, fmt.Errorf("create pull request: %w", err)
	}

	return pr, nil
}

func (svc *service) GetByID(ctx context.Context, prID pullrequest.ID) (*pullrequest.PullRequest, error) {

	pr, err := svc.prRepo.GetByID(ctx, prID)
	if err != nil {
		return nil, fmt.Errorf("get pull request by id %s: %w", prID, err)
	}

	return pr, nil
}

func (svc *service) UpdateStatus(ctx context.Context, cmd UpdateStatusCommand) error {

	if err := svc.prRepo.UpdateStatus(ctx, cmd.PullRequestID, cmd.Status); err != nil {
		return fmt.Errorf(
			"update pull request status (id: %s, status: %s): %w",
			cmd.PullRequestID,
			cmd.Status,
			err,
		)
	}

	return nil
}

func (svc *service) AssignReviewers(ctx context.Context, cmd AssignReviewersCommand) error {

	if err := svc.prRepo.AssignReviewers(ctx, cmd.PullRequestID, cmd.Reviewers); err != nil {
		return fmt.Errorf(
			"assign reviewers to pull request %s: reviewers=%v: %w",
			cmd.PullRequestID,
			cmd.Reviewers,
			err,
		)
	}

	return nil
}

func (svc *service) ReassignReviewers(ctx context.Context, cmd ReassignReviewerCommand) error {

	if err := svc.prRepo.ReassignReviewers(ctx, cmd.PullRequestID, cmd.OldReviewerID, cmd.NewReviewerID); err != nil {
		return fmt.Errorf(
			"reassign reviewer for pull request %s: from=%s to=%s: %w",
			cmd.PullRequestID,
			cmd.OldReviewerID,
			cmd.NewReviewerID,
			err,
		)
	}

	return nil
}

func (svc *service) ListByUserID(ctx context.Context, userID user.ID) ([]*pullrequest.PullRequest, error) {

	prs, err := svc.prRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list pull requests by user id %s: %w", userID, err)
	}

	return prs, nil
}
