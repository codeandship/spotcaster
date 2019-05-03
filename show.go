package spotcaster

import "fmt"

type MetaData struct {
	FeedURL       string `json:"feedUrl"`
	TotalEpisodes int    `json:"totalEpisodes"`
	Starts        int    `json:"starts"`
	Streams       int    `json:"streams"`
	Listeners     int    `json:"listeners"`
	Followers     int    `json:"followers"`
	Published     bool   `json:"published"`
}

var (
	metaTemplate = `Spotify Stats:
Starts:    %d
Streams:   %d
Listeners: %d
Followers: %d`
	metaMarkdownTemplate = "*Spotify Stats*\n```\nStarts:    %d\nStreams:   %d\nListeners: %d\nFollowers: %d```"
)

func (m MetaData) String() string {
	return fmt.Sprintf(metaTemplate, m.Starts, m.Streams, m.Listeners, m.Followers)
}

func (m MetaData) Markdown() string {
	return fmt.Sprintf(metaMarkdownTemplate, m.Starts, m.Streams, m.Listeners, m.Followers)
}
