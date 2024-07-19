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

	file := make(chan []byte)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(c net.Conn, f chan<- []byte, wg *sync.WaitGroup) {
		b, err := io.ReadAll(c)
		if err != nil {
			fmt.Println(err)
		}
		f <- b
		fmt.Println(b)
		close(f)
	}(con, file, wg)

	txt, err := getTransScript(file)

	if err != nil {
		con.Write([]byte("internal Error"))
		con.Close()
		return
	}

	wg.Add(1)
	go func(c net.Conn, wg *sync.WaitGroup) {
		defer wg.Done()
		n, errN := c.Write([]byte(txt))
		fmt.Println(txt)
		fmt.Println(n, errN)
		con.Close()
	}(con, wg)

	wg.Wait()

	fmt.Println("came here finally")

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
