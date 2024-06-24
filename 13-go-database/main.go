package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
)

const Version = "1.0.1"
const Directory = "./"
const DatabaseName = "users"
const ACCESS_PERMISSION_CODE = 0755
const WRITE_PERMISSION_CODE = 0644

type (
	Logger interface {
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Warn(string, ...interface{})
		Info(string, ...interface{})
		Debug(string, ...interface{})
		Trace(string, ...interface{})
	}
	Driver struct {
		mutex   sync.Mutex
		mutexes map[string]*sync.Mutex
		dir     string
		log     Logger
	}
)

type Options struct {
	Logger
}

func NewContext(dir string, options *Options) (*Driver, error) {
	dir = filepath.Clean(dir)
	opts := Options{}
	if options != nil {
		opts = *options
	}
	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger((lumber.INFO))
	}

	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
		log:     opts.Logger,
	}
	if _, err := os.Stat(dir); err != nil {
		opts.Logger.Debug("Using database '%s', Database already exists.\n", dir)
		return &driver, nil
	}
	opts.Logger.Debug("Creating database at directory: '%s'\n", dir)
	return &driver, os.MkdirAll(dir, ACCESS_PERMISSION_CODE)
}

func (driver *Driver) Write(collection, resource string, v interface{}) error {
	if collection == "" {
		return fmt.Errorf("missing collection - no place to save records")
	}
	if resource == "" {
		return fmt.Errorf("missing resource - no name - unable to save resource")
	}

	mutex := driver.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(driver.dir, collection)
	finalPath := filepath.Join(dir, resource+".json")
	tempPath := finalPath + ".tmp"
	if err := os.MkdirAll(dir, ACCESS_PERMISSION_CODE); err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	bytes = append(bytes, byte('\n'))
	if err := os.WriteFile(tempPath, bytes, WRITE_PERMISSION_CODE); err != nil {
		return err
	}

	return os.Rename(tempPath, finalPath)
}

func (driver *Driver) Read(collection, resource string, v interface{}) error {
	if collection == "" {
		return fmt.Errorf("missing collection - unable to read")
	}
	if resource == "" {
		return fmt.Errorf("missing resource - no name - unable to read resource")
	}

	record := filepath.Join(driver.dir, collection, resource)
	if _, err := stat(record); err != nil {
		return err
	}

	bytes, err := os.ReadFile(record + ".json")
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, &v)
}

func (driver *Driver) ReadAll(collection string) ([]string, error) {
	if collection == "" {
		return nil, fmt.Errorf("missing collection - unable to read")
	}

	dir := filepath.Join(driver.dir, collection)
	if _, err := stat(dir); err != nil {
		return nil, err
	}

	files, _ := os.ReadDir(dir)
	var records []string
	for _, file := range files {
		bytes, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		records = append(records, string(bytes))
	}
	return records, nil
}

func (driver *Driver) Delete(collection, resource string) error {
	path := filepath.Join(collection, resource)
	mutex := driver.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(driver.dir, path)
	switch fi, err := stat(dir); {
	case fi == nil, err != nil:
		return fmt.Errorf("unable to find file: %v", path)
	case fi.Mode().IsDir():
		return os.RemoveAll(dir)
	case fi.Mode().IsRegular():
		return os.RemoveAll(dir + ".json")
	}

	return nil
}

func (driver *Driver) getOrCreateMutex(collection string) *sync.Mutex {
	driver.mutex.Lock()
	defer driver.mutex.Unlock()

	mutex, exists := driver.mutexes[collection]
	if !exists {
		mutex = &sync.Mutex{}
		driver.mutexes[collection] = mutex
	}
	return mutex
}

func stat(path string) (fileInfo os.FileInfo, err error) {
	fileInfo, err = os.Stat(path)
	if err == nil {
		return fileInfo, nil // Return fileInfo if file exists
	}
	if os.IsNotExist(err) {
		// Check for file with .json extension if it does not exist
		return os.Stat(path + ".json")
	}
	return nil, err // Return original error if not a NotExist error
}

type Address struct {
	City     string
	Country  string
	AreaCode json.Number
}

type User struct {
	Name    string
	Age     json.Number
	Contact string
	Company string
	Address Address
}

func main() {
	dir := Directory

	db, err := NewContext(dir, nil)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	employees := []User{
		{Name: "John", Age: "23", Contact: "12312312", Company: "Canva", Address: Address{City: "Sydney", Country: "Australia", AreaCode: "123412"}},
		{Name: "Tim", Age: "33", Contact: "12341", Company: "Atlassian", Address: Address{City: "Auckland", Country: "New Zealand", AreaCode: "1243213"}},
		{Name: "Bob", Age: "13", Contact: "45354", Company: "Canva", Address: Address{City: "Melbourne", Country: "Australia", AreaCode: "4309"}},
		{Name: "Jimmy", Age: "43", Contact: "653434564", Company: "Spark", Address: Address{City: "Auckland", Country: "New Zealand", AreaCode: "43214"}},
		{Name: "Tommy", Age: "53", Contact: "3456345", Company: "Canva", Address: Address{City: "Auckland", Country: "New Zealand", AreaCode: "2341243"}},
		{Name: "Joe", Age: "22", Contact: "12321", Company: "Octopus", Address: Address{City: "Auckland", Country: "New Zealand", AreaCode: "123"}},
	}

	for _, value := range employees {
		err := db.Write(DatabaseName, value.Name, User{
			Name: value.Name, Age: value.Age, Contact: value.Contact, Company: value.Company, Address: value.Address,
		})
		if err != nil {
			fmt.Println(err)
		}
	}

	records, err := db.ReadAll(DatabaseName)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println(records)

	allUsers := []User{}
	for _, user := range records {
		employee := User{}
		if err := json.Unmarshal([]byte(user), &employee); err != nil {
			fmt.Println("Error: ", err)
		}
		allUsers = append(allUsers, employee)
	}
	fmt.Println(allUsers)

	// if err := db.Delete(DatabaseName, "John"); err != nil {
	// 	fmt.Println("Error: ", err)
	// }
	// if err := db.Delete(DatabaseName, ""); err != nil {
	// 	fmt.Println("Error: ", err)
	// }
}
