package rbq

var GlobalVar = &Global{}

type Global struct {
	botQQ                    int64
	botNickname              string
	canSendImg               bool
	canSendRecord            bool
	onlineClients            []*OnlineClient
	friendList               FriendList
	unidirectionalFriendList UnidirectionalFriendList
}

// GetBotQQ 获取当前机器人的QQ号
func (g *Global) GetBotQQ() int64 {
	return g.botQQ
}

// GetBotNickname 获取当前机器人的昵称
func (g *Global) GetBotNickname() string {
	return g.botNickname
}

// CanSendImg 获取当前机器人是否可以发送图片
func (g *Global) CanSendImg() bool {
	return g.canSendImg
}

// CanSendRecord 获取当前机器人是否可以发送语音
func (g *Global) CanSendRecord() bool {
	return g.canSendRecord
}

func (g *Global) GetOnlineClients() []*OnlineClient {
	return g.onlineClients
}

func (g *Global) GetFriendList() FriendList {
	return g.friendList
}

type FriendList []*FriendInfo

// Search 通过用户ID查找好友信息
func (f FriendList) Search(userId int64) *FriendInfo {
	// 从 FriendList 中查找 给定的 userId
	for _, friend := range f {
		if friend.UserId == userId {
			return friend
		}
	}
	return nil
}

type UnidirectionalFriendList []*UnidirectionalFriendInfo

// Search 通过用户ID查找好友信息
func (f UnidirectionalFriendList) Search(userId int64) *UnidirectionalFriendInfo {
	// 从 UnidirectionalFriendList 中查找 给定的 userId
	for _, friend := range f {
		if friend.UserId == userId {
			return friend
		}
	}
	return nil
}
