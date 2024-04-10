build:
	go mod tidy
	go build -o myapp cmd/yaxkcdro/main.go

install-deps:
	go get -u github.com/bbalet/stopwords 
	go get -u github.com/kljensen/snowball 
	go get -u gopkg.in/yaml.v3
	go get -u github.com/cheggaaa/pb/v3
