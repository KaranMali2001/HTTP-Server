package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

// 1) read file 8 byte at time
// 2) read complete line at time
//   - there are three ways to do this
//     1) using scanner interface less control because it itself detect the line and
//     2) using string and manually detecting \n
//     3) using buffer rather than string because strings are not mutaable and it will create copy
//
// 3) create a getLinesChannel so that we can reuse it later
// 4) instead of reading through file listen from http stream
func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)
	go func() {
		defer f.Close()
		defer close(out)
		buf := make([]byte, 8)
		var line []byte
		for {
			n, err := f.Read(buf)
			if err != nil {
				if err == io.EOF {
					//here we need to check that if file has any left over content
					if len(line) > 0 {
						out <- string(line)
					}
					// fmt.Print("File reading Done ")
					break
				}

				log.Fatalf("error while reading file %v\n", err)
				break
			}
			//now we got some bytes and n is number of byte
			// first we will check do we have any bytes or not
			if n > 0 {
				//if we have some bytes , we will get index of \n
				//store those bytes in some variable
				data := buf[:n] // we are doing this because if we dont have exact 8 byte then we might get some random character
				//we need to loop over the data
				for len(data) > 0 {

					i := bytes.IndexByte(data, '\n')
					//indexByte return -1 if it does not find given string so check that
					if i == -1 {
						// we havent found a new line so just append to line
						// line = append(line, data) -> compile time error because append expect single element and data is array of byte
						line = append(line, data...)
						break
					}
					line = append(line, data[:i]...)
					out <- string(line)
					line = line[:0]
					data = data[i+1:]
				}
			}
		}

	}()
	return out
}

func main() {
	listenr, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println("error while reading from tcp stream", err)
	}
	fmt.Println("server is listening at ", listenr.Addr())
	defer listenr.Close()
	for {
		conn, err := listenr.Accept()
		if err != nil {
			log.Fatal("error while accepting the lister", err)
		}
		for line := range getLinesChannel(conn) {
			fmt.Println("Line is", line)
		}
	}
}
