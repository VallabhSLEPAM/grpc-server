package application

import (
	"fmt"
	"math/rand/v2"
	"time"
)

type ResiliencyService struct {
}

func (r *ResiliencyService) GenerateResiliency(minDelay, maxDelay int, statusCode []uint32) (string, uint32) {
	// generate random delay between min and max delays
	delay := rand.IntN(int(maxDelay-minDelay+1)) + int(minDelay)
	delayInSecond := time.Duration(delay) * time.Second
	time.Sleep(delayInSecond)

	idx := rand.IntN(len(statusCode))
	str := fmt.Sprintf("The time now is %v, execution is delayed by %v seconds", time.Now().Format("10:00:00.000"), delay)
	return str, statusCode[idx]
}
