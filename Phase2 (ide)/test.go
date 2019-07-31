package main
import (
	"fmt"
	"encoding/gob"
	"bytes"
)

type Student struct {
	Name string
	Age int32
}

func main() {

	fmt.Println("Gob Example")

	studentEncode := Student{Name:"Ketan",Age:30}

	var b bytes.Buffer
	e := gob.NewEncoder(&b)
	if err := e.Encode(studentEncode); err != nil {
		panic(err)
	}
	fmt.Println("Encoded Struct ", b)

	var studentDecode Student
	d := gob.NewDecoder(&b)
	if err := d.Decode(&studentDecode); err != nil {
		panic(err)
	}

	fmt.Println("Decoded Struct ", studentDecode.Name,"\t",studentDecode.Age)


}