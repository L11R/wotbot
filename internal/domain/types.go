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

type Stat struct {
	ID         int        `db:"id"`
	UserID     int        `db:"user_id"`
	Name       string     `db:"name"`
	Value      string     `db:"value"`
	HtmlID     string     `db:"html_id"`
	TrendImage []byte     `db:"trend_img"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
}