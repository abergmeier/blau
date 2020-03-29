package main

import (
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
)

func main() {
	playB, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	p := &pb.Player{}
	err = proto.Unmarshal(os.Stdin, p)
	if err != nil {
		panic(err)
	}

	rules.CalcPoints()
}
