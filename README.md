# cloudgo-database
## cloudgo-database
基本没写什么自己的代码，老师都写好了大部分，慌。  
主要看下老师代码做了些什么： 
```
type SQLExecer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// DaoSource Data Access Object Source
type DaoSource struct {
	// if DB, each statement execute sql with random conn.
	// if Tx, all statements use the same conn as the Tx's connection
	SQLExecer
}
```
这段定义了DaoSource和SQLExecer，这样定义后，外部程序将不再过多区分TX和DB操作的差别。。因为不管是TX还是DB，都能用来初始化DaoSource。  
然后是下面这段
```
func (*UserInfoAtomicService) Save(u *UserInfo) error {
	tx, err := mydb.Begin()
	checkErr(err)

	dao := userInfoDao{tx}
	err = dao.Save(u)

	if err == nil {
		tx.Commit()
	} else {
		tx.Rollback()
	}
	return nil
}

```
dao.Save()并不知道它处理的是DB还是TX，在UserInfoAtomicService这一层才明确。  
代码里还有好多错误处理，搞数据库相关错误处理看来要很谨慎。  

在cloudgo-data-orm中是使用xorm改写的程序，没有了亲自写Dao的过程，直接在service通过很简单的代码完成了对数据库的操作。
通过xorm，传入一个结构指针，xorm就能通过结构本身的信息来生成与结构相对应的table，或者查询结果通过这个结构呈现。
但是结构成员本身的信息以及它的成员的名字都是属于程序语言中的元信息，这怎么提取的呢。
这就利用了go的反射技术。
虽然这样很方便，但是据老师说在大数据处理量和高负荷下，orm方式操作数据库会导致较大的性能损耗，对多线程的支持也不完美。
所以如果要求性能，还是要自己做Dao。
附简单性能测试，100的同时并发数，看起来并没有造成很大的性能差异，可能还需要再大一点的负荷：
non-orm:
```
Server Software:        
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /service/userinfo?userid=1
Document Length:        95 bytes

Concurrency Level:      100
Time taken for tests:   0.355 seconds
Complete requests:      1000
Failed requests:        0
Non-2xx responses:      1000
Total transferred:      227000 bytes
HTML transferred:       95000 bytes
Requests per second:    2820.03 [#/sec] (mean)
Time per request:       35.461 [ms] (mean)
Time per request:       0.355 [ms] (mean, across all concurrent requests)
Transfer rate:          625.14 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.4      0       2
Processing:     1   33  24.1     28     125
Waiting:        1   33  24.1     27     125
Total:          1   33  24.2     28     126

Percentage of the requests served within a certain time (ms)
  50%     28
  66%     39
  75%     47
  80%     53
  90%     67
  95%     81
  98%     97
  99%    104
 100%    126 (longest request)

```
orm:
```
Server Software:        
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /service/userinfo?userid=1
Document Length:        100 bytes

Concurrency Level:      100
Time taken for tests:   0.353 seconds
Complete requests:      1000
Failed requests:        0
Non-2xx responses:      1000
Total transferred:      233000 bytes
HTML transferred:       100000 bytes
Requests per second:    2833.21 [#/sec] (mean)
Time per request:       35.296 [ms] (mean)
Time per request:       0.353 [ms] (mean, across all concurrent requests)
Transfer rate:          644.67 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.5      0       3
Processing:     1   34  27.2     27     144
Waiting:        1   34  27.2     27     144
Total:          1   34  27.4     27     144

Percentage of the requests served within a certain time (ms)
  50%     27
  66%     41
  75%     48
  80%     54
  90%     71
  95%     86
  98%    113
  99%    124
 100%    144 (longest request)

```

sql模板：
就是老师给的代码，老师写了，我一时也不知道怎么改可以改的更好，干脆讲讲这些代码做了些什么
以Select为例：
```
func (sqlt *SQLTemplate) Select(selectQuery string, rowMapper RowMapperCallback, args ...interface{}) error {

	// https://stackoverflow.com/questions/24878264/how-can-i-build-varidics-of-type-interface-in-go
	rows, err := sqlt.Query(selectQuery, args...)

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err := rowMapper(rows)
		if err != nil {
			return err
		}
	}

	return nil
}
```
selectQuery要求一个sql模板字符串，args是模板里面占位符填的东西，rowMapper比较有意思，Select方法只会闷头接受数据，得到的结果自己也不会处理
而是交给rowMapper这个函数去解决，这个rowMapper就是再上一层的userInfoDao的内容了，他会将返回的结果按照需求处理好。   
所以这个sqlt在Dao下增加了一层Template，将Dao的行为进一步泛化。  
但是你还是要写Dao,还要再额外写个相关的rowMapper 回调函数，换来的是不用和什么QueryRow，Exec，Prepare方法打交道，使用更直观的spltempalte中的方法
或者这样做结构会更清晰？

我想，如果要进行更加复杂的数据库处理，这个会更加更加方便，可以用template封装起一些复杂的操作。
