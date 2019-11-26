package core

import (
	"testing"
)

func TestAddGoodCIDRs(t *testing.T) {
	f := func(s string, l *List) {
		CIDR := s
		err := l.InsertCIDR(CIDR)
		if err != nil {
			t.Error("Insert", CIDR, "failed", "error:", err)
		}
	}
	black := NewList()
	f("13.56.0.0/16", black)
	f("23.20.0.0/14", black)
	f("0.0.0.0/0", black)
	f("0.0.0.0/32", black)
	f("127.0.0.1/32", black)
	f("52.46.80.0/21", black)
	f("192.168.1.0/24", black)
	f("10.10.10.8/8", black)
}

func TestDeleteGoodCIDRs(t *testing.T) {
	f := func(s string, l *List) {
		CIDR := s
		err := l.DeleteCIDR(CIDR)
		if err != nil {
			t.Error("Delete", CIDR, "failed", "error:", err)
		}
	}
	black := NewList()
	f("13.56.0.0/16", black)
	f("23.20.0.0/14", black)
	f("0.0.0.0/0", black)
	f("0.0.0.0/32", black)
	f("127.0.0.1/32", black)
	f("52.46.80.0/21", black)
	f("192.168.1.0/24", black)
	f("10.10.10.8/8", black)
}

func TestAddBadCIDRs(t *testing.T) {
	f := func(s string, l *List) {
		CIDR := s
		err := l.InsertCIDR(CIDR)
		if err == nil {
			t.Error("Inserted", CIDR, "should fail, but not")
		}
	}
	black := NewList()
	f("sdfsdf", black)
	f("2342345", black)
	f("nLKJlkj654765", black)
	f(`&""^%""$&^%$##$)9&*&^)`, black)
	f(`192.168.999.1/24`, black)
	f(`192.168.1.1/96`, black)
	f(`192.168.1.1.1/24`, black)
	f("128.256.3.0/33", black)
}

func TestDeleteBadCIDRs(t *testing.T) {
	f := func(s string, l *List) {
		CIDR := s
		err := l.DeleteCIDR(CIDR)
		if err == nil {
			t.Error("Deleted", CIDR, "should fail, but not")
		}
	}
	black := NewList()
	f("sdfsdf", black)
	f("2342345", black)
	f("nLKJlkj654765", black)
	f(`&""^%""$&^%$##$)9&*&^)`, black)
	f(`192.168.999.1/24`, black)
	f(`192.168.1.1/96`, black)
	f(`192.168.1.1.1/24`, black)
	f("128.256.3.0/33", black)
}

func TestLookUpGoodIP(t *testing.T) {
	black := NewList()
	CIDR := "192.0.2.0/24"
	err := black.InsertCIDR(CIDR)
	if err != nil {
		t.Error("Insert", CIDR, "failed", "error:", err)
	}
	f := func(IP, CIDR string, black *List) {
		ok, err := black.LookUpIP(IP)
		if err != nil {
			t.Error("LookUP for IP", IP, "in CIDR", CIDR, "failed", "error:", err)
		}
		if !ok {
			t.Error("LookUP for IP", IP, "in CIDR", CIDR, "failed")
		}
	}
	f("192.0.2.0", CIDR, black)
	f("192.0.2.1", CIDR, black)
	f("192.0.2.111", CIDR, black)
	f("192.0.2.255", CIDR, black)
}

func TestLookUpBadIP(t *testing.T) {
	black := NewList()
	CIDR := "192.0.2.0/24"
	err := black.InsertCIDR(CIDR)
	if err != nil {
		t.Error("Insert", CIDR, "failed", "error:", err)
	}
	f := func(IP, CIDR string, black *List) {
		ok, err := black.LookUpIP(IP)
		if err == nil {
			t.Error("LookUP for IP", IP, "in CIDR", CIDR, "should fail!")
		}
		if ok {
			t.Error("LookUP for IP", IP, "in CIDR", CIDR, "should fail!")
		}
	}
	f("192.0.2.0.1", CIDR, black)
	f("6766754654654654", CIDR, black)
	f("vksjdlskjdfghsldfkj", CIDR, black)
	f("192.0.2.0/24", CIDR, black)
}

func TestLookUpGoodIPMiss(t *testing.T) {
	black := NewList()
	CIDR := "192.0.2.0/24"
	err := black.InsertCIDR(CIDR)
	if err != nil {
		t.Error("Insert", CIDR, "failed", "error:", err)
	}
	f := func(IP, CIDR string, black *List) {
		ok, err := black.LookUpIP(IP)
		if err != nil {
			t.Error("LookUP for IP", IP, "in CIDR", CIDR, "failed", "error:", err)
		}
		if ok {
			t.Error("LookUP for IP", IP, "in CIDR", CIDR, "succeeded but should MISS")
		}
	}
	f("193.0.2.0", CIDR, black)
	f("192.0.3.1", CIDR, black)
	f("192.1.2.111", CIDR, black)
	f("127.0.0.1", CIDR, black)
	f("0.0.0.0", CIDR, black)
	f("8.8.8.8", CIDR, black)
}

func TestAllIPs(t *testing.T) {
	black := NewList()
	CIDR := "0.0.0.0/0"
	err := black.InsertCIDR(CIDR)
	if err != nil {
		t.Error("Insert", CIDR, "failed", "error:", err)
	}
	f := func(IP, CIDR string, black *List) {
		ok, err := black.LookUpIP(IP)
		if err != nil {
			t.Error("LookUP for IP", IP, "in CIDR", CIDR, "failed", "error:", err)
		}
		if !ok {
			t.Error("LookUP for IP", IP, "in CIDR", CIDR, "failed")
		}
	}

	f("192.0.2.0", CIDR, black)
	f("192.0.2.1", CIDR, black)
	f("192.0.2.111", CIDR, black)
	f("192.0.2.255", CIDR, black)
	f("193.0.2.0", CIDR, black)
	f("192.0.3.1", CIDR, black)
	f("192.1.2.111", CIDR, black)
	f("127.0.0.1", CIDR, black)
	f("0.0.0.0", CIDR, black)
	f("8.8.8.8", CIDR, black)
	f("255.255.255.255", CIDR, black)
}

func TestAddCheckDeleteCheck(t *testing.T) {
	black := NewList()
	CIDR := "192.0.2.0/24"
	err := black.InsertCIDR(CIDR)
	if err != nil {
		t.Error("Insert", CIDR, "failed", "error:", err)
	}
	f := func(IP, CIDR string, black *List) {
		ok, err := black.LookUpIP(IP)
		if err != nil {
			t.Error("LookUP for IP", IP, "in CIDR", CIDR, "failed", "error:", err)
		}
		if !ok {
			t.Error("LookUP for IP", IP, "in CIDR", CIDR, "failed")
		}
	}
	nf := func(IP, CIDR string, black *List) {
		ok, err := black.LookUpIP(IP)
		if err != nil {
			t.Error("LookUP for IP", IP, "in CIDR", CIDR, "failed", "error:", err)
		}
		if ok {
			t.Error("LookUP for IP", IP, "in CIDR", CIDR, "succeeded but should MISS")
		}
	}
	f("192.0.2.0", CIDR, black)
	f("192.0.2.1", CIDR, black)
	f("192.0.2.111", CIDR, black)
	f("192.0.2.255", CIDR, black)
	nf("192.0.4.0", CIDR, black)
	nf("192.0.4.1", CIDR, black)
	nf("192.0.4.111", CIDR, black)
	nf("192.0.4.255", CIDR, black)
	CIDR = "192.0.4.0/24"
	err = black.InsertCIDR(CIDR)
	if err != nil {
		t.Error("Insert", CIDR, "failed", "error:", err)
	}
	f("192.0.2.0", CIDR, black)
	f("192.0.2.1", CIDR, black)
	f("192.0.2.111", CIDR, black)
	f("192.0.2.255", CIDR, black)
	f("192.0.4.0", CIDR, black)
	f("192.0.4.1", CIDR, black)
	f("192.0.4.111", CIDR, black)
	f("192.0.4.255", CIDR, black)
	CIDR = "192.0.2.0/24"
	err = black.DeleteCIDR(CIDR)
	if err != nil {
		t.Error("Delete", CIDR, "failed", "error:", err)
	}
	nf("192.0.2.0", CIDR, black)
	nf("192.0.2.1", CIDR, black)
	nf("192.0.2.111", CIDR, black)
	nf("192.0.2.255", CIDR, black)
	f("192.0.4.0", CIDR, black)
	f("192.0.4.1", CIDR, black)
	f("192.0.4.111", CIDR, black)
	f("192.0.4.255", CIDR, black)
}
