package listner

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/praveenmahasena/aiserver/internal/transcribe"
)

type Listner struct {
	Port string
}

func New(p string) Listner {
	return Listner{
		Port: p,
	}
}

func (l Listner) Run() error {
	li, liErr := net.Listen("tcp", l.Port)

	if liErr != nil {
		return liErr
	}

	for {
		con, conErr := li.Accept()
		if conErr != nil {
			continue
		}

		go handleCon(con)
	}

	//return nil
}

func handleCon(con net.Conn) {

	wg := sync.WaitGroup{}
	wg.Add(2)
	defer wg.Wait()

	file := make(chan []byte)

	go func(c net.Conn, f chan<- []byte, wg *sync.WaitGroup) {
		b, err := io.ReadAll(con)
		if err != nil {
			fmt.Println(err)
			return
		}
		f <- b
	}(con, file, &wg)

	close(file)

	txt, err := getTransScript(file)
	fmt.Println(err)

	go func(c net.Conn, wg *sync.WaitGroup) {
		defer wg.Done()
		n, errN := io.WriteString(c, txt)
		fmt.Println(n, errN)
	}(con, &wg)

	con.Close()
}

func getTransScript(file <-chan []byte) (string, error) {
	t := transcribe.New(<-file)
	if upLoadErr := t.UploadMediaFile(); upLoadErr != nil {
		return "", upLoadErr
	}

	if transErr := t.TranscribeRes(); transErr != nil {
		return "", transErr
	}

	return t.GetTransStr()
}
