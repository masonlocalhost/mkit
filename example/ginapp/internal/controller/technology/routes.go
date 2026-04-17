package technology

import "github.com/go-chi/chi/v5"

func (cc *Controller) Register(router chi.Router) {
	router.Get("/{id}", cc.GetTechnology)
	router.Get("/", cc.ListTechnologies)
	router.Get("/vendor", cc.ListTechnologyVendors)
	router.Get("/cpe-type", cc.ListTechnologyCPETypes)
}
