//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer szenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"time"
)

func producer(stream Stream) (tweets []*Tweet) {
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			return tweets
		}

		tweets = append(tweets, tweet)
	}
}

func consumer(tweets []*Tweet) {
	for _, t := range tweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func asyncProducer(stream Stream) <-chan *Tweet {
	out := make(chan *Tweet)

	go func() {
		defer close(out)

		for {
			tweet, err := stream.Next()
			if err == ErrEOF {
				return
			}

			out <- tweet
		}
	}()

	return out
}

func asyncConsumer(in <-chan *Tweet, done chan struct{}) {
	for t := range in {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}

	done <- struct{}{}
}

func main() {
	var (
		start  = time.Now()
		stream = GetMockStream()
		in     = asyncProducer(stream)
		done   = make(chan struct{})
	)

	go asyncConsumer(in, done)
	<-done

	fmt.Printf("Process took %s\n", time.Since(start))
}
