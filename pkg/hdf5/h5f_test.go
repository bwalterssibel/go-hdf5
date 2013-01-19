package hdf5

import (
	"os"
	"testing"
)

func TestFile(t *testing.T) {
	f, err := CreateFile(FNAME, F_ACC_TRUNC)
	if err != nil {
		t.Fatalf("CreateFile failed: %s", err)
	}
	defer os.Remove(FNAME)
	defer f.Close()

	if fileName := f.FileName(); fileName != FNAME {
		t.Fatalf("FileName() have %v, want %v", fileName, FNAME)
	}
	// The file is also the root group
	if name := f.Name(); name != "/" {
		t.Fatalf("Name() have %v, want %v", name, FNAME)
	}
	if err := f.Flush(F_SCOPE_GLOBAL); err != nil {
		t.Fatalf("Flush() failed: %s", err)
	}
	if !IsHdf5(FNAME) {
		t.Fatalf("IsHdf5 returned false")
	}

	f2 := f.File()
	fName := f.FileName()
	f2Name := f2.FileName()
	if fName != f2Name {
		t.Fatalf("f2 FileName() have %v, want %v", f2Name, fName)
	}

	// Test a Group
	groupName := "test"
	g, err := f.CreateGroup(groupName)
	if err != nil {
		t.Fatalf("CreateGroup() failed: %s", err)
	}
	if name := g.Name(); name != "/"+groupName {
		t.Fatalf("Group Name() have %v, want /%v", name, groupName)
	}

	g2, err := f.OpenGroup(groupName)
	if err != nil {
		t.Fatalf("OpenGroup() failed: %s", err)
	}
	if name := g2.Name(); name != "/"+groupName {
		t.Fatalf("Group Name() have %v, want /%v", name, groupName)
	}

	// Test a Dataset
	ds, err := CreateDataSpace(S_SCALAR)
	if err != nil {
		t.Fatalf("CreateDataspace failed: %s", err)
	}
	dsetName := "test_dataset"
	dset, err := f.CreateDataset(dsetName, T_NATIVE_INT, ds, P_DEFAULT)
	if err != nil {
		t.Fatalf("CreateDataset failed: %s", err)
	}
	if name := dset.Name(); name != "/"+dsetName {
		t.Fatalf("Dataset Name() have %v, want /%v", name, dsetName)
	}
	dFile := dset.File()
	if dFile.Name() != f.Name() {
		t.Fatalf("Dataset File() have %v, want %v", dFile.Name(), f.Name())
	}
}

func TestClosedFile(t *testing.T) {
	f, err := CreateFile(FNAME, F_ACC_TRUNC)
	if err != nil {
		t.Fatalf("CreateFile failed: %s", err)
	}
	fName := f.Name()
	f2 := f.File()
	f.Close()

	f2Name := f2.FileName()
	if f2Name != FNAME {
		t.Fatalf("f2 FileName() have %v, want %v", f2Name, fName)
	}
	f2.Close()

	os.Remove(FNAME)
	f3 := f.File()
	if f3 != nil {
		t.Fatalf("expected file to be nil")
	}

}
