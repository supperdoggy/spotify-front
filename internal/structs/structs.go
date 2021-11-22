package structs

import globalStructs "github.com/supperdoggy/spotify-web-project/spotify-globalStructs"

type GetAllSongsRespose struct {
	Songs []globalStructs.Song `json:"songs"`
	Error string `json:"error"`
}
