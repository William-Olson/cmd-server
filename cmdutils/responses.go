package cmdutils

// SlackResponse is the response returned as JSON
type SlackResponse struct {
	ResponseType string            `json:"response_type"`
	Text         string            `json:"text"`
	Attachments  []SlackAttachment `json:"attachments"`
}

// SlackAttachment is optionally attached to SlackResponses
type SlackAttachment struct {
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	Text      string `json:"text"`
}
