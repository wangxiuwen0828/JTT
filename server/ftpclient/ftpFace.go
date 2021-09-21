package ftpclient

import (
	"fmt"
	"gitee.com/ictt/JTTM/config"
	"gitee.com/ictt/JTTM/tools/logs"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func FtpFace() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logs.BeeLogger.Error("ftp face new error: %s", err)
		fmt.Printf("ftp face new error: %s\n", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fmt.Println("event:", event)
				if event.Op&fsnotify.Create == fsnotify.Create {
					//log.Println("modified file:", event.Name)
					//time.Sleep(1* time.Second)

					//fmt.Println(string(pic))
					time.Sleep(time.Second * time.Duration(config.FaceWriteTime))
					pic := readpic(event.Name)

					picName := strings.Split(event.Name, "/")
					picftpName := "faces/" + picName[len(picName) - 1]

					//picName := s
					msg := Message{
						FileName: picftpName,
						SaveMsg: pic,
					}
					FtpMsg <- msg
					time.Sleep(time.Second * 1 )
					fmt.Println("pic name ", picftpName)
					logs.BeeLogger.Info("pic name: %s", picftpName)

					err = os.Remove(event.Name)
					if err != nil {
						fmt.Println("remove pic fail", err)
						logs.BeeLogger.Error("ftp face remove pic error: %s", err)
						return
					}

				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logs.BeeLogger.Error("ftp face event error: %s", err)
				fmt.Printf("ftp face event error: %s\n", err)
			}
		}
	}()

	err = watcher.Add(config.FacePath)
	if err != nil {
		logs.BeeLogger.Error("ftp face add error: %s", err)
		fmt.Printf("ftp face add error: %s\n", err)
	}
	<-done
}

func readpic(path string)string {
	//读取到file中，再利用ioutil将file直接读取到[]byte中, 这是最优
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("read file fail", err)
		logs.BeeLogger.Error("ftp face open picture error: %s", err)
		return ""
	}

	defer f.Close()

	pic, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("read to pic fail", err)
		logs.BeeLogger.Error("ftp face open read to pic error: %s", err)
		return  ""
	}
	return string(pic)
}
