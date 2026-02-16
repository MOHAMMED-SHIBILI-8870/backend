package controllers

import (
	"backend/config"
	"backend/models"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login page",
	})

}

func ShowDashboard(c *gin.Context) {
    var totalUsers, totalProducts, totalOrders int64
    var totalRevenue float64

    // Count totals
    config.DB.Model(&models.User{}).Count(&totalUsers)
    config.DB.Model(&models.Product{}).Count(&totalProducts)
    config.DB.Model(&models.Order{}).Count(&totalOrders)

    // Total revenue
    config.DB.Model(&models.OrderItem{}).
        Joins("JOIN orders ON orders.id = order_items.order_id").
        Where("orders.status = ?", "delivered").
        Select("COALESCE(SUM(order_items.total_price),0)").Scan(&totalRevenue)

    // Daily revenue chart
    type RevenueChart struct {
        Day     string  `json:"day"`
        Revenue float64 `json:"revenue"`
    }

    var revenueChart []RevenueChart
    config.DB.Model(&models.OrderItem{}).
        Select("to_char(orders.created_at, 'YYYY-MM-DD') as day, COALESCE(SUM(order_items.total_price),0) as revenue").
        Joins("JOIN orders ON orders.id = order_items.order_id").
        Where("orders.status = ?", "delivered").
        Group("day").
        Order("day").
        Scan(&revenueChart)

    // Convert chart data to JSON
    revenueJSON, err := json.Marshal(revenueChart)
    if err != nil {
        revenueJSON = []byte("[]")
    }

    // Render dashboard template
    c.HTML(http.StatusOK, "dashboard.html", gin.H{
        "title":          "Admin Dashboard",
        "total_users":    totalUsers,
        "total_products": totalProducts,
        "total_orders":   totalOrders,
        "total_revenue":  totalRevenue,
        "revenue_json":   string(revenueJSON),
        "Active":         "dashboard",
    })
}



func ShowUsersPage(c *gin.Context) {
	var users []models.User
	if err := config.DB.Order("id ASC").Find(&users).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "users.html", gin.H{"error": "Failed to fetch users"})
		return
	}

	c.HTML(http.StatusOK, "users.html", gin.H{
		"title":  "Manage Users",
		"users":  users,
		"Active": "users",
	})
}
func ShowEditUserPage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.String(http.StatusNotFound, "User not found")
		return
	}

	c.HTML(http.StatusOK, "edit_user.html", gin.H{
		"title": "Edit User",
		"user":  user,
	})
}

// ---------------- PRODUCTS ----------------
func ShowProductsPage(c *gin.Context) {
	var products []models.Product
	if err := config.DB.Order("id ASC").Find(&products).Error; err != nil {
		products = []models.Product{}
	}

	c.HTML(http.StatusOK, "products.html", gin.H{
		"title":    "Manage Products",
		"products": products,
		"Active":   "products",
	})
}

// ---------------- CREATE PRODUCT PAGE ----------------
func ShowCreateProductPage(c *gin.Context) {
	c.HTML(http.StatusOK, "add_product.html", gin.H{
		"title": "Add Product",
	})
}

// ---------------- EDIT PRODUCT PAGE ----------------
func ShowEditProductPage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid product ID")
		return
	}

	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		c.String(http.StatusNotFound, "Product not found")
		return
	}

	c.HTML(http.StatusOK, "edit_product.html", gin.H{
		"title":   "Edit Product",
		"product": product,
	})
}

// ---------------- ORDERS ----------------
func ShowOrdersPage(c *gin.Context) {
	var orders []models.Order
	if err := config.DB.Preload("User").Preload("OrderItems").Order("id ASC").Find(&orders).Error; err != nil {
		orders = []models.Order{}
	}
	c.HTML(http.StatusOK, "orders.html", gin.H{
		"title":  "Manage Orders",
		"orders": orders,
		"Active": "orders",
	})
}

// -------MIDDLEWARE
func MethodOverride() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost {
			if method := c.PostForm("_method"); method != "" {
				c.Request.Method = method
			}
		}
		c.Next()
	}
}

func getUserIDFromContext(c *gin.Context) (uint, bool) {
	userIDValue, exists := c.Get("userId")
	if !exists {
		return 0, false
	}

	switch v := userIDValue.(type) {
	case uint:
		return v, true
	case int:
		return uint(v), true
	case float64:
		return uint(v), true
	case string:
		id, err := strconv.Atoi(v)
		if err != nil {
			return 0, false
		}
		return uint(id), true
	default:
		return 0, false
	}
}

// ---------------- PROFILE PAGE ----------------
func ShowAdminProfilePage(c *gin.Context) {
	adminID, ok := getUserIDFromContext(c)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var admin models.User
	if err := config.DB.First(&admin, adminID).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to load admin details")
		return
	}

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"title":  "Admin Profile",
		"admin":  admin,
		"Active": "profile",
	})
}

// ---------------- EDIT PAGE ----------------
func ShowEditAdminProfilePage(c *gin.Context) {
	adminID, ok := getUserIDFromContext(c)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var admin models.User
	if err := config.DB.First(&admin, adminID).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to load admin data")
		return
	}

	c.HTML(http.StatusOK, "edit_admin_profile.html", gin.H{
		"title":  "Edit Profile",
		"admin":  admin,
		"Active": "profile",
	})
}

// ---------------- UPDATE ----------------
func UpdateAdminProfile(c *gin.Context) {
	adminID, ok := getUserIDFromContext(c)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var admin models.User
	if err := config.DB.First(&admin, adminID).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to fetch admin")
		return
	}

	fullName := strings.TrimSpace(c.PostForm("full_name"))
	email := strings.TrimSpace(c.PostForm("email"))

	if fullName == "" || email == "" {
		c.HTML(http.StatusBadRequest, "edit_admin_profile.html", gin.H{
			"title":  "Edit Profile",
			"admin":  admin,
			"error":  "Fields cannot be empty",
			"Active": "profile",
		})
		return
	}

	admin.FullName = fullName
	admin.Email = email

	if err := config.DB.Save(&admin).Error; err != nil {
		c.String(http.StatusInternalServerError, "Update failed")
		return
	}

	c.Redirect(http.StatusFound, "/view/profile")
}

func ShowRevenuePage(c *gin.Context) {
	var totalRevenue float64

	err := config.DB.
		Model(&models.OrderItem{}).
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.status = ?", "completed").
		Select("COALESCE(SUM(order_items.total_price), 0)").
		Scan(&totalRevenue).Error

	if err != nil {
		totalRevenue = 0
	}

	c.HTML(http.StatusOK, "revenue.html", gin.H{
		"title":        "Revenue",
		"totalRevenue": totalRevenue,
		"Active":       "revenue",
	})
}
