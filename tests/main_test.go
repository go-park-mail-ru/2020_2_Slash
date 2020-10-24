package tests

import (
	"github.com/go-park-mail-ru/2020_2_Slash/app/user"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

type TestCase struct {
	name    string
	reqBody map[string]interface{}
	resBody interface{}
	status  int
	user    *user.User
}

var url = "/api/v1"

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}
