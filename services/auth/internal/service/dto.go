package service

type CreateAccountRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateAccountResponse struct {
	AccountId string `json:"accountId"`
}

type UpdateAccountRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateAccountResponse struct {
	AccountId string `json:"accountId"`
}

type GetAccountResponse struct {
	AccountId string `json:"accountId"`
}
