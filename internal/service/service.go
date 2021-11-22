package service

import (
	"encoding/json"
	"github.com/supperdoggy/spotify-web-project/spotify-front/internal/structs"
	globalStructs "github.com/supperdoggy/spotify-web-project/spotify-globalStructs"

	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

type IService interface {
	GetAllSongs() (*structs.GetAllSongsRespose, error)
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
