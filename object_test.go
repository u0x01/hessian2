// Copyright 2016-2019 Alex Stocks, Wongoo, Yincheng Fang
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

package hessian

import (
	"reflect"
	"testing"
)

type Department struct {
	Name string
}

func (Department) JavaClassName() string {
	return "com.bdt.info.Department"
}

type WorkerInfo struct {
	unexportedFiled   string
	Name              string
	Addrress          string
	Age               int
	Salary            float32
	Payload           map[string]int32
	FamilyMembers     []string `hessian:"familyMembers1"`
	FamilyPhoneNumber string   // default convert to => familyPhoneNumber
	Dpt               Department
}

func (WorkerInfo) JavaClassName() string {
	return "com.bdt.info.WorkerInfo"
}

func TestEncEmptyStruct(t *testing.T) {
	var (
		w   WorkerInfo
		err error
		e   *Encoder
		d   *Decoder
		res interface{}
	)

	e = NewEncoder()
	w = WorkerInfo{
		Name:              "Trump",
		Addrress:          "W,D.C.",
		Age:               72,
		Salary:            21000.03,
		Payload:           map[string]int32{"Number": 2017061118},
		FamilyMembers:     []string{"m1", "m2", "m3"},
		FamilyPhoneNumber: "010-12345678",
		// Dpt: Department{
		// 	Name: "Adm",
		// },
	}
	e.Encode(w)

	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	if err != nil {
		t.Errorf("Decode() = %+v", err)
	}
	t.Logf("decode(%v) = %v, %v\n", w, res, err)

	reflect.DeepEqual(w, res)
}

func TestEncStruct(t *testing.T) {
	var (
		w   WorkerInfo
		err error
		e   *Encoder
		d   *Decoder
		res interface{}
	)

	e = NewEncoder()
	w = WorkerInfo{
		Name:              "Trump",
		Addrress:          "W,D.C.",
		Age:               72,
		Salary:            21000.03,
		Payload:           map[string]int32{"Number": 2017061118},
		FamilyMembers:     []string{"m1", "m2", "m3"},
		FamilyPhoneNumber: "010-12345678",
		Dpt: Department{
			Name: "Adm",
		},
		unexportedFiled: "you cannot see me!",
	}
	wCopy := w
	wCopy.unexportedFiled = ""

	err = e.Encode(w)
	if err != nil {
		t.Errorf("Encode() = %+v", err)
	}

	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	if err != nil {
		t.Errorf("Decode() = %+v", err)
	}
	t.Logf("decode(%+v) = %+v, %v\n", w, res, err)

	w2, ok := res.(*WorkerInfo)
	if !ok {
		t.Fatalf("res:%T is not of type WorkerInfo", w2)
	}

	if !reflect.DeepEqual(wCopy, *w2) {
		t.Fatalf("w != w2:\n%#v\n!=\n%#v", wCopy, w2)
	}
}

type UserName struct {
	FirstName string
	LastName  string
}

func (UserName) JavaClassName() string {
	return "com.bdt.info.UserName"
}

type Person struct {
	UserName
	Age int32
	Sex bool
}

func (Person) JavaClassName() string {
	return "com.bdt.info.Person"
}

type JOB struct {
	Title   string
	Company string
}

func (JOB) JavaClassName() string {
	return "com.bdt.info.JOB"
}

type Worker struct {
	Person
	CurJob JOB
	Jobs   []JOB
}

func (Worker) JavaClassName() string {
	return "com.bdt.info.Worker"
}

func TestIssue6(t *testing.T) {
	name := UserName{
		FirstName: "John",
		LastName:  "Doe",
	}
	person := Person{
		UserName: name,
		Age:      18,
		Sex:      true,
	}

	worker := &Worker{
		Person: person,
		CurJob: JOB{Title: "cto", Company: "facebook"},
		Jobs: []JOB{
			JOB{Title: "manager", Company: "google"},
			JOB{Title: "ceo", Company: "microsoft"},
		},
	}

	e := NewEncoder()
	err := e.Encode(worker)
	if err != nil {
		t.Fatalf("encode(worker:%#v) = error:%s", worker, err)
	}
	bytes := e.Buffer()

	d := NewDecoder(bytes)
	res, err := d.Decode()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("type of decode object:%v", reflect.TypeOf(res))

	worker2, ok := res.(*Worker)
	if !ok {
		t.Fatalf("res:%#v is not of type Worker", res)
	}

	if !reflect.DeepEqual(worker, worker2) {
		t.Fatalf("worker:%#v != worker2:%#v", worker, worker2)
	}
}
