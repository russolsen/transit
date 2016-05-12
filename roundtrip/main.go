package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	//	"reflect"
	"github.com/russolsen/transit"
)

var logf, _ = os.OpenFile("/tmp/log.txt", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)

func ReadTransit(jsd *json.Decoder) interface{} {
	decoder := transit.NewJsonDecoder(jsd)

	fmt.Fprintf(logf, "Reading...")

	value, err := decoder.Decode()

	fmt.Fprintf(logf, "Value read:\n%v\n", value)

	//fmt.Printf("The value read is: %v[%v]\n", value, reflect.TypeOf(value))

	if err == io.EOF {
		return err
	} else if err != nil {
		fmt.Fprintf(logf, "Error reading Transit data: %v\n", err)
		return io.EOF
	}

	return value
}

func WriteTransit(value interface{}) {
	encoder := transit.NewEncoder(os.Stdout, false)

	fmt.Fprintf(logf, "Writing...")
	err := encoder.Encode(value)

	if err != nil {
		fmt.Fprintf(logf, "Error writing Transit data: %v\n", err)
		return
	}
}

func main() {

	jsd := json.NewDecoder(os.Stdin)

	for x := ReadTransit(jsd); x != io.EOF; x = ReadTransit(jsd) {
		WriteTransit(x)
		os.Stdout.Sync()
	}
	fmt.Fprintf(logf, "Done!")
}
