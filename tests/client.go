package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"
)

// URL возвращает адрес сервиса
func URL() string {
	if url := os.Getenv("HTTP_URL"); url != "" {
		return url
	}
	return "http://localhost:8080"
}

// RandomSuffix генерирует случайную строку для уникальности данных в тестах
func RandomSuffix() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%d", r.Intn(100000))
}

// DoPostJSON выполняет POST запрос
func DoPostJSON(t *testing.T, path string, body any) (*http.Response, []byte) {
	t.Helper()

	var buf io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body for POST %s: %v", path, err)
		}
		buf = bytes.NewReader(b)
	}

	req, err := http.NewRequest(http.MethodPost, URL()+path, buf)
	if err != nil {
		t.Fatalf("new POST request %s: %v", path, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do POST %s: %v", path, err)
	}

	data, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("read body for POST %s: %v", path, err)
	}

	return resp, data
}

// DoGet выполняет GET запрос
func DoGet(t *testing.T, path string) (*http.Response, []byte) {
	t.Helper()

	req, err := http.NewRequest(http.MethodGet, URL()+path, nil)
	if err != nil {
		t.Fatalf("new GET request %s: %v", path, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do GET %s: %v", path, err)
	}

	data, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("read body for GET %s: %v", path, err)
	}

	return resp, data
}

type ErrorDTO struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponseDTO struct {
	Error ErrorDTO `json:"error"`
}

type TeamMemberDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type TeamDTO struct {
	TeamName string          `json:"team_name"`
	Members  []TeamMemberDTO `json:"members"`
}

type AddTeamResponse struct {
	Team TeamDTO `json:"team"`
}

type UserDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type SetIsActiveResponse struct {
	User UserDTO `json:"user"`
}

type PullRequestDTO struct {
	PullRequestID     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
	CreatedAt         *string  `json:"createdAt,omitempty"`
	MergedAt          *string  `json:"mergedAt,omitempty"`
}

type CreatePRResponse struct {
	PR PullRequestDTO `json:"pr"`
}

type MergePRResponse struct {
	PR PullRequestDTO `json:"pr"`
}

type PRShortDTO struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

type GetReviewResponseDTO struct {
	UserID       string       `json:"user_id"`
	PullRequests []PRShortDTO `json:"pull_requests"`
}

type ReassignResponseDTO struct {
	PR         PullRequestDTO `json:"pr"`
	ReplacedBy string         `json:"replaced_by"`
}

// --- Helpers ---

func DecodeJSON(t *testing.T, data []byte, v any, context string) {
	t.Helper()
	if err := json.Unmarshal(data, v); err != nil {
		t.Fatalf("decode json (%s) failed: %v, body=%s", context, err, string(data))
	}
}

func DecodeError(t *testing.T, data []byte) ErrorResponseDTO {
	t.Helper()
	var er ErrorResponseDTO
	DecodeJSON(t, data, &er, "error response")
	return er
}

func AssertStatus(t *testing.T, resp *http.Response, expected int, body []byte) {
	t.Helper()
	if resp.StatusCode != expected {
		t.Fatalf("expected status %d, got %d. body: %s", expected, resp.StatusCode, string(body))
	}
}

func AssertErrorCode(t *testing.T, data []byte, expectedCode string) {
	t.Helper()
	errResp := DecodeError(t, data)
	if errResp.Error.Code != expectedCode {
		t.Fatalf("expected error code %s, got %s. Msg: %s", expectedCode, errResp.Error.Code, errResp.Error.Message)
	}
}
