# with
SPANDigital's implementation of the Functional Options Pattern using Go Generics

### Usage

Please see the [samples/](samples/) directory for examples of usage.

### TLDR 

- use generic function type [/with.go#L7](`with.Func\[O any\]`) as the return type of your `With..` high order functions.
- 


See you the "With.." functions are built below, and with.Build in your constructor.

```go
package catfacts

import (
	"encoding/json"
	"github.com/spandigital/with"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	DEFAULT_CAT_API_BASE             = "https://cat-fact.herokuapp.com"
	DEFAULT_CAT_API_RANDOM_FACT_PATH = "/facts/random"
	DEFAULT_CAT_API_LIMIT            = 1
)

type CatFacts interface {
	GetRandomFacts(limit *int) (catFacts []*CatFact, err error)
}

type CatFact struct {
	Id        string    `json:"_id"`
	V         int       `json:"__v"`
	Text      string    `json:"text"`
	UpdatedAt time.Time `json:"updatedAt"`
	Deleted   bool      `json:"deleted"`
	Source    string    `json:"source"`
	SentCount int       `json:"sentCount"`
}

type Options struct {
	HttpClient     *http.Client
	Base           string
	RandomFactPath string
	Limit          int
}

func WithBase(endpoint string) with.Func[Options] {
	return func(options *Options) (err error) {
		options.Base = endpoint
		return
	}
}

func WithHttpClient(httpClient *http.Client) with.Func[Options] {
	return func(options *Options) (err error) {
		options.HttpClient = httpClient
		return
	}
}

func WithResourceElement(path string) with.Func[Options] {
	return func(options *Options) (err error) {
		options.RandomFactPath = path
		return
	}
}

func WithLimit(limit int) with.Func[Options] {
	return func(options *Options) (err error) {
		options.Limit = limit
		return
	}
}

type catFacts struct {
	httpClient    *http.Client
	randomFactUrl string
	defaultLimit  int
}

func NewCatFacts(withOptions ...with.Func[Options]) (newCatFacts *catFacts, err error) {
	var options *Options
	if options, err = with.Build(&Options{
		HttpClient:     http.DefaultClient,
		Base:           DEFAULT_CAT_API_BASE,
		RandomFactPath: DEFAULT_CAT_API_RANDOM_FACT_PATH,
		Limit:          DEFAULT_CAT_API_LIMIT,
	}, nil, withOptions...); err == nil {
		var randomFactUrl string
		if randomFactUrl, err = url.JoinPath(options.Base, options.RandomFactPath); err == nil {
			newCatFacts = &catFacts{
				httpClient:    options.HttpClient,
				randomFactUrl: randomFactUrl,
				defaultLimit:  options.Limit,
			}
		}
	}
	return
}

func (c *catFacts) GetRandomFacts(limit *int) (catFacts []*CatFact, err error) {
	limitParam := c.defaultLimit
	if limit != nil {
		limitParam = *limit
	}
	var response *http.Response
	if getURL, err := url.Parse(c.randomFactUrl); err == nil {
		values := getURL.Query()
		values.Add("amount", strconv.Itoa(limitParam))
		getURL.RawQuery = values.Encode()
		var decodeTarget any
		//if limit is one then it returns one cat object, otherwise it sends a list of cat abjects
		if limitParam == 1 {
			catFacts = []*CatFact{&CatFact{}}
			decodeTarget = catFacts[0]
		} else {
			decodeTarget = &catFacts
		}
		if response, err = c.httpClient.Get(getURL.String()); err == nil {
			defer response.Body.Close()
			err = json.NewDecoder(response.Body).Decode(decodeTarget)
		}
	}
	return
}
```