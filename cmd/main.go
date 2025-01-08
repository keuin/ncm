package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/keuin/ncm"
)

func main() {
	log.SetFlags(0)
	var (
		in  = os.Stdin
		out = os.Stdout
	)
	defer func() {
		_ = in.Close()
		_ = out.Close()
	}()
	if len(os.Args) == 2 || len(os.Args) == 3 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatalf("failed to open file `%s`: %v", os.Args[1], err)
		}
		in = f
	}
	if len(os.Args) == 3 {
		f, err := os.OpenFile(os.Args[2], os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatalf("failed to open file `%s`: %v", os.Args[2], err)
		}
		out = f
	}
	if len(os.Args) > 3 {
		fmt.Printf("Usage: %s <input_file> <output_file>\n", os.Args[0])
		os.Exit(1)
	}
	dec, err := ncm.NewDecoder(in)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Print(getBriefInfo(&dec.Metadata), ".", dec.Metadata.Format)
	_, err = io.Copy(out, dec)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func getBriefInfo(m *ncm.Metadata) string {
	var sb strings.Builder
	sb.WriteString(m.MusicName)
	sb.WriteString(" - ")
	for i := range m.Artist {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(m.Artist[i].Name)
	}
	return sb.String()
}
