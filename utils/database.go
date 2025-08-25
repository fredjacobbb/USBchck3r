package utils

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Database(device, lastModified string) {
	db, err := sql.Open("sqlite3", "./database/chck3r.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	tableCreation := `
		CREATE TABLE IF NOT EXISTS usb_devices (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			device_name TEXT NOT NULL,
			last_modified TEXT NOT NULL,
			file_size INTEGER,
			CONSTRAINT unique_device_time UNIQUE (device_name, last_modified)
		);
	`
	_, err = db.Exec(tableCreation)
	if err != nil {
		fmt.Println(err)
		return
	}

	insertDevice := `
		INSERT INTO usb_devices (device_name,last_modified,file_size) VALUES (?,?,?);
	`
	_, err = db.Exec(insertDevice, device, lastModified, 4000)
	if err != nil {

	}

	readEntries := `
	   SELECT * FROM usb_devices;
	`

	rows, err := db.Query(readEntries)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var device, lastModified, insertDevice string
		var id uint
		err := rows.Scan(&id, &insertDevice, &device, &lastModified)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(id, insertDevice, device, lastModified)
	}

}
