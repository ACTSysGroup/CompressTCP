package main

import (
    "fmt"
    "bufio"
    "os"
    "io"
    "math"
    "math/rand"
    "time"
    "compress/flate"
    "bytes"
    "sort"
    "strconv"
    "gopkg.in/yaml.v2"
)

type a1_counter struct {
    item    byte
    number  int64
}

type Length struct {
    RawLength       []int64
    CompressedLen   []int
}

type Entropy struct {
    RawEntropy  []float64
    EstimatedEntropy    []float64
    CompressedEntropy   []float64
}

type Test_result struct {
    FileNum []int
    Length  Length
    Entropy Entropy
    RelativeErr []float64
}

func algorithm1 (input []byte, epsilon float64, delta float64) (int, float64) {

    var entropy float64

    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
    g := int(math.Round(2*math.Log2(1/delta)))
    z := int(math.Round(32*math.Log2(float64(len(input))/epsilon/epsilon)))
    totalNum := g*z
    if totalNum > len(input) {
        z = len(input) / g
        totalNum = g*z
    }
    a1Counters := make([]a1_counter, totalNum)
    var a1CountersNum int

    a1Positions := make([]int, totalNum)
    a1Positions_shuffle := make([]int, len(input))
    for i:=0; i<len(input); i++ {
        a1Positions_shuffle[i] = i
    }
    r1.Shuffle(len(a1Positions_shuffle), func(i, j int) {
        a1Positions_shuffle[i], a1Positions_shuffle[j] = a1Positions_shuffle[j], a1Positions_shuffle[i]
    })
    for i:=0; i<totalNum; i++ {
        a1Positions[i] = a1Positions_shuffle[i]
    }
    sort.Ints(a1Positions)

    for index:=0; index<len(input); index++ {
        //if a1CountersNum < totalNum {
        //    if index == a1Positions[a1CountersNum]{
        //        for original_index:=0; original_index<totalNum; original_index++ {
        //            if a1Positions_shuffle[original_index] == a1Positions[a1CountersNum] {
        //                a1Counters[original_index].item = input[index]
        //                a1CountersNum ++
        //                break
        //            }
        //        }
        //    }
        //}

        if a1CountersNum < totalNum {
            if index == a1Positions[a1CountersNum]{
                a1Counters[a1CountersNum].item = input[index]
                a1CountersNum ++
            }
        }

        for i:=0; i<a1CountersNum; i++ {
            if a1Counters[i].item == input[index] {
                a1Counters[i].number += 1
            }
        }
    }

    a1_S := make([]float64, g)
    for i:=0; i<g; i++ {
        for j:=0; j<z; j++ {
            if a1Counters[i*z+j].number > 1 {
                a1_S[i] += float64(len(input)) * (float64(a1Counters[i*z+j].number) * math.Log2(float64(a1Counters[i*z+j].number)) - float64(a1Counters[i*z+j].number - 1) * math.Log2(float64(a1Counters[i*z+j].number - 1)))
            }
        }
        a1_S[i] = a1_S[i] / float64(z)
    }
    sort.Float64s(a1_S)

    entropy = math.Log2(float64(len(input))) - a1_S[g/2] / float64(len(input))

    return totalNum, entropy
}

func main() {
    buf_len, _ := strconv.Atoi(os.Args[1])
    epsilon, _ := strconv.ParseFloat(os.Args[2], 64)
    delta, _ := strconv.ParseFloat(os.Args[3], 64)

    var f_out *os.File
    var err error
    f_out, err = os.Create("./result.yaml")
    if err != nil {
        panic(err)
    }
    defer f_out.Close()

    FileNum := make([]int, 11)
    RawLength := make([]int64, 11)
    CompressedLen := make([]int, 11)
    RawEntropy := make([]float64, 11)
    CompressedEntropy := make([]float64, 11)
    EstimatedEntropy := make([]float64, 11)
    RelativeErr := make([]float64, 11)

    fmt.Println("|File Number\t|Raw Entropy\t|Raw Length\t|Compressed Entropy\t|Compressed Length\t|A1 Entropy\t|A1 Buf Len\t|Relative Error\t|")

    for total_times:=0; total_times<1; total_times++ {
    for iter := 0; iter < 11; iter++ {
        var f *os.File
        switch iter {
        case 0:
            f, err = os.Open("../cantrbry/alice29.txt")
        case 1:
            f, err = os.Open("../cantrbry/asyoulik.txt")
        case 2:
            f, err = os.Open("../cantrbry/plrabn12.txt")
        case 3:
            f, err = os.Open("../cantrbry/sum")
        case 4:
            f, err = os.Open("../cantrbry/ptt5")
        case 5:
            f, err = os.Open("../cantrbry/cp.html")
        case 6:
            f, err = os.Open("../cantrbry/xargs.1")
        case 7:
            f, err = os.Open("../cantrbry/fields.c")
        case 8:
            f, err = os.Open("../cantrbry/lcet10.txt")
        case 9:
            f, err = os.Open("../cantrbry/grammar.lsp")
        case 10:
            f, err = os.Open("../cantrbry/kennedy.xls")

        }
        if err != nil {
            panic(err)
        }
        defer f.Close()
        raw_buf := make([]byte, buf_len)

        reader := bufio.NewReader(f)
        buf := make([]byte, 1)

        var compress_buf bytes.Buffer
        writer, _ := flate.NewWriter(&compress_buf, 9)
        writer.Write(buf)

        var counter [256]int64
        var sieveCounter [256]int64
        var elephantFlag [256]bool
        var elephantInterval [256]int64
        s1 := rand.NewSource(time.Now().UnixNano())
        r1 := rand.New(s1)
        var streamIndex int64

//        epsilon := 0.01
//        delta := 0.1

        for {
            _,err = reader.Read(buf)

            if err != nil {
                if err != io.EOF {
                    fmt.Println(err)
                }
                break
            }
            if streamIndex == int64(buf_len) {break}

            raw_buf[streamIndex] = buf[0]
            counter[buf[0]] += 1

            ifSample := r1.Intn(100)

            if ifSample == 50 {
                if sieveCounter[buf[0]] > 0 && !elephantFlag[buf[0]] {
                    elephantFlag[buf[0]] = true
                    elephantInterval[buf[0]] = sieveCounter[buf[0]]
                    sieveCounter[buf[0]] += 1
                } else {
                    sieveCounter[buf[0]] += 1
                }
            } else {
                if sieveCounter[buf[0]] != 0 {
                    sieveCounter[buf[0]] += 1
                }
            }

            streamIndex += 1
        }

        raw_buf = raw_buf[0:streamIndex]

        writer.Write(raw_buf)
        writer.Close()
        var compress_counter [256]int64
        for i:=0; i < compress_buf.Len(); i++ {
            b, _ := compress_buf.ReadByte()
            compress_counter[b] += 1
        }
        var compress_entropy float64
        for i := 0; i < 256; i++ {
            if compress_counter[i] != 0 {
                p := float64(compress_counter[i]) / float64(compress_buf.Len())
                compress_entropy -= p * math.Log2(p)
            }
        }

        var entropy float64
        for i := 0; i < 256; i++ {
            if counter[i] != 0 {
                p := float64(counter[i]) / float64(streamIndex)
                entropy -= p * math.Log2(p)
            }
        }

        var SE float64
        var X [256]float64
        for i := 0; i < 256; i++ {
            if elephantFlag[i] {
                totalCount := sieveCounter[i] + elephantInterval[i]
                SE += float64(totalCount) * math.Log2(float64(totalCount))
            } else {
                if sieveCounter[i] > 1 {
                    X[i] = float64(streamIndex) * (float64(sieveCounter[i]) * math.Log2(float64(sieveCounter[i])) - float64(sieveCounter[i] - 1) * math.Log2(float64(sieveCounter[i] - 1)))
                } else if sieveCounter[i] == 1 {
                    X[i] = 0.00001
                }
            }
        }

        a1_len, a1_entropy := algorithm1(raw_buf, epsilon, delta)

        fmt.Printf("|%-10d\t|%-10.4f\t|%-10d\t|%-10.4f\t\t|%-10d\t\t|%-10.4f\t|%-10d\t|%-10.4f\t|\n", iter+1, entropy, streamIndex, compress_entropy, compress_buf.Len(), a1_entropy, a1_len, math.Abs(entropy - a1_entropy)/entropy)
        FileNum[iter] = iter+1
        RawLength[iter] = streamIndex
        CompressedLen[iter] = compress_buf.Len()
        RawEntropy[iter] = entropy
        CompressedEntropy[iter] = compress_entropy
        EstimatedEntropy[iter] = a1_entropy
        RelativeErr[iter] = math.Abs(entropy - a1_entropy)/entropy
    }

    result := Test_result{
        FileNum:FileNum,
        Length:Length{
            RawLength:RawLength,
            CompressedLen:CompressedLen,
        },
        Entropy:Entropy{
            RawEntropy:RawEntropy,
            EstimatedEntropy:EstimatedEntropy,
            CompressedEntropy:CompressedEntropy,
        },
        RelativeErr:RelativeErr,
    }
    d, _ := yaml.Marshal(&result)
    fmt.Fprint(f_out, string(d))

    }
}
