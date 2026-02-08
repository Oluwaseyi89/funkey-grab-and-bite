package v1

import (
	"net/http"
	"strconv"
	"strings"

	"funkey-grab-and-bite/funkey-bite-api/internal/handlers"

	"github.com/gin-gonic/gin"

	"funkey-grab-and-bite/funkey-bite-api/internal/services"
)

type MenuHandler struct {
	menuService services.MenuService
}

func NewMenuHandler(menuService services.MenuService) *MenuHandler {
	return &MenuHandler{
		menuService: menuService,
	}
}

// GetMenu returns all available menu items
// @Summary Get all menu items
// @Description Get all available menu items
// @Tags menu
// @Accept json
// @Produce json
// @Success 200 {array} models.MenuItem
// @Router /menu [get]
func (h *MenuHandler) GetMenu(c *gin.Context) {
	menuItems, err := h.menuService.GetMenuItems()
	if err != nil {
		handlers.ErrorWithDetails(c, http.StatusInternalServerError,
			"MENU_FETCH_FAILED", "Failed to fetch menu items", err.Error())
		return
	}

	handlers.Success(c, menuItems)
}

// GetMenuItem returns a specific menu item by ID
// @Summary Get menu item by ID
// @Description Get a specific menu item by its ID
// @Tags menu
// @Accept json
// @Produce json
// @Param id path int true "Menu Item ID"
// @Success 200 {object} models.MenuItem
// @Failure 404 {object} map[string]string
// @Router /menu/{id} [get]
func (h *MenuHandler) GetMenuItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid menu item ID",
		})
		return
	}

	menuItem, err := h.menuService.GetMenuItemByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Menu item not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch menu item",
			})
		}
		return
	}

	c.JSON(http.StatusOK, menuItem)
}

// GetCategories returns all menu categories
// @Summary Get all menu categories
// @Description Get all available menu categories
// @Tags menu
// @Accept json
// @Produce json
// @Success 200 {array} models.MenuCategory
// @Router /categories [get]
func (h *MenuHandler) GetCategories(c *gin.Context) {
	categories, err := h.menuService.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch categories",
		})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetMenuByCategory returns menu items by category
// @Summary Get menu items by category
// @Description Get all menu items in a specific category
// @Tags menu
// @Accept json
// @Produce json
// @Param category_id query int true "Category ID"
// @Success 200 {array} models.MenuItem
// @Router /menu/category [get]
func (h *MenuHandler) GetMenuByCategory(c *gin.Context) {
	categoryIDStr := c.Query("category_id")
	if categoryIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "category_id query parameter is required",
		})
		return
	}

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid category ID",
		})
		return
	}

	menuItems, err := h.menuService.GetMenuItemsByCategory(categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch menu items by category",
		})
		return
	}

	c.JSON(http.StatusOK, menuItems)
}

// SearchMenu searches for menu items
// @Summary Search menu items
// @Description Search for menu items by name or description
// @Tags menu
// @Accept json
// @Produce json
// @Param query query string true "Search query"
// @Success 200 {array} models.MenuItem
// @Router /menu/search [get]
func (h *MenuHandler) SearchMenu(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "query parameter is required",
		})
		return
	}

	menuItems, err := h.menuService.SearchMenuItems(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to search menu items",
		})
		return
	}

	c.JSON(http.StatusOK, menuItems)
}

// GetFeaturedItems returns featured menu items (optional)
// @Summary Get featured menu items
// @Description Get featured/popular menu items
// @Tags menu
// @Accept json
// @Produce json
// @Success 200 {array} models.MenuItem
// @Router /menu/featured [get]
func (h *MenuHandler) GetFeaturedItems(c *gin.Context) {
	// For now, return all items. You can implement proper featured logic later
	menuItems, err := h.menuService.GetMenuItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch featured items",
		})
		return
	}

	// Limit to 6 featured items for now
	limit := 6
	if len(menuItems) > limit {
		menuItems = menuItems[:limit]
	}

	c.JSON(http.StatusOK, menuItems)
}
