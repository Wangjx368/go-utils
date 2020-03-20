package upload

import (
	"fmt"
	"io/ioutil"
	"openapiai/services/speech/common"
	"regexp"
	"strings"
	"testing"
)

func TestDoUploadFile(t *testing.T) {
	onePicPath1 := "../../a-z-animals/Aardvark/Aardvark_01.jpg"
	onePicPath2 := "../../a-z-animals/Aardvark/Aardvark_02.jpg"
	picList := []string{
		onePicPath1,
		onePicPath2,
	}

	for _, pic := range picList {
		doUpload(pic, "it's me")
	}
}

func doUpload(imgFile string, fnlAudioName string) {
	picBytes, err := ioutil.ReadFile(imgFile)
	if err != nil {
		fmt.Println("wav read error!")
		return
	}
	picURL := DoUploadFile(picBytes, "wav")
	common.WriteWithIo("./", "tts.list", fnlAudioName+"[>]"+picURL+"\n")
}

func TestWalkAnimals(t *testing.T) {
	animalDir := "../../laoshu"
	walkDir(animalDir)
}

func walkDir(path string) {
	//fmt.Println("current dir is : ", path)
	rd, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("read dir error!")
		return
	}
	for _, fi := range rd {
		if fi.IsDir() {
			newPath := path + "/" + fi.Name()
			walkDir(newPath)
		} else {
			picPath := path + "/" + fi.Name()
			fmt.Println("  file path is : ", picPath)

			fnlAudioName := geneFnlAudioName(picPath)
			fmt.Println("[fileName]:", fnlAudioName)
			doUpload(picPath, fnlAudioName)
		}
	}
}

func TestUploadExe(t *testing.T) {
	exeFilePath := "../../A_exe/AIScanner.exe"
	picBytes, err := ioutil.ReadFile(exeFilePath)
	if err != nil {
		fmt.Println("pic read error!")
		return
	}
	picURL := DoUploadFile(picBytes, "exe")
	fmt.Println(picURL)
}

func geneFnlAudioName(picPath string) (fnlAudioName string) {
	// 获取文件名
	tmpList := strings.Split(picPath, "/")
	tmpName := tmpList[len(tmpList)-1]
	// 去掉wav后缀
	audioName := tmpName[:len(tmpName)-4]
	// 只保留 数字 字母 以及'
	r, _ := regexp.Compile(`[^a-zA-z0-9\']+`)
	fnlAudioName = r.ReplaceAllString(audioName, " ")
	// 多个空格转换为一个空格
	r2, _ := regexp.Compile(`[ ]+`)
	fnlAudioName = r2.ReplaceAllString(fnlAudioName, " ")
	// 去掉首位空格
	return strings.TrimSpace(fnlAudioName)
}
