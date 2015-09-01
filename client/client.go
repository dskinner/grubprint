package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"dasa.cc/food/router"
	"dasa.cc/food/usda"
	"github.com/facebookgo/httpcontrol"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
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
		client = &http.Client{
			Transport: MemoryCacheTransport(&httpcontrol.Transport{MaxTries: 3}),
		}
	}

	c := &Client{}
	c.client = client
	c.BaseURL = &url.URL{Scheme: "http", Host: "localhost:8080", Path: "/api/"}
	c.Foods = &foodClient{c}

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
