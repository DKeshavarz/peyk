package service

import (
	"context"
	"time"

	"github.com/DKeshavarz/peyk/internal/domain"
)

type CodeGenerator interface {
	Generate() (string, error)
}

type ConnectionCodeRepository interface {
	Save(ctx context.Context, code *domain.ConnectionCode) error
}

type ConnectionUsecase interface {
	GenerateCode(ctx context.Context, sourceChatID string, platform domain.PlatformName) (code string, err error)
}

type connectionUsecase struct {
	codeGen  CodeGenerator
	codeRepo ConnectionCodeRepository
	ttl      time.Duration
	now      func() time.Time
}

func NewConnectionUsecase(codeGen CodeGenerator, codeRepo ConnectionCodeRepository, ttl time.Duration) ConnectionUsecase {
	return &connectionUsecase{
		codeGen:  codeGen,
		codeRepo: codeRepo,
		ttl:      ttl,
		now:      time.Now,
	}
}

func (u *connectionUsecase) GenerateCode(ctx context.Context, sourceChatID string, platform domain.PlatformName) (string, error) {

	code, err := u.codeGen.Generate()
	if err != nil {
		return "", err
	}

	cc := &domain.ConnectionCode{
		Code:      code,
		ChatID:    sourceChatID,
		Platform:  platform,
		ExpiresAt: u.now().Add(u.ttl),
	}

	if err := u.codeRepo.Save(ctx, cc); err != nil {
		return "", err
	}

	return code, nil
}
