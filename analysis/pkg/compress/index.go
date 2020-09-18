package compress

import (
    "bytes"
    "compress/lzw"
    "compress/zlib"
    "compress/gzip"
    "compress/flate"
)

func CompressLZW(raw []byte)[]byte{
    var b bytes.Buffer
    w:=lzw.NewWriter(&b,lzw.MSB,8)
    w.Write(raw)
    w.Close()
    return b.Bytes()
}

func CompressZLib(raw []byte)[]byte{
    var b bytes.Buffer
    w:=zlib.NewWriter(&b)
    w.Write(raw)
    w.Close()
    return b.Bytes()
}

func CompressDeflate(raw []byte,level int )[]byte{
    var b bytes.Buffer
    w,_:=flate.NewWriter(&b,level)
    w.Write(raw)
    w.Close()
    return b.Bytes()
}
func CompressGzip(raw []byte)[]byte{
    var b bytes.Buffer
    w:=gzip.NewWriter(&b)
    w.Write(raw)
    w.Close()
    return b.Bytes()
}

func FilterLarge(raw[]byte)bool{
    if len(raw)>500{
        return true
    }
    return false
}

