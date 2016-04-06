package keystore

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jwt"
)

const message = "despiteallobjections return"

var (
	urltoken, urlsecret string

	conf *jwt.Config
)

func TestMain(m *testing.M) {
	pub, priv, err := Keygen()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/oauth2/token", TokenHandler)
	mux.HandleFunc("/secret", func(w http.ResponseWriter, r *http.Request) {
		if err := VerifyRequest(r); err != nil {
			switch err := err.(type) {
			case Error:
				HandleError(w, err)
			default:
				http.Error(w, err.Error(), 500)
			}
			return
		}
		fmt.Fprint(w, message)
	})

	ts := httptest.NewServer(mux)
	urltoken = ts.URL + "/oauth2/token"
	urlsecret = ts.URL + "/secret"

	Set("keystore@test", pub)
	conf = &jwt.Config{
		Email:      "keystore@test",
		PrivateKey: priv,
		Scopes:     []string{},
		TokenURL:   urltoken,
	}

	ret := m.Run()
	ts.Close()
	os.Exit(ret)
}

func TestVerified(t *testing.T) {
	resp, err := conf.Client(oauth2.NoContext).Get(urlsecret)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	bin, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 || string(bin) != message {
		t.Fatalf("Bad response: %v %s", resp.StatusCode, bin)
	}
}

func TestNotVerified(t *testing.T) {
	resp, err := http.Get(urlsecret)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	bin, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode == 200 || string(bin) == message {
		t.Fatalf("Bad response: %v %s", resp.StatusCode, bin)
	}
}

func TestInvalidToken(t *testing.T) {
	var err error
	var resp *http.Response
	resp, err = conf.Client(oauth2.NoContext).Get(urlsecret)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatal("Expected 200")
	}

	now = func() time.Time { return time.Now().Add(2 * time.Hour) }
	defer func() { now = time.Now }()

	resp, err = conf.Client(oauth2.NoContext).Get(urlsecret)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 401 {
		t.Fatal("Expected 401")
	}
}
