package telegram

type UpdateResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

type IncomingMessage struct {
	Text     string `json:"text"`
	Document *File  `json:"document"`
	From     From   `json:"from"`
	Chat     Chat   `json:"chat"`
}

type File struct {
	ID   string `json:"file_id"`
	Name string `json:"file_name"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}

type FileResponse struct {
	Ok     bool     `json:"ok"`
	Result FilePath `json:"result"`
}

type FilePath struct {
	Path string `json:"file_path"`
}
