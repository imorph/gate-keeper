package core

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net"
	"testing"
	"time"
)

/*
 ******************************************************************
 Benchmarks.
 ******************************************************************
*/

func BenchmarkPCTrieHitIPv4UsingAWSRanges(b *testing.B) {
	benchmarkLookUpIP(b, "52.95.110.1", NewList())
}

func BenchmarkPCTrieHitIPv6UsingAWSRanges(b *testing.B) {
	benchmarkLookUpIP(b, "2620:107:300f::36b7:ff81", NewList())
}

func BenchmarkPCTrieMissIPv4UsingAWSRanges(b *testing.B) {
	benchmarkLookUpIP(b, "123.123.123.123", NewList())
}

func BenchmarkPCTrieHMissIPv6UsingAWSRanges(b *testing.B) {
	benchmarkLookUpIP(b, "2620::ffff", NewList())
}

func benchmarkLookUpIP(tb testing.TB, nn string, list *List) {
	configureListWithAWSRanges(tb, list)
	for n := 0; n < tb.(*testing.B).N; n++ {
		_, err := list.LookUpIP(nn)
		if err != nil {
			panic(err)
		}
	}
}

type AWSRanges struct {
	Prefixes     []Prefix     `json:"prefixes"`
	IPv6Prefixes []IPv6Prefix `json:"ipv6_prefixes"`
}

type Prefix struct {
	IPPrefix string `json:"ip_prefix"`
	Region   string `json:"region"`
	Service  string `json:"service"`
}

type IPv6Prefix struct {
	IPPrefix string `json:"ipv6_prefix"`
	Region   string `json:"region"`
	Service  string `json:"service"`
}

var awsRanges *AWSRanges
var ipV4AWSRangesIPNets []*net.IPNet
var ipV6AWSRangesIPNets []*net.IPNet

func loadAWSRanges() *AWSRanges {
	file, err := ioutil.ReadFile("./testdata/aws_ip_ranges.json")
	if err != nil {
		panic(err)
	}
	var ranges AWSRanges
	err = json.Unmarshal(file, &ranges)
	if err != nil {
		panic(err)
	}
	return &ranges
}

func configureListWithAWSRanges(tb testing.TB, list *List) {
	for _, prefix := range awsRanges.Prefixes {
		err := list.InsertCIDR(prefix.IPPrefix)
		if err != nil {
			panic(err)
		}
	}
}

func init() {
	awsRanges = loadAWSRanges()
	for _, prefix := range awsRanges.IPv6Prefixes {
		_, network, _ := net.ParseCIDR(prefix.IPPrefix)
		ipV6AWSRangesIPNets = append(ipV6AWSRangesIPNets, network)
	}
	for _, prefix := range awsRanges.Prefixes {
		_, network, _ := net.ParseCIDR(prefix.IPPrefix)
		ipV4AWSRangesIPNets = append(ipV4AWSRangesIPNets, network)
	}
	rand.Seed(time.Now().Unix())
}
