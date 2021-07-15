package wind

import (
	"fmt"
	"time"
)

var ctrlCh chan map[string]interface{}
var dir int
var callbacks []func(val int)

func init() {
	ticker := time.NewTicker(1 * time.Second)
	ctrlCh = make(chan map[string]interface{}, 4096)
	callbacks = make([]func(val int), 0)
	go func() {
		for {
			select {
			case req := <-ctrlCh:
				command := req["command"].(string)
				switch command {
				case "get":
					sendResponseWithTimeout(req["output"].(chan interface{}), dir, 1*time.Second)
				case "set":
					dir = req["value"].(int)
					for _, callback := range callbacks {
						callback(dir)
					}
				case "subscribe":
					callbacks = append(callbacks, req["callback"].(func(val int)))
				}
			case <-ticker.C:
				if dir > 0 {
					dir--
				} else if dir < 0 {
					dir++
				} else {
					break
				}
				for _, callback := range callbacks {
					callback(dir)
				}
			}
		}
	}()
}

func sendResponseWithTimeout(out chan interface{}, response interface{}, timeout time.Duration) {
	select {
	case out <- response:
	case <-time.After(timeout):
		fmt.Println("timeout")
	}
}

func GetDirection() int {
	readResponseFrom := make(chan interface{}, 1)
	select {
	case ctrlCh <- map[string]interface{}{
		"command": "get",
		"output":  readResponseFrom,
	}:
	case <-time.After(1 * time.Second):
		fmt.Println("timeout")
		return 0
	}

	select {
	case v := <-readResponseFrom:
		return v.(int)
	case <-time.After(1 * time.Second):
		fmt.Println("timeout")
		return 0
	}
}

func SetDirection(val int) {
	select {
	case ctrlCh <- map[string]interface{}{
		"command": "set",
		"value":   val,
	}:
	case <-time.After(1 * time.Second):
		fmt.Println("timeout")
	}
}

func Subscribe(callback func(val int)) {
	select {
	case ctrlCh <- map[string]interface{}{
		"command":  "subscribe",
		"callback": callback,
	}:
	case <-time.After(1 * time.Second):
		fmt.Println("timeout")
	}
}
