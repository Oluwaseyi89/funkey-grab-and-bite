package services

import (
	"context"
	"errors"
	"io"
	"testing"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

type memoryTeamStore struct {
	members map[int]*models.TeamMember
	nextID  int
}

func (m *memoryTeamStore) List(bool) ([]models.TeamMember, error) { return nil, nil }
func (m *memoryTeamStore) GetByID(id int) (*models.TeamMember, error) {
	if member, ok := m.members[id]; ok {
		copy := *member
		return &copy, nil
	}
	return nil, nil
}
func (m *memoryTeamStore) Create(input *models.TeamMemberInput) (*models.TeamMember, error) {
	m.nextID++
	member := &models.TeamMember{ID: m.nextID, Name: input.Name, Role: input.Role, Bio: input.Bio, PhotoURL: input.PhotoURL}
	m.members[member.ID] = member
	return member, nil
}
func (m *memoryTeamStore) Update(id int, input *models.TeamMemberInput) (*models.TeamMember, error) {
	member, ok := m.members[id]
	if !ok {
		return nil, nil
	}
	member.Name, member.Role, member.PhotoURL = input.Name, input.Role, input.PhotoURL
	return member, nil
}
func (m *memoryTeamStore) Delete(id int) (bool, error) {
	_, ok := m.members[id]
	delete(m.members, id)
	return ok, nil
}

type recordingStorage struct{ deleted []string }

func (r *recordingStorage) Put(context.Context, string, string, io.Reader, int64) error { return nil }
func (r *recordingStorage) Delete(_ context.Context, key string) error {
	r.deleted = append(r.deleted, key)
	return nil
}
func (r *recordingStorage) PublicBaseURL() string { return "https://cdn.example.com/assets" }

func TestTeamServiceNormalizesAndCleansUpPhotos(t *testing.T) {
	store := &memoryTeamStore{members: map[int]*models.TeamMember{}}
	images := &recordingStorage{}
	svc := NewTeamService(store, images)

	photo := "https://cdn.example.com/assets/team/abc.png"
	member, err := svc.Create(&models.TeamMemberInput{Name: "  Bisi Ade ", Role: " Head Chef ", PhotoURL: &photo, IsActive: true})
	if err != nil {
		t.Fatal(err)
	}
	if member.Name != "Bisi Ade" || member.Role != "Head Chef" || *member.PhotoURL != "team/abc.png" {
		t.Fatalf("expected trimmed fields and the photo stored as a key, got %+v", member)
	}

	newPhoto := "team/def.png"
	if _, err := svc.Update(member.ID, &models.TeamMemberInput{Name: "Bisi Ade", Role: "Head Chef", PhotoURL: &newPhoto}); err != nil {
		t.Fatal(err)
	}
	if len(images.deleted) != 1 || images.deleted[0] != "team/abc.png" {
		t.Fatalf("expected the replaced photo deleted, got %v", images.deleted)
	}

	if err := svc.Delete(member.ID); err != nil {
		t.Fatal(err)
	}
	if len(images.deleted) != 2 || images.deleted[1] != "team/def.png" {
		t.Fatalf("expected the photo deleted with the member, got %v", images.deleted)
	}

	inline := "data:image/png;base64,AAAA"
	var invalid *TeamValidationError
	if _, err := svc.Create(&models.TeamMemberInput{Name: "Tunde", Role: "Rider", PhotoURL: &inline}); !errors.As(err, &invalid) {
		t.Fatalf("inline image data must be rejected, got %v", err)
	}
	if _, err := svc.Create(&models.TeamMemberInput{Name: " ", Role: "Rider"}); !errors.As(err, &invalid) {
		t.Fatalf("blank name must be rejected, got %v", err)
	}
	if err := svc.Delete(404); !errors.Is(err, ErrTeamMemberNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
	if _, err := svc.Update(404, &models.TeamMemberInput{Name: "Tunde", Role: "Rider"}); !errors.Is(err, ErrTeamMemberNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}
