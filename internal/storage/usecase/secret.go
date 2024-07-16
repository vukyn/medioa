package usecase

import (
	"context"
	"fmt"
	"medioa/constants"
	secretModel "medioa/internal/secret/models"
	storageModel "medioa/internal/storage/models"
	"medioa/pkg/log"
	"regexp"

	"github.com/google/uuid"
	"github.com/vukyn/kuery/crypto"
	"golang.org/x/crypto/bcrypt"
)

func (u *usecase) CreateSecret(ctx context.Context, userId int64, params *storageModel.CreateSecretRequest) (*storageModel.CreateSecretResponse, error) {
	log := log.New("service", "CreateSecret")

	// check if username already exists
	foundSecret, err := u.secretSv.GetOne(ctx, &secretModel.RequestParams{
		Username: params.Username,
	})
	if err != nil {
		log.Error("service.secretSv.GetOne", err)
		return nil, err
	}
	if foundSecret != nil {
		return nil, fmt.Errorf("username already has a secret")
	}

	// check if pin code is valid
	pinCodePattern := `^\d{4}$`
	regex, _ := regexp.Compile(pinCodePattern)
	if !regex.MatchString(params.PinCode) {
		return nil, fmt.Errorf("pin code must be 4 digits")
	}

	isMaster := false
	if params.MasterKey == u.cfg.Secret.SecretKey {
		isMaster = true
	}

	// Save to database
	_id := uuid.New().String()
	accessToken := crypto.HashedToken()
	if _, err := u.secretSv.Create(ctx, userId, &secretModel.SaveRequest{
		UUID:        _id,
		Username:    params.Username,
		Password:    params.Password,
		PinCode:     params.PinCode,
		AccessToken: accessToken,
		Type:        constants.SECRET_TYPE_MEDIA,
		IsMaster:    isMaster,
	}); err != nil {
		log.Error("service.secretSv.Create", err)
		return nil, err
	}

	return &storageModel.CreateSecretResponse{
		UserId:      _id,
		AccessToken: accessToken,
	}, nil
}

func (u *usecase) RetrieveSecret(ctx context.Context, userId int64, params *storageModel.RetrieveSecretRequest) (*storageModel.RetrieveSecretResponse, error) {
	log := log.New("service", "RetrieveSecret")

	// check if username exists
	foundSecret, err := u.secretSv.GetOne(ctx, &secretModel.RequestParams{
		Username: params.Username,
	})
	if err != nil {
		log.Error("service.secretSv.GetOne", err)
		return nil, err
	}
	if foundSecret == nil {
		return nil, fmt.Errorf("username not found")
	}

	// check if password is correct
	if err := comparePassword(foundSecret.Password, params.Password); err != nil {
		log.Error("comparePassword", err)
		return nil, fmt.Errorf("password is incorrect")
	}

	accessToken := crypto.HashedToken()
	if _, err := u.secretSv.Update(ctx, userId, &secretModel.SaveRequest{
		UUID:        foundSecret.UUID,
		AccessToken: accessToken,
	}); err != nil {
		log.Error("service.secretSv.Update", err)
		return nil, err
	}

	return &storageModel.RetrieveSecretResponse{
		UserId:      foundSecret.UUID,
		AccessToken: accessToken,
	}, nil
}

func (u *usecase) ResetPinCode(ctx context.Context, userId int64, params *storageModel.ResetPinCodeRequest) (int64, error) {
	log := log.New("service", "ResetPinCode")

	// check if secret exists
	foundSecret, err := u.secretSv.GetOne(ctx, &secretModel.RequestParams{
		AccessToken: params.AccessToken,
	})
	if err != nil {
		log.Error("service.secretSv.GetOne", err)
		return 0, err
	}
	if foundSecret == nil {
		return 0, fmt.Errorf("secret not found")
	}

	if _, err := u.secretSv.Update(ctx, userId, &secretModel.SaveRequest{
		UUID:    foundSecret.UUID,
		PinCode: params.NewPinCode,
	}); err != nil {
		log.Error("service.secretSv.Update", err)
		return 0, err
	}

	return 1, nil
}

func comparePassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
