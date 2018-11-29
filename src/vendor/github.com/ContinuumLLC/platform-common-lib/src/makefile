all: dependencies build-all test cover

build-all: api-windows-build logger-build cassandra-build clar-build env-build exception-build json-build kafka-build protocol-build pluginUtils-build procParser-build web-build webClient-build wmi-build utils-build namedpipes-build setupUtils-build downloader-build instrumentation-build


dependencies:
	glide install

api-windows-build:
	GOOS=windows go install ./api/win/pdh

logger-build:
	go install ./logging

cassandra-build:
	go install ./cassandra

cassandra-mockgen:
	mockgen -package mock github.com/ContinuumLLC/platform-common-lib/src/cassandra DbConnector,Factory,BatchQueryExecutor > cassandra/mock/mocks.go

clar-build:
	go install ./clar

clar-mockgen:
	mockgen	-package mock github.com/ContinuumLLC/platform-common-lib/src/clar ServiceInit,ServiceInitFactory > clar/mock/mocks.go

env-build:
	go install ./env

env-mockgen:
	mockgen -package mock github.com/ContinuumLLC/platform-common-lib/src/env FactoryEnv,Env > env/mock/mocks.go

exception-build:
	go install ./exception

json-build:
	go install ./json

json-mockgen:
	mockgen -package mock github.com/ContinuumLLC/platform-common-lib/src/json FactoryJSON,DeserializerJSON,SerializerJSON > json/mock/mocks.go

kafka-build:
	go install ./kafka

kafka-mockgen:
	mockgen -package mock github.com/ContinuumLLC/platform-common-lib/src/kafka ProducerFactory,ConsumerFactory,ProducerService,ConsumerService,Limiter > kafka/mock/mocks.go

protocol-build:
	go install ./plugin/protocol ./plugin/protocol/http

instrumentation-build:
	go install ./instrumentation

protocol-mockgen:
	mockgen -package mock github.com/ContinuumLLC/platform-common-lib/src/plugin/protocol Client,Server,ClientFactory,ServerFactory,ResponseSerializer > plugin/protocol/mock/mocks.go

pluginUtils-build:
	go install ./pluginUtils

pluginUtils-mockgen:
	mockgen -package mock github.com/ContinuumLLC/platform-common-lib/src/pluginUtils IOReaderFactory,IOWriterFactory > ./pluginUtils/mock/mocks.go

procParser-build:
	go install ./procParser

procParser-mockgen:
	mockgen -package mock github.com/ContinuumLLC/platform-common-lib/src/procParser Parser,ParserFactory > ./procParser/mock/mocks.go

web-build:
	go install ./web

web-mockgen:
	mockgen -package mock github.com/ContinuumLLC/platform-common-lib/src/web Server,ServerFactory,Resource,RequestContext > web/mock/mocks.go

webClient-build:
	go install ./webClient

webClient-mockgen:
	mockgen -package mock github.com/ContinuumLLC/platform-common-lib/src/webClient ClientFactory,ClientService,HTTPClientFactory,HTTPClientService > webClient/mock/mocks.go

services-mockgen:
	mockgen -package mock github.com/ContinuumLLC/platform-common-lib/src/services/model VersionFactory,VersionService,VersionDependencies,HealthCheckServiceFactory,HealthCheckService,HealthCheckDalFactory,HealthCheckDal,HealthCheckDependencies > services/mock/mocks.go

config-mockgen:
	mockgen -package mock github.com/ContinuumLLC/platform-common-lib/src/config ConfigurationService > config/mock/mocks.go

namedpipes-mockgen:
	mockgen -package mock github.com/ContinuumLLC/platform-common-lib/src/namedpipes ServerPipe,ClientPipe > namedpipes/mock/mocks.go
	mockgen -package mock net Listener,Conn > namedpipes/mock/nMocks.go

sync-mockgen:
	mockgen -package mock github.com/ContinuumLLC/platform-common-lib/src/sync Service > sync/mock/mocks.go
	mockgen -package mock github.com/ContinuumLLC/platform-common-lib/src/sync/zookeeper Connection > sync/zookeeper/mock/mocks.go

exec-mockgen:
	mockgen -package mock github.com/ContinuumLLC/platform-common-lib/src/exec Command > exec/mock/mocks.go

wmi-build:
	GOOS=windows go install ./plugin/wmi

wmi-mockgen:
	mockgen -package wmiMock github.com/ContinuumLLC/platform-common-lib/src/plugin/wmi Wrapper > ./plugin/wmi/wmiMock/mocks_windows.go

messaging-mockgen:
	mockgen -package mock github.com/ContinuumLLC/platform-common-lib/src/messaging Service > messaging/mock/mocks.go

downloader-mockgen:
	mockgen -package mock github.com/ContinuumLLC/platform-common-lib/src/downloader Service > downloader/mock/mocks.go

checksum-mockgen:
	mockgen -package mock github.com/ContinuumLLC/platform-common-lib/src/checksum Service > checksum/mock/mocks.go

utils-build:
	go install ./utils

namedpipes-build:
	go install ./namedpipes

setupUtils-build:
	go install ./setupUtils

downloader-build:
	go install ./downloader


sqlite-mockgen:
	mockgen -package mock github.com/ContinuumLLC/platform-common-lib/src/db/sqlite Service > db/sqlite/mock/mocks.go


test:
	go test -coverprofile ./uuid/cover.out -covermode=count ./uuid
	go test -coverprofile ./logging/cover.out -covermode=count ./logging
	go test -coverprofile ./plugin/protocol/cover.out -covermode=count ./plugin/protocol
	go test -coverprofile ./plugin/protocol/http/cover.out -covermode=count ./plugin/protocol/http
	go test -coverprofile ./json/cover.out -covermode=count ./json
	go test -coverprofile ./web/cover.out -covermode=count ./web
	go test -coverprofile ./exception/cover.out -covermode=count ./exception
	go test -coverprofile ./kafka/cover.out -covermode=count ./kafka
	go test -coverprofile ./procParser/cover.out -covermode=count ./procParser
	go test -coverprofile ./clar/cover.out -covermode=count ./clar
	go test -coverprofile ./web/microService/cover.out -covermode=count ./web/microService
	go test -coverprofile ./web/rest/cover.out -covermode=count ./web/rest
	go test -coverprofile ./webClient/cover.out -covermode=count ./webClient
	go test -coverprofile ./utils/cover.out -covermode=count ./utils
	go test -coverprofile ./namedpipes/cover.out -covermode=count ./namedpipes
	go test -coverprofile ./sync/zookeeper/cover.out -covermode=count ./sync/zookeeper
	go test -coverprofile ./messaging/cover.out -covermode=count ./messaging
	go test -coverprofile ./downloader/http/cover.out -covermode=count ./downloader/http
	go test -coverprofile ./kafka/encode/cover.out -covermode=count ./kafka/encode
	go test -coverprofile ./entitlement/cover.out -covermode=count ./entitlement
	go test -coverprofile ./checksum/cover.out -covermode=count ./checksum
	go test -coverprofile ./cherwell/cover.out -covermode=count ./cherwell
	go test -coverprofile ./instrumentation/cover.out -covermode=count ./instrumentation

cover:
	echo "mode: count" > cover-all.out
	tail -n +2 ./uuid/cover.out >> cover-all.out
	tail -n +2 ./logging/cover.out >> cover-all.out
	tail -n +2 ./plugin/protocol/cover.out >> cover-all.out
	tail -n +2 ./plugin/protocol/http/cover.out >> cover-all.out
	tail -n +2 ./json/cover.out >> cover-all.out
	tail -n +2 ./web/cover.out >> cover-all.out
	tail -n +2 ./exception/cover.out >> cover-all.out
	tail -n +2 ./webClient/cover.out >> cover-all.out
	tail -n +2 ./kafka/cover.out >> cover-all.out
	tail -n +2 ./procParser/cover.out >> cover-all.out
	tail -n +2 ./clar/cover.out >> cover-all.out
	tail -n +2 ./web/microService/cover.out >> cover-all.out
	tail -n +2 ./web/rest/cover.out >> cover-all.out
	tail -n +2 ./sync/zookeeper/cover.out >> cover-all.out
	tail -n +2 ./messaging/cover.out >> cover-all.out
	tail -n +2 ./downloader/http/cover.out >> cover-all.out
	tail -n +2 ./kafka/encode/cover.out >> cover-all.out
	tail -n +2 ./entitlement/cover.out >> cover-all.out
	tail -n +2 ./checksum/cover.out >> cover-all.out
	tail -n +2 ./cherwell/cover.out >> cover-all.out
	tail -n +2 ./instrumentation/cover.out >> cover-all.out
	gocover-cobertura < cover-all.out > cover-cobertura.xml

cover-xml:
	@$(foreach package,$(packages), \
		gocov convert $(package)/cover.out | gocov-xml > $(package)/coverage.xml;)
