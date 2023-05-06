package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/insighted4/correios-cep/correios"
	"github.com/insighted4/correios-cep/pkg/errors"
	"github.com/insighted4/correios-cep/storage"
	"github.com/sirupsen/logrus"
)

func listAddressHandler(s storage.Storage, log logrus.FieldLogger) gin.HandlerFunc {
	type PaginationRequest struct {
		PerPage int `json:"per_page" form:"per_page"`
		Page    int `json:"page" form:"page"`
	}

	type ListAddressesRequest struct {
		PaginationRequest
		State string `json:"state" form:"state"`
	}

	const op errors.Op = "handler.handleListAddresses"
	return func(ctx *gin.Context) {
		var form ListAddressesRequest
		if err := ctx.ShouldBind(&form); err != nil {
			abortWithError(ctx, errors.E(op, errors.KindBadRequest, err), nil)
			return
		}

		params := storage.ListParams{
			Pagination: storage.NewPagination(form.PerPage, form.Page),
			State:      form.State,
		}

		result, err := s.ListAddresses(ctx, params)
		if err != nil {
			log.Errorf("failed to list addresses: %v", err)
			abortWithError(ctx, err, nil)
			return
		}

		ctx.JSON(http.StatusOK, result)
	}
}

func getAddressHandler(c correios.Correios, s storage.Storage, log logrus.FieldLogger) gin.HandlerFunc {
	const op errors.Op = "handler.handleListAddresses"

	getAddressFn := func(ctx context.Context, cep string) (*storage.Address, error) {
		addr, err := s.GetAddress(ctx, cep)
		if err == nil {
			return addr, nil
		}

		if errors.Is(err, errors.KindUnexpected) {
			return nil, err
		}

		addr, err = c.Lookup(ctx, cep)
		if err != nil && !errors.Is(err, errors.KindNotFound) {
			return nil, err
		}

		if err := s.CreateAddress(ctx, addr); err != nil {
			return nil, err
		}

		return addr, nil
	}

	return func(ctx *gin.Context) {
		cep := ctx.Param("cep")
		result, err := getAddressFn(ctx, cep)
		switch {
		case err == nil:
			ctx.JSON(http.StatusOK, result)
		case errors.Is(err, errors.KindNotFound):
			log.Infof("address not found: cep %s", cep)
			abortWithError(
				ctx, errors.E(op, errors.KindNotFound, "address not found"),
				fmt.Sprintf("CEP %s not found",
					cep),
			)
		case err != nil:
			log.Errorf("failed to get addresses: %v", err)
			abortWithError(ctx, err, nil)
		}
	}
}

func handleCreateAddress(ctx *gin.Context) {
	return
}

func handleUpdateAddress(ctx *gin.Context) {
	return
}
