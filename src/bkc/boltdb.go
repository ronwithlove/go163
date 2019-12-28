package main

import (
	"fmt"
	"log"

"github.com/boltdb/bolt"
)

func main() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		//创建一个桶
		b, err := tx.CreateBucket([]byte("MyBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		//写入数据
		if nil !=b{
			err:=b.Put([]byte("answer"),[]byte("42"))//写入key 和value
			if nil !=err{
				return err
			}
		}
		return nil
	})

	//读
	db.View(func(tx *bolt.Tx) error {
		//获取桶
		b:=tx.Bucket([]byte("MyBucket"))
		if nil !=b{
			value:=b.Get([]byte("answer"))//读key
			fmt.Printf("value: %s\n",value)
		}
		return nil
	})
}
