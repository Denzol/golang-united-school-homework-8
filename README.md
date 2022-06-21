# golang-united-school-homework-8

In the task you have to write console application for managing users list. It should accept there types of operation:
`add
list
findById
remove`
Users list should be stored in the json file. When you start your application and tries to perform some operations, existing file should be used or new one should be created if it does not exist.
Example of the json file (users.json):
`[{id: "1", email: «test@test.com», age: 31}, {id: "2", email: «test2@test.com», age: 41}]`
In the `main.go` file you can find a function called Perform(args Arguments, writer io.Writer) error.
You have to call this function from the main function and pass arguments from the console and os.Stdout stream. Perform function body you have to write by yourself :).
Arguments - is a `map[string]string` with the following fields:
`id, item, operation and fileName`. Create a separate type for the arguments map: `type Arguments map[string]string` to prevent duplication of `map[string]string`. Unit tests use Arguments type also.
Arguments should be passed via console flags:
`./main.go -operation «add» -item ‘{«id»: "1", «email»: «email@test.com», «age»: 23}’ -fileName «users.json»`
`-operation`, `-item`and `-fileName` are console flags. To parse them and build map you can take a Look at «flag» package: https://golang.org/pkg/flag/.
Pay attention that `-fileName` flag should be provided every time with the name of file where you store users!

#### Getting list of items:
Application has to retrieve list from the users.json file and print it to the `io.Writer` stream. Use writer from the argument of Perform function to print the result! It is important for passing unit tests. It can be smth like `writer.Write(bytes)`
File content: `[{«id»: "1", «email»: «email@test.com», «age»: 23}]`
Command: `./main -operation «list» -fileName «users.json»` (main is bult go application. Binary file after go build command)
Output to the console: `[{«id»: "1", «email»: «email@test.com», «age»: 23}]`
If file is empty then nothing should be printed to the console.
**Errors:** 
1. If `-operation` flag is missing, then error `-operation` flag has to be specified» has to be returned from Perform function. Package `errors` can be used for creating errors (https://golang.org/pkg/errors/).
2. If `—fileName` flag is missing, then error «-fileName flag has to be specified» has to be returned from Perform function.
3. If operation can not be handled, for example «abc» operation, then «Operation abcd not allowed!» error has to be return from the Perform function
All cases are covered by unit tests. If you want to be sure your solution works correct, just start `go test -v` command in the root of the repo
#### Adding new item:
For adding new item to the array inside users.json file, application should provide the following cmd command:
`./main -operhation «add» -item «{«id»: "1", «email»: «email@test.com», «age»: 23}» -fileName «users.json»`
`-item` - valid json object with the id, email and age fields
**Errors:**
1. All errors about operation and fileName flags mentioned above
2. If `-item` flag is not provided Perform function should return error «-item flag has to be specified»

#### Remove user
Application should allow removing user with the following command:
`./main -operation «remove» -id "2" -fileName «users.json»`
If user with id `"2"`, for example, does not exist, Perform functions should print message to the `io.Writer` «Item with id 2 not found».
If user with specified id exists, it should be removed from the users.json file.
**Errors:**
1. All errors about operation and fileName flags mentioned above
2. If `-id` flag is not provided error «-id flag has to be specified» should be returned from Perform function

#### Find by id
Application should allow finding user by id with the following command:
`./main -operation «findById» -id "1" -fileName «users.json»`
If user with specified id does not exist in the users.json file, then empty string has to be written to  the `io.Writer`
If user exists, then json object should be written in `io.Writer`
**Errors:**
1. All errors about operation and fileName flags mentioned above
2. If `-id` flag is not provided error «-id flag has to be specified» should be returned from Perform function
All cases of the task are covered by unit tests, So, you can check your solution during the implementation.

### Useful info:
1. For opening and creating file use `os` package and `OpenFile` function https://golang.org/pkg/os/
2. To simply read file use `ioutil` package and `ReadAll` function https://golang.org/pkg/io/ioutil/
3. To convert json string to the object use `encoding/json` package and `Unmarshal` function: https://golang.org/pkg/encoding/json/
4. To convert json array or object to string use json package and `Marshal` function.
5. Go does not have throw operator and try catch statement. Instead, it has multi return and allows returning error from a function: `func () (User, error) {}`
Take a look: https://medium.com/@hussachai/error-handling-in-go-a-quick-opinionated-guide-9199dd7c7f76
6. If you receive error in Perform function, just call panic function for exiting the execution and printing error

**Note that flags and operations names should be the same as mentioned above or unit tests will never pass.**
