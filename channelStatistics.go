package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	chans := []chan int{
		make(chan int, 1000),
		make(chan int, 1000),
		make(chan int, 1000),
		make(chan int, 1000),
	}
	go worker(chans[0], time.Duration(30*time.Millisecond), false, 0)
	go worker(chans[1], time.Duration(300*time.Millisecond), false, 0)
	go worker(chans[2], time.Duration(300*time.Millisecond), false, 0)
	go worker(chans[3], time.Duration(30*time.Millisecond), false, 0)
	throttle := time.Tick(time.Millisecond)
	roundRobinIndex := 0
	for {
		<-throttle
		c := chans[roundRobinIndex]
		shouldContinue := true
		for shouldContinue {
			rand.Seed(time.Now().Unix())
			select {
			case c <- 1:
				shouldContinue = false
			default:
				difference := float64(len(c)) - findMostEmpty(chans)
				if difference/float64(len(c)) > float64(0.8) {
					//drop the message, this queue is out of sync
					fmt.Printf("dropping message on channel %v\n", roundRobinIndex)
					shouldContinue = false
					//for now, just panic
					panic("asdf")
				} else {
					// wait a little and retry
					fmt.Printf("%v, %v, %v, %v\n", len(chans[0]), len(chans[1]), len(chans[2]), len(chans[3]))
					<-time.After(500 * time.Millisecond)
				}

				//mean, standardDeviation := findMeanAndStandardDeviation(chans)
				//fmt.Printf("standard deviation %v\n", standardDeviation)
				//if float64(len(c)) > standardDeviation*2+mean {
				//	//drop the message, this queue is out of sync
				//	fmt.Printf("dropping message, length of channel %v is %v, more than two standard deviations (%v) from the mean (%v)\n", roundRobinIndex, len(c), standardDeviation, mean)
				//	//shouldContinue = false
				//	//for now, just panic
				//	panic("asdf")
				//} else {
				//	// wait a little and retry
				//	fmt.Printf("channel %v has message length %v which is less than two standard deviations (%v) from the mean (%v), sleeping and trying again\n", roundRobinIndex, len(c), standardDeviation, mean)
				//	fmt.Printf("%v, %v, %v, %v\n", len(chans[0]), len(chans[1]), len(chans[2]), len(chans[3]))
				//	<-time.After(500 * time.Millisecond)
				//}
			}
		}
		fmt.Printf("length of channel %v is %v\n", roundRobinIndex, len(c))
		roundRobinIndex = (roundRobinIndex + 1) % len(chans)
	}
}
func findMostEmpty(chans []chan int) float64 {
	lowest := math.MaxFloat64
	for _, c := range chans {
		if lowest > float64(len(c)) {
			lowest = float64(len(c))
		}
	}
	return lowest
}
func findMedian(chans []chan int) float64 {
	values := []float64{}
	for _, c := range chans {
		values = append(values, float64(len(c)))
	}
	if len(values)%2 == 0 {
		middle := len(values) / 2
		total := values[middle-1] + values[middle]
		return total / float64(2)
	} else {
		middle := len(values) / 2
		return values[middle+1]
	}
}
func findMeanAndStandardDeviation(chans []chan int) (float64, float64) {
	mean := float64(0)
	for j := 0; j < len(chans); j++ {
		mean += float64(len(chans[j]))
	}
	mean /= float64(len(chans))
	variance := float64(0)
	for j := 0; j < len(chans); j++ {
		variance += math.Pow(float64(len(chans[j]))-mean, 2)
	}
	standardDeviation := math.Sqrt(variance / float64(len(chans)))
	return mean, standardDeviation
}

func worker(c chan int, sleepTime time.Duration, hang bool, id int) {
	for {
		<-c
		if hang {
			select {}
		}
		myOffset := rand.Int63n(int64(sleepTime))
		if rand.Intn(1)%2 == 0 {
			myOffset *= -1
		}
		mySleepTime := time.Duration(int64(sleepTime) + myOffset)
		<-time.After(mySleepTime)
	}
}
