package app

import (
	"PRService/internal/adapters/repo/transactor"
	"PRService/internal/domain/pullrequest"
	"PRService/internal/domain/team"
	"PRService/internal/domain/user"
	pullrequestusecase "PRService/internal/usecase/pullrequest"
	teamusecase "PRService/internal/usecase/team"
	userusecase "PRService/internal/usecase/user"
	"context"
	"fmt"
	"math/rand"
	"slices"
)

// Агрегированные в одну сущность сервисы. Точка взаимодействия с приложение.
type Services struct {
	User        userusecase.Service
	Team        teamusecase.Service
	PullRequest pullrequestusecase.Service
	tm          *transactor.Transactor
}

func NewServices(
	user userusecase.Service,
	team teamusecase.Service,
	pr pullrequestusecase.Service,
	tm *transactor.Transactor,
) *Services {
	return &Services{User: user, Team: team, PullRequest: pr, tm: tm}
}

func (svc *Services) CreateTeam(ctx context.Context, cmd teamusecase.CreateTeamAndUsersCommand) (*team.Team, []*user.User, error) {
	var (
		t     *team.Team
		users []*user.User
	)

	err := svc.tm.WithinTransaction(ctx, func(txCtx context.Context) error {
		// 1. Save team
		ids := make([]user.ID, 0, len(cmd.Members))
		for _, u := range cmd.Members {
			ids = append(ids, u.UserID)
		}
		createTeamCmd := teamusecase.CreateTeamCommand{
			Name:    cmd.Name,
			Members: ids,
		}

		var err error
		t, err = svc.Team.CreateTeam(txCtx, createTeamCmd)
		if err != nil {
			return fmt.Errorf("create team: %q: %w", createTeamCmd.Name, err)
		}

		// 2. Save users
		users = make([]*user.User, 0, len(cmd.Members))
		for _, u := range cmd.Members {
			if _, err := svc.User.CreateUser(txCtx, u); err != nil {
				return fmt.Errorf("create user: %w", err)
			}
			users = append(users, u)
		}
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return t, users, nil
}

func (svc *Services) CreatePR(ctx context.Context, cmd pullrequestusecase.CreatePRCommand) (*pullrequest.PullRequest, error) {

	var pr *pullrequest.PullRequest

	err := svc.tm.WithinTransaction(ctx, func(txCtx context.Context) error {

		if pr, _ = svc.PullRequest.GetByID(txCtx, pullrequest.ID(cmd.ID)); pr != nil {
			return pullrequest.ErrPullRequestExists
		}

		// 1. get user
		u, err := svc.User.GetByID(txCtx, cmd.Author)
		if err != nil {
			return fmt.Errorf("service: create pr: get author: %w", err)
		}

		if u.TeamName == "" {
			return team.ErrTeamNotFound
		}

		// 2. get team
		t, err := svc.Team.GetByName(txCtx, u.TeamName)
		if err != nil {
			return fmt.Errorf("service: create pr: get the author team: %w", err)
		}

		// 3. get active users
		activeMembers, err := svc.Team.GetActiveUsersInTeam(txCtx, t.Name)
		if err != nil {
			return fmt.Errorf("service: create pr: get active users in team: %q: %w", t.Name, err)
		}

		// 4. chose at most 2 reviewers
		candidates := make([]*user.User, 0, len(activeMembers))
		for _, member := range activeMembers {
			if member.UserID != cmd.Author {
				candidates = append(candidates, member)
			}
		}

		rand.Shuffle(len(candidates), func(i, j int) {
			candidates[i], candidates[j] = candidates[j], candidates[i]
		})

		reviewers := candidates[:min(2, len(candidates))]
		reviewersIDs := make([]user.ID, 0)
		for _, u := range reviewers {
			reviewersIDs = append(reviewersIDs, u.UserID)
		}

		pr = pullrequest.NewPullRequest(cmd.ID, cmd.Name, cmd.Author, pullrequest.OPEN, reviewersIDs)

		if err = svc.PullRequest.Save(txCtx, pr); err != nil {
			return fmt.Errorf("service: create pr: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return pr, nil
}

type ReassignReviewerCommand struct {
	PullRequestID pullrequest.ID
	OldReviewerID user.ID
	NewReviewerID user.ID // будет заполнен в сервисе
}

func (svc *Services) ReassignReviewer(ctx context.Context, cmd pullrequestusecase.ReassignReviewerCommand) (*pullrequest.PullRequest, user.ID, error) {

	var pr *pullrequest.PullRequest
	var uid user.ID

	err := svc.tm.WithinTransaction(ctx, func(txCtx context.Context) error {
		var err error

		pr, err = svc.PullRequest.GetByID(txCtx, cmd.PullRequestID)
		if err != nil {
			return fmt.Errorf("service: reassign reviewer: get pr: %w", err)
		}

		if pr.Author == "" {
			return pullrequest.ErrNoAuthor
		}

		oldIdx := -1
		for i, rid := range pr.Reviewers {
			if rid == cmd.OldReviewerID {
				oldIdx = i
				break
			}
		}
		if oldIdx == -1 {
			return pullrequest.ErrReviewerNotAssigned
		}

		author, err := svc.User.GetByID(txCtx, pr.Author)
		if err != nil {
			return fmt.Errorf("service: reassign reviewer: get author: %w", err)
		}

		users, err := svc.Team.GetActiveUsersInTeam(txCtx, author.TeamName)
		if err != nil {
			return fmt.Errorf("service: reassign reviewer: get team users: %w", err)
		}

		candidates := make([]user.ID, 0, len(users))
		for _, u := range users {
			uidIter := u.UserID

			if uidIter == pr.Author {
				continue
			}
			if uidIter == cmd.OldReviewerID {
				continue
			}
			alreadyAssigned := false

			if slices.Contains(pr.Reviewers, uidIter) {
				alreadyAssigned = true
			}

			if alreadyAssigned {
				continue
			}

			candidates = append(candidates, uidIter)
		}

		if len(candidates) == 0 {
			return pullrequest.ErrNoCandidate
		}

		newIdx := rand.Intn(len(candidates))
		newReviewerID := candidates[newIdx]
		cmd.NewReviewerID = newReviewerID

		if err := svc.PullRequest.ReassignReviewers(txCtx, cmd); err != nil {
			return fmt.Errorf("service: reassign reviewer: %w", err)
		}

		updated := *pr
		updated.Reviewers = make([]user.ID, len(pr.Reviewers))
		copy(updated.Reviewers, pr.Reviewers)
		updated.Reviewers[oldIdx] = newReviewerID

		pr = &updated
		uid = newReviewerID
		return nil
	})

	if err != nil {
		return nil, "", err
	}

	return pr, uid, nil
}
