package encrypt

import (
	"fmt"
	"math/rand"
)

func RandInt(length int) (generatedNum string) {
	for i := 0; i < length; i++ {		
		generatedNum = generatedNum + fmt.Sprint(rand.Intn(10))		
	}
	return
}