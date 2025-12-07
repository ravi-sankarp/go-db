**GO-DB** is an embedded database written in Go.

I built this project to dive deeper into database internals. Having used PostgreSQL for several years, I was always curious about how it works under the hood. So, I decided to write my own database to gain a better understanding of its internals. That's how this project began.

The database uses [slotted pages](https://www.interdb.jp/pg/pgsql01/03.html) to store data in binary format. All data is converted to binary and inserted into a database file.

The [Page Headers](https://github.com/ravi-sankarp/go-db/blob/master/core/internal/page_headers.go) file contains all the headers used in this process.

Currently, the database supports both reading and writing data.

Next steps include:

- Building a query parser for SQL queries
- Implementing B-Tree indexes
- Adding transaction handling
- Setting up a distributed database system
