package services

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/storage"
)

var ErrTeamMemberNotFound = errors.New("team member not found")

// TeamValidationError explains why a team member was rejected.
type TeamValidationError struct {
	Message string
}

func (e *TeamValidationError) Error() string { return e.Message }

type TeamService interface {
	List(activeOnly bool) ([]models.TeamMember, error)
	Create(input *models.TeamMemberInput) (*models.TeamMember, error)
	Update(id int, input *models.TeamMemberInput) (*models.TeamMember, error)
	Delete(id int) error
}

type teamStore interface {
	List(activeOnly bool) ([]models.TeamMember, error)
	GetByID(id int) (*models.TeamMember, error)
	Create(input *models.TeamMemberInput) (*models.TeamMember, error)
	Update(id int, input *models.TeamMemberInput) (*models.TeamMember, error)
	Delete(id int) (bool, error)
}

type teamService struct {
	repo   teamStore
	images storage.Storage // nil when storage is not configured
}

func NewTeamService(repo teamStore, images storage.Storage) TeamService {
	return &teamService{repo: repo, images: images}
}

func (s *teamService) List(activeOnly bool) ([]models.TeamMember, error) {
	return s.repo.List(activeOnly)
}

func (s *teamService) Create(input *models.TeamMemberInput) (*models.TeamMember, error) {
	if err := s.normalize(input); err != nil {
		return nil, err
	}
	return s.repo.Create(input)
}

func (s *teamService) Update(id int, input *models.TeamMemberInput) (*models.TeamMember, error) {
	if err := s.normalize(input); err != nil {
		return nil, err
	}

	current, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, ErrTeamMemberNotFound
	}

	updated, err := s.repo.Update(id, input)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, ErrTeamMemberNotFound
	}
	if current.PhotoURL != nil && (input.PhotoURL == nil || *input.PhotoURL != *current.PhotoURL) {
		s.deletePhoto(*current.PhotoURL)
	}
	return updated, nil
}

func (s *teamService) Delete(id int) error {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if current == nil {
		return ErrTeamMemberNotFound
	}
	if _, err := s.repo.Delete(id); err != nil {
		return err
	}
	if current.PhotoURL != nil {
		s.deletePhoto(*current.PhotoURL)
	}
	return nil
}

// normalize trims the text and stores uploaded photos by key, so a resolved
// public URL sent back by the admin app is saved as the key again.
func (s *teamService) normalize(input *models.TeamMemberInput) error {
	input.Name = strings.TrimSpace(input.Name)
	input.Role = strings.TrimSpace(input.Role)
	input.Bio = strings.TrimSpace(input.Bio)
	if len(input.Name) < 2 || len(input.Role) < 2 {
		return &TeamValidationError{Message: "Name and role are required"}
	}

	if input.PhotoURL != nil {
		photo := strings.TrimSpace(*input.PhotoURL)
		if strings.HasPrefix(strings.ToLower(photo), "data:") {
			return &TeamValidationError{Message: "Upload the photo instead of pasting image data"}
		}
		photo = storage.KeyFromPublicURL(s.images, photo)
		if photo == "" {
			input.PhotoURL = nil
		} else {
			input.PhotoURL = &photo
		}
	}
	return nil
}

// deletePhoto removes a replaced photo; failure only leaves an unused file.
func (s *teamService) deletePhoto(value string) {
	if s.images == nil || !storage.IsManagedKey(value) {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.images.Delete(ctx, value); err != nil {
		log.Printf("failed to delete team photo %s: %v", value, err)
	}
}
