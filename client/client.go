package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jwt"
	"grubprint.io/router"
	"grubprint.io/usda"
)

var apiRouter = router.New()

type Client struct {
	client *http.Client

	BaseURL *url.URL

	Foods     usda.FoodService
	Weights   usda.WeightService
	Nutrients usda.NutrientService
}

func New(client *http.Client) *Client {
	if client == nil {
		bin, err := ioutil.ReadFile("id_rsa")
		if err != nil {
			log.Println(err)
			bin = []byte{}
		}
		conf := &jwt.Config{
			Email:      "oauth2@keystore",
			PrivateKey: bin,
			Scopes:     []string{},
			TokenURL:   "http://localhost:8080/oauth2/token",
		}
		client = conf.Client(oauth2.NoContext)
		client.Transport = MemoryCacheTransport(client.Transport)
		// TODO proxy with &httpcontrol.Transport{MaxTries: 3})
	}

	c := &Client{}
	c.client = client
	c.BaseURL = &url.URL{Scheme: "http", Host: "localhost:8080", Path: "/api/"}
	c.Foods = &foodClient{c}
	c.Weights = &weightClient{c}
	c.Nutrients = &nutrientClient{c}

	return c
}

func (c *Client) url(name string, pairs ...string) (*url.URL, error) {
	r := apiRouter.Get(name)
	if r == nil {
		return nil, fmt.Errorf("No route named %q.", name)
	}

	url, err := r.URL(pairs...)
	if err != nil {
		return nil, err
	}

	// make path relative to resolve against baseURL
	url.Path = strings.TrimPrefix(url.Path, "/")
	return c.BaseURL.ResolveReference(url), nil
}

type foodClient struct {
	*Client
}

func (c *foodClient) ById(id string) (*usda.Food, error) {
	ctx := context.TODO()

	url, err := c.url(router.Food, "id", id)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := ctxhttp.Do(ctx, c.client, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s", b)
	}

	var m *usda.Food

	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}

	return m, nil
}

func (c *foodClient) Search(x string) ([]*usda.Food, error) {
	ctx := context.TODO()

	url, err := c.url(router.Foods, "q", x)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := ctxhttp.Do(ctx, c.client, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s", b)
	}

	var m []*usda.Food

	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}

	return m, nil
}

type weightClient struct {
	*Client
}

func (c *weightClient) ByFoodId(id string) ([]*usda.Weight, error) {
	ctx := context.TODO()

	url, err := c.url(router.Weights, "id", id)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := ctxhttp.Do(ctx, c.client, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s", b)
	}

	var m []*usda.Weight

	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}

	return m, nil
}

type nutrientClient struct {
	*Client
}

func (c *nutrientClient) ByFoodId(id string) ([]*usda.Nutrient, error) {
	ctx := context.TODO()

	url, err := c.url(router.Nutrients, "id", id)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := ctxhttp.Do(ctx, c.client, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s", b)
	}

	var m []*usda.Nutrient

	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}

	return m, nil
}
