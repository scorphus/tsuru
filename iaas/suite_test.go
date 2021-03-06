// Copyright 2014 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iaas

import (
	"testing"

	"launchpad.net/gocheck"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

type S struct{}

var _ = gocheck.Suite(&S{})

func (s *S) SetUpTest(c *gocheck.C) {
	iaasProviders = make(map[string]IaaS)
	RegisterIaasProvider("test-iaas", TestIaaS{})
	coll := collection()
	defer coll.Close()
	coll.RemoveAll(nil)
	tplColl := template_collection()
	defer tplColl.Close()
	tplColl.RemoveAll(nil)
}

type TestIaaS struct{}

func (TestIaaS) DeleteMachine(m *Machine) error {
	m.Status = "destroyed"
	return nil
}

func (TestIaaS) CreateMachine(params map[string]string) (*Machine, error) {
	params["should"] = "be in"
	m := Machine{
		Id:      params["id"],
		Status:  "running",
		Address: params["id"] + ".somewhere.com",
	}
	return &m, nil
}

type TestDescriberIaaS struct {
	TestIaaS
}

func (i TestDescriberIaaS) DeleteMachine(m *Machine) error {
	return i.TestIaaS.DeleteMachine(m)
}

func (i TestDescriberIaaS) CreateMachine(params map[string]string) (*Machine, error) {
	return i.TestIaaS.CreateMachine(params)
}

func (i TestDescriberIaaS) Describe() string {
	return "ahoy desc!"
}

type TestCustomizableIaaS struct {
	name string
	TestIaaS
}

func (i TestCustomizableIaaS) DeleteMachine(m *Machine) error {
	return i.TestIaaS.DeleteMachine(m)
}

func (i TestCustomizableIaaS) CreateMachine(params map[string]string) (*Machine, error) {
	return i.TestIaaS.CreateMachine(params)
}

func (i TestCustomizableIaaS) Clone(name string) IaaS {
	i.name = name
	return i
}
