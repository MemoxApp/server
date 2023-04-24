package mail

import (
	"fmt"
	"math/rand"
	"time"
)

// init 设置随机数种子
func init() {
	rand.Seed(time.Now().Unix())
}

// randomNumberStr 获取指定长度的随机数字文本
func randomNumberStr(length int) string {
	result := 0
	for i := 0; i < length; i++ {
		result = result*10 + rand.Int()%10
	}
	return fmt.Sprintf("%06d", result)
}

func newCode(id string, length int) Code {
	return Code{
		ID:   id,
		Code: randomNumberStr(length),
	}
}

func codeKey(id string) string {
	return fmt.Sprintf(keyVerifyCode, id)
}

func coolDownKey(id string) string {
	return fmt.Sprintf(keyVerifyCodeCoolDown, id)
}
