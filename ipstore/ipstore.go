package ipstore

import (
	"database/sql"
	//"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
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
	if os.IsNotExist(err) {
		log.Println("Not removing", db_filename, "because it does not exist")
  } else if err != nil {
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
	stmt, err := db.Prepare("select count(3) from hosts where name = ?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var count int
	stmt.QueryRow(hostname).Scan(&count)
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func DropHostIp(hostname string) (err error) {
	result, err := db.Exec("delete from hosts where name = '?'", hostname)
	affected, _ := result.RowsAffected()
	log.Printf("RowsAffected: %s", affected)
	return err
}

func SetHostIp(hostname, ip string) (err error) {
	exists, err := HostExists(hostname)
	if err != nil {
		return err
	}
	if exists {
		stmt, err := db.Prepare("update hosts set ip = '%s' where name = '%s'")
		if err != nil {
			return err
		}
    defer stmt.Close()
    stmt.QueryRow(ip, hostname)
	} else {
		stmt, err := db.Prepare("insert into hosts(name, ip) values('%s', '%s')")
		if err != nil {
			return err
		}
    defer stmt.Close()
    stmt.QueryRow(ip, hostname)
	}
	return
}

func GetHostIp(hostname string) (ip string, err error) {
	stmt, err := db.Prepare("select ip from hosts where name = ?")
	if err != nil {
		return
	}
	defer stmt.Close()
	stmt.QueryRow(hostname).Scan(&ip)
	return ip, nil
}
