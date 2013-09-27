package main

import (
    "fmt"
    "io"
    "os"
    "strings"
)

type rot13Reader struct {
    r io.Reader
}

func (reader *rot13Reader) Read(p []byte) (n int, err error) {
    n, err = reader.r.Read(p)
    
    for i:=0; i<len(p); i++ {
        switch {
        case 65 <= p[i] && p[i] <= 90: // upper case letter
        	newValue := ((p[i] + 13) % (65 + 26)) 
		if newValue < 65 {
			newValue += 65
		}
		p[i] = newValue
        case 97 <= p[i] && p[i] <= 122: // lower case letter
        	newValue := ((p[i] + 13) % (97 + 26))
		if newValue < 97 {
			newValue += 97
		}
		p[i] = newValue
        }
    }
    
    return
}

func main() {
    s := strings.NewReader("Lbh penpxrq gur pbqr!")
    r := rot13Reader{s}
    io.Copy(os.Stdout, &r)
    fmt.Println()
}
