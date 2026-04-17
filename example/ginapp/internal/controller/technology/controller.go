package technology

import (
	"mkit/example/ginapp/internal/controller"
	"mkit/example/ginapp/pkg/converter"
	"mkit/example/ginapp/pkg/dto"
	"mkit/pkg/error/serviceerror"
	chiutil "mkit/pkg/server/chi/util"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Controller struct {
	dep *controller.DependencyContainer
}

func NewController(dep *controller.DependencyContainer) *Controller {
	return &Controller{
		dep: dep,
	}
}

func (cc *Controller) GetTechnology(w http.ResponseWriter, r *http.Request) {
	var (
		req = &dto.GetTechnologyRequest{}
		ctx = r.Context()
	)

	req.ID = chi.URLParam(r, "id")
	if req.ID == "" {
		chiutil.HandleError(w, r, serviceerror.NewInvalidArgument(nil).SetMessage("id is required"))
		return
	}

	tech, err := cc.dep.TechnologyService.FirstByID(ctx, req.ID)
	if err != nil {
		chiutil.HandleError(w, r, err)
		return
	}

	chiutil.HandleSuccess(w, r, &dto.GetTechnologyResponse{
		Technology: converter.TechnologyToDTO(tech),
	})
}

func (cc *Controller) ListTechnologies(w http.ResponseWriter, r *http.Request) {
	var (
		req = &dto.ListTechnologiesRequest{}
		ctx = r.Context()
	)

	if err := chiutil.BindQuery(r, req); err != nil {
		chiutil.HandleError(w, r, serviceerror.NewInvalidArgument(err))
		return
	}

	sorts, err := chiutil.GetSortArray(req.SortBy, req.SortOrder, dto.ListTechnologiesOrderColumns)
	if err != nil {
		chiutil.HandleError(w, r, err)
		return
	}

	techs, total, err := cc.dep.TechnologyService.FindByFilters(
		ctx, req.Vendors, req.CPETypes, req.CreatedFrom, req.CreatedTo, req.Search,
		req.Limit, req.Offset, true, sorts,
	)
	if err != nil {
		chiutil.HandleError(w, r, err)
		return
	}

	chiutil.HandleSuccess(w, r, dto.ListTechnologiesResponse(
		converter.TechnologiesToDTOs(techs),
	), total, req.Offset, req.Limit)
}

func (cc *Controller) ListTechnologyVendors(w http.ResponseWriter, r *http.Request) {
	var (
		req = &dto.ListTechnologyVendorsRequest{}
		ctx = r.Context()
	)

	if err := chiutil.BindQuery(r, req); err != nil {
		chiutil.HandleError(w, r, serviceerror.NewInvalidArgument(err))
		return
	}

	result, err := cc.dep.TechnologyService.FindTechnologyVendors(ctx)
	if err != nil {
		chiutil.HandleError(w, r, err)
		return
	}

	chiutil.HandleSuccess(w, r, dto.ListTechnologyVendorsResponse(result))
}

func (cc *Controller) ListTechnologyCPETypes(w http.ResponseWriter, r *http.Request) {
	var (
		req = &dto.ListTechnologyCPETypesRequest{}
		ctx = r.Context()
	)

	if err := chiutil.BindQuery(r, req); err != nil {
		chiutil.HandleError(w, r, serviceerror.NewInvalidArgument(err))
		return
	}

	result, err := cc.dep.TechnologyService.FindTechnologyCPETypes(ctx)
	if err != nil {
		chiutil.HandleError(w, r, err)
		return
	}

	chiutil.HandleSuccess(w, r, dto.ListTechnologyCPETypesResponse(result))
}
