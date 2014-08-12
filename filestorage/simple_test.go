// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package filestorage_test

import (
	"bytes"
	"io/ioutil"

	"github.com/juju/testing"
	gc "launchpad.net/gocheck"

	"github.com/juju/utils/filestorage"
)

//---------------------------
// metadata storage

var _ = gc.Suite(&MetadataStorageSuite{})

type MetadataStorageSuite struct {
	testing.IsolationSuite
	original filestorage.Metadata
}

func (s *MetadataStorageSuite) SetUpTest(c *gc.C) {
	s.IsolationSuite.SetUpTest(c)
	s.original = filestorage.NewMetadata(nil)
	s.original.SetFile(0, "", "")
}

func (s *MetadataStorageSuite) TestMetadataStorageNewMetadataStorage(c *gc.C) {
	stor := filestorage.NewMetadataStorage()

	c.Check(stor, gc.NotNil)
}

func (s *MetadataStorageSuite) TestMetadataStorageDoc(c *gc.C) {
	stor := filestorage.NewMetadataStorage()
	id, err := stor.AddDoc(s.original)
	c.Assert(err, gc.IsNil)

	doc, err := stor.Doc(id)
	c.Assert(err, gc.IsNil)
	meta, ok := doc.(filestorage.Metadata)
	c.Assert(ok, gc.Equals, true)
	c.Check(meta, gc.DeepEquals, s.original)
}

func (s *MetadataStorageSuite) TestMetadataStorageMetadata(c *gc.C) {
	stor := filestorage.NewMetadataStorage()
	id, err := stor.AddDoc(s.original)
	c.Assert(err, gc.IsNil)

	meta, err := stor.Metadata(id)
	c.Assert(err, gc.IsNil)
	c.Check(meta, gc.DeepEquals, s.original)
}

func (s *MetadataStorageSuite) TestMetadataStorageListDocs(c *gc.C) {
	stor := filestorage.NewMetadataStorage()
	id, err := stor.AddDoc(s.original)
	c.Assert(err, gc.IsNil)

	list, err := stor.ListDocs()
	c.Assert(err, gc.IsNil)
	c.Assert(list, gc.HasLen, 1)
	c.Assert(list[0], gc.NotNil)
	meta, ok := list[0].(filestorage.Metadata)
	c.Assert(ok, gc.Equals, true)
	c.Check(meta.ID(), gc.Equals, id)
}

func (s *MetadataStorageSuite) TestMetadataStorageListMetadata(c *gc.C) {
	stor := filestorage.NewMetadataStorage()
	id, err := stor.AddDoc(s.original)
	c.Assert(err, gc.IsNil)

	list, err := stor.ListMetadata()
	c.Assert(err, gc.IsNil)
	c.Assert(list, gc.HasLen, 1)
	c.Assert(list[0], gc.NotNil)
	c.Check(list[0].ID(), gc.Equals, id)
}

func (s *MetadataStorageSuite) TestMetadataStorageAddDoc(c *gc.C) {
	stor := filestorage.NewMetadataStorage()
	list, err := stor.ListMetadata()
	c.Assert(err, gc.IsNil)
	c.Assert(list, gc.HasLen, 0)

	id, err := stor.AddDoc(s.original)

	meta, err := stor.Metadata(id)
	c.Assert(err, gc.IsNil)
	c.Check(meta, gc.DeepEquals, s.original)
}

func (s *MetadataStorageSuite) TestMetadataStorageRemoveDoc(c *gc.C) {
	stor := filestorage.NewMetadataStorage()
	id, err := stor.AddDoc(s.original)
	c.Assert(err, gc.IsNil)

	err = stor.RemoveDoc(id)
	c.Assert(err, gc.IsNil)
	_, err = stor.Metadata(id)
	c.Assert(err, gc.NotNil)
}

func (s *MetadataStorageSuite) TestMetadataStorageNew(c *gc.C) {
	stor := filestorage.NewMetadataStorage()

	meta := stor.New()
	c.Assert(meta.ID(), gc.Equals, "")
}

func (s *MetadataStorageSuite) TestMetadataStorageSetStored(c *gc.C) {
	stor := filestorage.NewMetadataStorage()
	id, err := stor.AddDoc(s.original)
	c.Assert(err, gc.IsNil)
	meta, err := stor.Metadata(id)
	c.Assert(err, gc.IsNil)
	c.Check(meta.Stored(), gc.Equals, false)

	err = stor.SetStored(meta)
	c.Assert(err, gc.IsNil)
	meta, err = stor.Metadata(id)
	c.Assert(err, gc.IsNil)
	c.Check(meta.Stored(), gc.Equals, true)
}

//---------------------------
// raw file storage

var _ = gc.Suite(&RawFileSuite{})

type RawFileSuite struct {
	testing.IsolationSuite
}

func (s *RawFileSuite) TestRawFileStorageNewRawFileStorage(c *gc.C) {
	stor, err := filestorage.NewRawFileStorage(c.MkDir())
	c.Assert(err, gc.IsNil)

	c.Check(stor, gc.NotNil)
}

func (s *RawFileSuite) TestRawFileStorageFile(c *gc.C) {
	stor, err := filestorage.NewRawFileStorage(c.MkDir())
	c.Assert(err, gc.IsNil)
	data := bytes.NewBufferString("spam")
	err = stor.AddFile("eggs", data, 4)
	c.Assert(err, gc.IsNil)

	file, err := stor.File("eggs")
	c.Assert(err, gc.IsNil)
	content, err := ioutil.ReadAll(file)
	c.Assert(err, gc.IsNil)
	c.Check(string(content), gc.Equals, "spam")
}

func (s *RawFileSuite) TestRawFileStorageAddFile(c *gc.C) {
	stor, err := filestorage.NewRawFileStorage(c.MkDir())
	c.Assert(err, gc.IsNil)
	data := bytes.NewBufferString("spam")

	_, err = stor.File("eggs")
	c.Check(err, gc.NotNil)

	err = stor.AddFile("eggs", data, 4)
	c.Assert(err, gc.IsNil)
	file, err := stor.File("eggs")
	c.Assert(err, gc.IsNil)
	content, err := ioutil.ReadAll(file)
	c.Assert(err, gc.IsNil)
	c.Check(string(content), gc.Equals, "spam")
}

func (s *RawFileSuite) TestRawFileStorageRemoveFile(c *gc.C) {
	stor, err := filestorage.NewRawFileStorage(c.MkDir())
	c.Assert(err, gc.IsNil)
	data := bytes.NewBufferString("spam")
	err = stor.AddFile("eggs", data, 4)
	c.Assert(err, gc.IsNil)

	err = stor.RemoveFile("eggs")
	c.Check(err, gc.IsNil)
	_, err = stor.File("eggs")
	c.Check(err, gc.NotNil)
}
