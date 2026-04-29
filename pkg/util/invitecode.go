// @Title 邀请码生成器
// @Description
// @Author 蔺保仲 2020/04/20
// @Update 蔺保仲 2020/04/20
package util

var (
	invitecode = NewHashID("Harry19921115")
)

func GetInviteCode(uid uint64) (string, error) {
	return invitecode.Encode(uid)
}

// Decode 根据邀请码,获取用户id

func GetIDFromInviteCode(inviteCode string) (uint64, error) {
	return invitecode.Decode(inviteCode)
}
