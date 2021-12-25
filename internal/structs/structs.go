package structs

import globalStructs "github.com/supperdoggy/spotify-web-project/spotify-globalStructs"

type GetAllSongsRespose struct {
	Songs []globalStructs.Song `json:"songs"`
	Error string               `json:"error"`
}

type UploadSongRequest struct {
	Name        string `json:"name"`
	Album       string `json:"album"`
	Band        string `json:"band"`
	ReleaseDate string `json:"release_date"`
	SongData    string `json:"song_data"`
}

type AuthReq struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Login     bool   `json:"login"`
}

type AuthResp struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
	Error  string `json:"error"`
}
