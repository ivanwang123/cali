.PHONY: build gui full-build

build:
	go build -o cali.exe && copy cali.exe "C:\Users\ivanw\Documents\Go\bin"

full-build:
	go build -o cali.exe -ldflags -H=windowsgui && copy cali.exe "C:\Users\ivanw\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup"

gui:
	boltdbweb --db-name=cali.db --port=8080