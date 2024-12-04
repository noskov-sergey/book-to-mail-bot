package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"go.uber.org/zap"

	"github.com/noskov-sergey/book-to-mail-bot/lib/e"
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

	log *zap.Logger
}

func New(host string, token string, log *zap.Logger) *Client {
	return &Client{
		host:     host,
		basePath: NewBasePath(token),
		client:   http.Client{},
		log:      log.Named("telegram client"),
	}
}

func (c *Client) Updates(offset int, limit int) (updates []Update, err error) {
	defer func() { err = e.WrapIfErr("can't do update: %w", err) }()

	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	body, err := c.doRequest(c.makePath(getUpdateMessage), q)
	if err != nil {
		return nil, err
	}

	defer func() { _ = body.Close() }()

	data, err := io.ReadAll(body)
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

	body, err := c.doRequest(c.makePath(getFile), q)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var res FileResponse

	if err = json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return &res.Result, nil
}

func (c *Client) DownloadFile(p string) (data io.ReadCloser, err error) {
	filePath := path.Join(getFileMethod, c.basePath, p)

	body, err := c.doRequest(filePath, url.Values{})
	if err != nil {
		return nil, e.WrapErr("can't download file: %w", err)
	}

	c.log.Info("document has been downloaded", zap.String("document", p))

	return body, nil
}

func (c *Client) doRequest(path string, query url.Values) (data io.ReadCloser, err error) {
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

	return resp.Body, err
}

func NewBasePath(token string) string {
	return "bot" + token
}

func (c *Client) makePath(method string) string {
	return path.Join(c.basePath, method)
}
