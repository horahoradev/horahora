package routes

import (
	"context"
	"net/http"

	schedulerproto "github.com/horahoradev/horahora/scheduler/protocol"
	"github.com/labstack/echo/v4"
)

func (r *RouteHandler) getInferenceCategories(c echo.Context) error {
	profile, err := r.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	if !(profile.Rank >= 1) {
		return c.String(http.StatusForbidden, "Insufficient user status")
	}

	categories, err := r.s.GetInferenceCategories(context.Background(), &schedulerproto.Empty{})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, categories)
}

func (r *RouteHandler) addInferenceCategory(c echo.Context) error {
	profile, err := r.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	if !(profile.Rank >= 1) {
		return c.String(http.StatusForbidden, "Insufficient user status")
	}

	category := c.FormValue("category")
	tag := c.FormValue("tag")

	_, err := r.s.AddInferenceCategory(context.Background(), &schedulerproto.InferenceEntry{
		Category: category,
		Tag:      tag,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}
