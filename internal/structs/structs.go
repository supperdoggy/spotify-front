package structs

import "github.com/supperdoggy/webproject/globalStructs"

type GetAllSongsRespose struct {
	Songs []globalStructs.Song `json:"songs"`
	Error string `json:"error"`
}
