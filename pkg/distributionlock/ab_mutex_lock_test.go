package distributionlock

/*import (
	"fmt"
	"mtcloud.com/mtstorage/pkg/journal"
	"sync"
	"testing"
	"time"
)

///////////////////////test cod ///////////
var data = 0
var wgAdd sync.WaitGroup
var wgDec sync.WaitGroup

var metexKey = "metexKey"
var shareAddKey = "shareAddKey"
var shareDecKey = "shareDecKey"

func Inc() {
	defer wgAdd.Done()
	journal.JournalObject.Lock(metexKey, shareAddKey, journal.AddExpire, journal.AddShareExpire)
	data++
	fmt.Printf("inc data: %d \n", data)
	time.Sleep(1 * time.Second)
}

func Dec() {
	defer wgDec.Done()
	for true {
		if err := journal.JournalObject.Lock(metexKey, shareDecKey, journal.DeleteExpire, journal.DeleteShareExpire); err == nil {
			break
		}
	}
	data--
	fmt.Printf("dec data: %d \n", data)
	time.Sleep(1 * time.Second)
}

var count = 10

func DistributeLock() {
	journal.Init()
	wgAdd.Add(count)
	wgDec.Add(count)

	data = 0
	//让添加获取lock.
	Inc()
	for i := 0; i < count-1; i++ {
		go Inc()
	}

	for i := 0; i < count; i++ {
		go Dec()
	}

	wgAdd.Wait()
	fmt.Printf("add finish \n")
	time.Sleep(30 * time.Second)

	fmt.Printf("add begin unlock \n")
	for i := 0; i < count; i++ {
		journal.JournalObject.UnLock(metexKey, shareAddKey)
	}
	fmt.Printf("add unlock finish \n")

	wgDec.Wait()
	fmt.Printf("dec begin unlock \n")
	for i := 0; i < count; i++ {
		journal.JournalObject.UnLock(metexKey, shareDecKey)
	}
	fmt.Printf("dec unlock finish \n")
}

func Inc_1() {
	defer wgAdd.Done()

	for true {
		if err := journal.JournalObject.Lock(metexKey, shareAddKey, journal.AddExpire, journal.AddShareExpire); err == nil {
			break
		}
	}

	defer journal.JournalObject.UnLock(metexKey, shareAddKey)
	data++
	fmt.Printf("inc data: %d \n", data)
	time.Sleep(1 * time.Second)
}

func Dec_1() {
	defer wgDec.Done()
	for true {
		if err := journal.JournalObject.Lock(metexKey, shareDecKey, journal.DeleteExpire, journal.DeleteShareExpire); err == nil {
			break
		}
	}
	defer journal.JournalObject.UnLock(metexKey, shareDecKey)

	data--
	fmt.Printf("dec data: %d \n", data)
	time.Sleep(1 * time.Second)
}

func DistributeLock_1() {
	journal.Init()
	wgAdd.Add(count)
	wgDec.Add(count)

	data = 0
	//让添加获取lock.
	for i := 0; i < count; i++ {
		go Inc_1()
	}

	for i := 0; i < count; i++ {
		go Dec_1()
	}

	wgAdd.Wait()
	wgDec.Wait()

	fmt.Printf("finish data: %d \n", data)

}

//////////////////////////////

func TestLock(t *testing.T) {
	DistributeLock()
	//DistributeLock_1()
}*/
