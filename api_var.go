package rbq

// GetDeviceList 登录设备的 model 参数取值，和 SetLoginDevice 的 model 参数取值的一个特定值
// 下列值只是一些特殊值，调用 api 时可以自由传入任意值
const (
	APIModelIphoneX    = "iPhone11,2"
	APIModelIphoneXR   = "iPhone11,8"
	APIModelIphone11   = "iPhone12,1"
	APIModelIphone12   = "iPhone13,2"
	APIModelIpadPro    = "iPad8,1"
	APIModelIpadMini   = "iPad11,2"
	APIModelIpadAir4   = "iPad13,2"
	APIModelAppleWatch = "Apple Watch"
)

// todo GetRecord 的 outFormat 参数取值, 由于 get_record 接口暂未实现，所以这里先注释掉
//const (
//	APIRecordFormatMp3  = "mp3"
//	APIRecordFormatAmr  = "amr"
//	APIRecordFormatWma  = "wma"
//	APIRecordFormatM4a  = "m4a"
//	APIRecordFormatSpx  = "spx"
//	APIRecordFormatOgg  = "ogg"
//	APIRecordFormatWav  = "wav"
//	APIRecordFormatFlac = "flac"
//)

//var apiRecordFormatTypes = []string{
//	APIRecordFormatMp3,
//	APIRecordFormatAmr,
//	APIRecordFormatWma,
//	APIRecordFormatM4a,
//	APIRecordFormatSpx,
//	APIRecordFormatOgg,
//	APIRecordFormatWav,
//	APIRecordFormatFlac,
//}

// GetGroupHonorInfo 的 honerType 参数取值
const (
	APIHonerTypeTalkative    = "talkative"     // 龙王
	APIHonerTypePerformer    = "performer"     // 群聊之火
	APIHonerTypeLegend       = "legend"        // 群聊炽焰
	APIHonerTypeStrongNewbie = "strong_newbie" // 冒尖小春笋
	APIHonerTypeEmotion      = "emotion"       // 快乐之源
	APIHonerTypeAll          = "all"           // 所有荣誉
)
