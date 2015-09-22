package ipstore

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db          *sql.DB
	GlobalToken = "global"
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
		"create table hosts (id integer not null primary key, name text, ip text, token text)",
		//"delete from hosts",
	}

	for _, sql := range initialize_sqls {
		_, err = db.Exec(sql)
		if err != nil {
			log.Printf("ipstore: %q: %s\n", err, sql)
			return
		}
	}

	return
}

func Close() (err error) {
	return db.Close()
}

// HostExists validates whether there is a record for the hostname in the GlobalToken context
func HostExists(hostname string) (ret_val bool, err error) {
	return HostExistsToken(hostname, GlobalToken)
}

// HostExistsToken validates whether there is a record for the hostname
func HostExistsToken(hostname, token string) (ret_val bool, err error) {
	stmt, err := db.Prepare("select count(1) from hosts where name = ? and token = ?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var count int
	stmt.QueryRow(hostname, token).Scan(&count)
	if count > 0 {
		return true, nil
	}
	return false, nil
}

// DropHostIp removes the record for hostname
func DropHostIp(hostname string) (err error) {
	return DropHostIpToken(hostname, GlobalToken)
}

// DropHostIpToken removes the record for hostname
func DropHostIpToken(hostname, token string) (err error) {
	result, err := db.Exec("delete from hosts where name = ? and token = ?", hostname, token)
	affected, _ := result.RowsAffected()
	log.Printf("RowsAffected: %s", affected)
	return err
}

// SetHostIp sets the ip for a host. Either by updating or inserting the record
func SetHostIp(hostname, ip string) (err error) {
	return SetHostIpToken(hostname, ip, GlobalToken)
}

// SetHostIpToken sets host ip, in context of the provided token
func SetHostIpToken(hostname, ip, token string) (err error) {
	exists, err := HostExistsToken(hostname, token)
	if err != nil {
		return err
	}
	if exists {
		stmt, err := db.Prepare("update hosts set ip = ? where name = ? and token = ?")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(ip, hostname, token)
		if err != nil {
			return err
		}
	} else {
		stmt, err := db.Prepare("insert into hosts(name, ip, token) values(?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(hostname, ip, token)
		if err != nil {
			return err
		}
	}
	return
}

// GetHostIp gets the ip address for hostname
func GetHostIp(hostname string) (ip string, err error) {
	return GetHostIpToken(hostname, GlobalToken)
}

// GetHostIpToken gets the ip address for hostname in the context of token
func GetHostIpToken(hostname, token string) (ip string, err error) {
	rows, err := db.Query("select ip from hosts where name = ? and token = ?", hostname, token)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&ip)
		log.Println(ip)
		if len(ip) > 0 {
			//return ip, nil
		}
	}
	return ip, nil
}
