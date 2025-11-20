package app

import (
	"PRService/internal/domain/pullrequest"
	"PRService/internal/domain/team"
	"PRService/internal/domain/user"
	pullrequest_usecase "PRService/internal/usecase/pullrequest"
	team_usecase "PRService/internal/usecase/team"
	user_usecase "PRService/internal/usecase/user"
	"context"
	"fmt"
	"math/rand"
)

// Агрегированные в одну сущность сервисы. Точка взаимодействия с приложение.
type Services struct {
	User        user_usecase.Service
	Team        team_usecase.Service
	PullRequest pullrequest_usecase.Service
}

func NewServices(user user_usecase.Service, team team_usecase.Service, pr pullrequest_usecase.Service) *Services {
	return &Services{User: user, Team: team, PullRequest: pr}
}

func (svc *Services) CreateTeam(ctx context.Context, cmd team_usecase.CreateTeamAndUsersCommand) (*team.Team, []*user.User, error) {
	// TODO: создать здесь транзакцию. Либо создаем и team и все users, либо ничего

	// 1. Make cmd
	members := cmd.Members
	ids := make([]user.ID, 0, len(members))
	for _, u := range members {
		ids = append(ids, u.UserID)
	}
	createTeamCmd := team_usecase.CreateTeamCommand{
		Name:    cmd.Name,
		Members: ids,
	}
	// 2. Save team
	t, err := svc.Team.CreateTeam(ctx, createTeamCmd)
	if err != nil {
		return nil, nil, fmt.Errorf("create team: %q: %w", createTeamCmd.Name, err)
	}

	// 3. Save users
	users := make([]*user.User, 0, len(members))
	for _, u := range members {
		_, err := svc.User.CreateUser(ctx, u)
		if err != nil {
			return nil, nil, fmt.Errorf("create team: %w", err)
		}
		users = append(users, u)
	}

	return t, users, nil
}

func (svc *Services) CreatePR(ctx context.Context, cmd pullrequest_usecase.CreatePRCommand) (*pullrequest.PullRequest, error) {
	// TODO: делать все операции под общей tx

	if pr, _ := svc.PullRequest.GetByID(ctx,pullrequest.ID(cmd.ID)); pr != nil {
		return nil, pullrequest.ErrPullRequestExists
	}
	
	// 1. get user
	u, err := svc.User.GetByID(ctx, cmd.Author)
	if err != nil {
		return nil, fmt.Errorf("service: create pr: get author: %w", err)
	}

	if u.TeamName == "" {
		return nil, team.ErrTeamNotFound
	}

	// 2. get team
	t, err := svc.Team.GetByName(ctx, u.TeamName)
	if err != nil {
		return nil, fmt.Errorf("service: create pr: get the author team: %w", err)
	}

	// 3. get active users
	activeMembers, err := svc.Team.GetActiveUsersInTeam(ctx, t.Name)
	if err != nil {
		return nil, fmt.Errorf("service: create pr: get active users in team: %q: %w", t.Name, err)
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

	pr := pullrequest.NewPullRequest(cmd.ID, cmd.Name, cmd.Author, pullrequest.OPEN, reviewersIDs)
	
	if err = svc.PullRequest.Save(ctx, pr); err != nil {
		return nil, fmt.Errorf("service: create pr: %w", err)
	}
	return pr, nil
}


type ReassignReviewerCommand struct {
	PullRequestID  pullrequest.ID
	OldReviewerID  user.ID
	NewReviewerID  user.ID // будет заполнен в сервисе
}

func (svc *Services) ReassignReviewer(ctx context.Context, cmd pullrequest_usecase.ReassignReviewerCommand) (*pullrequest.PullRequest, user.ID, error) {
	pr, err := svc.PullRequest.GetByID(ctx, cmd.PullRequestID)
	if err != nil {
		return nil, "", fmt.Errorf("service: reassign reviewer: get pr: %w", err)
	}

	if pr.Author == "" {
		return nil, "", pullrequest.ErrNoAuthor
	}

	oldIdx := -1
	for i, rid := range pr.Reviewers {
		if rid == cmd.OldReviewerID {
			oldIdx = i
			break
		}
	}
	if oldIdx == -1 {
		return nil, "", pullrequest.ErrReviewerNotAssigned
	}

	author, err := svc.User.GetByID(ctx, pr.Author)
	if err != nil {
		return nil, "", fmt.Errorf("service: reassign reviewer: get author: %w", err)
	}

	users, err := svc.Team.GetActiveUsersInTeam(ctx, author.TeamName)
	if err != nil {
		return nil, "", fmt.Errorf("service: reassign reviewer: get team users: %w", err)
	}

	candidates := make([]user.ID, 0, len(users))
	for _, u := range users {
		uid := u.UserID

		if uid == pr.Author {
			continue
		}
		if uid == cmd.OldReviewerID {
			continue
		}
		alreadyAssigned := false
		for _, rid := range pr.Reviewers {
			if rid == uid {
				alreadyAssigned = true
				break
			}
		}
		if alreadyAssigned {
			continue
		}

		candidates = append(candidates, uid)
	}

	if len(candidates) == 0 {
		return nil, "", pullrequest.ErrNoCandidate
	}

	newIdx := rand.Intn(len(candidates))
	newReviewerID := candidates[newIdx]
	cmd.NewReviewerID = newReviewerID

	if err := svc.PullRequest.ReassignReviewers(ctx, cmd); err != nil {
		return nil, "", fmt.Errorf("service: reassign reviewer: %w", err)
	}

	updated := *pr
	updated.Reviewers = make([]user.ID, len(pr.Reviewers))
	copy(updated.Reviewers, pr.Reviewers)
	updated.Reviewers[oldIdx] = newReviewerID

	return &updated, newReviewerID, nil
}
