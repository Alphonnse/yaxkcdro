default:
	go mod tidy
	go build -o myapp main.go

install-deps:
	go get -u github.com/bbalet/stopwords 
	go get -u github.com/kljensen/snowball 

bench:
	go test -bench=. -benchmem -benchtime=10s -count=5

test:
	go test ./...
