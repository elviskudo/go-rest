package handlers

import (
	"go-rest/internal/database"
	"go-rest/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateRole creates a new role
// CreateRole godoc
// @Summary      Create a role
// @Description  Create a new role
// @Tags         rbac
// @Accept       json
// @Produce      json
// @Param        role  body      models.Role  true  "Role JSON"
// @Success      201   {object}  models.Role
// @Failure      400   {object}  gin.H
// @Failure      500   {object}  gin.H
// @Security     BearerAuth
// @Router       /rbac/roles [post]
func CreateRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, role)
}

// GetRoles godoc
// @Summary      List roles
// @Description  Get all roles
// @Tags         rbac
// @Produce      json
// @Success      200  {array}   models.Role
// @Failure      500  {object}  gin.H
// @Security     BearerAuth
// @Router       /rbac/roles [get]
func GetRoles(c *gin.Context) {
	var roles []models.Role
	if err := database.DB.Preload("Permissions").Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// CreatePermission godoc
// @Summary      Create a permission
// @Description  Create a new permission
// @Tags         rbac
// @Accept       json
// @Produce      json
// @Param        permission  body      models.Permission  true  "Permission JSON"
// @Success      201         {object}  models.Permission
// @Failure      400         {object}  gin.H
// @Failure      500         {object}  gin.H
// @Security     BearerAuth
// @Router       /rbac/permissions [post]
func CreatePermission(c *gin.Context) {
	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&permission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, permission)
}

// GetPermissions godoc
// @Summary      List permissions
// @Description  Get all permissions
// @Tags         rbac
// @Produce      json
// @Success      200  {array}   models.Permission
// @Failure      500  {object}  gin.H
// @Security     BearerAuth
// @Router       /rbac/permissions [get]
func GetPermissions(c *gin.Context) {
	var permissions []models.Permission
	if err := database.DB.Find(&permissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

// AssignPermissionsToRole godoc
// @Summary      Assign permissions to role
// @Description  Assign a list of permissions to a role
// @Tags         rbac
// @Accept       json
// @Produce      json
// @Param        id     path      string  true  "Role ID"
// @Param        input  body      object  true  "Permission IDs"
// @Success      200    {object}  gin.H
// @Failure      400    {object}  gin.H
// @Failure      404    {object}  gin.H
// @Failure      500    {object}  gin.H
// @Security     BearerAuth
// @Router       /rbac/roles/{id}/permissions [post]
func AssignPermissionsToRole(c *gin.Context) {
	roleID := c.Param("id")
	var input struct {
		PermissionIDs []string `json:"permission_ids"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var role models.Role
	if err := database.DB.First(&role, "id = ?", roleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	var permissions []models.Permission
	if err := database.DB.Where("id IN ?", input.PermissionIDs).Find(&permissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Model(&role).Association("Permissions").Replace(permissions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permissions assigned successfully"})
}

// AssignRoleToUser godoc
// @Summary      Assign role to user
// @Description  Assign a role to a user
// @Tags         rbac
// @Accept       json
// @Produce      json
// @Param        id     path      string  true  "User ID"
// @Param        input  body      object  true  "Role ID"
// @Success      200    {object}  gin.H
// @Failure      400    {object}  gin.H
// @Failure      404    {object}  gin.H
// @Failure      500    {object}  gin.H
// @Security     BearerAuth
// @Router       /rbac/users/{id}/role [post]
func AssignRoleToUser(c *gin.Context) {
	userID := c.Param("id")
	var input struct {
		RoleID string `json:"role_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.RoleID = uuid.MustParse(input.RoleID)
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role assigned successfully"})
}
