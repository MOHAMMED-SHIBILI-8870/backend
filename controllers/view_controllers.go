package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK,"login.html",gin.H{
		"title":"Login page",
	})

}


func ShowDashboard(c *gin.Context){
	
	c.HTML(http.StatusOK,"dashboard.html",gin.H{
		"title":"Admin Dashboard",
	})
}

