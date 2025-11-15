package pullrequest_usecase

import (
	"PRService/internal/domain/pullrequest"
	"PRService/internal/domain/user"
	"context"
)

type Service interface {
    Create(ctx context.Context, cmd CreatePRCommand) (*pullrequest.PullRequest, error)
    GetByID(ctx context.Context, id pullrequest.ID) (*pullrequest.PullRequest, error)
    UpdateStatus(ctx context.Context, cmd UpdateStatusCommand) error
    AssignReviewers(ctx context.Context, cmd AssignReviewersCommand) error
    ReassignReviewers(ctx context.Context, cmd ReassignReviewerCommand) error
    ListByUserID(ctx context.Context, id user.ID) ([]*pullrequest.PullRequest, error)
}

type pullrequestService struct {
	prRepo pullrequest.Repo
}

func (svc *pullrequestService) Create(ctx context.Context, pr *pullrequest.PullRequest) error {
	return svc.prRepo.Save(ctx, pr)
}

func (svc *pullrequestService) GetByID(ctx context.Context, prID pullrequest.ID) (*pullrequest.PullRequest, error) {
	return svc.prRepo.GetByID(ctx, prID)
}

func (svc *pullrequestService) UpdateStatus(ctx context.Context,  cmd UpdateStatusCommand) error {
	return svc.prRepo.UpdateStatus(ctx, cmd.PullRequestID, cmd.Status)
}

func (svc *pullrequestService) AssignReviewers(ctx context.Context, cmd AssignReviewersCommand) error {
	return svc.prRepo.AssignReviewers(ctx, cmd.PullRequestID, cmd.Reviewers)
}

func (svc *pullrequestService) ReassignReviewers(ctx context.Context, cmd ReassignReviewerCommand) error {
	return svc.prRepo.ReassignReviewers(ctx, cmd.PullRequestID, cmd.OldReviewerID, cmd.NewReviewerID)
}

func (svc *pullrequestService) ListByUserID(ctx context.Context, userID user.ID) ([]*pullrequest.PullRequest, error) {
	return svc.prRepo.ListByUserID(ctx, userID)
}

