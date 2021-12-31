make proto:
	cd protocol/; protoc --go_out=. *.proto; cd ../;