# Usage:
# 'make dep' and 'make webtools' to install dependencies.
# 'make clean' to clear all work files
# 'make' to build css and js into static/
# 'make serve' to start dev webserver

NODE_VER = 16

JSFILES = index.js helpers.js data.js
JSFILES2 = Index.svelte Accounts.svelte Txns.svelte AccountForm.svelte TxnForm.svelte
JSFILES3 = Tablinks.svelte Journal.svelte Report.svelte SummaryRpt.svelte Setup.svelte 
JSFILES4 = UserLogin.svelte UserSignup.svelte UserPassword.svelte UserDel.svelte
JSFILES5 = SetupBooks.svelte BookForm.svelte SetupCurrencies.svelte CurrencyForm.svelte SetupUser.svelte

SRCS = ab.go util.go web.go user.go
SRCS2 = db.go dbdata.go dbrptdata.go

all: ab static/style.css static/bundle.js

dep:
	sudo apt install curl software-properties-common
	curl -fsSL https://deb.nodesource.com/setup_$(NODE_VER).x | sudo bash -
	sudo apt install -y nodejs
	sudo npm --force install -g npx

depgo:
	go env -w GO111MODULE=auto
	#go get -u github.com/mattn/go-sqlite3
	#go get -u golang.org/x/crypto/bcrypt
	#go get -u github.com/xuri/excelize

webtools:
	npm install --save-dev tailwindcss
	npm install --save-dev postcss-cli
	npm install --save-dev cssnano
	npm install --save-dev svelte
	npm install --save-dev rollup
	npm install --save-dev rollup-plugin-svelte
	npm install --save-dev @rollup/plugin-node-resolve

static/style.css: twsrc.css
	#npx tailwind build twsrc.css -o twsrc.o 1>/dev/null
	#npx postcss twsrc.o > static/style.css
	npx tailwind -i twsrc.css -o static/style.css 1>/dev/null

static/bundle.js: $(JSFILES) $(JSFILES2) $(JSFILES3) $(JSFILES4) $(JSFILES5)
	npx rollup -c

ab: $(SRCS) $(SRCS2)
	go build -o ab $(SRCS) $(SRCS2) $(SRCS3)

importer: importer.go dbdata.go db.go util.go
	go build -o importer importer.go dbdata.go db.go util.go

clean:
	rm -rf ab static/*.js static/*.css static/*.map

serve:
	python -m SimpleHTTPServer

