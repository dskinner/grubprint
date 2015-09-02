package keystore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

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
	handleError := func(w http.ResponseWriter, err error) {
		switch err := err.(type) {
		case Error:
			http.Error(w, err.Error(), err.Code())
		default:
			http.Error(w, err.Error(), 500)
		}
	}
	mux.HandleFunc("/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		x := r.FormValue("assertion")
		if err := Verify(x); err != nil {
			handleError(w, err)
			return
		}
		if err := json.NewEncoder(w).Encode(oauth2.Token{AccessToken: x}); err != nil {
			handleError(w, err)
			return
		}
	})
	mux.HandleFunc("/secret", func(w http.ResponseWriter, r *http.Request) {
		if err := VerifyRequest(r); err != nil {
			handleError(w, err)
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
