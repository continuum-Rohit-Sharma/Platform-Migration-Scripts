package main

import "log"

//"github.com/ContinuumLLC/platform-common-lib/src/testapps/kafka/Test"

func main() {
	// var wait sync.WaitGroup
	// wait.Add(1)
	//timer := time.NewTimer(time.Second * 2)
	// f, _ := os.OpenFile("Confluent", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//log.SetOutput(f)
	log.Println("Start Job")
	// defer f.Close()
	func() {
		go testSaramaConsumerPullHandlerLoop()
	}()
	testSaramaProducerPushLoop("message#-")
	//	<-timer.C
	//time.Sleep(time.Duration(1) * time.Second)

	//		i++
	//	}
	// wait.Wait()

	// testProducer()
	// testConsumer()
}
func testProducer() {
	var result bool

	result = testNewSaramaProducer()
	if result == false {
		log.Println("Cannot create new Producer")
		return
	}
	log.Println("Producer created sucessfully")

	result = testNewSaramaProducerNilBrokerAddress()
	if result == false {
		log.Println("Producer Validation test Failed")
		return
	}
	log.Println("Producer Validated Successfully")

	result = testNewSaramaProducerEmptyBrokerAddress()
	if result == false {
		log.Println("Producer Validation test Failed")
		return
	}
	log.Println("Producer Validated Successfully")

	/*	result = TestSaramaProducerConnect()
		if result == false {
			log.Println("Cannot connect to Producer")
			return
		}

		log.Println("Producer connected sucessfully")

		result = TestSaramaProducerIsConnected()
		if result == false {
			log.Println("Cannot check connection status for Producer")
			return
		}

		log.Println("Producer connection checked sucessfully")*/

	result = testSaramaProducerPush()

	if result == false {
		log.Println("Cannot Push data to Producer")
		return
	}

	log.Println("Producer data pushed sucessfully")
	log.Println("All producer tests executed sucessfully")

}

func testConsumer() {
	var result bool

	result = testNewSaramaConsumer()
	if result == false {
		log.Println("Cannot create new Consumer")
		return
	}

	log.Println("Consumer created sucessfully")

	result = testNewSaramaConsumerNilBrokerAddress()
	if result == false {
		log.Println("Consumer Validation test failed")
		return
	}

	log.Println("Consumer validated sucessfully")

	result = testNewSaramaConsumerEmptyBrokerAddress()
	if result == false {
		log.Println("Consumer Validation test failed")
		return
	}

	log.Println("Consumer validated sucessfully")

	result = testNewSaramaConsumerEmptyGroupID()
	if result == false {
		log.Println("Consumer Validation test failed")
		return
	}
	log.Println("Consumer validated sucessfully")
	/*result = TestSaramaConsumerConnect()
	if result == false {
		log.Println("Cannot connect to Consumer")
		return
	}

	log.Println("Consumer connected sucessfully")

	result = TestSaramaConsumerIsConnected()
	if result == false {
		log.Println("Cannot check connection status for Consumer")
		return
	}

	log.Println("Consumer connection checked sucessfully")

	result = TestSaramaConsumerPull()

	if result == false {
		log.Println("Cannot pull data to Consumer")
		return
	}
	log.Println("Consumer data pulled sucessfully")*/

	result = testSaramaConsumerPullHandler()

	if result == false {
		log.Println("Cannot pull data to Consumer")
		return
	}
	log.Println("ConsumerHandler data pulled sucessfully")

	log.Println("All tests executed sucessfully")

}
