package tree

import "testing"

import (
    "os"
    "bytes"
    "math/rand"
    "encoding/binary"
)

import (
  "github.com/timtadh/data-structures/types"
)

func init() {
    if urandom, err := os.Open("/dev/urandom"); err != nil {
        return
    } else {
        buf := make([]byte, 8)
        if _, err := urandom.Read(buf); err == nil {
            buf_reader := bytes.NewReader(buf)
            if seed, err := binary.ReadVarint(buf_reader); err == nil {
                rand.Seed(seed)
            }
        }
        urandom.Close()
    }
}

func randstr(length int) types.String {
    if urandom, err := os.Open("/dev/urandom"); err != nil {
        panic(err)
    } else {
        slice := make([]byte, length)
        if _, err := urandom.Read(slice); err != nil {
            panic(err)
        }
        urandom.Close()
        return types.String(slice)
    }
    panic("unreachable")
}

func TestPutHasGetRemoveBucket(t *testing.T) {

    type record struct {
        key types.String
        value types.String
    }

    records := make([]*record, 400)
    var tree *AvlTree
    var err error
    var val interface{}
    var updated bool

    ranrec := func() *record {
        return &record{ randstr(20), randstr(20) }
    }

    for i := range records {
        r := ranrec()
        records[i] = r
        tree, updated = tree.Put(r.key, types.String(""))
        if updated {
            t.Error("should have not been updated")
        }
        tree, updated = tree.Put(r.key, r.value)
        if !updated {
            t.Error("should have been updated")
        }
        if tree.Size() != (i+1) {
            t.Error("size was wrong", tree.Size(), i+1)
        }
    }

    for _, r := range records {
        if has := tree.Has(r.key); !has {
            t.Error("Missing key")
        }
        if has := tree.Has(randstr(12)); has {
            t.Error("Table has extra key")
        }
        if val, err := tree.Get(r.key); err != nil {
            t.Error(err, val.(types.String), r.value)
        } else if !(val.(types.String)).Equals(r.value) {
            t.Error("wrong value")
        }
    }

    for i, x := range records {
        if tree, val, err = tree.Remove(x.key); err != nil {
            t.Error(err)
        } else if !(val.(types.String)).Equals(x.value) {
            t.Error("wrong value")
        }
        for _, r := range records[i+1:] {
            if has := tree.Has(r.key); !has {
                t.Error("Missing key")
            }
            if has := tree.Has(randstr(12)); has {
                t.Error("Table has extra key")
            }
            if val, err := tree.Get(r.key); err != nil {
                t.Error(err)
            } else if !(val.(types.String)).Equals(r.value) {
                t.Error("wrong value")
            }
        }
        if tree.Size() != (len(records) - (i+1)) {
            t.Error("size was wrong", tree.Size(), (len(records) - (i+1)))
        }
    }
}

func TestIterators(t *testing.T) {
    var data []int = []int{
        1, 5, 7, 9, 12, 13, 17, 18, 19, 20,
    }
    var order []int = []int{
        6, 1, 8, 2, 4 , 9 , 5 , 7 , 0 , 3 ,
    }
    var tree *AvlTree
    var updated bool

    for j := range order {
        if tree, updated = tree.Put(types.Int(data[order[j]]), order[j]); updated {
            t.Error("should have not been updated")
        }
    }

    j := 0
    for k, v, next := tree.Iterate()(); next != nil; k, v, next = next() {
        if !k.Equals(types.Int(data[j])) {
            t.Error("Wrong key")
        }
        if v.(int) != j {
            t.Error("Wrong value")
        }
        j += 1
    }

    j = 0
    for k, next := tree.Keys()(); next != nil; k, next = next() {
        if !k.Equals(types.Int(data[j])) {
            t.Error("Wrong key")
        }
        j += 1
    }

    j = 0
    for v, next := tree.Values()(); next != nil; v, next = next() {
        if v.(int) != j {
            t.Error("Wrong value")
        }
        j += 1
    }
}
