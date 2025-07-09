package handlers

import (
	"code.haedhutner.dev/mvv/LastMUD/services/auth/internal/service"
	"code.haedhutner.dev/mvv/LastMUD/shared/httputils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type GetAccountHandler struct {
	userService *service.UserService
}

func NewGetAccountHandler() *GetAccountHandler {
	return &GetAccountHandler{}
}

func (*GetAccountHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	log.Println("GET: Found AccountId: ", vars["accountId"])

	err := httputils.EncodeBody(resp, service.GetAccountResponse{AccountId: uuid.New().String()})

	if err != nil {
		httputils.WriteUnhandledError(resp, err)
		return
	}

	resp.WriteHeader(http.StatusOK)
}

type CreateAccountHandler struct {
	userService *service.UserService
}

func NewCreateAccountHandler() *CreateAccountHandler {
	return &CreateAccountHandler{}
}

func (*CreateAccountHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	err := httputils.DecodeBody(req.Body, &service.CreateAccountRequest{})

	if err != nil {
		httputils.WriteUnhandledError(resp, err)
		return
	}

	err = httputils.EncodeBody(resp, service.CreateAccountResponse{AccountId: uuid.New().String()})

	if err != nil {
		httputils.WriteUnhandledError(resp, err)
		return
	}

	resp.WriteHeader(http.StatusOK)
}

type UpdateAccountHandler struct {
	userService *service.UserService
}

func NewUpdateAccountHandler() *UpdateAccountHandler {
	return &UpdateAccountHandler{}
}

func (*UpdateAccountHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	logger := httputils.GetLogger(req)
	vars := mux.Vars(req)

	logger.Info("PUT: Found AccountId: ", vars["accountId"])

	body := service.UpdateAccountRequest{}
	err := httputils.DecodeBody(req.Body, &body)

	if err != nil {
		httputils.WriteUnhandledError(resp, err)
		return
	}

	logger.Infof("PUT Body: %+v", body)
	resp.WriteHeader(http.StatusOK)
}

type DeleteAccountHandler struct {
	userService *service.UserService
}

func NewDeleteAccountHandler() *DeleteAccountHandler {
	return &DeleteAccountHandler{}
}

func (*DeleteAccountHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	log.Println("DELETE: Found AccountId: ", vars["accountId"])

	resp.WriteHeader(http.StatusOK)
}
