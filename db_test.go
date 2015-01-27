package db

import (
	"testing"

	// For testing the storage of bcrypt password hashes
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/sha256"
	"github.com/xyproto/permissions2"
	"io"
)

const (
	listname     = "testlist"
	setname      = "testset"
	hashmapname  = "testhashmap"
	keyvaluename = "testkeyvalue"
	testdata1    = "abc123"
	testdata2    = "def456"
	testdata3    = "ghi789"
)

func TestLocalConnection(t *testing.T) {
	Verbose = true

	//err := TestConnection() // locally
	err := TestConnectionHost("travis:@127.0.0.1/") // for travis-ci
	//err := TestConnectionHost("go:go@/main") // laptop
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestList(t *testing.T) {
	Verbose = true

	//host := New() // locally
	host := NewHost("travis:@127.0.0.1/") // for travis-ci
	//host := NewHost("go:go@/main") // laptop

	defer host.Close()
	list := NewList(host, listname)
	list.Clear()
	if err := list.Add(testdata1); err != nil {
		t.Errorf("Error, could not add item to list! %s", err.Error())
	}
	items, err := list.GetAll()
	if err != nil {
		t.Errorf("Error when retrieving list! %s", err.Error())
	}
	if len(items) != 1 {
		t.Errorf("Error, wrong list length! %v", len(items))
	}
	if (len(items) > 0) && (items[0] != testdata1) {
		t.Errorf("Error, wrong list contents! %v", items)
	}
	if err := list.Add(testdata2); err != nil {
		t.Errorf("Error, could not add item to list! %s", err.Error())
	}
	if err := list.Add(testdata3); err != nil {
		t.Errorf("Error, could not add item to list! %s", err.Error())
	}
	items, err = list.GetAll()
	if err != nil {
		t.Errorf("Error when retrieving list! %s", err.Error())
	}
	if len(items) != 3 {
		t.Errorf("Error, wrong list length! %v", len(items))
	}
	item, err := list.GetLast()
	if err != nil {
		t.Errorf("Error, could not get last item from list! %s", err.Error())
	}
	if item != testdata3 {
		t.Errorf("Error, expected %s, got %s with GetLast()!", testdata3, item)
	}
	items, err = list.GetLastN(2)
	if err != nil {
		t.Errorf("Error, could not get last N items from list! %s", err.Error())
	}
	if len(items) != 2 {
		t.Errorf("Error, wrong list length! %v", len(items))
	}
	if items[0] != testdata2 {
		t.Errorf("Error, expected %s, got %s with GetLast()!", testdata2, items[0])
	}
	err = list.Remove()
	if err != nil {
		t.Errorf("Error, could not remove list! %s", err.Error())
	}

	// Check that list qualifies for the IList interface
	var _ IList = list
}

func TestSet(t *testing.T) {
	Verbose = true

	//host := New() // locally
	host := NewHost("travis:@127.0.0.1/") // for travis-ci
	//host := NewHost("go:go@/main") // laptop

	defer host.Close()
	set := NewSet(host, setname)
	set.Clear()
	if err := set.Add(testdata1); err != nil {
		t.Errorf("Error, could not add item to set! %s", err.Error())
	}
	items, err := set.GetAll()
	if err != nil {
		t.Errorf("Error when retrieving set! %s", err.Error())
	}
	if len(items) != 1 {
		t.Errorf("Error, wrong set length! %v", len(items))
	}
	if (len(items) > 0) && (items[0] != testdata1) {
		t.Errorf("Error, wrong set contents! %v", items)
	}
	if err := set.Add(testdata2); err != nil {
		t.Errorf("Error, could not add item to set! %s", err.Error())
	}
	if err := set.Add(testdata3); err != nil {
		t.Errorf("Error, could not add item to set! %s", err.Error())
	}
	// Add an element twice. This is a set, so the element should only appear once.
	if err := set.Add(testdata3); err != nil {
		t.Errorf("Error, could not add item to set! %s", err.Error())
	}
	items, err = set.GetAll()
	if err != nil {
		t.Errorf("Error when retrieving set! %s", err.Error())
	}
	if len(items) != 3 {
		t.Errorf("Error, wrong set length! %v\n%v\n", len(items), items)
	}
	err = set.Remove()
	if err != nil {
		t.Errorf("Error, could not remove set! %s", err.Error())
	}

	// Check that set qualifies for the ISet interface
	var _ ISet = set
}

func TestHashMap(t *testing.T) {
	Verbose = true

	//host := New() // locally
	host := NewHost("travis:@127.0.0.1/") // for travis-ci
	//host := NewHost("go:go@/main") // laptop

	defer host.Close()
	hashmap := NewHashMap(host, hashmapname)
	hashmap.Clear()

	username := "bob"
	key := "password"
	value := "hunter1"

	// Get key that doesn't exist yet
	item, err := hashmap.Get("ownerblabla", "keyblabla")
	if err == nil {
		t.Errorf("Key found, when it should be missing! %s", err.Error())
	}

	if err := hashmap.Set(username, key, value); err != nil {
		t.Errorf("Error, could not set value in hashmap! %s", err.Error())
	}
	// Once more, with the same data
	if err := hashmap.Set(username, key, value); err != nil {
		t.Errorf("Error, could not set value in hashmap! %s", err.Error())
	}
	items, err := hashmap.GetAll()
	if err != nil {
		t.Errorf("Error when retrieving elements! %s", err.Error())
	}
	if len(items) != 1 {
		t.Errorf("Error, wrong element length! %v", len(items))
	}
	if (len(items) > 0) && (items[0] != username) {
		t.Errorf("Error, wrong elementid! %v", items)
	}
	item, err = hashmap.Get(username, key)
	if err != nil {
		t.Errorf("Error, could not fetch value from hashmap! %s", err.Error())
	}
	if item != value {
		t.Errorf("Error, expected %s, got %s!", value, item)
	}
	err = hashmap.Remove()
	if err != nil {
		t.Errorf("Error, could not remove hashmap! %s", err.Error())
	}

	// Check that hashmap qualifies for the IHashMap interface
	var _ IHashMap = hashmap
}

func TestKeyValue(t *testing.T) {
	Verbose = true

	//host := New() // locally
	host := NewHost("travis:@127.0.0.1/") // for travis-ci
	//host := NewHost("go:go@/main") // laptop

	defer host.Close()
	keyvalue := NewKeyValue(host, keyvaluename)
	keyvalue.Clear()

	key := "password"
	value := "hunter1"

	if err := keyvalue.Set(key, value); err != nil {
		t.Errorf("Error, could not set value in keyvalue! %s", err.Error())
	}
	// Twice
	if err := keyvalue.Set(key, value); err != nil {
		t.Errorf("Error, could not set value in keyvalue! %s", err.Error())
	}
	item, err := keyvalue.Get(key)
	if err != nil {
		t.Errorf("Error, could not fetch value from keyvalue! %s", err.Error())
	}
	if item != value {
		t.Errorf("Error, expected %s, got %s!", value, item)
	}
	err = keyvalue.Remove()
	if err != nil {
		t.Errorf("Error, could not remove keyvalue! %s", err.Error())
	}
	// Check that keyvalue qualifies for the IKeyValue interface
	var _ IKeyValue = keyvalue
}

func TestHashStorage(t *testing.T) {
	Verbose = true

	//host := New() // locally
	host := NewHost("travis:@127.0.0.1/") // for travis-ci
	//host := NewHost("go:go@/main") // laptop

	defer host.Close()
	hashmap := NewHashMap(host, hashmapname)
	hashmap.Clear()

	username := "bob"
	key := "password"
	password := "hunter1"

	// bcrypt test

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	value := string(passwordHash)

	if err := hashmap.Set(username, key, value); err != nil {
		t.Errorf("Error, could not set value in hashmap! %s", err.Error())
	}
	item, err := hashmap.Get(username, key)
	if err != nil {
		t.Errorf("Unable to retrieve from hashmap! %s\n", err.Error())
	}
	if item != value {
		t.Errorf("Error, got a different value back (bcrypt)! %s != %s\n", value, item)
	}

	// sha256 test

	hasher := sha256.New()
	io.WriteString(hasher, password+permissions.RandomCookieFriendlyString(30)+username)
	passwordHash = hasher.Sum(nil)
	value = string(passwordHash)

	if err := hashmap.Set(username, key, value); err != nil {
		t.Errorf("Error, could not set value in hashmap! %s", err.Error())
	}
	item, err = hashmap.Get(username, key)
	if err != nil {
		t.Errorf("Unable to retrieve from hashmap! %s\n", err.Error())
	}
	if item != value {
		t.Errorf("Error, got a different value back (sha256)! %s != %s\n", value, item)
	}

	err = hashmap.Remove()
	if err != nil {
		t.Errorf("Error, could not remove hashmap! %s", err.Error())
	}
}

func TestTwoFields(t *testing.T) {
	test, test23, ok := twoFields("test1@test2@test3", "@")
	if ok && ((test != "test1") || (test23 != "test2@test3")) {
		t.Error("Error in twoFields functions")
	}
}
