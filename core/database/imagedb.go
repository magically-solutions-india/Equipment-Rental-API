package database

import (
	"time"
	"log"
	"github.com/remony/Equipment-Rental-API/core/router"
	"strings"
)

type Images struct {
	Images []Image `json:"images"`
	Total  int     `json:"total"`
}

type Image struct {
	Title      string           `json:"title"`
	Location   string           `json:"location"`
	Date_added time.Time        `json:"date_added"`
	Size       size             `json:"size"`
}

type size struct {
	Large  string        `json:"large"`
	Medium string        `json:"medium"`
	Small  string        `json:"small"`
	Thumb  string        `json:"thumb"`
}

func DoesImageExist(api router.API, code string) bool {
	var exist bool
	err := api.Context.Session.QueryRow("SELECT EXISTS (SELECT 1 FROM images WHERE file_name LIKE '%" + code + "%');").Scan(&exist)
	if (err != nil) {
		log.Println(err)
	}
	// If it exists return true
	if exist {
		return true
	}
	// Otherwise return false
	return false

}

func GetImage(api router.API, id int) []Image {
	var images = []Image{}
	stmt, err := api.Context.Session.Prepare("CALL getImage(?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var image Image
		err := rows.Scan(
			&image.Location,
			&image.Title,
			&image.Date_added,
		)

		if err != nil {
			panic(err)
		}

		values := strings.Split(image.Location, ".")
		// set storage locations of other sizes of images
		if (len(values) == 2) {
			if (values[1] != "gif") {
				image.Size.Large = "/data/" + values[0] + "_large" + "." + values[1]
				image.Size.Medium = "/data/" + values[0] + "_medium" + "." + values[1]
				image.Size.Small = "/data/" + values[0] + "_small" + "." + values[1]
				image.Size.Thumb = "/data/" + values[0] + "_thumb" + "." + values[1]
			} else {
				image.Size.Large = "/data/" + values[0] + "." + values[1]
				image.Size.Medium = "/data/" + values[0] + "." + values[1]
				image.Size.Small = "/data/" + values[0] + "." + values[1]
				image.Size.Thumb = "/data/" + values[0] + "." + values[1]
			}

		}

		images = append(images, image)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return images
}

func GetAllImages(api router.API) Images {
	var images = []Image{}
	stmt, err := api.Context.Session.Prepare("SELECT file_name, title, date_added FROM images ORDER BY date_added DESC")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var image Image
		err := rows.Scan(
			&image.Location,
			&image.Title,
			&image.Date_added,
		)

		if err != nil {
			panic(err)
		}
		values := strings.Split(image.Location, ".")

		image.Size.Large = "/data/" + values[0] + "_large" + "." + values[1]
		image.Size.Medium = "/data/" + values[0] + "_medium" + "." + values[1]
		image.Size.Small = "/data/" + values[0] + "_small" + "." + values[1]
		image.Size.Thumb = "/data/" + values[0] + "_thumb" + "." + values[1]
		images = append(images, image)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return Images{Images:images, Total:len(images)}
}

// IsImageAvilable queries the database if the filename already exists
func IsImageAvailable(api router.API, token string) bool {
	var exist bool
	err := api.Context.Session.QueryRow("SELECT EXISTS (SELECT 1 FROM images WHERE file_name = ?)", token).Scan(&exist)
	if (err != nil) {
		log.Println(err)
	}
	// If it exists return true
	if exist {
		return true
	}
	// Otherwise return false
	return false
}

func AddImageLocationToDb(api router.API, filename string, title string, original_name string, token string) bool {
	stmt, err := api.Context.Session.Prepare("CALL addImage(?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(filename, title, original_name, token)

	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	return true
}
func DeleteImage(api router.API, title string) bool {
	stmt, err := api.Context.Session.Prepare("CALL DeleteImage(?)")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(title)

	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	return true
}

func AddImageToProduct(api router.API, filename string, token string, pid string) bool {
	stmt, err := api.Context.Session.Prepare("CALL AddAnotherImage(?, ?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(filename, filename, filename, token, pid)

	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	return true
}


