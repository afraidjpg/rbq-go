<div style="text-align: center">

<h1>rbq-go</h1>
<h3>一个以go-cqhttp为后端的插件开发框架</h3>
<img src="asset/logo.png" alt="logo" width="400"/>

![](https://img.shields.io/badge/go-v1.18%2B-blue)
[![](https://img.shields.io/badge/go--cqhttp-v1.0.0-orange)](https://github.com/Mrs4s/go-cqhttp)

</div>
<h2>写在前面</h2>

这是一个个人学习性质的项目，不能提供任何可靠性保证，包括但不限于：**代码安全与质量，相关文档，issue回复的及时性等**。

基于上述理由，本项目的更新速度可能会十分缓慢（就是懒）。如果star高了，用的人多了，项目变得活跃了，那我应该会变得很积极（手动狗头）

一个人的力量是有限的，所以如果你愿意使用，并在使用过程中有任何问题、新的feature/代码修改建议，欢迎提出issue，我会尽可能的回复并改进。

<span style="font-size:18px">PS: rbq = **r**o**b**ot **q**q</span>

<h2>简介</h2>

rbq-go是一个以go-cqhttp为后端的插件开发框架，它的目标是提供一个简单、易用的插件开发环境，通过封装go-cqhttp的api，使插件开发者可以更简单的与qq交互。

项目灵感来自于许久以前写着玩儿的时候，发现要写一堆代码去调用后端，感觉很麻烦，于是有了写这个项目的念头，我但愿这个项目能够帮助到你。

为了保持项目简单易懂，项目会 **尽量避免使用过于复杂的设计**
，降低学习成本。同时保留高级功能的空间，使之可以胜任复杂场景。期望是不再需要理解以及再封装onebot的文档，即可简单的开发插件

<h2>安装</h2>

```
go get -u github.com/afraidjpg/rbq-go
```

<h2>使用</h2>

示例代码请 [点击这里](./example)

```go
package main

import (
	"fmt"
	"github.com/afraidjpg/rbq-go"
)

// YourPlugin 这是你的插件本体
func YourPlugin(ctx *rbq.Context) {
	msg := ctx.GetMessage()
	if msg == "hello world" {
		ctx.Reply("welcome")
		return
	}

	ctx.Reply("test") // 发送 "test"

	ctx.AddText("boo") // 添加消息但不立刻发送
	ctx.AddText("bar")
	ctx.Reply() // 发送 "boobar"

	ctx.AddText("hello") // 添加消息但不立刻发送
	ctx.AddText(" ")
	ctx.Reply("world") // 发送 "hello world"
}

// main 方法
func main() {
	bot := rbq.NewApp()          // 新建app
	pld := bot.GetPluginLoader() // 获取插件装载器
	// 添加单个插件
	pld.BindPlugin(YourPlugin, nil)

	// 创建一个插件组，并向插件组中加入插件
	gp := pld.Group("gp", nil)
	gp.BindPlugin(YourPlugin, nil)

	// 添加插件中发生的错误
	errs := gp.GetErrors()
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
			return
		}
	}
	// 不给参数，自动设置为 127.0.0.1:8080（go-cqhttp的默认正向websocket端口）
	bot.Run("")
}
```

<h2>运行</h2>

- 下载 [go-cqhttp](https://github.com/Mrs4s/go-cqhttp/releases)
- 启动 `go-cqhttp` 并使用选择正向websocket连接
- 运行你的程序

<h2>插件选项</h2>

`BindPlugin`的第二个参数是一个`*PluginOption`类型的变量，可以为插件设置一些选项，例如：

```go
package main

import (
	"fmt"
	"github.com/afraidjpg/rbq-go"
)

var SimpleOption = &rbq.PluginOption{
	Name: "example_reply", // 插件名称
	FilterFunc: []rbq.PluginFilterFunc{
		func(ctx *rbq.Context) bool {
			return ctx.IsGroup() // 只回复群消息
		},
	},
	RecoverFunc: func(ctx *rbq.Context, err any) {
		fmt.Println("插件运行错误:", err)
	}, // 当插件运行错误的时候执行的逻辑
}

func YourPlugin(ctx *rbq.Context) {
	ctx.Reply("test")
}

// 然后在绑定插件的时候使用
func main() {
	bot := rbq.NewApp()
	pld := bot.GetPluginLoader()
	pld.BindPlugin(YourPlugin, SimpleOption) // 该插件将应用上述设置

	gp := pld.Group("gp", &rbq.PluginGroupOption{SimpleOption})
	gp.BindPlugin(YourPlugin, nil) // gp组内的插件选项如果未设置，则会直接应用gp组的选项
	// ...
}

```

<h2>CQ码</h2>

如何获取CQ码与发送CQ码：

```go
package yourplugin

import (
	"fmt"
	"github.com/afraidjpg/rbq-go"
)

func YourPlugin(ctx *rbq.Context) {
	// 获取该条消息中的所有CQ码，返回类型为 []CQCodeInterface，需要手动转换类型
	cq := ctx.GetAllCQCode()
	fmt.Println(cq)

	// 获取具体的某种CQ码的方式为 ctx.GetCQX(), 其中X为CQ码的类型
	// 需要注意的是，一条消息中允许存在多个相同类型的CQ码时，会返回[]*CQX
	// 如果不会返回多个，则返回*CQX 
	// 例如 ctx.GetCQAt() 获取at码
	ats := ctx.GetCQAt() // 返回 []*CQAt，获取该条消息中的所有at
	for _, at := range ats {
		fmt.Println(at.GetQQ()) // 获取at的qq号
	}

	// 发送CQ码的方式为 ctx.AddCQX(), X代表CQ码类型，且首字母大写
	// 如 At、Image、Face...
	ctx.AddCQAt(1234567)
	ctx.AddCQFace(22)
	// ...

	ctx.Reply() // 发送消息
}

```

CQ码的支持与 go-cqhttp 同步，可以[点击这里](https://github.com/Mrs4s/go-cqhttp)查看README中对CQ码的描述

|    CQ码    | 进度  | 收发限制 | 备注                                                                    |
|:---------:|:---:|:----:|-----------------------------------------------------------------------|
|   face    | 完成  | 收&发  |                                                                       |                                                                       |
|  record   | 完成  | 收&发  | ⚠️可能存在问题，具体查看该[issue](https://github.com/Mrs4s/go-cqhttp/issues/1749) |
|   video   | 完成  | 收&发  | 如果没有后缀，通常为mp4                                                         |
|    at     | 完成  | 收&发  |                                                                       |
|   share   | 完成  | 收&发  | ⚠️对群时发送/接受时可能会失败                                                      |
|   music   | 完成  | 收&发  |                                                                       |
|   reply   | 完成  | 收&发  |                                                                       |
|  forward  | 完成  |  收   |                                                                       |
|   node    | 完成  |  发   |                                                                       |
|    xml    | 未实现 |      |                                                                       |
|   json    | 未实现 |      |                                                                       |
|   image   | 完成  | 收&发  |                                                                       |
|  redbag   | 完成  |  收   |                                                                       |
|   poke    | 完成  |  发   |                                                                       |
| cardimage | 完成  |  发   |                                                                       |
|    tts    | 完成  |  发   |                                                                       |

<h2>API</h2>

实现了对 go-cqhttp 的 api 调用，可以通过 `rbq.Api` 获取 `APIWrapper` 对象，然后调用对应的方法即可

|      分类      |          名称          |             方法名             |         go-cqhttp 接口名          |    进度     |
|:------------:|:--------------------:|:---------------------------:|:------------------------------:|:---------:|
|    bot账号     |       获取登录号信息        |        GetLoginInfo         |         get_login_info         |    完成     |
|    bot账号     |       设置登录号资料        |        SetQQProfile         |         set_qq_profile         |    完成     |
|    bot账号     |       获取企点账号信息       |    GetQidianAccountInfo     |    qidian_get_account_info     |    完成     |
|    bot账号     |        获取在线机型        |        GetDeviceList        |        _get_model_show         |    完成     |
|    bot账号     |        设置在线机型        |       SetOnlineDevice       |        _set_model_show         |    完成     |
|    bot账号     |        获取在线设备        |      GetOnlineClients       |       get_online_clients       |    完成     |
|     好友信息     |       获取陌生人信息        |       GetStrangerInfo       |       get_stranger_info        |    完成     |
|     好友信息     |        获取好友列表        |        GetFriendList        |        get_friend_list         |    完成     |
|     好友信息     |       获取单向好友列表       | GetUnidirectionalFriendList | get_unidirectional_friend_list |    完成     |
|     好友操作     |         删除好友         |        DeleteFriend         |         delete_friend          |    完成     |
|     好友操作     |         删除单向         | DeleteUnidirectionalFriend  |  delete_unidirectional_friend  |    完成     |
|      消息      |        发送私聊消息        |       SendPrivateMsg        |        send_private_msg        |    完成     |
|      消息      |        发送群消息         |        SendGroupMsg         |         send_group_msg         |    完成     |
|      消息      |         发送消息         |           SendMsg           |            send_msg            |    完成     |
|      消息      |         获取消息         |           GetMsg            |            get_msg             |    完成     |
|      消息      |         撤回消息         |          DeleteMsg          |           delete_msg           |    完成     |
|      消息      |        标记消息已读        |        MarkMsgAsRead        |        mark_msg_as_read        |    完成     |
|      消息      |       获取合并转发内容       |        GetforwardMsg        |        get_forward_msg         |    完成     |
|      消息      |    发送合并转发 ( 群聊 )     |     SendGroupForwardMsg     |     send_group_forward_msg     |    完成     |
|      消息      |    发送合并转发 ( 好友 )     |    SendPrivateForwardMsg    |    send_private_forward_msg    |    完成     |
|      消息      |      获取群消息历史记录       |     GetGroupMsgHistory      |     get_group_msg_history      |    完成     |
|      图片      |        获取图片信息        |          GetImage           |           get_image            |    完成     |
|      图片      |      检查是否可以发送图片      |        CanSendImage         |         can_send_image         |    完成     |
|      图片      |        图片 OCR        |          OcrImage           |           ocr_image            |    完成     |
|      语音      |         获取语音         |          GetRecord          |           get_record           | cqhttp未实现 |
|      语音      |      检查是否可以发送语音      |        CanSendRecord        |        can_send_record         |    完成     |
|      处理      |       处理加好友请求        |     SetFriendAddRequest     |     set_friend_add_request     |    完成     |
|      处理      |      处理加群请求／邀请       |     SetGroupAddRequest      |     set_group_add_request      |    完成     |
|     群信息      |        获取群信息         |        GetGroupInfo         |         get_group_info         |    完成     |
|     群信息      |        获取群列表         |        GetGroupList         |         get_group_list         |    完成     | 
|     群信息      |       获取群成员信息        |     GetGroupMemberInfo      |     get_group_member_info      |    完成     |
|     群信息      |       获取群成员列表        |     GetGroupMemberList      |     get_group_member_list      |    完成     |
|     群信息      |       获取群荣誉信息        |      GetGroupHonorInfo      |      get_group_honor_info      |    完成     |
|     群信息      |       获取群系统消息        |      GetGroupSystemMsg      |      get_group_system_msg      |    完成     |
|     群信息      |       获取精华消息列表       |      GetEssenceMsgList      |      get_essence_msg_list      |    完成     |
|     群信息      |    获取群 @全体成员 剩余次数    |     GetGroupAtAllRemain     |    get_group_at_all_remain     |    完成     |
|     群设置      |         设置群名         |        SetGroupName         |         set_group_name         |    完成     |
|     群设置      |        设置群头像         |      SetGroupPortrait       |       set_group_portrait       |    完成     |
|     群设置      |        设置群管理员        |        SetGroupAdmin        |        set_group_admin         |    完成     |
|     群设置      |    设置群名片 ( 群备注 )     |        SetGroupCard         |         set_group_card         |    完成     |
|     群设置      |       设置群专属头衔        |    SetGroupSpecialTitle     |    set_group_special_title     |    完成     |
|     群操作      |        群单人禁言         |         SetGroupBan         |         set_group_ban          |    完成     |
|     群操作      |        群全员禁言         |      SetGroupWholeBan       |      set_group_whole_ban       |    完成     |
|     群操作      |       群匿名用户禁言        |    SetGroupAnonymousBan     |    set_group_anonymous_ban     | cqhttp未实现 |
|     群操作      |        设置精华消息        |        SetEssenceMsg        |        set_essence_msg         |    完成     |
|     群操作      |        删除精华消息        |      DeleteEssenceMsg       |       delete_essence_msg       |    完成     |
|     群操作      |         群打卡          |        SendGroupSign        |        send_group_sign         |    完成     |
|     群操作      |        群设置匿名         |      SetGroupAnonymous      |      set_group_anonymous       |    完成     |
|     群操作      |        发送群公告         |       SendGroupNotice       |       _send_group_notice       |    完成     |
|     群操作      |        获取群公告         |       GetGroupNotice        |       _get_group_notice        |    完成     |
|     群操作      |         群组踢人         |        SetGroupKick         |         set_group_kick         |    完成     |
|     群操作      |         退出群组         |        SetGroupLeave        |        set_group_leave         |    完成     |
|      文件      |        上传群文件         |       UploadGroupFile       |       upload_group_file        |    完成     |
|      文件      |        删除群文件         |       DeleteGroupFile       |       delete_group_file        |    完成     |
|      文件      |       创建群文件文件夹       |    CreateGroupFileFolder    |    create_group_file_folder    |    完成     |
|      文件      |       删除群文件文件夹       |    DeleteGroupFileFolder    |      delete_group_folder       |    完成     |
|      文件      |      获取群文件系统信息       |   GetGroupFileSystemInfo    |   get_group_file_system_info   |    完成     |
|      文件      |      获取群根目录文件列表      |      GetGroupRootFiles      |      get_group_root_files      |    完成     |
|      文件      |      获取群子目录文件列表      |    GetGroupFilesByFolder    |   get_group_files_by_folder    |    完成     |
|      文件      |      获取群文件资源链接       |       GetGroupFileUrl       |       get_group_file_url       |    完成     |
|      文件      |        上传私聊文件        |      UploadPrivateFile      |      upload_private_file       |    完成     |
| Go-CqHttp 相关 |      获取 Cookies      |         GetCookies          |          get_cookies           | cqhttp未实现 |
| Go-CqHttp 相关 |    获取 CSRF Token     |        GetCsrfToken         |         get_csrf_token         | cqhttp未实现 |
| Go-CqHttp 相关 |     获取 QQ 相关接口凭证     |       GetCredentials        |        get_credentials         | cqhttp未实现 |
| Go-CqHttp 相关 |        获取版本信息        |       GetVersionInfo        |        get_version_info        |    完成     |
| Go-CqHttp 相关 |         获取状态         |          GetStatus          |           get_status           |    完成     |
| Go-CqHttp 相关 |         清理缓存         |         CleanCache          |          clean_cache           |    完成     |
| Go-CqHttp 相关 |       重载事件过滤器        |      ReloadEventFilter      |      reload_event_filter       |    完成     |
| Go-CqHttp 相关 |      下载文件到缓存目录       |        DownloadFile         |         download_file          |    完成     |
| Go-CqHttp 相关 |       检查链接安全性        |       CheckUrlSafely        |        check_url_safely        |    完成     |
| Go-CqHttp 相关 |  获取中文分词 ( 隐藏 API )   |        getWordSlices        |        .get_word_slices        |    完成     |
| Go-CqHttp 相关 | 对事件执行快速操作 ( 隐藏 API ) |    handleQuickOperation     |    .handle_quick_operation     |    完成     |

<h2 id="global_var"> 全局变量 </h2>
