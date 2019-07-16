package main

import (
    "fmt"
    "os"
    "sync"

    "github.com/qiniu/api.v7/auth/qbox"
    "github.com/qiniu/api.v7/storage"
    "github.com/qiniu/x/rpc.v7"
)

var (
    AK     string
    SK     string
)

func init() {
    AK = os.Getenv("Q_AK")
    SK = os.Getenv("Q_SK")
}

func main() {
    if AK == "" || SK == "" {
        fmt.Println("Please set Q_AK, Q_SK environment variable")
    }

    if len(os.Args) < 3 {
        fmt.Printf("Usage: %s <src_bucket> <dest_bucket>\n", os.Args[0])
        return
    }

    srcBkt, dstBkt := os.Args[1], os.Args[2]

    mac := qbox.NewMac(AK, SK)
    manager := storage.NewBucketManager(mac, nil)
    fileCh, err := manager.ListBucket(srcBkt,"", "", "")
    if err != nil {
        fmt.Println(err)
        return
    }

    moveOps := make([]string, 0, 1000)
    opCh := make(chan string)
    exitCh := make(chan struct{})
    wg := sync.WaitGroup{}

    wg.Add(1)
    go func(ch <-chan string, exitCh chan struct{}) {
        for {
            select {
            case moveOp := <-opCh:
                if moveOps == nil {
                    moveOps = make([]string, 0, 1000)
                }
                moveOps = append(moveOps, moveOp)

                if len(moveOps) == 1000 {
                    go batchMove(manager, moveOps)
                    moveOps = nil
                }
            case <-exitCh:
                batchMove(manager, moveOps)
                wg.Done()
                return
            }
        }

    }(opCh, exitCh)

    for item := range fileCh {
        moveOp := storage.URIMove(srcBkt, item.Item.Key, dstBkt, item.Item.Key, true)
        opCh <- moveOp
    }
    exitCh <- struct{}{}
    wg.Wait()
}


func batchMove(manager *storage.BucketManager, ops []string) {
    ret, err := manager.Batch(ops)
    if err != nil {
        if _, ok := err.(*rpc.ErrorInfo); ok {
            for _, r := range ret {
                if r.Code != 200 {
                    fmt.Println(r.Data.Error)
                }
            }
        } else {
            fmt.Printf("batch error: %v\n", err)
        }
    }
}