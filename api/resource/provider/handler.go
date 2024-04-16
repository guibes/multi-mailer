package provider

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	e "myapp/api/resource/common/err"
	l "myapp/api/resource/common/log"
	ctxUtil "myapp/util/ctx"
	validatorUtil "myapp/util/validator"
)

type API struct {
	logger     *zerolog.Logger
	validator  *validator.Validate
	repository *Repository
}

func New(logger *zerolog.Logger, validator *validator.Validate, db *gorm.DB) *API {
	return &API{
		logger:     logger,
		validator:  validator,
		repository: NewRepository(db),
	}
}

// List godoc
//
//	@summary		List providers
//	@description	List providers
//	@tags			providers
//	@accept			json
//	@produce		json
//	@success		200	{array}		DTO
//	@failure		500	{object}	err.Error
//	@router			/providers [get]
func (a *API) List(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())

	providers, err := a.repository.List()
	if err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	}

	if len(providers) == 0 {
		fmt.Fprint(w, "[]")
		return
	}

	if err := json.NewEncoder(w).Encode(providers.ToDto()); err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}
}

// Create godoc
//
//	@summary		Create provider
//	@description	Create provider
//	@tags			providers
//	@accept			json
//	@produce		json
//	@param			body	body	Form	true	"Provider form"
//	@success		201
//	@failure		400	{object}	err.Error
//	@failure		422	{object}	err.Errors
//	@failure		500	{object}	err.Error
//	@router			/providers [post]
func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())

	form := &Form{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.BadRequest(w, e.RespJSONDecodeFailure)
		return
	}

	if err := a.validator.Struct(form); err != nil {
		respBody, err := json.Marshal(validatorUtil.ToErrResponse(err))
		if err != nil {
			a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
			e.ServerError(w, e.RespJSONEncodeFailure)
			return
		}

		e.ValidationErrors(w, respBody)
		return
	}

	newProvider := form.ToModel()
	newProvider.ID = uuid.New()

	book, err := a.repository.Create(newProvider)
	if err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataInsertFailure)
		return
	}

	a.logger.Info().Str(l.KeyReqID, reqID).Str("id", book.ID.String()).Msg("new provider created")
	w.WriteHeader(http.StatusCreated)
}

// Read godoc
//
//	@summary		Read provider
//	@description	Read provider
//	@tags			providers
//	@accept			json
//	@produce		json
//	@param			id	path		string	true	"Provider ID"
//	@success		200	{object}	DTO
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		500	{object}	err.Error
//	@router			/providers/{id} [get]
func (a *API) Read(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		e.BadRequest(w, e.RespInvalidURLParamID)
		return
	}

	book, err := a.repository.Read(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	}

	dto := book.ToDto()
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}
}

// Update godoc
//
//	@summary		Update provider
//	@description	Update provider
//	@tags			providers
//	@accept			json
//	@produce		json
//	@param			id		path	string	true	"Provider ID"
//	@param			body	body	Form	true	"Provider form"
//	@success		200
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		422	{object}	err.Errors
//	@failure		500	{object}	err.Error
//	@router			/providers/{id} [put]
func (a *API) Update(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		e.BadRequest(w, e.RespInvalidURLParamID)
		return
	}

	form := &Form{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.BadRequest(w, e.RespJSONDecodeFailure)
		return
	}

	if err := a.validator.Struct(form); err != nil {
		respBody, err := json.Marshal(validatorUtil.ToErrResponse(err))
		if err != nil {
			a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
			e.ServerError(w, e.RespJSONEncodeFailure)
			return
		}

		e.ValidationErrors(w, respBody)
		return
	}

	provider := form.ToModel()
	provider.ID = id

	rows, err := a.repository.Update(provider)
	if err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataUpdateFailure)
		return
	}
	if rows == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	a.logger.Info().Str(l.KeyReqID, reqID).Str("id", id.String()).Msg("book updated")
}

// Delete godoc
//
//	@summary		Delete provider
//	@description	Delete provider
//	@tags			providers
//	@accept			json
//	@produce		json
//	@param			id	path	string	true	"Provider ID"
//	@success		200
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		500	{object}	err.Error
//	@router			/providers/{id} [delete]
func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		e.BadRequest(w, e.RespInvalidURLParamID)
		return
	}

	rows, err := a.repository.Delete(id)
	if err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataRemoveFailure)
		return
	}
	if rows == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	a.logger.Info().Str(l.KeyReqID, reqID).Str("id", id.String()).Msg("provider deleted")
}
