package kttc

import "encoding/json"

type Response struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type StatsByBattles map[string]*Stats

type Stats struct {
	WN8                 float64 `json:"WN8"`
	WTR                 int     `json:"WG"`
	Battles             int     `json:"BT"`
	Wins                int     `json:"BW"`
	Losses              int     `json:"BL"`
	Ties                int     `json:"BD"`
	Winrate             float64 `json:"PW"`
	AverageLevel        float64 `json:"LVL"`
	AverageBattlesLevel float64 `json:"LVLB"`
	Damaged             float64 `json:"DMG"`
	Defended            float64 `json:"TNK"`
	Exp                 int     `json:"EAV"`
	Spotting            float64 `json:"SPT"`
	BaseCaptured        float64 `json:"CPT"`
	BaseDefended        float64 `json:"DEF"`
	HitsPercentage      float64 `json:"HTP,string"`
	Survived            float64 `json:"LIV"`
	KD                  float64 `json:"KDES"`
	Max                 string  `json:"MAX"`
	Date                string  `json:"DATE"`
	FullDate            string  `json:"FULLDATE"`
	Deltas              *Deltas `json:"DELTA"`
}

type Deltas struct {
	WN8            Delta `json:"WN8"`
	WTR            Delta `json:"WG"`
	Winrate        Delta `json:"PW"`
	Damaged        Delta `json:"DMG"`
	Defended       Delta `json:"TNK"`
	Exp            Delta `json:"EAV"`
	Spotted        Delta `json:"SPT"`
	Destroyed      Delta `json:"DST"`
	BaseCaptured   Delta `json:"CPT"`
	BaseDefended   Delta `json:"DEF"`
	Survived       Delta `json:"LIV"`
	KD             Delta `json:"KDES"`
	HitsPercentage Delta `json:"HTP"`
}

type Delta struct {
	Value float64 `json:"value"`
	Diff  string  `json:"diff"`
}
