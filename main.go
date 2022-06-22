package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type Arguments map[string]string

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type List struct {
	Users []User
}

func Perform(args Arguments, writer io.Writer) error {
	const k = "flag has to be specified"
	if args["operation"] == "" {
		err := fmt.Errorf("-operation %s", k)
		return err
	}
	if args["item"] == "" && args["operation"] == "add" {
		err := fmt.Errorf("-item %s", k)
		return err
	}
	if args["id"] == "" && args["operation"] == "remove" {
		err := fmt.Errorf("-id %s", k)
		return err
	}
	if args["id"] == "" && args["operation"] == "findById" {
		err := fmt.Errorf("-id %s", k)
		return err
	}
	if args["operation"] != "" && args["operation"] != "list" && args["operation"] != "add" && args["operation"] != "remove" && args["operation"] != "findById" {
		err := fmt.Errorf("Operation %s not allowed!", args["operation"])
		return err
	}
	if args["fileName"] == "" {
		err := fmt.Errorf("-fileName %s", k)
		return err
	}
	if args["operation"] == "list" {
		dataIn, err := ioutil.ReadFile(args["fileName"])
		if err != nil {
			log.Fatal("Cannot load settings:", err)
		}
		writer.Write(dataIn)
	}

	if args["operation"] == "add" {
		var luser User
		err := json.Unmarshal([]byte(args["item"]), &luser)
		if err != nil {
			log.Fatal("JSON unmarshaling failed:", err)
		}
		file, err := os.OpenFile(args["fileName"], os.O_RDONLY|os.O_CREATE, 0777)
		if err != nil {
			log.Fatal("Cannot load settings:", err)
		}
		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		var res []User
		err = json.Unmarshal(data, &res)
		for i := 0; i < len(res); i++ {
			if res[i].Id == luser.Id {
				dataIn := "Item with id 1 already exists"
				writer.Write([]byte(dataIn))
			}
		}
		res = append(res, luser)
		data, err = json.Marshal(&res)
		if err != nil {
			log.Fatal("JSON marshaling failed:", err)
		}
		err = ioutil.WriteFile(args["fileName"], data, 0777)
		if err != nil {
			log.Fatal("Cannot write updated settings file:", err)
		}
	}

	if args["operation"] == "remove" {
		file, err := os.OpenFile(args["fileName"], os.O_RDONLY|os.O_CREATE, 0777)
		if err != nil {
			log.Fatal("Cannot load settings:", err)
		}
		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		var res []User
		err = json.Unmarshal(data, &res)
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < len(res); i++ {
			if res[i].Id == "1" {
				res[i] = res[len(res)-1]
				res[len(res)-1] = User{}
				res = res[:len(res)-1]
				fileData, err := json.Marshal(&res)
				if err != nil {
					log.Fatal("JSON marshaling failed:", err)
				}
				err = ioutil.WriteFile(args["fileName"], fileData, 0777)
				if err != nil {
					log.Fatal("Cannot write updated settings file:", err)
				}
			}
		}
		writer.Write([]byte("Item with id 2 not found"))

	}
	if args["operation"] == "findById" {
		file, err := os.OpenFile(args["fileName"], os.O_RDONLY|os.O_CREATE, 0777)
		if err != nil {
			log.Fatal("Cannot load settings:", err)
		}
		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		var res []User
		err = json.Unmarshal(data, &res)
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < len(res); i++ {
			if res[i].Id == args["id"] {
				dataOut, err := json.Marshal(&res[i])
				if err != nil {
					log.Fatal("JSON marshaling failed:", err)
				}
				writer.Write(dataOut)
			}
		}
	}
	return nil
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
