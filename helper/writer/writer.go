package writer

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"sync"
	"time"

	"github.com/fahmiaz411/bookla/helper/constant"
)

var PayloadActivity = map[int64]string{}

var wg = &sync.WaitGroup{}

func StartScheduling() {
	
	if len(PayloadActivity) == constant.ZeroValue {
		time.Sleep(5 * time.Millisecond)
		StartScheduling()
		return
	}
	fmt.Println(len(PayloadActivity))
	wg.Add(len(PayloadActivity))

	for k, v := range PayloadActivity {
		key := k
		val := v
		delete(PayloadActivity, key)

		go func(){
			defer wg.Done()
			ioutil.WriteFile("cache/activities/" + fmt.Sprint(key) + ".json", []byte(val), fs.ModePerm)
		}()		
	}

	wg.Wait()

	StartScheduling()

	return
}