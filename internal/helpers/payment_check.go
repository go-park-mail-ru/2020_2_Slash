package helpers

import (
	// nolint: gosec
	"crypto/sha1"
	"encoding/hex"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"strconv"
	"strings"
)

type PaymentData struct {
	notificationType string
	operationID      string
	amount           string
	currency         string
	datetime         string
	sender           string
	codepro          string
	Label            string
	unaccepted       string
	sha1hash         string
}

func NewPaymentDataFromContext(cntx echo.Context) *PaymentData {
	return &PaymentData{
		notificationType: cntx.FormValue("notification_type"),
		operationID:      cntx.FormValue("operation_id"),
		amount:           cntx.FormValue("amount"),
		currency:         cntx.FormValue("currency"),
		datetime:         cntx.FormValue("datetime"),
		sender:           cntx.FormValue("sender"),
		codepro:          cntx.FormValue("codepro"),
		Label:            cntx.FormValue("label"),
		unaccepted:       cntx.FormValue("unaccepted"),
		sha1hash:         cntx.FormValue("sha1_hash"),
	}
}

func (pd *PaymentData) GetParametersString(notificationSecret string) string {
	return strings.Join([]string{pd.notificationType,
		pd.operationID, pd.amount, pd.currency, pd.datetime,
		pd.sender, pd.codepro, notificationSecret, pd.Label}, "&")
}

func (pd *PaymentData) CheckPaymentData() *errors.Error {
	if customErr := pd.CheckHash(); customErr != nil {
		return customErr
	}
	if customErr := pd.CheckUnaccepted(); customErr != nil {
		return customErr
	}
	if customErr := pd.CheckCodepro(); customErr != nil {
		return customErr
	}
	if customErr := pd.CheckLabel(); customErr != nil {
		return customErr
	}
	return nil
}

func getNotificationSecret(filename string) (string, error) {
	// nolint: gosec
	notificationSecretBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(notificationSecretBytes), nil
}

func (pd *PaymentData) CheckHash() *errors.Error {
	notificationSecret, err := getNotificationSecret("secret.key")
	if err != nil {
		return errors.Get(consts.CodeReadKeyFileError)
	}

	parametersString := pd.GetParametersString(notificationSecret)

	// nolint: gosec
	hash := sha1.Sum([]byte(parametersString))
	hashSlice := hash[:]
	hexString := hex.EncodeToString(hashSlice)

	if hexString != pd.sha1hash {
		return errors.Get(consts.CodeWrongPaymentHash)
	}

	return nil
}

func (pd *PaymentData) CheckUnaccepted() *errors.Error {
	unaccepted, err := strconv.ParseBool(pd.unaccepted)
	if err != nil {
		return errors.Get(consts.CodeParseUnacceptedError)
	}
	if unaccepted {
		return errors.Get(consts.CodeUnacceptedPayment)
	}

	return nil
}

func (pd *PaymentData) CheckCodepro() *errors.Error {
	codePro, err := strconv.ParseBool(pd.codepro)
	if err != nil {
		return errors.Get(consts.CodeParseCodeProError)
	}
	if codePro {
		return errors.Get(consts.CodeProtectedPayment)
	}
	return nil
}

func (pd *PaymentData) CheckLabel() *errors.Error {
	if pd.Label == "" {
		return errors.Get(consts.CodeEmptyLabelError)
	}
	return nil
}
