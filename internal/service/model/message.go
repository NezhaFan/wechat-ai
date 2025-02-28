package model

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
)

func getMessages(uid string) (messages []RequestMessage) {
	filename := "./chat/" + uid
	hasfile := existsFile(filename)
	if !hasfile {
		return
	}

	f, err := openFile(filename)
	if err != nil {
		log.Println("[ERROR] 无法打开文件：" + err.Error())
		return
	}
	defer f.Close()

	// 读取
	bs, _ := io.ReadAll(f)
	end := bytes.LastIndexByte(bs, ',')
	bs = bs[:end]
	messagesBytes := make([]byte, len(bs)+2)
	messagesBytes[0] = '['
	messagesBytes[len(messagesBytes)-1] = ']'
	copy(messagesBytes[1:], bs)
	err = json.Unmarshal(messagesBytes, &messages)
	if err != nil {
		log.Println("[ERROR] 无法解析文件：" + err.Error())
	}
	return
}

func addMessages(uid, role, msg string) {
	f, err := openFile("./chat/" + uid)
	if err != nil {
		return
	}
	defer f.Close()
	b, _ := json.Marshal(RequestMessage{Role: role, Content: msg})
	f.Write(b)
	f.Write([]byte(",\n"))
}

// 判断文件是否存在
func existsFile(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// 打开文件 (文件不存在则创建)
func openFile(filename string) (*os.File, error) {
	dir := filepath.Dir(filename)
	if !existsFile(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	return os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0775)
}
