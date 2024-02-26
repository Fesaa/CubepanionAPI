package models

type LeaderboardSubmission struct {
	Uuid          string           `json:"uuid"`
	UnixTimeStamp uint64           `json:"unix_time_stamp"`
	Game          string           `json:"game"`
	Entries       []LeaderboardRow `json:"entries"`
}
