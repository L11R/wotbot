package domain

import "time"

type User struct {
	ID          int        `db:"id"`
	TelegramID  int        `db:"telegram_id"`
	Nickname    *string    `db:"nickname"`
	WargamingID *int       `db:"wargaming_id"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

type StatType string

const (
	TrendStat   StatType = "trend"
	VehicleStat StatType = "vehicle"
)

type Stat struct {
	ID        int        `db:"id"`
	UserID    int        `db:"user_id"`
	Type      StatType   `db:"type"`
	Name      string     `db:"name"`
	Value     *string    `db:"value"`
	HtmlID    string     `db:"html_id"`
	Image     []byte     `db:"img"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
