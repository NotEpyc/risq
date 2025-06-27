package startup

import (
	"context"
	"fmt"
	"time"

	"risq_backend/internal/user"
	"risq_backend/pkg/logger"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, input StartupOnboardingInput, userID uuid.UUID) (*Startup, error)
	CreateWithUserLink(ctx context.Context, input StartupOnboardingInput, userID uuid.UUID) (*Startup, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Startup, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Startup, error)
	Update(ctx context.Context, startup *Startup) error
}

type service struct {
	repo     Repository
	userRepo user.Repository
}

func NewService(repo Repository, userRepo user.Repository) Service {
	return &service{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *service) Create(ctx context.Context, input StartupOnboardingInput, userID uuid.UUID) (*Startup, error) {
	logger.Infof("Creating startup: %s for user: %s", input.Name, userID)

	// Parse founded date
	foundedDate, err := time.Parse("2006-01-02", input.FoundedDate)
	if err != nil {
		return nil, fmt.Errorf("invalid founded date format: %w", err)
	}

	startup := &Startup{
		ID:           uuid.New(),
		Name:         input.Name,
		Description:  input.Description,
		Industry:     input.Industry,
		FundingStage: input.FundingStage,
		Location:     input.Location,
		FoundedDate:  foundedDate,
		TeamSize:     input.TeamSize,
		Website:      input.Website,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, startup); err != nil {
		logger.Errorf("Failed to create startup: %v", err)
		return nil, err
	}

	logger.Infof("Successfully created startup: %s", startup.ID)
	return startup, nil
}

func (s *service) CreateWithUserLink(ctx context.Context, input StartupOnboardingInput, userID uuid.UUID) (*Startup, error) {
	logger.Infof("Creating startup with user link: %s for user: %s", input.Name, userID)

	// First create the startup
	startup, err := s.Create(ctx, input, userID)
	if err != nil {
		return nil, err
	}

	// Update the user to link to this startup
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		logger.Errorf("Failed to get user for linking: %v", err)
		return startup, nil // Don't fail startup creation if user link fails
	}

	user.StartupID = &startup.ID
	if err := s.userRepo.Update(ctx, user); err != nil {
		logger.Errorf("Failed to link user to startup: %v", err)
		// Continue without failing - the startup was created successfully
	}

	logger.Infof("Successfully created startup with user link: %s", startup.ID)
	return startup, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*Startup, error) {
	logger.Debugf("Getting startup by ID: %s", id)

	startup, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Errorf("Failed to get startup by ID: %v", err)
		return nil, err
	}

	return startup, nil
}

func (s *service) GetByUserID(ctx context.Context, userID uuid.UUID) (*Startup, error) {
	logger.Debugf("Getting startup by user ID: %s", userID)

	startup, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		logger.Errorf("Failed to get startup by user ID: %v", err)
		return nil, err
	}

	return startup, nil
}

func (s *service) Update(ctx context.Context, startup *Startup) error {
	logger.Infof("Updating startup: %s", startup.ID)

	if err := s.repo.Update(ctx, startup); err != nil {
		logger.Errorf("Failed to update startup: %v", err)
		return err
	}

	logger.Infof("Successfully updated startup: %s", startup.ID)
	return nil
}
