package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

const (
	opAdd      = "add"
	opList     = "list"
	opRemove   = "remove"
	opFindByID = "findById"
)

type Arguments map[string]string

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func Perform(args Arguments, writer io.Writer) (err error) {

	if err = checkArguments(args); err != nil {
		return
	}
	switch args["operation"] {
	case opList:
		var dataIn []byte
		if dataIn, err = ioutil.ReadFile(args["fileName"]); err != nil {
			log.Fatal("Cannot load settings:", err)
		}
		if _, err = writer.Write(dataIn); err != nil {
			return
		}
	case opAdd:
		if err = doAdd(args["fileName"], args["item"], writer); err != nil {
			log.Fatal(err)
		}
	case opRemove:
		if err = doRemove(args["fileName"], args["id"], writer); err != nil {
			log.Fatal(err)
		}
	case opFindByID:
		if err = doFindByID(args["fileName"], args["id"], writer); err != nil {
			log.Fatal(err)
		}
	default:
		err = fmt.Errorf("Operation %s not allowed!", args["operation"])
	}
	return
}

func doFindByID(fileName, id string, writer io.Writer) (err error) {

	var users []User
	if users, err = loadUsers(fileName); err != nil {
		return
	}
	for _, user := range users {
		if user.ID == id {
			var dataOut []byte
			if dataOut, err = json.Marshal(user); err != nil {
				return fmt.Errorf("JSON marshaling failed: %s", err.Error())
			}
			_, err = writer.Write(dataOut)
			break
		}
	}
	return
}

func doRemove(fileName, id string, writer io.Writer) (err error) {

	var users []User
	if users, err = loadUsers(fileName); err != nil {
		return
	}
	for idx, user := range users {
		if user.ID == id {
			var fileData []byte
			users = append(users[:idx], users[idx+1:]...)
			if fileData, err = json.Marshal(&users); err != nil {
				return fmt.Errorf("JSON marshaling failed: %s", err.Error())
			}
			if err = ioutil.WriteFile(fileName, fileData, 0600); err != nil {
				return fmt.Errorf("Cannot write updated settings file: %s", err.Error())
			}
			break
		}
	}
	_, err = writer.Write([]byte(fmt.Sprintf("Item with id %s not found", id)))
	return
}

func doAdd(fileName, item string, writer io.Writer) (err error) {

	var users []User
	if users, err = loadUsers(fileName); err != nil {
		return
	}
	var argUser User
	if err = json.Unmarshal([]byte(item), &argUser); err != nil {
		return fmt.Errorf("JSON unmarshalling failed: %s", err.Error())
	}
	for _, user := range users {
		if user.ID == argUser.ID {
			_, err = writer.Write([]byte("Item with id 1 already exists"))
			return
		}
	}
	var data []byte
	users = append(users, argUser)
	if data, err = json.Marshal(&users); err != nil {
		return fmt.Errorf("JSON marshaling failed: %s", err.Error())
	}
	err = ioutil.WriteFile(fileName, data, 0777)
	if err != nil {
		err = fmt.Errorf("Cannot write updated settings file: %s", err.Error())
	}
	return
}

func loadUsers(fileName string) (users []User, err error) {

	var data []byte
	if data, err = ioutil.ReadFile(fileName); err == nil {
		if err = json.Unmarshal(data, &users); err != nil {
			err = fmt.Errorf("JSON unmarshalling failed: %s", err.Error())
			return
		}
	}
	if err != nil && errors.Is(err, os.ErrNotExist) {
		err = fmt.Errorf("open file error: %s", err.Error())
	}
	return users, nil
}

func checkArguments(args Arguments) (err error) {

	if args["fileName"] == "" {
		return errors.New("-fileName flag has to be specified")
	}
	if args["operation"] == "" {
		return errors.New("-operation flag has to be specified")
	}
	if (args["operation"] == "remove" || args["operation"] == "findById") && args["id"] == "" {
		return errors.New("-id flag has to be specified")
	}
	if args["item"] == "" && args["operation"] == "add" {
		return errors.New("-item flag has to be specified")
	}
	return err
}

func main() {
	args := Arguments{}
	id := flag.String("id", "", "")
	item := flag.String("item", "", "")
	operation := flag.String("operation", "", "")
	fileName := flag.String("fileName", "", "")
	flag.Parse()
	args["id"] = *id
	args["item"] = *item
	args["operation"] = *operation
	args["fileName"] = *fileName
	buf := os.Stdout
	err := Perform(args, buf)
	if err != nil {
		panic(err)
	}
}
