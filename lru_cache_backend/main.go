package main

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type KeyPair struct {
	Key       int       `json:"key"`
	Value     int       `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

type lruCache struct {
	capacity int
	list     *list.List
	elements map[int]*list.Element
	mu       sync.Mutex
}

func (cache *lruCache) get(key int) (KeyPair, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	if node, ok := cache.elements[key]; ok {
		value := node.Value.(*list.Element).Value.(KeyPair)
		value.Timestamp = time.Now()
		node.Value.(*list.Element).Value = value
		cache.list.MoveToFront(node)
		return value, nil
	}
	return KeyPair{}, errors.New("key not found")
}

func (cache *lruCache) set(val *KeyPair) {
	defer cache.mu.Unlock()
	cache.mu.Lock()

	if node, ok := cache.elements[val.Key]; ok {
		cache.list.MoveToFront(node)
		node.Value.(*list.Element).Value = KeyPair{
			Key:       val.Key,
			Value:     val.Value,
			Timestamp: time.Now(),
		}
		return
	} else {
		if cache.capacity == cache.list.Len() {
			idx := cache.list.Back().Value.(*list.Element).Value.(KeyPair).Key
			delete(cache.elements, idx)
			cache.list.Remove(cache.list.Back())
			return
		}
	}

	node := &list.Element{
		Value: KeyPair{
			Key:       val.Key,
			Value:     val.Value,
			Timestamp: time.Now(),
		},
	}

	pointer := cache.list.PushFront(node)
	cache.elements[val.Key] = pointer

}

// func (cache *lruCache) remove(key int) {
// 	cache.mu.Lock()
// 	defer cache.mu.Unlock()
// 	if node, ok := cache.elements[key]; ok {
// 		delete(cache.elements, key)
// 		cache.list.Remove(node)
// 	}
// }
// func (cache *lruCache) top() interface{} {
// 	if cache.list.Len() != 0 {
// 		return cache.list.Front().Value.(*list.Element).Value.(KeyPair).Value
// 	} else {
// 		fmt.Println("cache Empty")
// 		return -1
// 	}
// }

func (cache *lruCache) cleanEvery5Sec() {

	cache.mu.Lock()
	defer cache.mu.Unlock()
	for element := cache.list.Back(); element != nil; element = element.Prev() {
		if time.Since(element.Value.(*list.Element).Value.(KeyPair).Timestamp) > 3*time.Second {
			fmt.Printf("time since %d", time.Since(element.Value.(*list.Element).Value.(KeyPair).Timestamp)*time.Second)
			delete(cache.elements, element.Value.(*list.Element).Value.(KeyPair).Key)
			cache.list.Remove(element)
		} else {
			fmt.Printf("time since %d", time.Since(element.Value.(*list.Element).Value.(KeyPair).Timestamp))

		}
	}
}

func (cache *lruCache) GetAll() []KeyPair {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	keysAndValues := make([]KeyPair, 0, cache.capacity) // Preallocate for efficiency
	for element := cache.list.Front(); element != nil; element = element.Next() {
		value := element.Value.(*list.Element).Value.(KeyPair)
		keysAndValues = append(keysAndValues, value)
	}
	return keysAndValues
}

func main() {
	var wg sync.WaitGroup
	cache := &lruCache{
		capacity: 1024,
		list:     list.New(),
		elements: make(map[int]*list.Element),
	}
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "lruapi",
		AppName:       "Test App v1.0.1",
	})
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: false,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	api := app.Group("/api")

	api.Get("/all/", func(c *fiber.Ctx) error {
		out := cache.GetAll()
		return c.JSON(out)

	})

	api.Get("/:key", func(c *fiber.Ctx) error {
		fmt.Println("inside get")
		key_s := c.Params("key")
		key, err := strconv.Atoi(key_s)
		if err != nil {
			fmt.Println("Error:", err)
			return errors.New("invalid key")
		}

		resultCh := make(chan KeyPair)
		errCh := make(chan error)
		wg.Add(1)
		go func() {
			defer wg.Done()
			val, err := cache.get(key)

			if err != nil {
				errCh <- err
			}
			resultCh <- val
		}()

		select {
		case value := <-resultCh:
			fmt.Println(value)
			jsonData, err := json.Marshal(value)
			if err != nil {
				log.Printf("Error marshalling value to JSON: %s", err.Error())
				return err
			}
			fmt.Println(string(jsonData))
			var out KeyPair
			json.Unmarshal(jsonData, &out)
			return c.JSON(out)
		case err := <-errCh:
			log.Printf("Error fetching value for key %d: %s", key, err.Error())
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching value for key ")
		}

	})

	api.Post("/", func(c *fiber.Ctx) error {
		fmt.Println("inside post")
		input := new(KeyPair)

		if err := c.BodyParser(input); err != nil {
			return err
		}

		cache.set(input)

		fmt.Println("new cache added")
		return c.Status(fiber.StatusAccepted).JSON(
			map[string]interface{}{
				"code": "20",
				"msg":  "new cache added",
				"err":  nil,
			},
		)

	})

	ticker := time.NewTicker(5 * time.Second)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range ticker.C {
			cache.cleanEvery5Sec()
		}
	}()

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	fmt.Println("Top value in the cache:", cache.top())
	// }()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := app.Listen(":5000"); err != nil {
			fmt.Println("Had error")
			panic(err.Error())
		}
	}()
	wg.Wait()
}
