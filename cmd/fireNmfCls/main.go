package main

import (
	"flag"
	"log"
	"strconv"
	"time"

	"github.com/nimezhu/fire"

	"github.com/boltdb/bolt"
	. "github.com/nimezhu/ice"
	"github.com/nimezhu/ice/encoding/tsv"
)

var help = flag.Bool("help", false, "print help")
var kstar = flag.Int("k", 5, "k number")
var fn = flag.String("t", "", "filename")
var out = flag.String("o", "", "output database filename")

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}
func generateNames(k int, prefix string) []string {
	s := make([]string, k)
	for i := 0; i < k; i++ {
		s[i] = prefix + strconv.Itoa(i)
	}
	return s
}
func float64arr(a []int) []float64 {
	r := make([]float64, len(a))
	for i := range a {
		r[i] = float64(a[i])
	}
	return r
}
func main() {

	flag.Parse()
	if *help {
		flag.PrintDefaults()
		return
	}
	k := *kstar
	t := new(Table)
	t.LoadTsv(*fn)
	r, c := t.Dims()
	mat, rowIds, rowGroup, colIds, colGroup, W, H, _, rowLDA, colLDA := fire.SortLabeledNMF(t.Dense(), t.Rows(), t.Cols(), k)
	db, err := bolt.Open(*out, 0600, &bolt.Options{Timeout: 1 * time.Second})
	checkErr(err)
	defer db.Close()
	tx, err := db.Begin(true)
	checkErr(err)
	defer tx.Commit()
	//defer tx.Rollback()

	// Use the transaction...
	bucketName := "bucket"
	bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
	checkErr(err)
	outdb := &Bucket{Bucket: bucket, Encoder: tsv.Encode, Decoder: tsv.Decode}
	outdb.Add(&Table{[]string{"colGroup"}, colIds, 1, c, float64arr(colGroup), "colGroup", "colGroup"}, "colGroup")
	outdb.Add(&Table{[]string{"rowGroup"}, rowIds, 1, r, float64arr(rowGroup), "rowGroup", "rowGroup"}, "rowGroup")
	outdb.Add(&Table{colIds, generateNames(k, "E"), c, k, H.RawMatrix().Data, "H", "H"}, "H")
	outdb.Add(&Table{generateNames(k, "E"), rowIds, k, r, W.RawMatrix().Data, "W", "W"}, "W")
	outdb.Add(&Table{colIds, rowIds, c, r, mat.RawMatrix().Data, "mat", "mat"}, "mat")
	outdb.Add(&Table{generateNames(k, "Group"), rowIds, k, r, rowLDA.RawMatrix().Data, "rowLDA", "rowLDA"}, "rowLDA")
	outdb.Add(&Table{colIds, generateNames(k, "Group"), k, c, colLDA.RawMatrix().Data, "colLDA", "colLDA"}, "colLDA")
	log.Println("done")

}
