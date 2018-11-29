package main

import (
	"fmt"

	"github.com/ContinuumLLC/platform-common-lib/src/webClient"
)

func main() {

	TestWebClientInvalidHTTPMethod()
	TestWebClientEmptyContentType()
	TestWebClientNilData()
	TestWebClientInvalidMessageType()
	TestWebClientEmptyURLSuffix()
	TestWebClientInvalidURLSuffix()
	TestWebClientDataPost()
}

func TestWebClientInvalidHTTPMethod() bool {

	message := new(webClient.Message)
	message.Method = 2
	message.ContentType = "json/string"
	message.Data = []byte(`{"message":"Hello World"}`)
	//message.MessageType = webClient.Broker
	message.URLSuffix = "/broker/post"

	httpCommandFactory := new(webClient.HTTPCommandFactoryImpl)
	httpClientFactory := new(webClient.ClientFactoryImpl)
	client := httpClientFactory.GetClientService(httpCommandFactory)

	_, err := client.Do(message)
	if err != nil {
		fmt.Println(err)
		fmt.Println("InvalidHTTPMethod Validated Successfully")
		return true
	}
	fmt.Println("InvalidHTTPMethod Validation failed")
	return false
}

func TestWebClientEmptyContentType() bool {
	message := new(webClient.Message)
	message.Method = webClient.Post
	message.ContentType = ""
	message.Data = []byte(`{"message":"Hello World"}`)
	//message.MessageType = webClient.Broker
	message.URLSuffix = "/broker/post"

	httpCommandFactory := new(webClient.HTTPCommandFactoryImpl)
	httpClientFactory := new(webClient.ClientFactoryImpl)
	client := httpClientFactory.GetClientService(httpCommandFactory)

	_, err := client.Do(message)
	if err != nil {
		fmt.Println(err)
		fmt.Println("EmptyContentType Validated Successfully")
		return true
	}
	fmt.Println("EmptyContentType Validation failed")
	return false
}

func TestWebClientNilData() bool {
	message := new(webClient.Message)
	message.Method = webClient.Post
	message.ContentType = "json/string"
	message.Data = nil
	//message.MessageType = webClient.Broker
	message.URLSuffix = "/broker/post"

	httpCommandFactory := new(webClient.HTTPCommandFactoryImpl)
	httpClientFactory := new(webClient.ClientFactoryImpl)
	client := httpClientFactory.GetClientService(httpCommandFactory)

	_, err := client.Do(message)
	if err != nil {
		fmt.Println(err)
		fmt.Println("NilData Successfully")
		return true
	}
	fmt.Println("NilData Successfully")
	return false
}

func TestWebClientInvalidMessageType() bool {
	message := new(webClient.Message)
	message.Method = webClient.Post
	message.ContentType = "json/string"
	message.Data = []byte(`{"message":"Hello World"}`)
	//message.MessageType = 0
	message.URLSuffix = "/broker/post"

	httpCommandFactory := new(webClient.HTTPCommandFactoryImpl)
	httpClientFactory := new(webClient.ClientFactoryImpl)
	client := httpClientFactory.GetClientService(httpCommandFactory)

	_, err := client.Do(message)
	if err != nil {
		fmt.Println(err)
		fmt.Println("InvalidMessageType Validated Successfully")
		return true
	}
	fmt.Println("InvalidMessageType Validation failed")
	return false
}

func TestWebClientEmptyURLSuffix() bool {
	message := new(webClient.Message)
	message.Method = webClient.Post
	message.ContentType = "json/string"
	message.Data = []byte(`{"message":"Hello World"}`)
	//message.MessageType = webClient.Broker
	message.URLSuffix = ""

	httpCommandFactory := new(webClient.HTTPCommandFactoryImpl)
	httpClientFactory := new(webClient.ClientFactoryImpl)
	client := httpClientFactory.GetClientService(httpCommandFactory)

	_, err := client.Do(message)
	if err != nil {
		fmt.Println(err)
		fmt.Println("EmptyURLSuffix Validated Successfully")
		return true
	}
	fmt.Println("EmptyURLSuffix Validation failed")
	return false
}

func TestWebClientInvalidURLSuffix() bool {
	message := new(webClient.Message)
	message.Method = webClient.Post
	message.ContentType = "json/string"
	message.Data = []byte(`{"message":"Hello World"}`)
	//message.MessageType = webClient.Broker
	message.URLSuffix = "/Invalid/Route"

	httpCommandFactory := new(webClient.HTTPCommandFactoryImpl)
	httpClientFactory := new(webClient.ClientFactoryImpl)
	client := httpClientFactory.GetClientService(httpCommandFactory)

	resp, err := client.Do(message)
	if err != nil {
		fmt.Println("InvalidURLSuffix Validation failed")
		fmt.Println(err)
		return false
	}
	if resp.StatusCode == 404 {
		fmt.Println("InvalidURLSuffix Validated Successfully")
		return true
	}
	fmt.Println("InvalidURLSuffix Validation failed")
	return false
}

func TestWebClientDataPost() bool {
	message := new(webClient.Message)
	message.Method = webClient.Post
	message.ContentType = "json/string"
	message.Data = []byte(`{"message":"Hello World"}`)
	//message.MessageType = webClient.Broker
	message.URLSuffix = "/broker/post"

	httpCommandFactory := new(webClient.HTTPCommandFactoryImpl)
	httpClientFactory := new(webClient.ClientFactoryImpl)
	client := httpClientFactory.GetClientService(httpCommandFactory)

	resp, err := client.Do(message)

	if err != nil {
		fmt.Println("DataPost failed")
		fmt.Println(err)
		return false
	}
	if resp.StatusCode != 200 {
		fmt.Println("DataPost failed, server not running")
		fmt.Println(err)
		return false
	}
	fmt.Println("Data Posted Successfully")
	return true
}
