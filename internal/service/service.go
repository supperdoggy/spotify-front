package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/floyernick/fleep-go"
	structsBack "github.com/supperdoggy/spotify-web-project/spotify-back/shared/structs"
	"github.com/supperdoggy/spotify-web-project/spotify-front/internal/structs"
	globalStructs "github.com/supperdoggy/spotify-web-project/spotify-globalStructs"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

type IService interface {
	GetAllSongs() (*structs.GetAllSongsRespose, error)
	UploadNewSong(req *structs.UploadSongRequest) error
}

type Service struct {
	logger *zap.Logger
}

func NewService(l *zap.Logger) IService {
	return &Service{logger: l}
}

func (s *Service) GetAllSongs() (*structs.GetAllSongsRespose, error) {
	resp, err := http.Get("http://localhost:8080/allsongs")
	if err != nil {
		s.logger.Error("error making request to backend", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("error reading body", zap.Error(err))
		return nil, err
	}

	var unmarshal []globalStructs.Song

	if err := json.Unmarshal(data, &unmarshal); err != nil {
		s.logger.Error("error unmarshalling body", zap.Error(err), zap.String("data", string(data)))
		return nil, err
	}

	if len(unmarshal) == 0 {
		return &structs.GetAllSongsRespose{Songs: unmarshal, Error: "not found"}, nil
	}

	return &structs.GetAllSongsRespose{Songs: unmarshal}, nil
}

func (s *Service) UploadNewSong(req *structs.UploadSongRequest) error {
	if req.SongData == "" || len(req.SongData) == 0 || req.ReleaseDate == "" || req.Album == "" || req.Band == "" || req.Name == "" {
		return errors.New("you need to fill all the fields")
	}

	info, err := fleep.GetInfo([]byte(req.SongData))
	if err != nil {
		return err
	}

	if !info.IsAudio() {
		return errors.New("file must be audio")
	}
	// TODO properly parse time
	//release, err := time.Parse("", req.ReleaseDate)
	//if err != nil {
	//	return err
	//}
	release := time.Now()

	reqToBack := structsBack.CreateNewSongReq{
		SongData: []byte(req.SongData),
		Song:     globalStructs.Song{
			Name: req.Name,
			Band: req.Band,
			Album: req.Album,
			ReleaseDate: release,
		},
	}

	data, err := json.Marshal(reqToBack)
	if err != nil {
		return errors.New("error marshalling req")
	}

	resp, err := http.Post("http://localhost:8080/api/v1/newsong", http.DetectContentType(data), bytes.NewReader(data))
	if err != nil {
		return err
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	s.logger.Info("data", zap.Any("data", string(data)))

	return nil
}
