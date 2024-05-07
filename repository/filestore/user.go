package filestore

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.dev/constant"
	"go.dev/entity"
	"os"
	"strconv"
	"strings"
)

type Filestore struct {
	Filepath          string
	serializationMode string
}

func New(path string, serializationMode string) Filestore {
	return Filestore{
		Filepath: path, serializationMode: serializationMode,
	}
}

func (f Filestore) Save(u entity.User) {
	f.writeUserToFile(u)
}

func (f Filestore) Load() []entity.User {

	var uStore []entity.User

	file, err := os.Open(f.Filepath)
	if err != nil {
		fmt.Println("can't open file", err)
		return nil
	}
	var data = make([]byte, 10240)
	_, oErr := file.Read(data)
	if oErr != nil {
		fmt.Println("can't read from file")
		return nil
	}
	var datastr = string(data)
	userSlice := strings.Split(datastr, "\n")

	for _, u := range userSlice {

		var userStruct = entity.User{}

		if u[0] != '{' && u[len(u)-1] != '}' {
			continue
		}

		switch f.serializationMode {

		case constant.ManDarAvardi_serializationMode:
			var dErr error

			userStruct, dErr = deserilizeFromManDaravardi(u)
			if dErr != nil {
				fmt.Println("can't deserilize user record to user struct", dErr)
				return nil
			}

		case constant.Json_serializationMode:

			uErr := json.Unmarshal([]byte(u), &userStruct)
			if uErr != nil {
				fmt.Println("can't deserilize user record to user struct from jason mode", uErr)
				return nil
			}

		default:
			fmt.Println("")
		}

	}
	return uStore
}

func (f Filestore) writeUserToFile(user entity.User) {
	var file *os.File
	file, err := os.OpenFile(f.Filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println("can't open or creat file")
		return
	}
	var data []byte
	defer file.Close()
	if f.serializationMode == constant.ManDarAvardi_serializationMode {

		data = []byte(fmt.Sprintf("ID:%d ,name:%s ,email:%s ,password:%s\n",
			user.ID, user.Name, user.Email, user.Password))

	} else if f.serializationMode == constant.Json_serializationMode {

		data, err = json.Marshal(user)
		if err != nil {
			fmt.Println("can't marshal user struct to json", err)
			return
		}
		data = append(data, []byte("\n")...)

	} else {
		fmt.Println("invalid serialization Mode")
		return
	}

	file.Write(data)
}

func deserilizeFromManDaravardi(userstr string) (entity.User, error) {
	if userstr == "" {
		return entity.User{}, errors.New("user string is empty")
	}
	var user = entity.User{}
	// fmt.Println("Line of file:",index,"User:",u)
	userField := strings.Split(userstr, ",")
	for _, Field := range userField {
		value := strings.Split(Field, ":")
		fieldNmae := strings.ReplaceAll(value[0], " ", "")
		fieldValue := strings.ReplaceAll(value[1], " ", "")

		switch fieldNmae {
		case "ID":
			id, err := strconv.Atoi(fieldValue)
			if err != nil {
				return entity.User{}, errors.New("strconv error")
			}
			user.ID = id
		case "name":
			user.Name = fieldValue
		case "email":
			user.Email = fieldValue
		case "password":
			user.Password = fieldValue
		}

	}

	return user, nil
}
