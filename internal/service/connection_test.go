package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DKeshavarz/peyk/internal/domain"
)

// --- test ---
func TestGenerateCode(t *testing.T) {
	tests := []struct {
		name         string
		genCode      string
		genErr       error
		repoErr      error
		sourceChatID string
		wantErr      bool
	}{
		{
			name:         "success",
			genCode:      "ABC123",
			sourceChatID: "chat-1",
			wantErr:      false,
		},
		{
			name:    "generator fails",
			genErr:  errors.New("gen failed"),
			wantErr: true,
		},
		{
			name:         "repo fails",
			genCode:      "XYZ999",
			repoErr:      errors.New("db error"),
			sourceChatID: "chat-2",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := &fakeGen{code: tt.genCode, err: tt.genErr}
			repo := &fakeRepo{err: tt.repoErr}

			u := NewConnectionUsecase(gen, repo, time.Minute)

			code, err := u.GenerateCode(context.Background(), tt.sourceChatID)

			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !tt.wantErr && code != tt.genCode {
				t.Fatalf("expected code %s, got %s", tt.genCode, code)
			}

			if !tt.wantErr {
				if repo.saved == nil {
					t.Fatalf("expected code to be saved")
				}
				if repo.saved.SourceChat != tt.sourceChatID {
					t.Fatalf("wrong source chat saved")
				}
			}
		})
	}
}

// --- fakes ---

type fakeGen struct {
	code string
	err  error
}

func (f *fakeGen) Generate() (string, error) {
	return f.code, f.err
}

type fakeRepo struct {
	saved *domain.ConnectionCode
	err   error
}

func (r *fakeRepo) Save(_ context.Context, c *domain.ConnectionCode) error {
	if r.err != nil {
		return r.err
	}
	r.saved = c
	return nil
}
