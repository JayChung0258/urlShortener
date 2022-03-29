# urlshortener
## Table of Content
* Install
* Usage
* Ideas
* Database
* 3rd Party libs
* Unit test(unfinished)
* Reflection

## Install

### Using Docker to connect (Recommended)
This project uses Docker. Go check them out if you don't have them locally installed.
```
cd project_location/urlshortener-master
```
```
docker compose up
```

After building up the container, open your cmd/terminal or other database connection tool.

You can check the database state under this command.

**cmd/terminal**
```
psql -h localhost -p 5432 -U postgres urlshortener
```

```
\l #checking all database under this port
\d #checking the realtions in urls
```

### Using localhost to connect
My environment include below:
* Postgresql: 14.2
* Go: 1.18
* Homebrew: 3.4.3

First, Install Go on official website.

And we should insatll Homebrew for install and start Postgresql services. 
(I'll recommend you to use Docker if your OS is Windows)

**terminal**
```
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
# install Homebrew
```
```
brew install postgresql
# install postgresql
```
```
brew start services postgresql
# start the postgresql database with default port 5432
```
```
psql postgresql
# log in to the database interactive mode
```

## Usage
You can use cmd/terminal or Postman to use urlshortener. Here we introduce the way on cmd/terminal.


**Create short url**
```
curl -X POST -H "Content-Type:application/json" http://localhost/api/v1/urls -d '{
"url": "<original_url>",
"expireAt": "2021-02-08T09:20:41Z"
}'

#Response
{"id": "<url_id>", "shortUrl": "http://localhost/<uri_id>"}
```

**GET short url**
```
curl -L -X GET http://localhost/<url_id>

#Response
"<original_url>"
```

**GET all data**
```
curl -L -X GET http://localhost/urls

#Response
{
"url": "<original_url>",
"expireAt": "2021-02-08T09:20:41Z",
"id": "<url_id>", 
"shortUrl": "http://localhost/<uri_id>"
}, ... all data
```





## Ideas
一開始看到這個題目，第一個是先去找縮網址的網站，看看它的運作原理，在做這題目之前，我以為縮網址的運作原理是用數學計算，後來發現題目可以使用database，才發現原來是要把原網址跟縮網址pair再一起，放在資料庫裡面等待檢索。

因為在寫這份作業之前，對於資料庫的了解不太深，僅有在前端call過後端做好的API而已，因此這次使用Postgresql也是一個挑戰，只不過因為畢業論文也是使用Postgresql，因此剛好來學習使用它。

而往下看看到curl指令，第一時間就想到了這是要使用http的作業，第一個想法是先去載Postman，因為等等就會用到了，比起curl，用Postman測試起來方便多了，剛好學校最近又在教tcp/udp的protocol，倒是這次作業讓我也把port的很多小細節釐清了。

簡單來說呢，我的想法是使用Go語言建立router去監聽一個port(這裡是用port:80)，而建立與Postgresql database的連線(這裡是port: 5432)，而CRUD就在這個main.go進行，因為對於Go語言沒有到很熟悉，因此邊學邊用順便參照js的寫法完成了這次作業。

但完成之後隔天睡起來又想到，阿別人要怎麼使用我的程式呢，於是我把它整包丟到另外一台電腦上，發現光環境的設定就快不行了，搞了超久才弄好，於是我就想到，如果打包起來丟過來不就好了嗎，於是這是我第一次學習docker，琢磨了一天，成功包起來使用compose的方式啟動我的url image 和 postgres image讓另外一台電腦可以無痛運行。

## Database
原本想使用本地記憶體像是Redis之類的，應該會比連線資料庫簡單且快速，但因為論文有使用Postgresql的需求，因此就順便在這次的作業中實現了。

## Standard libs
* fmt
fmt用來做基本字串串接以及print out
* net/http
官方http庫做router處理
* os
os來讀取.env檔案
* json
json來encode/decode檔案
* time
來解析傳進來的時間以判斷expire
* net/http/httptest
生成mock request 和 writer
* regexp
正則表達式去比對測試SQL syntax
* testing
原生的Go test 

## 3rd Party libs
* 資料庫處理: go-gorm
gorm V2可以以簡潔的方式處理資料庫的連接，以及強大的ORM支援，在Go檔案中就可以處理SQL syntax
```
db, err = gorm.Open(postgres.Open(dsn_urlShortener), &gorm.Config{})
# 開啟資料庫
```
* 連線處理: gorilla/mux
使用mux代替原生的http router，可以更方便的批量創造request HandlerFunc，而且與官方http.serveMux相容
```
// create Router
router := mux.NewRouter()
// Listener
http.ListenAndServe(":80", router)
```
* Mock SQL database: DATA-DOG/go-sqlmock
使用這個sqlmock可以在unit test裡面生成固定格式的假資料庫，這樣測試可以不會影響到真的的資料庫
```
//set up the mock sql connection
	testDB, mock, err := sqlmock.New()
    
// uses "gorm.io/driver/postgres" library
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 testDB,
		PreferSimpleProtocol: true,
	})
```

## Unit Test (Unfinished)
Unit Test因為時間的關係，目前只有做出`getURL()`的測試，這裡測試的重點在於SQL syntax，不在於http連線，以`getURL()`為例子，我們預期會得到的SQL syntax為 
`SELECT * FROM "urls" WHERE "urls"."id"=<id>`，下圖為測試的示意圖。
![](https://i.imgur.com/dsrSb43.png)

**操作方法(目前還未支援docker裡run test)**
```
cd project/urlshortener
```
```
go test -v
```


還沒補完的方法我認為都是造樣造句即可，專注於SQL syntax的解析與配對，這個問題其實解滿久的，差不多有兩天的時間都在想unit test怎麼寫，因為很少接觸，而Mux + Gorm + postgresql又讓難度更高了，查遍了stackoverflow也沒有甚麼好想法最後硬著頭皮發了一篇文詢問，再加上之前爬文看到的想法拼拼湊湊寫出來了。

[stackoverflow上的提問](https://stackoverflow.com/questions/71645815/how-to-unit-testing-with-gorm-mux-postgresql/71648556?noredirect=1#comment126631348_71648556)

### Reflection
這次有機會寫到Dcard的作業純屬意外，滑著ig看到Dcard有在徵實習生的推播，就點進來了，要不然真的不知道Dcard也有在找實習生。而這個作業屬於後端的實作CRUD，對我來說也滿新鮮的，在學校沒什麼機會可以做到完整的操作，而上一份實習也是大一的事情了，就想說來試試看，雖然我對Go以及資料庫都沒有到很熟練，但花了這幾天的過程也讓我知道哪裡還可以加強以及要成為後端應該要往哪裡點技能樹。
