package etcdao

import (
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func getKapi() (client.KeysAPI, error) {
	cfg := client.Config{
		Endpoints:               []string{"http://127.0.0.1:2379"},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: 5 * time.Second,
	}

	c, err := client.New(cfg)

	if err != nil {
		return nil, err
	}

	return client.NewKeysAPI(c), nil
}

func TestReflectionScalars(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration")
	}

	kapi, err := getKapi()
	if err != nil {
		t.Fatal(err)
	}

	kapi.Set(context.Background(), "/testStr/testInt", "1", nil)

	defer func() {
		kapi.Delete(context.Background(), "/testStr", &client.DeleteOptions{Recursive: true})
	}()

	intV := 0
	err = ReadObject(kapi, context.Background(), "/testStr/testInt", &intV)
	if err != nil {
		t.Fatal(err)
	}
	if intV != 1 {
		t.Fatalf("intV != 1 but = %d", intV)
	}

	strV := ""

	err = ReadObject(kapi, context.Background(), "/testStr/testInt", &strV)
	if err != nil {
		t.Fatal(err)
	}
	if strV != "1" {
		t.Fatalf("strV != 1 but = %s", strV)
	}

}

func TestReflectionStruct(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration")
	}

	kapi, err := getKapi()
	if err != nil {
		t.Fatal(err)
	}

	kapi.Set(context.Background(), "/testStr/A", "bar", nil)
	kapi.Set(context.Background(), "/testStr/A", "1000", nil)
	kapi.Set(context.Background(), "/testStr/B", "BBB", nil)
	kapi.Set(context.Background(), "/testStr/C", "", nil)
	kapi.Set(context.Background(), "/testStr/D", "false", nil)
	kapi.Set(context.Background(), "/testStr/E", "1", nil)

	defer func() {
		kapi.Delete(context.Background(), "/testStr", &client.DeleteOptions{Recursive: true})
	}()

	var str1 struct {
		A       int
		B       string
		C, D, E bool
	}

	err = ReadObject(kapi, context.Background(), "/testStr", &str1)

	if str1.A != 1000 {
		t.Errorf("str1.A !=1000 str1.A = %v", str1.A)
	}

	if str1.B != "BBB" {
		t.Errorf("str1.A !=1000 str1.A = %v", str1.A)
	}

	if str1.C != false {
		t.Errorf("str1.C !=false str1.C = %v", str1.C)
	}

	if str1.D != false {
		t.Errorf("str1.D !=false str1.D = %v", str1.C)
	}

	if str1.E != true {
		t.Errorf("str1.E !=true str1.E = %v", str1.C)
	}

}

func TestReflectionInnerStruct(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration")
	}

	kapi, err := getKapi()
	if err != nil {
		t.Fatal(err)
	}

	kapi.Set(context.Background(), "/testStrIn/AEL", "10", nil)
	kapi.Set(context.Background(), "/testStrIn/B/BEL", "20", nil)
	kapi.Set(context.Background(), "/testStrIn/B/BEL2", "BEL2VAL", nil)

	defer func() {
		kapi.Delete(context.Background(), "/testStrIn", &client.DeleteOptions{Recursive: true})
	}()

	type B struct {
		BEL  int
		BEL2 string
	}

	type A struct {
		AEL int
		B   B
	}

	a := A{}

	err = ReadObject(kapi, context.Background(), "/testStrIn", &a)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if a.AEL != 10 {
		t.Errorf("a.AEL !=10 a.AEL = %v", a.AEL)
	}

	if a.B.BEL != 20 {
		t.Errorf("a.B.BEL !=20 a.B.BEL = %v", a.B.BEL)
	}

	if a.B.BEL2 != "BEL2VAL" {
		t.Errorf("a.B.BEL2 !='BEL2VAL' a.B.BEL2 = %v", a.B.BEL2)
	}

}

func TestSliceInt(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration")
	}

	kapi, err := getKapi()
	if err != nil {
		t.Fatal(err)
	}

	kapi.Set(context.Background(), "/testArr/0", "10", nil)
	kapi.Set(context.Background(), "/testArr/1", "20", nil)
	kapi.Set(context.Background(), "/testArr/2", "30", nil)
	defer func() {
		kapi.Delete(context.Background(), "/testArr", &client.DeleteOptions{Recursive: true})
	}()

	var a []int

	err = ReadObject(kapi, context.Background(), "/testArr", &a)

	if err != nil {
		t.Fatal(err)
	}

	if a == nil {
		t.Fatalf("a = nil")
	}
	if len(a) != 3 {
		t.Fatalf("len(a) != 3 a=%v", a)
	}

	for i := 0; i < 3; i++ {
		tv := 10 + i*10
		if a[i] != tv {
			t.Fatalf("a[%d] != %d a[%d]=%d", i, tv, i, a[0])
		}
	}

}
