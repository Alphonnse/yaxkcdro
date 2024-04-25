build:
	go mod tidy
	go build -o xkcd cmd/yaxkcdro/main.go

install-deps:
	go get -u github.com/bbalet/stopwords 
	go get -u github.com/kljensen/snowball 
	go get -u gopkg.in/yaml.v3
	go get -u github.com/cheggaaa/pb/v3
	go install golang.org/x/perf/cmd/benchstat@latest

test:
	go test -bench=FindComicsByStringUsingIndex ./pkg/database -benchmem -count=6 | tee benchWithIndex.txt
	go test -bench=FindComicsByStringNotUsingIndex ./pkg/database -benchmem -count=6 | tee benchWithoutIndex.txt 
	benchstat benchWithIndex.txt benchWithoutIndex.txt | tee benchCompareResoult.txt
