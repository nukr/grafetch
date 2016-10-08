package grafetch

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// GraphQLQuery ...
type GraphQLQuery struct {
	Query     string `json:"query"`
	Variables string `json:"variables"`
}

// Grafetch ...
type Grafetch struct {
	Queries []GraphQLQuery
	URL     string
	Headers []map[string]string
}

// NewReader ...
func (g Grafetch) newReader(i int) io.Reader {
	bs, _ := json.Marshal(g.Queries[i])
	return strings.NewReader(string(bs))
}

// SetHeader ...
func (g *Grafetch) SetHeader(key, value string) {
	g.Headers = append(g.Headers, map[string]string{key: value})
}

// Fetch ...
func (g Grafetch) Fetch(data interface{}) error {
	for i := range g.Queries {
		client := &http.Client{}
		req, err := http.NewRequest("POST", g.URL, g.newReader(i))
		if err != nil {
			return err
		}
		for _, header := range g.Headers {
			for k, v := range header {
				req.Header.Add(k, v)
			}
		}
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if resp.StatusCode >= 400 {
			return errors.New(resp.Status)
		}
		json.Unmarshal(body, data)
	}
	return nil
}

// SetQuery ...
func (g *Grafetch) SetQuery(q GraphQLQuery) {
	g.Queries = append(g.Queries, q)
}

// New ...
func New(url string) Grafetch {
	return Grafetch{
		URL: url,
	}
}
