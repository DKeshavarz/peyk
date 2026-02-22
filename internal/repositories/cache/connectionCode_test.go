package cache

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/DKeshavarz/peyk/internal/domain"
	"github.com/DKeshavarz/peyk/internal/infra/cache"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectionCodeRepository_Save(t *testing.T) {
	c := cache.New()
	repo := NewConnectionCodeRepository(c)
	ctx := context.Background()

	tests := []struct {
		name    string
		code    *domain.ConnectionCode
		wantErr bool
	}{
		{
			name: "successful save",
			code: &domain.ConnectionCode{
				Code:      "TEST-123",
				ChatID:    "chat_1",
				ExpiresAt: time.Now().Add(5 * time.Minute),
			},
			wantErr: false,
		},
		{
			name:    "nil code",
			code:    nil,
			wantErr: true,
		},
		{
			name: "expired code",
			code: &domain.ConnectionCode{
				Code:      "EXPIRED-123",
				ChatID:    "chat_2",
				ExpiresAt: time.Now().Add(-5 * time.Minute),
			},
			wantErr: true,
		},
		{
			name: "zero TTL code",
			code: &domain.ConnectionCode{
				Code:      "ZERO-123",
				ChatID:    "chat_3",
				ExpiresAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "empty code string",
			code: &domain.ConnectionCode{
				Code:      "",
				ChatID:    "chat_4",
				ExpiresAt: time.Now().Add(5 * time.Minute),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Save(ctx, tt.code)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify it was saved correctly
				if tt.code != nil && tt.code.Code != "" {
					saved, err := repo.FindByCode(ctx, tt.code.Code)
					assert.NoError(t, err)
					assert.Equal(t, tt.code.Code, saved.Code)
					assert.Equal(t, tt.code.ChatID, saved.ChatID)
					assert.WithinDuration(t, tt.code.ExpiresAt, saved.ExpiresAt, time.Second)
				}
			}
		})
	}
}

func TestConnectionCodeRepository_Save_Overwrite(t *testing.T) {
	c := cache.New()
	repo := NewConnectionCodeRepository(c)
	ctx := context.Background()

	code1 := &domain.ConnectionCode{
		Code:      "SAME-CODE",
		ChatID:    "chat_1",
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	err := repo.Save(ctx, code1)
	require.NoError(t, err)

	code2 := &domain.ConnectionCode{
		Code:      "SAME-CODE",
		ChatID:    "chat_2",
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	err = repo.Save(ctx, code2)
	assert.NoError(t, err)

	saved, err := repo.FindByCode(ctx, "SAME-CODE")
	assert.NoError(t, err)
	assert.Equal(t, "SAME-CODE", saved.Code)
	assert.Equal(t, "chat_2", saved.ChatID)
	assert.WithinDuration(t, code2.ExpiresAt, saved.ExpiresAt, time.Second)
}

func TestConnectionCodeRepository_FindByCode(t *testing.T) {
	c := cache.New()
	repo := NewConnectionCodeRepository(c)
	ctx := context.Background()

	testCodes := []*domain.ConnectionCode{
		{
			Code:      "FIND-001",
			ChatID:    "chat_1",
			ExpiresAt: time.Now().Add(5 * time.Minute),
		},
		{
			Code:      "FIND-002",
			ChatID:    "chat_2",
			ExpiresAt: time.Now().Add(10 * time.Minute),
		},
	}

	for _, tc := range testCodes {
		err := repo.Save(ctx, tc)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		code      string
		wantErr   bool
		checkFunc func(t *testing.T, found *domain.ConnectionCode)
	}{
		{
			name:    "existing code",
			code:    "FIND-001",
			wantErr: false,
			checkFunc: func(t *testing.T, found *domain.ConnectionCode) {
				assert.Equal(t, "FIND-001", found.Code)
				assert.Equal(t, "chat_1", found.ChatID)
			},
		},
		{
			name:    "another existing code",
			code:    "FIND-002",
			wantErr: false,
			checkFunc: func(t *testing.T, found *domain.ConnectionCode) {
				assert.Equal(t, "FIND-002", found.Code)
				assert.Equal(t, "chat_2", found.ChatID)
			},
		},
		{
			name:    "non-existent code",
			code:    "NON-EXISTENT",
			wantErr: true,
		},
		{
			name:    "empty code string",
			code:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := repo.FindByCode(ctx, tt.code)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, found)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, found)
				if tt.checkFunc != nil {
					tt.checkFunc(t, found)
				}
			}
		})
	}
}

func TestConnectionCodeRepository_FindByCode_Expired(t *testing.T) {
	c := cache.New()
	repo := NewConnectionCodeRepository(c)
	ctx := context.Background()

	// Save a code that expires soon
	code := &domain.ConnectionCode{
		Code:      "EXPIRING",
		ChatID:    "chat_expire",
		ExpiresAt: time.Now().Add(2 * time.Second),
	}

	err := repo.Save(ctx, code)
	require.NoError(t, err)

	// Find it immediately - should work
	found, err := repo.FindByCode(ctx, "EXPIRING")
	assert.NoError(t, err)
	assert.NotNil(t, found)

	// Wait for expiration
	time.Sleep(3 * time.Second)

	// Try to find after expiration - cache should auto-delete
	_, err = repo.FindByCode(ctx, "EXPIRING")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "code not found")
}

func TestConnectionCodeRepository_Delete(t *testing.T) {
	c := cache.New()
	repo := NewConnectionCodeRepository(c)
	ctx := context.Background()

	tests := []struct {
		name     string
		setup    func() string // returns code to delete
		wantErr  bool
		validate func(t *testing.T, code string)
	}{
		{
			name: "delete existing code",
			setup: func() string {
				code := &domain.ConnectionCode{
					Code:      "DELETE-001",
					ChatID:    "chat_delete",
					ExpiresAt: time.Now().Add(5 * time.Minute),
				}
				err := repo.Save(ctx, code)
				require.NoError(t, err)
				return "DELETE-001"
			},
			wantErr: false,
			validate: func(t *testing.T, code string) {
				// Verify it's gone
				_, err := repo.FindByCode(ctx, code)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "code not found")
			},
		},
		{
			name: "delete non-existent code",
			setup: func() string {
				return "NON-EXISTENT"
			},
			wantErr: false, // Should return nil for non-existent
			validate: func(t *testing.T, code string) {
				// Should still be gone (no error)
				_, err := repo.FindByCode(ctx, code)
				assert.Error(t, err)
			},
		},
		{
			name: "delete empty code",
			setup: func() string {
				return ""
			},
			wantErr: false,
		},
		{
			name: "delete twice",
			setup: func() string {
				code := &domain.ConnectionCode{
					Code:      "DELETE-TWICE",
					ChatID:    "chat_twice",
					ExpiresAt: time.Now().Add(5 * time.Minute),
				}
				err := repo.Save(ctx, code)
				require.NoError(t, err)
				return "DELETE-TWICE"
			},
			wantErr: false,
			validate: func(t *testing.T, code string) {
				// First delete
				err := repo.Delete(ctx, code)
				assert.NoError(t, err)

				// Second delete - should not error
				err = repo.Delete(ctx, code)
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codeToDelete := tt.setup()

			err := repo.Delete(ctx, codeToDelete)

			if tt.wantErr {
				assert.Error(t, err)

			} else {
				assert.NoError(t, err)
			}

			if tt.validate != nil {
				tt.validate(t, codeToDelete)
			}
		})
	}
}

func TestConnectionCodeRepository_DataIntegrity(t *testing.T) {
	c := cache.New()
	repo := NewConnectionCodeRepository(c)
	ctx := context.Background()

	// Test with various data
	testCases := []struct {
		name    string
		code    *domain.ConnectionCode
		wantErr bool
	}{
		{
			name: "special characters in code",
			code: &domain.ConnectionCode{
				Code:      "TEST-@#$%-123",
				ChatID:    "chat_special",
				ExpiresAt: time.Now().Add(5 * time.Minute),
			},
		},
		{
			name: "unicode in source chat",
			code: &domain.ConnectionCode{
				Code:      "UNICODE-001",
				ChatID:    "چت_فارسی_123",
				ExpiresAt: time.Now().Add(5 * time.Minute),
			},
		},
		{
			name: "very long code",
			code: &domain.ConnectionCode{
				Code:      "THIS-IS-A-VERY-LONG-CODE-THAT-MIGHT-CAUSE-ISSUES-1234567890",
				ChatID:    "chat_long",
				ExpiresAt: time.Now().Add(5 * time.Minute),
			},
		},
		{
			name: "very long source chat",
			code: &domain.ConnectionCode{
				Code:      "LONG-SRC-001",
				ChatID:    "this-is-a-very-long-source-chat-id-that-might-cause-issues-1234567890",
				ExpiresAt: time.Now().Add(5 * time.Minute),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Save
			err := repo.Save(ctx, tc.code)
			assert.NoError(t, err)

			// Find
			found, err := repo.FindByCode(ctx, tc.code.Code)
			assert.NoError(t, err)
			assert.Equal(t, tc.code.Code, found.Code)
			assert.Equal(t, tc.code.ChatID, found.ChatID)
			assert.WithinDuration(t, tc.code.ExpiresAt, found.ExpiresAt, time.Second)

			// Marshal/Unmarshal test
			data, err := json.Marshal(found)
			assert.NoError(t, err)

			var unmarshaled domain.ConnectionCode
			err = json.Unmarshal(data, &unmarshaled)
			assert.NoError(t, err)
			assert.Equal(t, tc.code.Code, unmarshaled.Code)
		})
	}
}
