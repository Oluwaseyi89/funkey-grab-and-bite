package v1

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/services"
)

var teamFieldLabels = map[string]string{
	"name": "Name", "role": "Role", "bio": "Bio", "photoUrl": "Photo", "displayOrder": "Display order",
}

type TeamHandler struct {
	teamService services.TeamService
	validate    *validator.Validate
}

func NewTeamHandler(teamService services.TeamService) *TeamHandler {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
	})
	return &TeamHandler{teamService: teamService, validate: validate}
}

func respondTeamError(c *gin.Context, err error, fallback string) {
	var invalid *services.TeamValidationError
	switch {
	case errors.As(err, &invalid):
		c.JSON(http.StatusBadRequest, gin.H{"error": invalid.Message})
	case errors.Is(err, services.ErrTeamMemberNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "Team member not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": fallback})
	}
}

// ListPublic lists the team members shown on the contact page
// @Summary List team members
// @Tags public
// @Produce json
// @Success 200 {array} models.TeamMember
// @Router /team [get]
func (h *TeamHandler) ListPublic(c *gin.Context) {
	h.list(c, true)
}

// ListAll lists every team member, including hidden ones (admin)
// @Summary List all team members
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.TeamMember
// @Router /admin/team [get]
func (h *TeamHandler) ListAll(c *gin.Context) {
	h.list(c, false)
}

func (h *TeamHandler) list(c *gin.Context, activeOnly bool) {
	members, err := h.teamService.List(activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch team members"})
		return
	}
	c.JSON(http.StatusOK, members)
}

// Create adds a team member (admin)
// @Summary Create team member
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param member body models.TeamMemberInput true "Team member"
// @Success 201 {object} models.TeamMember
// @Router /admin/team [post]
func (h *TeamHandler) Create(c *gin.Context) {
	input, ok := h.bind(c)
	if !ok {
		return
	}
	member, err := h.teamService.Create(input)
	if err != nil {
		respondTeamError(c, err, "Failed to add team member")
		return
	}
	c.JSON(http.StatusCreated, member)
}

// Update replaces a team member (admin)
// @Summary Update team member
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Team member ID"
// @Param member body models.TeamMemberInput true "Team member"
// @Success 200 {object} models.TeamMember
// @Router /admin/team/{id} [put]
func (h *TeamHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team member ID"})
		return
	}
	input, ok := h.bind(c)
	if !ok {
		return
	}
	member, err := h.teamService.Update(id, input)
	if err != nil {
		respondTeamError(c, err, "Failed to update team member")
		return
	}
	c.JSON(http.StatusOK, member)
}

// Delete removes a team member (admin)
// @Summary Delete team member
// @Tags admin
// @Security BearerAuth
// @Param id path int true "Team member ID"
// @Success 200 {object} map[string]string
// @Router /admin/team/{id} [delete]
func (h *TeamHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team member ID"})
		return
	}
	if err := h.teamService.Delete(id); err != nil {
		respondTeamError(c, err, "Failed to delete team member")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Team member removed"})
}

func (h *TeamHandler) bind(c *gin.Context) (*models.TeamMemberInput, bool) {
	var input models.TeamMemberInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return nil, false
	}
	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationMessage(err, teamFieldLabels)})
		return nil, false
	}
	return &input, true
}
