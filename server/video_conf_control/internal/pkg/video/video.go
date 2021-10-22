package video

import (
	"context"
	logs "github.com/sirupsen/logrus"
	"net/http"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/config"
	"github.com/volcengine/volc-sdk-golang/service/vod"
	"github.com/volcengine/volc-sdk-golang/service/vod/models/request"
)

// only choose wanted url, not all param
type response struct {
	Result result `json:"result"`
}

type result struct {
	Data data `json:"Data"`
}

type data struct {
	Status       int        `json:"Status"`
	VideoID      string     `json:"VideoID"`
	PlayInfoList []playInfo `json:"PlayInfoList"`
}

type playInfo struct {
	MainPlayUrl string `json:"MainPlayUrl"`
}

var vodClient *vod.Vod

func Init() {
	vodClient = vod.NewInstance()
	vodClient.SetAccessKey(config.Config.PostProcessingAK)
	vodClient.SetSecretKey(config.Config.PostProcessingSK)
}

func GetVideoURL(ctx context.Context, vids []string) map[string]string {
	res := make(map[string]string)
	for _, vid := range vids {
		query := &request.VodGetPlayInfoRequest{
			Vid:        vid,
			Format:     "",
			Codec:      "",
			Definition: "",
			FileType:   "",
			LogoType:   "",
			Ssl:        "",
		}
		// 发起请求并获取响应
		resp, code, err := vodClient.GetPlayInfo(query)
		if err != nil || code != http.StatusOK {
			logs.Errorf("get paly info failed,vid:%s,error:%s", vid, err)
			continue
		}
		if len(resp.GetResult().PlayInfoList) == 0 {
			res[vid] = ""
		} else {
			res[vid] = resp.GetResult().PlayInfoList[0].MainPlayUrl
		}
	}
	return res
}

func DeleteRecord(ctx context.Context, vid string) (bool, error) {
	query := &request.VodDeleteMediaRequest{
		Vids: vid,
	}
	_, code, err := vodClient.DeleteMedia(query)
	if err != nil || code != http.StatusOK {
		return false, err
	}
	return true, nil
}
