package telegram

import (
	"book-to-mail-bot/lib/e"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strconv"
)

const (
	getUpdateMessage  = "getUpdates"
	sendMessageMethod = "sendMessage"
	getFile           = "getFile"
	getFileMethod     = "file"
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

	data, err := c.doRequest(c.makePath(getUpdateMessage), q)
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

	_, err := c.doRequest(c.makePath(sendMessageMethod), q)
	if err != nil {
		return e.WrapErr("can't send message: %w", err)
	}

	return nil
}

func (c *Client) GetFileLink(fileID string) (file *FilePath, err error) {
	defer func() { err = e.WrapIfErr("can't get url to download file: %w", err) }()
	q := url.Values{}
	q.Add("file_id", fileID)

	data, err := c.doRequest(c.makePath(getFile), q)
	if err != nil {
		return nil, err
	}

	var res FileResponse

	if err = json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return &res.Result, nil
}

func (c *Client) DownloadFile(path string) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can't download file: %w", err) }()

	q := url.Values{}

	filePath := filepath.Join(getFileMethod, c.basePath, path)

	res, err := c.doRequest(filePath, q)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) doRequest(path string, query url.Values) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can't do request: %w", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path,
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

func (c *Client) makePath(method string) string {
	return path.Join(c.basePath, method)
}
