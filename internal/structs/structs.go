package structs

import globalStructs "github.com/supperdoggy/spotify-web-project/spotify-globalStructs"

type GetAllSongsRespose struct {
	Songs []globalStructs.Song `json:"songs"`
	Error string `json:"error"`
}

type UploadSongRequest struct {
	Name        string    `json:"name"`
	Album       string    `json:"album"`
	Band        string    `json:"band"`
	ReleaseDate string `json:"release_date"`
	SongData string `json:"song_data"`
}
