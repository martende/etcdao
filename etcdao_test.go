package etcdao


import (
	"testing"
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)


func TestReflectionScalars(t*testing.T) {
	if  testing.Short() {
		t.Skip("Integration")
	}

	etcd, err := New("/test", []string{"http://127.0.0.1:2379"})
	if err != nil {
		t.Fatal(err)
	}

	intV := 0
	err = etcd.ReadObject(kapi,context.Background()"testInt",&intV)
	if err != nil {
		t.Fatal(err)
	}
	if intV != 1 {
		t.Fatalf("intV != 1 but = %d",intV)
	}

	strV := ""

	err = etcd.ReadObject(kapi,context.Background()"testInt",&strV)
	if err != nil {
		t.Fatal(err)
	}
	if strV != "1" {
		t.Fatalf("intV != 1 but = %s",strV)
	}

}

func TestReflectionStruct(t*testing.T) {
	if  testing.Short() {
		t.Skip("Integration")
	}
	etcd, err := New("/test", []string{"http://127.0.0.1:2379"})
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	etcd.WriteKey("/testStr/A","1000")
	etcd.WriteKey("/testStr/B","BBB")
	etcd.WriteKey("/testStr/C","")
	etcd.WriteKey("/testStr/D","false")
	etcd.WriteKey("/testStr/E","1")

	defer func() {
		kapi,root := etcd.Kapi()
		kapi.Delete(context.Background(),root + "/testStr",&client.DeleteOptions{Recursive:true})
	}()

	var str1 struct {
		A int
		B string
		C,D,E bool
	}

	err = etcd.ReadObject(kapi,context.Background()"testStr",&str1)

	if str1.A != 1000 {
		t.Errorf("str1.A !=1000 str1.A = %v",str1.A )
	}

	if str1.B != "BBB" {
		t.Errorf("str1.A !=1000 str1.A = %v",str1.A )
	}

	if str1.C != false {
		t.Errorf("str1.C !=false str1.C = %v",str1.C )
	}

	if str1.D != false {
		t.Errorf("str1.D !=false str1.D = %v",str1.C )
	}

	if str1.E != true {
		t.Errorf("str1.E !=true str1.E = %v",str1.C )
	}

}


func TestReflectionInnerStruct(t*testing.T) {
	if  testing.Short() {
		t.Skip("Integration")
	}
	etcd, err := New("/test", []string{"http://127.0.0.1:2379"})
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	etcd.WriteKey("/testStrIn/AEL","10")
	etcd.WriteKey("/testStrIn/B/BEL","20")
	etcd.WriteKey("/testStrIn/B/BEL2","BEL2VAL")

	defer func() {
		kapi,root := etcd.Kapi()
		kapi.Delete(context.Background(),root + "/testStrIn",&client.DeleteOptions{Recursive:true})
	}()

	type B struct {
		BEL int
		BEL2 string
	}

	type A struct {
		AEL int
		B B
	}

	a := A{}

	err = etcd.ReadObject(kapi,context.Background(),"/test/testStrIn",&a)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if a.AEL != 10 {
		t.Errorf("a.AEL !=10 a.AEL = %v",a.AEL )
	}

	if a.B.BEL != 20 {
		t.Errorf("a.B.BEL !=20 a.B.BEL = %v",a.B.BEL )
	}

	if a.B.BEL2 != "BEL2VAL" {
		t.Errorf("a.B.BEL2 !='BEL2VAL' a.B.BEL2 = %v",a.B.BEL2 )
	}

}

func TestSliceInt(t*testing.T) {
	if  testing.Short() {
		t.Skip("Integration")
	}
	etcd, err := New("/test", []string{"http://127.0.0.1:2379"})
	if err != nil {
		t.Fatal(err)
	}
	etcd.WriteKey("/testArr/0","10")
	etcd.WriteKey("/testArr/1","20")
	etcd.WriteKey("/testArr/2","30")
	defer func() {
		kapi,root := etcd.Kapi()
		kapi.Delete(context.Background(),root + "/testArr",&client.DeleteOptions{Recursive:true})
	}()

	var a []int

	err = etcd.ReadObject(kapi,context.Background(),"/test/testArr",&a)

	if err != nil {
		t.Fatal(err)
	}

	if a == nil {
		t.Fatalf("a = nil")
	}
	if len(a) != 3 {
		t.Fatalf("len(a) != 3 a=%v",a)
	}

	for i := 0 ; i < 3 ; i++ {
		tv := 10 + i * 10
		if a[i] != tv {
			t.Fatalf("a[%d] != %d a[%d]=%d",i,tv,i,a[0])
		}
	}

}