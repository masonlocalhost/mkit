package technology

import "github.com/gin-gonic/gin"

func (cc *Controller) Register(router gin.IRouter) {
	router.GET("/:id", cc.GetTechnology)
	router.GET("", cc.ListTechnologies)
	router.GET("/vendor", cc.ListTechnologyVendors)
	router.GET("/cpe-type", cc.ListTechnologyCPETypes)
}
