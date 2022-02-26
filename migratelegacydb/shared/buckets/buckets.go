package buckets

import (
	"encoding/binary"
	"time"

	"github.com/FallenTaters/bbucket"
	"go.etcd.io/bbolt"
)

const path = "streepjes.bbolt"

var DB *bbolt.DB

func Init() func() error {
	database, err := bbolt.Open(path, 0666, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
	DB = database

	Orders = bbucket.New(DB, ordersBucketName)
	Users = bbucket.New(DB, usersBucketName)
	Members = bbucket.New(DB, membersBucketName)
	Categories = bbucket.New(DB, categoriesBucketName)
	Products = bbucket.New(DB, productsBucketName)

	return DB.Close
}

var (
	Orders     bbucket.Bucket
	Users      bbucket.Bucket
	Members    bbucket.Bucket
	Categories bbucket.Bucket
	Products   bbucket.Bucket

	ordersBucketName     = []byte("orders")
	usersBucketName      = []byte("users")
	membersBucketName    = []byte("members")
	categoriesBucketName = []byte("categories")
	productsBucketName   = []byte("products")
)

func Itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
