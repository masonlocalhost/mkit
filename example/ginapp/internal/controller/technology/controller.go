package technology

import (
	"mkit/example/ginapp/internal/controller"
	"mkit/example/ginapp/pkg/converter"
	"mkit/example/ginapp/pkg/dto"
	"mkit/pkg/error/serviceerror"
	ginutil "mkit/pkg/server/gin/util"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	dep *controller.DependencyContainer
}

func NewController(dep *controller.DependencyContainer) *Controller {
	return &Controller{
		dep: dep,
	}
}

func (cc *Controller) GetTechnology(c *gin.Context) {
	var (
		req = &dto.GetTechnologyRequest{}
		ctx = c.Request.Context()
	)
	if err := c.ShouldBindUri(req); err != nil {
		ginutil.HandleError(c, serviceerror.NewInvalidArgument(err))

		return
	}

	tech, err := cc.dep.TechnologyService.FirstByID(ctx, req.ID)
	if err != nil {
		ginutil.HandleError(c, err)

		return
	}

	ginutil.HandleSuccess(c, &dto.GetTechnologyResponse{
		Technology: converter.TechnologyToDTO(tech),
	})
}

func (cc *Controller) ListTechnologies(c *gin.Context) {
	var (
		req = &dto.ListTechnologiesRequest{}
		ctx = c.Request.Context()
		// user = &model.User{ID: "mock user"}
	)
	if err := c.ShouldBind(req); err != nil {
		ginutil.HandleError(c, serviceerror.NewInvalidArgument(err))

		return
	}

	sorts, err := ginutil.GetSortArray(req.SortBy, req.SortOrder, dto.ListTechnologiesOrderColumns)
	if err != nil {
		ginutil.HandleError(c, err)

		return
	}

	techs, total, err := cc.dep.TechnologyService.FindByFilters(
		ctx, req.Vendors, req.CPETypes, req.CreatedFrom, req.CreatedTo, req.Search,
		req.Limit, req.Offset, true, sorts,
	)
	if err != nil {
		ginutil.HandleError(c, err)

		return
	}

	ginutil.HandleSuccess(c, dto.ListTechnologiesResponse(
		converter.TechnologiesToDTOs(techs),
	), total, req.Offset, req.Limit)
}

func (cc *Controller) ListTechnologyVendors(c *gin.Context) {
	var (
		req = &dto.ListTechnologyVendorsRequest{}
		// user = &model.User{ID: "mock user"}
		ctx = c.Request.Context()
	)

	if err := c.ShouldBind(req); err != nil {
		ginutil.HandleError(c, serviceerror.NewInvalidArgument(err))

		return
	}

	r, err := cc.dep.TechnologyService.FindTechnologyVendors(ctx)
	if err != nil {
		ginutil.HandleError(c, err)

		return
	}

	ginutil.HandleSuccess(c, dto.ListTechnologyVendorsResponse(r))
}

func (cc *Controller) ListTechnologyCPETypes(c *gin.Context) {
	var (
		req = &dto.ListTechnologyCPETypesRequest{}
		// user = &model.User{ID: "mock user"}
		ctx = c.Request.Context()
	)

	if err := c.ShouldBind(req); err != nil {
		ginutil.HandleError(c, serviceerror.NewInvalidArgument(err))

		return
	}

	r, err := cc.dep.TechnologyService.FindTechnologyCPETypes(ctx)
	if err != nil {
		ginutil.HandleError(c, err)

		return
	}

	ginutil.HandleSuccess(c, dto.ListTechnologyCPETypesResponse(r))
}
