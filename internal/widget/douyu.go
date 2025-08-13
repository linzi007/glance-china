package widget

import (
	"context"
	"fmt"
	"time"
	
	"github.com/glance-china/internal/service"
)

// DouyuLiveWidget 斗鱼直播组件
type DouyuLiveWidget struct {
	ChineseWidget
	Rooms         []DouyuRoom `yaml:"rooms"`
	Limit         int         `yaml:"limit"`
	ShowOffline   bool        `yaml:"show-offline"`
	CollapseAfter int         `yaml:"collapse-after"`
}

type DouyuRoom struct {
	RoomID string `yaml:"room-id"`
	Name   string `yaml:"name"`
}

func NewDouyuLiveWidget() *DouyuLiveWidget {
	return &DouyuLiveWidget{
		ChineseWidget: ChineseWidget{
			BaseWidget: BaseWidget{
				Type: "douyu-live",
			},
			Region:    "cn",
			APISource: "douyu",
		},
		Limit:         10,
		ShowOffline:   false,
		CollapseAfter: 5,
	}
}

func (d *DouyuLiveWidget) GetData(ctx context.Context, config Config) (interface{}, error) {
	serviceManager := ctx.Value("serviceManager").(*service.ServiceManager)
	
	client, err := serviceManager.GetClient("douyu")
	if err != nil {
		return nil, err
	}
	
	douyuClient := client.(*service.DouyuClient)
	
	var roomIDs []string
	for _, room := range d.Rooms {
		roomIDs = append(roomIDs, room.RoomID)
	}
	
	streams, err := douyuClient.GetLiveStreams(ctx, roomIDs)
	if err != nil {
		return nil, err
	}
	
	// 过滤离线直播间
	if !d.ShowOffline {
		var liveStreams []service.DouyuStreamData
		for _, stream := range streams {
			if stream.IsLive {
				liveStreams = append(liveStreams, stream)
			}
		}
		streams = liveStreams
	}
	
	// 限制数量
	if len(streams) > d.Limit {
		streams = streams[:d.Limit]
	}
	
	return map[string]interface{}{
		"streams":        streams,
		"show_offline":   d.ShowOffline,
		"collapse_after": d.CollapseAfter,
		"title":          d.getTitle(),
	}, nil
}

func (d *DouyuLiveWidget) GetCacheKey(config Config) string {
	return fmt.Sprintf("douyu-live:%d", len(d.Rooms))
}

func (d *DouyuLiveWidget) getTitle() string {
	if d.Title != "" {
		return d.Title
	}
	return "斗鱼直播"
}

func (d *DouyuLiveWidget) Validate(config Config) error {
	if len(d.Rooms) == 0 {
		return fmt.Errorf("至少需要配置一个直播间")
	}
	
	for _, room := range d.Rooms {
		if room.RoomID == "" {
			return fmt.Errorf("直播间ID不能为空")
		}
	}
	
	return nil
}
