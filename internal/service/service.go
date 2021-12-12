package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/floyernick/fleep-go"
	"github.com/gin-gonic/gin"
	structs2 "github.com/supperdoggy/spotify-web-project/spotify-auth/shared/structs"
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
	Auth(req structs.AuthReq) (resp structs.AuthResp, err error)
	CheckToken(c *gin.Context) bool
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

func (s *Service) Auth(req structs.AuthReq) (resp structs.AuthResp, err error) {
	if req.Email == "" || req.Password == "" {
		resp.Error = "you need to fill all the fields"
		return resp, errors.New(resp.Error)
	}

	if req.Login == false && (req.FirstName == "" || req.LastName == "") {
		resp.Error = "you need to fill all the fields"
		return resp, errors.New(resp.Error)
	}
	var url = "http://localhost:8080/login"
	if !req.Login {
		url = "http://localhost:8080/register"
	}
	marshalled, err := json.Marshal(req)
	if err != nil {
		s.logger.Error("error marshalling req", zap.Error(err))
		resp.Error = err.Error()
		return
	}

	respData, err := http.Post(url, "application/json", bytes.NewBuffer(marshalled))
	if err != nil {
		s.logger.Error("error making request to backend", zap.Error(err))
		resp.Error = err.Error()
		return
	}
	defer respData.Body.Close()

	data, err := ioutil.ReadAll(respData.Body)
	if err != nil {
		s.logger.Error("error reading body", zap.Error(err))
		resp.Error = err.Error()
		return
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		s.logger.Error("error unmarshalling data", zap.Error(err))
		resp.Error = err.Error()
		return resp, err
	}

	if resp.Error != "" {
		s.logger.Error("got error from backend", zap.Any("error", resp.Error))
		return resp, errors.New(resp.Error)
	}

	return
}

func (s *Service) CheckToken(c *gin.Context) bool {
	token, err := c.Cookie("t")
	if err != nil {
		c.Redirect(http.StatusPermanentRedirect, "http://localhost:8081/auth")
		return false
	}

	data, err := s.GetToken(token)
	if err != nil {
		s.logger.Error("error when getting token", zap.Error(err))
		c.Redirect(http.StatusPermanentRedirect, "http://localhost:8081/auth")
		return false
	}

	if !data.OK {
		c.Redirect(http.StatusPermanentRedirect, "http://localhost:8081/auth")
		return false
	}
	return true
}

func (s *Service) GetToken(token string) (resp structs2.CheckTokenResp, err error){
	req := structs2.CheckTokenReq{Token: token}
	data, err := json.Marshal(req)
	if err != nil {
		s.logger.Error("error marshalling token", zap.Error(err))
		return resp, err
	}

	respRow, err := http.Post("http://localhost:8083/api/v1/check_token", "application/json", bytes.NewBuffer(data))
	if err != nil {
		s.logger.Error("error making request to auth")
		return
	}
	defer respRow.Body.Close()

	data, err = ioutil.ReadAll(respRow.Body)
	if err != nil {
		s.logger.Error("error reading body", zap.Error(err))
		return
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		s.logger.Error("error unmarshalling data", zap.Error(err))
		return
	}

	if resp.Error != "" {
		s.logger.Error("got error from auth", zap.Any("error", resp.Error))
		return
	}
	return
}
