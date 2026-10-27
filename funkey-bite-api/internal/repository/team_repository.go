package repository

import (
	"database/sql"
	"fmt"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

const teamMemberColumns = `id, name, role, bio, photo_url, display_order, is_active, created_at, updated_at`

// List returns team members in display order; activeOnly hides hidden ones.
func (r *TeamRepository) List(activeOnly bool) ([]models.TeamMember, error) {
	rows, err := r.db.Query(`SELECT `+teamMemberColumns+` FROM team_members
		WHERE ($1 = false OR is_active = true)
		ORDER BY display_order, id`, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to list team members: %w", err)
	}
	defer rows.Close()

	members := []models.TeamMember{}
	for rows.Next() {
		member, err := scanTeamMember(rows)
		if err != nil {
			return nil, err
		}
		members = append(members, *member)
	}
	return members, rows.Err()
}

func (r *TeamRepository) GetByID(id int) (*models.TeamMember, error) {
	member, err := scanTeamMember(r.db.QueryRow(`SELECT `+teamMemberColumns+` FROM team_members WHERE id = $1`, id))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return member, err
}

func (r *TeamRepository) Create(input *models.TeamMemberInput) (*models.TeamMember, error) {
	var id int
	err := r.db.QueryRow(`
		INSERT INTO team_members (name, role, bio, photo_url, display_order, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		input.Name, input.Role, input.Bio, input.PhotoURL, input.DisplayOrder, input.IsActive,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to create team member: %w", err)
	}
	return r.GetByID(id)
}

// Update replaces a team member; it returns nil if there is none with that id.
func (r *TeamRepository) Update(id int, input *models.TeamMemberInput) (*models.TeamMember, error) {
	result, err := r.db.Exec(`
		UPDATE team_members SET
			name = $1, role = $2, bio = $3, photo_url = $4,
			display_order = $5, is_active = $6, updated_at = $7
		WHERE id = $8`,
		input.Name, input.Role, input.Bio, input.PhotoURL,
		input.DisplayOrder, input.IsActive, time.Now(), id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update team member: %w", err)
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return nil, nil
	}
	return r.GetByID(id)
}

// Delete removes a team member and reports whether one existed.
func (r *TeamRepository) Delete(id int) (bool, error) {
	result, err := r.db.Exec(`DELETE FROM team_members WHERE id = $1`, id)
	if err != nil {
		return false, fmt.Errorf("failed to delete team member: %w", err)
	}
	rows, _ := result.RowsAffected()
	return rows > 0, nil
}

func scanTeamMember(row rowScanner) (*models.TeamMember, error) {
	var member models.TeamMember
	var photo sql.NullString
	err := row.Scan(&member.ID, &member.Name, &member.Role, &member.Bio, &photo,
		&member.DisplayOrder, &member.IsActive, &member.CreatedAt, &member.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan team member: %w", err)
	}
	member.PhotoURL = nullStringPtr(photo)
	return &member, nil
}
