package driver

import (
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
)

type Datastore struct {
	ds     *object.Datastore
	driver *Driver
}

func (d *Driver) NewDatastore(ref *types.ManagedObjectReference) *Datastore {
	return &Datastore{
		ds:     object.NewDatastore(d.client.Client, *ref),
		driver: d,
	}
}

// If name is an empty string, returns the default datastore (is exists)
func (d *Driver) FindDatastore(name string) (*Datastore, error) {
	ds, err := d.finder.DatastoreOrDefault(d.ctx, name)
	if err != nil {
		return nil, err
	}
	return &Datastore{
		ds:     ds,
		driver: d,
	}, nil
}

func (ds *Datastore) Info(params ...string) (*mo.Datastore, error) {
	var p []string
	if len(params) == 0 {
		p = []string{"*"}
	} else {
		p = params
	}
	var info mo.Datastore
	err := ds.ds.Properties(ds.driver.ctx, ds.ds.Reference(), p, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (ds *Datastore) FileExists(path string) bool {
	_, err := ds.ds.Stat(ds.driver.ctx, path)
	return err == nil
}

func (ds *Datastore) Name() string {
	return ds.ds.Name()
}

func (ds *Datastore) ResolvePath(path string) string {
	return ds.ds.Path(path)
}

func (ds *Datastore) UploadFile(src, dst string) error {
	p := soap.DefaultUpload
	return ds.ds.UploadFile(ds.driver.ctx, src, dst, &p)
}

func (ds *Datastore) Delete(path string) error {
	dc, err := ds.driver.finder.Datacenter(ds.driver.ctx, ds.ds.DatacenterPath)
	if err != nil {
		return err
	}
	fm := ds.ds.NewFileManager(dc, false)
	return fm.Delete(ds.driver.ctx, path)
}

// Cuts out the datastore prefix
// Example: "[datastore1] file.ext" --> "file.ext"
func RemoveDatastorePrefix(path string) string {
	res := object.DatastorePath{}
	if hadPrefix := res.FromString(path); hadPrefix {
		return res.Path
	} else {
		return path
	}
}
