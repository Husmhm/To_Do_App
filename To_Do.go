package main

import (
	"bufio"
	"flag"
	"fmt"
	"go.dev/constant"
	"go.dev/contract"
	"go.dev/repository/filestore"
	"go.dev/repository/memorystore"
	"go.dev/service/task"
	"os"

	// "path/filepath"

	// "os/user"
	"crypto/md5"
	"encoding/hex"
	"go.dev/entity"
	"strconv"
)

var (
	userstorage       []entity.User
	authenticatedUser *entity.User

	taskstorage     []entity.Task
	categorystorage []entity.Category

	serializationMode string
)

const (
	userstoragepath = "user.txt"
)

func main() {
	taskMemoryRepo := memorystore.NewTaskStore()
	taskService := task.NewService(taskMemoryRepo)

	fmt.Println("Hello to do app")
	sm := flag.String("serialization-mode", constant.Json_serializationMode, "serialization mode to write file")
	switch *sm {
	case constant.ManDarAvardi_serializationMode:
		serializationMode = constant.ManDarAvardi_serializationMode
	default:
		serializationMode = constant.Json_serializationMode
	}
	var userFilestor = filestore.New(userstoragepath, serializationMode)

	command := flag.String("command", "no-command", "command to run")
	flag.Parse()

	// loaduser storage from file
	// loadUserStorageFromFile(serializationMode)

	users := userFilestor.Load()
	userstorage = append(userstorage, users...)

	runcommand(userFilestor, *command)
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("please enter another command")
		scanner.Scan()
		*command = scanner.Text()
		runcommand(userFilestor, *command)
	}
}

func runcommand(store contract.UserWriteStore, command string) {
	if command != "register-user" && command != "exit" && command != "login" && authenticatedUser == nil {
		login()
		if authenticatedUser == nil {
			return
		}
	}

	// var store userWriteStore
	// store = filestore{}

	switch command {
	case "create-task":
		createTask()
	case "create-category":
		createCategory()
	case "register-user":
		registerUser(filestore.Filestore{
			Filepath: userstoragepath,
		})
	case "exit":
		os.Exit(0)
	case "login":
		login()
	case "list-task":
		//listTask()
	default:
		fmt.Println("no-command")
	}

}
func createTask() {
	scanner := bufio.NewScanner(os.Stdin)
	var title, catagory, dodate string

	fmt.Println("please enter the task title")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("please enter the task category ID")
	scanner.Scan()
	catagory = scanner.Text()
	catagory_id, err := strconv.Atoi(catagory)
	if err != nil {
		fmt.Printf("category-id is not valid integer ,%v\n", err)
		return
	}

	fmt.Println("please enter the task dodate")
	scanner.Scan()
	dodate = scanner.Text()

}
func createCategory() {
	scanner := bufio.NewScanner(os.Stdin)
	var title, color string

	fmt.Println("please enter the category title")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("please enter the category color")
	scanner.Scan()
	color = scanner.Text()

	category := Category{
		ID:     len(categorystorage) + 1,
		Title:  title,
		Color:  color,
		UserID: authenticatedUser.ID,
	}
	categorystorage = append(categorystorage, category)

}

func registerUser(store contract.UserWriteStore) {
	scanner := bufio.NewScanner(os.Stdin)
	var id, name, email, password string

	fmt.Println("please enter your Name")
	scanner.Scan()
	name = scanner.Text()

	fmt.Println("please enter your email")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("please enter your password")
	scanner.Scan()
	password = scanner.Text()
	id = email
	fmt.Println("user:", id, email, password)

	user := entity.User{
		ID:       len(userstorage) + 1,
		Name:     name,
		Email:    email,
		Password: hashThePassword(password),
	}

	userstorage = append(userstorage, user)
	fmt.Printf("userstorage: %v\n", userstorage)

	store.Save(user)
}

func login() {
	fmt.Println("*****login process*****")
	scanner := bufio.NewScanner(os.Stdin)
	var email, password string

	fmt.Println("please enter your email")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("please enter your password")
	scanner.Scan()
	password = scanner.Text()

	for _, user := range userstorage {
		if user.Email == email && user.Password == hashThePassword(password) {
			fmt.Println("you are login")
			authenticatedUser = &user

			break
		}
	}
	if authenticatedUser == nil {
		fmt.Println("email or password incorrect")
	}

}

func hashThePassword(password string) string {
	hash := md5.Sum([]byte(password))

	return hex.EncodeToString(hash[:])
}
