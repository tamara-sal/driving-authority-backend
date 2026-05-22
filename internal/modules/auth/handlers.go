package auth

import (
	"net/http"

	"driving-authority-backend/internal/http/middleware"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	svc *Service
}

func NewHandlers(svc *Service) *Handlers {
	return &Handlers{svc: svc}
}

// Register godoc
// @Summary      Register a new citizen account
// @Description  Creates a citizen user and returns a JWT access token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      RegisterInput  true  "Registration payload"
// @Success      201   {object}  AuthOutput
// @Failure      400   {object}  map[string]string
// @Router       /auth/register [post]
func (h *Handlers) Register(c *gin.Context) {
	var in RegisterInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := h.svc.Register(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, out)
}

// Login godoc
// @Summary      Login
// @Description  Authenticates a user and returns a JWT access token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      LoginInput  true  "Login credentials"
// @Success      200   {object}  AuthOutput
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Router       /auth/login [post]
func (h *Handlers) Login(c *gin.Context) {
	var in LoginInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := h.svc.Login(c.Request.Context(), in)
	if err != nil {
		if err == ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

// BootstrapAdmin godoc
// @Summary      Bootstrap first admin account
// @Description  Creates an admin user when BOOTSTRAP_ADMIN_SECRET is configured and the request secret matches.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      BootstrapAdminInput  true  "Bootstrap payload"
// @Success      201   {object}  AuthOutput
// @Failure      400   {object}  map[string]string
// @Failure      403   {object}  map[string]string
// @Router       /auth/bootstrap-admin [post]
func (h *Handlers) BootstrapAdmin(c *gin.Context) {
	var in BootstrapAdminInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := h.svc.BootstrapAdmin(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, out)
}

func (h *Handlers) VerifyEmail(c *gin.Context) {
	var in VerifyEmailInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.VerifyEmail(c.Request.Context(), in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "email verified"})
}

func (h *Handlers) ForgotPassword(c *gin.Context) {
	var in ForgotPasswordInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.svc.ForgotPassword(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handlers) Me(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	u, err := h.svc.GetUser(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":    u.ID.Hex(),
		"email": u.Email,
		"name":  userDisplayName(u),
		"phone": u.Phone,
		"role":  string(u.Role),
	})
}

func (h *Handlers) SeedDemo(c *gin.Context) {
	out, err := h.svc.SeedDemoUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handlers) ResetPassword(c *gin.Context) {
	var in ResetPasswordInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.ResetPassword(c.Request.Context(), in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "password reset successful"})
}
