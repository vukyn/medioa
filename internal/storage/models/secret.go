package models

type CreateSecretRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	PinCode   string `json:"pin_code"`
	MasterKey string `json:"master_key"`
}

type CreateSecretResponse struct {
	UserId      string `json:"user_id"`
	AccessToken string `json:"access_token"`
}

type RetrieveSecretRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RetrieveSecretResponse struct {
	UserId    string `json:"user_id"`
	AccessKey string `json:"access_key"`
}

type ResetPinCodeRequest struct {
	AccessToken string `json:"access_token"`
	NewPinCode  string `json:"new_pin_code"`
}
