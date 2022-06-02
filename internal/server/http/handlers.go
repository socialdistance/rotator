package internalhttp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"rotator/internal/app"
)

type ServerHandlers struct {
	app *app.App
}

func NewServerHandlers(a *app.App) *ServerHandlers {
	return &ServerHandlers{app: a}
}

func ResponseError(w http.ResponseWriter, code int, err error) {
	data, err := json.Marshal(ErrorDto{
		Success: false,
		Error:   err.Error(),
	})

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Failed to marshall error dto"))
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func ParsingData(r *http.Request, dto interface{}) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("error read body: %w", err)
	}

	err = json.Unmarshal(data, dto)
	if err != nil {
		return fmt.Errorf("error unmarshall body: %w", err)
	}

	return nil
}

func (s *ServerHandlers) AddBannerToSlot(w http.ResponseWriter, r *http.Request) {
	var dto BannerToSlotDto

	err := ParsingData(r, &dto)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err)
		return
	}

	err = s.app.AddBannerToSlot(r.Context(), dto.BannerID, dto.SlotID)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (s *ServerHandlers) RemoveBannerToSlot(w http.ResponseWriter, r *http.Request) {
	var dto BannerToSlotDto

	err := ParsingData(r, &dto)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err)
		return
	}

	err = s.app.RemoveBannerToSlot(r.Context(), dto.BannerID, dto.SlotID)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (s *ServerHandlers) CountTransition(w http.ResponseWriter, r *http.Request) {
	var dto CountTransitionDto

	err := ParsingData(r, &dto)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err)
		return
	}

	err = s.app.CountTransition(r.Context(), dto.BannerID, dto.SlotID, dto.SocialGroupID)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (s *ServerHandlers) ChooseBanner(w http.ResponseWriter, r *http.Request) {
	var dto ChooseBannerDto

	err := ParsingData(r, &dto)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err)
		return
	}

	bannerID, err := s.app.ChooseBanner(r.Context(), dto.SlotID, dto.SocialGroupID)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err)
		return
	}

	response := map[string]int64{
		"banner_id": bannerID,
	}

	res, err := json.Marshal(response)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
