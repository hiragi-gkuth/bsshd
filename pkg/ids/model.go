package ids

import (
	"net"
	"sync"
	"time"

	"github.com/hiragi-gkuth/bitris-analyzer/pkg/threshold"
)

// ModelData は，検知モデルデータ
type ModelData struct {
	sync.RWMutex
	th           *threshold.Threshold
	mask         int
	entirePeriod time.Duration
	divisions    int
}

// 検知モデル本体，直接扱うのではなくGetModelで取得して使う感じ
var m *ModelData

// InitModel は，検知モデルの初期化を行う
func InitModel(subnetMask int, entirePeriod time.Duration, divisions int) {
	if m != nil {
		return
	}
	m = NewModel(subnetMask, entirePeriod, divisions)

}

// NewModel は空の新しい見地モデルを返す
func NewModel(subnetMask int, entirePeriod time.Duration, divisions int) *ModelData {
	// 初期はハードコーディングしておく
	th := threshold.New(subnetMask, entirePeriod, divisions)
	th.BaseThreshold = 1.0
	return &ModelData{
		th:           th,
		mask:         subnetMask,
		entirePeriod: entirePeriod,
		divisions:    divisions,
	}
}

// GetModel は，検知モデルを取得する
func GetModel() *ModelData {
	if m == nil {
		panic("Ids model isn't initialized yet, call InitModel() first.")
	}
	return m
}

func (rcv *ModelData) Write(th *threshold.Threshold) {
	rcv.Lock()
	rcv.th = th
	rcv.Unlock()
}

// Search は，最適な閾値を返す
func (rcv *ModelData) Search(ip net.IP) float64 {
	now := time.Now()

	rcv.Lock()
	defer rcv.Unlock()
	model := rcv.th
	threshold := model.BaseThreshold

	// 時間帯に対応するしきい値があるなら格納
	t, ok := model.OnTime.Get(now)
	if ok {
		threshold = t
	}

	// IPアドレスに対応するしきい値があるなら格納
	t, ok = model.OnIP.Get(ip)
	if ok {
		threshold = t
	}

	// IPアドレスに対応し，かつ時間帯に対応するしきい値があるなら格納

	if onTime, ok := model.OnIPTime.GetByIP(ip); ok {
		t, ok = onTime.Get(now)
		if ok {
			threshold = t
		}
	}
	return threshold
}
