package utils

import (
	"io/ioutil"
	"os"
)

func VerifyFileExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}

	return exist
}

func VerifyPathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func GetCurPathFileNames(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		names = append(names, file.Name())
	}

	return names, nil
}

func GetFileNames(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, file := range files {
		if file.IsDir() {
			fileNameList,err := GetFileNames(file.Name())
			if err != nil{
				return nil,err
			}

			names = append(names, fileNameList...)
		}else{
			names = append(names, file.Name())
		}
	}

	return names, nil
}

func Mythical(a string)  {

}