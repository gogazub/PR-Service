package tests

import (
	"fmt"
	"net/http"
	"testing"
)

// 1. Тест на создание команды и получение (Happy Path + Error)
func TestTeamLifecycle(t *testing.T) {
	suffix := RandomSuffix()
	teamName := "team_" + suffix
	user1 := "u1_" + suffix
	user2 := "u2_" + suffix

	// --- A. Create Team ---
	reqBody := TeamDTO{
		TeamName: teamName,
		Members: []TeamMemberDTO{
			{UserID: user1, Username: "Alice", IsActive: true},
			{UserID: user2, Username: "Bob", IsActive: true},
		},
	}

	resp, data := DoPostJSON(t, "/team/add", reqBody)
	AssertStatus(t, resp, http.StatusCreated, data)

	var addResp AddTeamResponse
	DecodeJSON(t, data, &addResp, "create team")

	if addResp.Team.TeamName != teamName {
		t.Errorf("expected team name %s, got %s", teamName, addResp.Team.TeamName)
	}

	// --- B. Try Create Duplicate Team (expect 400 TEAM_EXISTS) ---
	respErr, dataErr := DoPostJSON(t, "/team/add", reqBody)
	AssertStatus(t, respErr, http.StatusBadRequest, dataErr)
	AssertErrorCode(t, dataErr, "TEAM_EXISTS")

	// --- C. Get Team ---
	respGet, dataGet := DoGet(t, fmt.Sprintf("/team/get?team_name=%s", teamName))
	AssertStatus(t, respGet, http.StatusOK, dataGet)

	var getResp AddTeamResponse
	DecodeJSON(t, dataGet, &getResp, "get team")
	if len(getResp.Team.Members) != 2 {
		t.Errorf("expected 2 members, got %d", len(getResp.Team.Members))
	}

	// --- D. Get Non-Existent Team (expect 404 NOT_FOUND) ---
	respNotFound, dataNotFound := DoGet(t, "/team/get?team_name=unknown_team")
	AssertStatus(t, respNotFound, http.StatusNotFound, dataNotFound)
}

// 2. Тест на работу с пользователями (Active status)
func TestUserOperations(t *testing.T) {
	suffix := RandomSuffix()
	teamName := "users_team_" + suffix
	userTarget := "target_" + suffix

	// Setup: создаем команду c пользователем
	setupBody := TeamDTO{
		TeamName: teamName,
		Members: []TeamMemberDTO{
			{UserID: userTarget, Username: "TargetUser", IsActive: true},
		},
	}
	DoPostJSON(t, "/team/add", setupBody)

	// --- A. Set User Inactive ---
	activeReq := map[string]interface{}{
		"user_id":   userTarget,
		"is_active": false,
	}
	resp, data := DoPostJSON(t, "/users/setIsActive", activeReq)
	AssertStatus(t, resp, http.StatusOK, data)

	var userResp SetIsActiveResponse
	DecodeJSON(t, data, &userResp, "set active")
	if userResp.User.IsActive != false {
		t.Error("expected user to be inactive")
	}

	// --- B. Set Unknown User (expect 404) ---
	notFoundReq := map[string]interface{}{
		"user_id":   "nobody_" + suffix,
		"is_active": false,
	}
	respNF, _ := DoPostJSON(t, "/users/setIsActive", notFoundReq)
	AssertStatus(t, respNF, http.StatusNotFound, nil)
}

// 3. Тест полного цикла Pull Request (Create -> Merge)
func TestPullRequestLifecycle(t *testing.T) {
	suffix := RandomSuffix()
	teamName := "pr_team_" + suffix
	author := "author_" + suffix
	reviewer1 := "rev1_" + suffix
	reviewer2 := "rev2_" + suffix
	prID := "pr_" + suffix

	// 1. Создаем команду (Автор + 2 Ревьювера)
	DoPostJSON(t, "/team/add", TeamDTO{
		TeamName: teamName,
		Members: []TeamMemberDTO{
			{UserID: author, Username: "Author", IsActive: true},
			{UserID: reviewer1, Username: "R1", IsActive: true},
			{UserID: reviewer2, Username: "R2", IsActive: true},
		},
	})

	// 2. Создаем PR
	prReq := map[string]string{
		"pull_request_id":   prID,
		"pull_request_name": "New Feature",
		"author_id":         author,
	}
	resp, data := DoPostJSON(t, "/pullRequest/create", prReq)
	AssertStatus(t, resp, http.StatusCreated, data)

	var createResp CreatePRResponse
	DecodeJSON(t, data, &createResp, "create pr")

	// Проверяем, что назначились ревьюверы
	if len(createResp.PR.AssignedReviewers) != 2 {
		t.Errorf("expected 2 reviewers, got %d", len(createResp.PR.AssignedReviewers))
	}
	if createResp.PR.Status != "OPEN" {
		t.Errorf("expected status OPEN, got %s", createResp.PR.Status)
	}

	// 3. Проверяем Duplicate PR (409 PR_EXISTS)
	respDup, dataDup := DoPostJSON(t, "/pullRequest/create", prReq)
	AssertStatus(t, respDup, http.StatusConflict, dataDup)
	AssertErrorCode(t, dataDup, "PR_EXISTS")

	// 4. Проверяем, что у ревьювера появился этот PR в списке (/users/getReview)
	// Берем первого назначенного
	assignedRev := createResp.PR.AssignedReviewers[0]
	respRev, dataRev := DoGet(t, fmt.Sprintf("/users/getReview?user_id=%s", assignedRev))
	AssertStatus(t, respRev, http.StatusOK, dataRev)

	var reviewList GetReviewResponseDTO
	DecodeJSON(t, dataRev, &reviewList, "get review list")
	found := false
	for _, pr := range reviewList.PullRequests {
		if pr.PullRequestID == prID {
			found = true
			break
		}
	}
	if !found {
		t.Error("created PR not found in reviewers list")
	}

	// 5. Merge PR
	mergeReq := map[string]string{"pull_request_id": prID}
	respMerge, dataMerge := DoPostJSON(t, "/pullRequest/merge", mergeReq)
	AssertStatus(t, respMerge, http.StatusOK, dataMerge)

	var mergeResp MergePRResponse
	DecodeJSON(t, dataMerge, &mergeResp, "merge pr")
	if mergeResp.PR.Status != "MERGED" {
		t.Errorf("expected status MERGED, got %s", mergeResp.PR.Status)
	}
}

// 4. Сложный сценарий: Reassign Reviewer
func TestReassignReviewer(t *testing.T) {
	suffix := RandomSuffix()
	teamName := "reassign_team_" + suffix
	author := "auth_" + suffix
	r1 := "r1_" + suffix
	r2 := "r2_" + suffix
	r3 := "r3_" + suffix
	prID := "pr_swap_" + suffix

	// 1. Команда из 4 человек (1 автор, 3 возможных ревьювера)
	DoPostJSON(t, "/team/add", TeamDTO{
		TeamName: teamName,
		Members: []TeamMemberDTO{
			{UserID: author, Username: "Author", IsActive: true},
			{UserID: r1, Username: "R1", IsActive: true},
			{UserID: r2, Username: "R2", IsActive: true},
			{UserID: r3, Username: "R3", IsActive: true},
		},
	})

	// 2. Создаем PR
	respCreate, dataCreate := DoPostJSON(t, "/pullRequest/create", map[string]string{
		"pull_request_id":   prID,
		"pull_request_name": "Swap Test",
		"author_id":         author,
	})
	AssertStatus(t, respCreate, http.StatusCreated, dataCreate)

	var createResp CreatePRResponse
	DecodeJSON(t, dataCreate, &createResp, "create pr response")

	// Берем того, кто был назначен системой
	if len(createResp.PR.AssignedReviewers) == 0 {
		t.Fatal("no reviewers assigned automatically, cannot test reassign")
	}
	// Берем первого попавшегося назначенного ревьювера для замены
	oldReviewerID := createResp.PR.AssignedReviewers[0]

	// 3. Делаем Reassign
	reassignReq := map[string]string{
		"pull_request_id": prID,
		"old_user_id":     oldReviewerID,
	}

	resp, data := DoPostJSON(t, "/pullRequest/reassign", reassignReq)
	AssertStatus(t, resp, http.StatusOK, data)

	var reassignResp ReassignResponseDTO
	DecodeJSON(t, data, &reassignResp, "reassign")

	// Проверяем, что вернулся новый ревьювер
	if reassignResp.ReplacedBy == "" {
		t.Error("expected replaced_by to be set")
	}
	if reassignResp.ReplacedBy == oldReviewerID {
		t.Error("reviewer should be different")
	}

	// Проверяем список ревьюверов в PR
	isOldPresent := false
	isNewPresent := false
	for _, r := range reassignResp.PR.AssignedReviewers {
		if r == oldReviewerID {
			isOldPresent = true
		}
		if r == reassignResp.ReplacedBy {
			isNewPresent = true
		}
	}

	if isOldPresent {
		t.Errorf("old reviewer %s should be removed", oldReviewerID)
	}
	if !isNewPresent {
		t.Errorf("new reviewer %s should be present", reassignResp.ReplacedBy)
	}
}

// 5. Тест ошибки NO_CANDIDATE при переназначении
func TestReassignNoCandidate(t *testing.T) {
	suffix := RandomSuffix()
	teamName := "nocand_team_" + suffix
	author := "au_" + suffix
	r1 := "rr1_" + suffix
	r2 := "rr2_" + suffix
	// Всего 3 человека. 1 Автор, 2 Ревьювера. Больше никого нет.
	// Если попытаться заменить ревьювера, некем заменить.

	DoPostJSON(t, "/team/add", TeamDTO{
		TeamName: teamName,
		Members: []TeamMemberDTO{
			{UserID: author, Username: "Author", IsActive: true},
			{UserID: r1, Username: "R1", IsActive: true},
			{UserID: r2, Username: "R2", IsActive: true},
		},
	})

	prID := "pr_nocand_" + suffix
	respCreate, dataCreate := DoPostJSON(t, "/pullRequest/create", map[string]string{
		"pull_request_id":   prID,
		"pull_request_name": "Test",
		"author_id":         author,
	})
	AssertStatus(t, respCreate, http.StatusCreated, dataCreate)

	var createResp CreatePRResponse
	DecodeJSON(t, dataCreate, &createResp, "create pr response")

	if len(createResp.PR.AssignedReviewers) == 0 {
		t.Fatal("no reviewers assigned")
	}
	targetReviewer := createResp.PR.AssignedReviewers[0]

	resp, data := DoPostJSON(t, "/pullRequest/reassign", map[string]string{
		"pull_request_id": prID,
		"old_user_id":     targetReviewer,
	})

	AssertStatus(t, resp, http.StatusConflict, data)
	AssertErrorCode(t, data, "NO_CANDIDATE")
}
