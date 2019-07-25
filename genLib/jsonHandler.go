// jsonHandler.go

package genLib

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

// // read json datas from file to given interface / structure
// // i.e: err := ReadJson(filename, &person)
// func ReadJson(filename string, interf interface{}) error {
// 	textFileBytes, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		return err
// 	}
// 	err = json.Unmarshal(textFileBytes, &interf)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // Write json datas to file from given interface / structure
// // i.e: err := WriteJson(filename, &person)
// func WriteJson(filename string, interf interface{}) error {
// 	var out bytes.Buffer
// 	jsonData, err := json.Marshal(&interf)
// 	if err != nil {
// 		return err
// 	}
// 	err = json.Indent(&out, jsonData, "", "\t")
// 	if err != nil {
// 		return err
// 	}
// 	f, err := os.Create(filename)
// 	defer f.Close()
// 	if err != nil {
// 		return err
// 	}
// 	_, err = fmt.Fprintln(f, string(out.Bytes()))
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// JsonRead: datas from file to given interface / structure
// i.e: err := ReadJson(filename, &person)
// remember to put upper char at left of variables names to be saved.
func JsonRead(filename string, interf interface{}) (err error) {
	var textFileBytes []byte
	if textFileBytes, err = ioutil.ReadFile(filename); err == nil {
		err = json.Unmarshal(textFileBytes, &interf)
	}
	return err
}

// JsonWrite: datas to file from given interface / structure
// i.e: err := WriteJson(filename, &person)
// remember to put upper char at left of variables names to be saved.
func JsonWrite(filename string, interf interface{}) (err error) {
	var out bytes.Buffer
	var jsonData []byte
	if jsonData, err = json.Marshal(&interf); err == nil {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			if err = ioutil.WriteFile(filename, out.Bytes(), 0644); err == nil {
			}
		}
	}
	return err
}
