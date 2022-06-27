package msg

type PrivateMsg struct {
	UserId     int64  `json:"user_id"`
	GroupId    int64  `json:"group_id"`
	Message    string `json:"message"`
	AutoEscape bool   `json:"auto_escape"`
}

type GroupMsg struct {
	GroupId    int64  `json:"group_id"`
	Message    string `json:"message"`
	AutoEscape bool   `json:"auto_escape"`
}

//
//
//type Sender struct {
//	to int64
//	message string
//}
//
//func (s *Sender) SetTo(to int64)  {
//	s.to = to
//}
//
//func (s Sender) SetMessage(message string) {
//	s.message = message
//}
//
//func (s Sender) Reply(message string) {
//	s.SetMessage(message)
//}