package main

import (
	"bytes"
	"github.com/Comcast/gots/packet"
	"github.com/Comcast/gots/pes"
	"github.com/Comcast/gots/psi"
	"io"
	"log"
	"os"
	"strconv"
)
func cut(stream []byte, start int, end int){
	final := make([]byte,0)
	flag := 0
	i := 0
	var pat psi.PAT
	file ,_ :=os.Create(strconv.Itoa(start)+"-"+strconv.Itoa(end)+".ts")
	for i < len(stream){
		chank, err := packet.FromBytes(stream[i:i+packet.PacketSize])
		if err != nil{
			log.Println(err)
			return
		}

		if chank.IsPAT() {
			pat, err = psi.ReadPAT(bytes.NewReader(chank[:]))
			if err != nil {
				log.Println("pat", err)
				return
			}
		}
		if chank.PayloadUnitStartIndicator() {
			if ok, _ := psi.IsPMT(chank, pat); !ok {
				data, err := packet.PESHeader(chank)
				if err != nil {
					log.Println("pes", err,i)

					if flag == 0{
						final = append(final,stream[i:i+packet.PacketSize]...)
					}
					i = i + packet.PacketSize
					continue
				}

				pusi, err := pes.NewPESHeader(data)
				if err != nil {
					log.Println("pusi", err)
					i = i + packet.PacketSize
					continue
				}

				log.Println("PTS", pusi.PTS())
				if pusi.PTS() >= uint64(start) && pusi.PTS() <= uint64(end) {
					flag = 1
				} else {
					flag = 0
				}
			}
		}
		if flag == 0{
			final = append(final,stream[i:i+packet.PacketSize]...)
		}
		i = i + packet.PacketSize
	}
	n,err := file.Write(final)
	if err!=nil{
		log.Println(err)
	}
	if n != len(final){
		log.Println("n!=len(b)")
	}
	defer file.Close()
}

func main(){
	file := os.Args[1]
	start,_ := strconv.Atoi(os.Args[2])
	end,_ := strconv.Atoi(os.Args[3])
	data, err := os.Open(file)
	if err != nil{
		log.Println(err)
		return
	}
	defer data.Close()

	buf := make([]byte, 64)
	stream := make([]byte, 0)
	for {
		n, err := data.Read(buf)
		if err == io.EOF{
			log.Println(err)
			break
		}
		stream = append(stream,buf[:n]...)
	}
	go cut(stream,start,end)
	for{

	}
}
