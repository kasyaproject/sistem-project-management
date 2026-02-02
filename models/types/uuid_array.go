package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type UUIDArray []uuid.UUID

// function untuk mengkonversi data dari database ke dalam UUIDArray, karena pada DB disimpan dalam bentuk string
func (a *UUIDArray) Scan(value interface{}) error {
	// Variable untuk menampung data uuid sebagai string
	var str string

	// Memeriksa tipe data dari value yang diterima dan di konversi ke string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return errors.New("failed to parse UUIDArray: unsupport data type")
	}

	// Menghapus kurung kurawal dan memisahkan string berdasarkan koma
	str = strings.Trim(str, "{}")
	// Mengonversi string menjadi array
	parts := strings.Split(str, ",")

	// Mengonversi array string menjadi array UUID
	*a = make(UUIDArray, 0, len(parts))
	for _, s := range parts {
		// Membersihkan spasi dan tanda kutip
		s = strings.TrimSpace(strings.Trim(s, `"`))
		// Melewati string kosong
		if s == "" {
			continue
		}

		// Mengonversi string menjadi UUID
		u, err := uuid.Parse(s)
		// Memeriksa error konversi
		if err != nil {
			return fmt.Errorf("invalid UUID in Array : %v", err)
		}
		// Menambahkan UUID ke array
		*a = append(*a, u)
	}

	return nil
}

// function untuk mengkonversi data dari UUIDArray ke dalam bentuk yang dapat disimpan di database, yaitu string
func (a UUIDArray) Value() (driver.Value, error) {
	// Jika array kosong, kembalikan string kosong dalam format PostgreSQL
	if len(a) == 0 {
		return "{}", nil
	}

	// Mengonversi setiap UUID dalam array ke string dan membentuk format PostgreSQL
	postgreFormat := make([]string, 0, len(a))
	// Looping setiap UUID dalam array dan tambahkan "" di sekelilingnya
	for _, value := range a {
		postgreFormat = append(postgreFormat, fmt.Sprintf(`"%s"`, value.String()))
	}

	// Menggabungkan semua string menjadi satu dalam format PostgreSQL
	return "{" + strings.Join(postgreFormat, ",") + "}", nil
}

// middleware gorm untuk mendefinisikan tipe data custom UUIDArray, setiap kali ada typedata uuid[]
func (UUIDArray) GormDataType() string {
	return "uuid[]"
}
