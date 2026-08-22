// Copyright (c) 2026 The BFE Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build ignore

package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/bfenetworks/bfe-access-pb/b2log"
)

func payloadOfSize(size int) []byte {
	// Use repeating pattern that never contains the magic bytes 0xA7 0xBE 0xAE 0xB0.
	payload := make([]byte, size)
	for i := range payload {
		payload[i] = byte(i % 251)
	}
	return payload
}

func writeRecord(buf *bytes.Buffer, payload []byte, corruptMagic bool, compressLen uint32, uncompressLen uint32) {
	var header b2log.Header
	header.MagicNumber = b2log.MAGIC_NUMBER
	if corruptMagic {
		header.MagicNumber = 0xDEADBEEF
	}
	header.Version = b2log.HEADER_VERSION
	if uncompressLen != 0 {
		header.UnCompressLen = uncompressLen
	} else {
		header.UnCompressLen = uint32(len(payload))
	}
	header.CompressLen = compressLen
	header.TimeStamp = 0

	hbuf := new(bytes.Buffer)
	_ = binary.Write(hbuf, binary.LittleEndian, header)
	buf.Write(hbuf.Bytes())
	buf.Write(payload)
}

func main() {
	_ = os.MkdirAll("test_data", 0755)

	// pb_access_1.log: 8 normal records
	buf := new(bytes.Buffer)
	for i := 0; i < 8; i++ {
		writeRecord(buf, []byte(fmt.Sprintf("record-%d", i)), false, 0, 0)
	}
	if err := os.WriteFile("test_data/pb_access_1.log", buf.Bytes(), 0644); err != nil {
		panic(err)
	}

	// pb_access_2.log: first record magic broken, then 7 normal records => 7 parsed
	buf.Reset()
	writeRecord(buf, []byte("bad-record"), true, 0, 0)
	for i := 0; i < 7; i++ {
		writeRecord(buf, []byte(fmt.Sprintf("record-%d", i)), false, 0, 0)
	}
	if err := os.WriteFile("test_data/pb_access_2.log", buf.Bytes(), 0644); err != nil {
		panic(err)
	}

	// pb_access_3.log: first record compress_len != 0, then 7 normal records => 7 parsed
	// Parser bypasses HEADER_SIZE + compressLen bytes for compressed record.
	buf.Reset()
	writeRecord(buf, payloadOfSize(100), false, 100, 0)
	for i := 0; i < 7; i++ {
		writeRecord(buf, []byte(fmt.Sprintf("record-%d", i)), false, 0, 0)
	}
	if err := os.WriteFile("test_data/pb_access_3.log", buf.Bytes(), 0644); err != nil {
		panic(err)
	}

	// pb_access_4.log: first record compress_len > 100K, then 507 normal records => 507 parsed
	// compressLen is capped to MAX_RECORD_LEN (100K), so payload size must be 100K.
	buf.Reset()
	writeRecord(buf, payloadOfSize(100*1024), false, uint32(100*1024+1), 0)
	for i := 0; i < 507; i++ {
		writeRecord(buf, []byte(fmt.Sprintf("record-%d", i)), false, 0, 0)
	}
	if err := os.WriteFile("test_data/pb_access_4.log", buf.Bytes(), 0644); err != nil {
		panic(err)
	}

	// pb_access_5.log: first record uncompress_len > 100K, then 507 normal records => 508 parsed
	// Parser truncates to MAX_RECORD_LEN and returns it as a valid record.
	buf.Reset()
	writeRecord(buf, payloadOfSize(100*1024), false, 0, uint32(100*1024+1))
	for i := 0; i < 507; i++ {
		writeRecord(buf, []byte(fmt.Sprintf("record-%d", i)), false, 0, 0)
	}
	if err := os.WriteFile("test_data/pb_access_5.log", buf.Bytes(), 0644); err != nil {
		panic(err)
	}

	fmt.Println("test data generated")
}
