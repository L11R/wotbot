package domain

import "fmt"

var (
	// Error that could occur during database querying
	ErrInternalDatabase = fmt.Errorf("internal database error")
	// Error that could occure during Wargaming API call
	ErrInternalWargaming = fmt.Errorf("internal Wargaming API error")
	// Error that could occur during XVM stats call
	ErrInternalXVM = fmt.Errorf("internal XVM stats error")
	// Error that could occur during KTTC stats call
	ErrInternalKTTC = fmt.Errorf("internal KTTC stats error")
	// Error that occurs if user passed wrong data on input
	ErrBotBadRequest = fmt.Errorf("bot bad request")
	// Error that occurs if player not found
	ErrPlayerNotFound = fmt.Errorf("player not found")
	// Error that occurs if user not found
	ErrUserNotFound = fmt.Errorf("user not found")
	// Error that occurs if user didn't save nickname
	ErrNicknameNotSaved = fmt.Errorf("nickname not saved")
	// Error that occurs if trend image not found
	ErrTrendImageNotFound = fmt.Errorf("trend image not found")
)
