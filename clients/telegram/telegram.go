package telegram

import (
	"book-to-mail-bot/lib/e"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	getUpdateMessage  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: NewBasePath(token),
		client:   http.Client{},
	}
}

func (c *Client) Updates(offset int, limit int) (updates []Update, err error) {
	defer func() { err = e.WrapIfErr("can't do update: %w", err) }()

	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdateMessage, q)
	if err != nil {
		return nil, err
	}

	var res UpdateResponse

	if err = json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessage(ChatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(ChatID))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.WrapErr("can't send message: %w", err)
	}

	return nil
}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can't do request: %w", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}

func NewBasePath(token string) string {
	return "bot" + token
}
