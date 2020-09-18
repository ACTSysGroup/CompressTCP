package mpcap

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"tcpanalysis/pkg/compress"
	"time"
)

var methods = []CompressStrategy{
    {
        Name: "middle-deflate-th500",
        Compressor: func(raw []byte) []byte {
            return compress.CompressDeflate(raw, 4)
        },
        Filter: compress.FilterLarge,
    },
    {
        Name: "huffman-th500",
        Compressor: func(raw []byte) []byte {
            return compress.CompressDeflate(raw, -2)
        },
        Filter: compress.FilterLarge,
    },
    {
        Name: "full-deflate-th500",
        Compressor: func(raw []byte) []byte {
            return compress.CompressDeflate(raw, 9)
        },
        Filter: compress.FilterLarge,
    },
    {
        Name:       "lzw-th500",
        Compressor: compress.CompressLZW,
        Filter: compress.FilterLarge,
    },
    {
        Name:       "gzip-th500",
        Compressor: compress.CompressGzip,
        Filter: compress.FilterLarge,
    },
    {
        Name:       "zlib-th500",
        Compressor: compress.CompressZLib,
        Filter: compress.FilterLarge,
    },
    {
        Name: "middle-deflate",
        Compressor: func(raw []byte) []byte {
            return compress.CompressDeflate(raw, 4)
        },
    },
    {
        Name: "huffman",
        Compressor: func(raw []byte) []byte {
            return compress.CompressDeflate(raw, -2)
        },
    },
    {
        Name: "full-deflate",
        Compressor: func(raw []byte) []byte {
            return compress.CompressDeflate(raw, 9)
        },
    },
    {
        Name:       "lzw",
        Compressor: compress.CompressLZW,
    },
    {
        Name:       "gzip",
        Compressor: compress.CompressGzip,
    },
    {
        Name:       "zlib",
        Compressor: compress.CompressZLib,
    },
}

type TCPPayload []byte

type CompressStrategy struct {
	Name        string
	Compressor  func(raw []byte) []byte
	SizeCounter uint64
	Elapsed     int64
    Filter      func(raw []byte) bool
}

func LoadPayload(fp string) {

	if handle, err := pcap.OpenOffline(fp); err != nil {
		panic(err)
	} else {
		defer handle.Close()
		packetStore := gopacket.NewPacketSource(handle, handle.LinkType())
		rawSize := 0
		for packet := range packetStore.Packets() {
			data := packet.Data()
			rs := len(data)
			for idx := range methods {
                flt :=  methods[idx].Filter
                if flt==nil || flt(data){
                    start := time.Now()
                    comp := methods[idx].Compressor(data)
                    elapsed := time.Since(start)
                    methods[idx].SizeCounter += uint64(len(comp))
                    methods[idx].Elapsed += elapsed.Microseconds()
                } else{
                    methods[idx].SizeCounter += uint64(rs)
                }
			}
			rawSize += rs
		}

		for idx := range methods {
			c := methods[idx].SizeCounter
			e := methods[idx].Elapsed
			fmt.Printf("method: %25s\t size: %d\t ratio: %f elapsed: %d\n", methods[idx].Name, c, float64(c)/float64(rawSize), e)
		}
	}
}
