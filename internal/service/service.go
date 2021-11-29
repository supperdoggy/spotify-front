package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/floyernick/fleep-go"
	structsBack "github.com/supperdoggy/spotify-web-project/spotify-back/shared/structs"
	"github.com/supperdoggy/spotify-web-project/spotify-front/internal/structs"
	globalStructs "github.com/supperdoggy/spotify-web-project/spotify-globalStructs"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type IService interface {
	GetAllSongs() (resp structs.GetAllSongsRespose, err error)
	UploadNewSong(req *structs.UploadSongRequest) error
}

type Service struct {
	logger *zap.Logger
}

func NewService(l *zap.Logger) IService {
	return &Service{logger: l}
}

func (s *Service) GetAllSongs() (resp structs.GetAllSongsRespose, err error) {
	rawResp, err := http.Get("http://localhost:8080/allsongs")
	if err != nil {
		s.logger.Error("error making request to backend", zap.Error(err))
		resp.Error = err.Error()
		return resp, err
	}
	defer rawResp.Body.Close()

	data, err := ioutil.ReadAll(rawResp.Body)
	if err != nil {
		s.logger.Error("error reading body", zap.Error(err))
		resp.Error = err.Error()
		return resp, err
	}

	var unmarshal structs.GetAllSongsRespose

	if err := json.Unmarshal(data, &unmarshal); err != nil {
		s.logger.Error("error unmarshalling body", zap.Error(err), zap.String("data", string(data)))
		resp.Error = err.Error()
		return resp, err
	}

	if len(unmarshal.Songs) == 0 {
		return unmarshal, nil
	}

	return unmarshal, nil
}

func (s *Service) UploadNewSong(req *structs.UploadSongRequest) error {
	if len(req.SongData) == 0 || req.ReleaseDate == "" || req.Album == "" || req.Band == "" || req.Name == "" {
		return errors.New("you need to fill all the fields")
	}

	b64data := req.SongData[strings.IndexByte(req.SongData, ',')+1:]
	data, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		return err
	}

	info, err := fleep.GetInfo(data)
	if err != nil {
		return err
	}

	if !info.IsAudio() {
		return errors.New("file must be audio")
	}

	release, err := time.Parse("2006-01-02T15:04", req.ReleaseDate)
	if err != nil {
		return err
	}

	reqToBack := structsBack.CreateNewSongReq{
		SongData: data,
		Song: globalStructs.Song{
			Name:        req.Name,
			Band:        req.Band,
			Album:       req.Album,
			ReleaseDate: release,
		},
	}

	data, err = json.Marshal(reqToBack)
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

	return nil
}
