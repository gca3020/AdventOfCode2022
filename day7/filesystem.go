package main

type item struct {
	name   string
	parent *Directory
}

type Directory struct {
	subdirs map[string]*Directory
	files   []*File
	item
}

type File struct {
	size int
	item
}

func NewDirectory(name string, parent *Directory) *Directory {
	return &Directory{make(map[string]*Directory, 0), make([]*File, 0), item{name: name, parent: parent}}
}

func (d *Directory) Size() int {
	size := 0
	for _, subdir := range d.subdirs {
		size += subdir.Size()
	}
	for _, file := range d.files {
		size += file.size
	}
	return size
}

func (d *Directory) AddSubdir(name string) {
	d.subdirs[name] = NewDirectory(name, d)
}

func (d *Directory) AddFile(name string, size int) {
	d.files = append(d.files, NewFile(name, size, d))
}

func NewFile(name string, size int, parent *Directory) *File {
	return &File{size, item{name, parent}}
}

func (f *File) Size() int {
	return f.size
}
