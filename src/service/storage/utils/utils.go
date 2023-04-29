package utils

import (
	"fmt"
	"regexp"
)

// FileNamePattern 文件名匹配格式：仅支持文件名 16 进制值 (16位或32位)，后缀名为 png、jpg、jpeg、gif、webp 或无后缀的图片的上传
const FileNamePattern = "([0-9a-fA-F]{16}|[0-9a-fA-F]{32})(\\.(png|jpg|jpeg|gif|webp)|)"
const IdPattern = "[0-9a-fA-F]{24}"

func CheckFileName(path string) bool {
	r, _ := regexp.Compile("^" + FileNamePattern + "$")
	return r.MatchString(path)
}

func PickupReferences(content string) []string {
	r, _ := regexp.Compile(fmt.Sprintf("\\${(%s)}", IdPattern))
	subMatches := r.FindAllStringSubmatch(content, -1)
	fmt.Printf("%+v", subMatches)
	if len(subMatches) == 0 {
		return nil
	}
	var references []string
	for _, subMatch := range subMatches {
		references = append(references, subMatch[1])
	}
	return references
}

func GeneratePath(userID, path string) string {
	return fmt.Sprintf("%s/%s", userID, path)
}
func GenerateResourcePath(userID, path string) string {
	return fmt.Sprintf("resources/%s/%s", userID, path)
}
