package professor

import (
	"io"
	"encoding/binary"
	"bytes"	
	"fmt"
)

// Version
const (
	REQUEST = 0x01
	RESPONSE = 0x81
)

// Flags
const (
	NO_FLAGS    = 0x00
	COMPRESSION = 0x01
	TRACING		= 0x02
)

// Opcodes
const (	
 	ERROR 			= 0x00
	STARTUP 		= 0x01
	READY 			= 0x02
	AUTHENTICATE 	= 0x03
	CREDENTIALS 	= 0x04
	OPTIONS 		= 0x05
	SUPPORTED 		= 0x06
	QUERY 			= 0x07
	RESULT 			= 0x08
	PREPARE 		= 0x09
	EXECUTE 		= 0x0A
	REGISTER 		= 0x0B
	EVENT 			= 0x0C	
)

// Consistencies
const (
	ANY 			= 0x0000
 	ONE				= 0x0001
	TWO				= 0x0002
	THREE			= 0x0003
	QUORUM 			= 0x0004
	ALL             = 0x0005
	LOCAL_QUORUM    = 0x0006
	EACH_QUORUM     = 0x0007
)

// Results
const (
	VOID			= 0x0001
	ROWS			= 0x0002
	SET_KEYSPACE	= 0x0003
	PREPARED		= 0x0004
	SCHEMA_CHANGE 	= 0x0005
)

type Header struct {
	version	uint8
	flags 	uint8
	stream	uint8
	opcode 	uint8
	length	uint32
}

type Frame struct {
	header 	*Header
	body 	interface{}
}

type LongString string

type Query struct {
	cql 		LongString
	consistency uint16
}

type Result struct {
	kind uint32
	body interface{}
}

func Options(stream uint8) *Frame {
	header := &Header{version:REQUEST, flags:NO_FLAGS, stream:stream, opcode:OPTIONS, length:0}
	return &Frame{header:header, body:nil}
}

func Startup(stream uint8, options map[string]string) *Frame {
	header := &Header{version:REQUEST, flags:NO_FLAGS, stream:stream, opcode:STARTUP}
	return &Frame{header:header, body:options}	
}

func Cql(stream uint8, query *Query) *Frame {
	header := &Header{version:REQUEST, flags:NO_FLAGS, stream:stream, opcode:QUERY}	
	return &Frame{header:header, body: query}	
}

func send(writer io.Writer, frame *Frame) {	
	buf := new(bytes.Buffer)

	switch content := frame.body.(type) {
		case map[string]string:			
			writeMap(buf, content)			
		case *Query:
			writeLongString(buf, content.cql)
			writeShort(buf, content.consistency)
		case nil:
			// ignore (e.g. for OPTIONS request)
	}

	//log.Println("Send buffer:",buf.Bytes())

	frame.header.length = uint32(buf.Len())	
	
	binary.Write(writer, binary.BigEndian, frame.header)
	if frame.header.length > 0 {
		buf.WriteTo(writer)
	}
		
}

func recv(reader io.Reader) (*Frame, error) {
	var header Header
	
	binary.Read(reader, binary.BigEndian, &header.version)
	binary.Read(reader, binary.BigEndian, &header.flags)
	binary.Read(reader, binary.BigEndian, &header.stream)
	binary.Read(reader, binary.BigEndian, &header.opcode)
	binary.Read(reader, binary.BigEndian, &header.length)
		
	frame := &Frame{header:&header}
	
	//log.Printf("%+v",frame.header)

	switch header.opcode {
		case RESULT:
			frame.body = readResult(reader)
		case ERROR:
			code := readInt(reader)
			message := readString(reader)			
			return frame, fmt.Errorf("Error from remote peer: code = %d; message = %s", code, message)			
		case SUPPORTED:
			frame.body = readMultiMap(reader)		
		case READY:
			// ignore
	}

	return frame, nil

}

func readLength(reader io.Reader) (n uint16) {	
	binary.Read(reader, binary.BigEndian, &n)
	return
}

func writeStringLength(buf *bytes.Buffer, s string) {
	n := uint16(len(s))	
	binary.Write(buf, binary.BigEndian, n)	
}

func writeMapLength(buf *bytes.Buffer, dict map[string]string) {
	n := uint16(len(dict))	
	binary.Write(buf, binary.BigEndian, n)	
}

func readInt(reader io.Reader) (i int32) {
	binary.Read(reader, binary.BigEndian, &i)
	return		
}

func writeShort(buf *bytes.Buffer, s uint16) {	
	binary.Write(buf, binary.BigEndian, s)
}

func readString(reader io.Reader) string {
	n := readLength(reader)	
	buf := make([]byte, n)
	reader.Read(buf)
	return string(buf)		
}

func writeString(buf *bytes.Buffer, s string) {	
	writeStringLength(buf, s)	
	buf.WriteString(s)	
}

func writeLongString(buf *bytes.Buffer, ls LongString) {
	s := string(ls)
	n := uint32(len(s))
	binary.Write(buf, binary.BigEndian, n)
	buf.WriteString(s)	
}

func writeMap(buf *bytes.Buffer, dict map[string]string) {
	writeMapLength(buf, dict)
	for k, v := range dict {		
		writeString(buf, k)
		writeString(buf, v)				
	}	
}

func readStringList(reader io.Reader) []string {
	n := readLength(reader)
	list := make([]string, n)
	
	for i := uint16(0); i < n; i++ {
		value := readString(reader)
		list = append(list, value)		
	}
	
	return list
}

func readMap(reader io.Reader) map[string]string {
	n := readLength(reader)
	dict := make(map[string]string, n)
	
	for i := uint16(0); i < n; i++ {
		key := readString(reader)
		value := readString(reader)
		dict[key] = value
	}

	return dict
}

func readMultiMap(reader io.Reader) map[string][]string {
	n := readLength(reader)
	multiMap := make(map[string][]string, n)
	
	for i := uint16(0); i < n; i++ {
		key := readString(reader)
		multiMap[key] = readStringList(reader)
	}

	return multiMap
}

func readResult(reader io.Reader) (result *Result) {
	result = new(Result)
	result.kind = uint32(readInt(reader))
	switch result.kind {
		case VOID:
		case ROWS:
		case SET_KEYSPACE:
			keyspace := readString(reader)			
			result.body = keyspace
		case SCHEMA_CHANGE:
		case PREPARED:
	}
	return
}
