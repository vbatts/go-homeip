package ipstore

import (
	"database/sql"
	//"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	//"os"
)

var (
	db *sql.DB
)

func InitFilename(db_filename string) (err error) {
	/*
			  1) check if db_filename exists.
		     a) if not, create it
		     b) if so, validate it
		     c) OR just os.Remove the filename and create new everytime
		    2) open the db_filename to db
	*/
	err = os.Remove(db_filename)
	if err != nil {
		return
	}

	db, err = sql.Open("sqlite3", db_filename)
	if err != nil {
		return
	}
	initialize_sqls := []string{
		"create table hosts (id integer not null primary key, name text, ip text)",
		"delete from hosts",
	}

	for _, sql := range initialize_sqls {
		_, err = db.Exec(sql)
		if err != nil {
			log.Printf("%q: %s\n", err, sql)
			return
		}
	}

	return
}

func Close() (err error) {
	return db.Close()
}

func HostExists(hostname string) (ret_val bool, err error) {
	rows, err := db.Query("select count(1) from hosts where name = '%s'", hostname)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var count int
		rows.Scan(&count)
		if count > 0 {
			return true, nil
		}
	}
	return false, nil
}

func SetHostIp(hostname, ip string) (err error) {
	exists, err := HostExists(hostname)
	if err != nil {
		return
	}
	if exists {
		_, err = db.Exec("update hosts set ip = '%s' where name = '%s'", ip, hostname)
	} else {
		_, err = db.Exec("insert into hosts(name, ip) values('%s', '%s')", ip, hostname)
	}
	return
}

func GetHostIp(hostname string) (ip string, err error) {
	rows, err := db.Query("select ip from hosts where name = '%s'", hostname)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&ip)
		return ip, nil
	}
	return
}
