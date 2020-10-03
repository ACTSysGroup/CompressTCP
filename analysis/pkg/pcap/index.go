package mpcap

import (
	"fmt"
	"sync"
	"tcpanalysis/common"
	"tcpanalysis/pkg/compress"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"gopkg.in/yaml.v2"
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
		Filter:     compress.FilterLarge,
	},
	{
		Name:       "gzip-th500",
		Compressor: compress.CompressGzip,
		Filter:     compress.FilterLarge,
	},
	{
		Name:       "zlib-th500",
		Compressor: compress.CompressZLib,
		Filter:     compress.FilterLarge,
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

func LoadPayload(dataset common.DatasetMeta) {
	var wg = sync.WaitGroup{}

	if handle, err := pcap.OpenOffline(dataset.File); err != nil {
		panic(err)
	} else {
		defer handle.Close()
		packetStore := gopacket.NewPacketSource(handle, handle.LinkType())
		rawSize := 0
		packetCounts := 0
		for packet := range packetStore.Packets() {
			packetCounts++
			data := packet.Data()
			rs := len(data)
			wg.Add(len(methods))

			doCompress := func(raw []byte, strategy *CompressStrategy, wg *sync.WaitGroup) uint64 {
				defer wg.Done()
				flt := strategy.Filter
				if flt == nil || flt(data) {
					start := time.Now()
					comp := strategy.Compressor(data)
					elapsed := time.Since(start)
					strategy.SizeCounter += uint64(len(comp))
					strategy.Elapsed += elapsed.Microseconds()
					return uint64(len(comp))
				} else {
					strategy.SizeCounter += uint64(len(data))
					return uint64(len(data))
				}
			}

			for idx := range methods {
				go doCompress(data, &methods[idx], &wg)
			}
			wg.Wait()
			rawSize += rs
		}
		GenerateResult(DisplayCtx{
			Dataset:      dataset,
			Methods:      &methods,
			RawSize:      uint64(rawSize),
			PacketCounts: uint64(packetCounts),
		})
	}
}

type DisplayCtx struct {
	Dataset      common.DatasetMeta
	Methods      *[]CompressStrategy
	RawSize      uint64
	PacketCounts uint64
}

func GenerateResult(ctx DisplayCtx) {
	type Result struct {
		Method         string
		CompressedSize uint64 `yaml:"compressed_size"`
		Ratio          float64
		Elapsed        int64
	}
	type Ext struct {
		DatasetSize  uint64 `yaml:"dataset_size"`
		PacketCounts uint64 `yaml:"packet_counte"`
	}
	type Output struct {
		Ext     Ext
		Dataset common.DatasetMeta
		Results []Result
	}

	var results = make([]Result, len(*ctx.Methods))

	for idx, me := range *ctx.Methods {
		results[idx].Method = me.Name
		results[idx].CompressedSize = me.SizeCounter
		results[idx].Elapsed = me.Elapsed
		results[idx].Ratio = float64(me.SizeCounter) / float64(ctx.RawSize)
	}

	var ext = Ext{
		DatasetSize:  ctx.RawSize,
		PacketCounts: ctx.PacketCounts,
	}
	var out = Output{
		Dataset: ctx.Dataset,
		Results: results,
		Ext:     ext,
	}
	outStr, _ := yaml.Marshal(out)
	fmt.Printf(string(outStr))
}
