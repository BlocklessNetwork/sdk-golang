package ipfs

import (
	"errors"
	"fmt"

	"github.com/txlabs/blockless-sdk-golang/jsonparser"
)

type Args struct {
	Name  string
	Value string
}

func NewArgs(name, value string) *Args {
	return &Args{
		Name:  name,
		Value: value,
	}
}

type IpfsOptions struct {
	Api  string
	Args []*Args
}

func (opts *IpfsOptions) PushArg(arg *Args) {
	opts.Args = append(opts.Args, arg)
}

func NewIpfsOptions(api string) *IpfsOptions {
	return &IpfsOptions{
		Api:  api,
		Args: []*Args{},
	}
}

func (opts *IpfsOptions) JsonString() string {
	var encoder = jsonparser.NewJSONEncoder()
	encoder.PushObject("")
	encoder.SetString("api", opts.Api)
	encoder.PushArray("args")
	for i := 0; i < len(opts.Args); i++ {
		var arg = opts.Args[i]
		encoder.PushObject("")
		encoder.SetString("name", arg.Name)
		encoder.SetString("value", arg.Value)
		encoder.PopObject()
	}
	encoder.PopArray()
	encoder.PopObject()
	return encoder.ToString()
}

func IpfsCreateDir(path string, parents bool) error {
	var opts = NewIpfsOptions("files/mkdir")
	opts.PushArg(NewArgs("arg", path))
	opts.PushArg(NewArgs("parents", fmt.Sprintf("%t", parents)))
	result, err := ipfsCommandResult(opts)
	if err != nil {
		return err
	}
	if result.statusCode != 200 {
		return errors.New(string(result.respBody))
	}
	return nil
}

func IpfsFileRemove(path string, recursive bool, force bool) error {
	var opts = NewIpfsOptions("files/rm")
	opts.PushArg(NewArgs("arg", path))
	opts.PushArg(NewArgs("recursive", fmt.Sprintf("%t", recursive)))
	opts.PushArg(NewArgs("force", fmt.Sprintf("%t", force)))
	result, err := ipfsCommandResult(opts)
	if err != nil {
		return err
	}
	if result.statusCode != 200 {
		return errors.New(string(result.respBody))
	}
	return nil
}

type File struct {
	Name  string
	Ftype int64
	Size  int64
	Hash  string
}

func (f *File) String() string {
	return fmt.Sprintf("Name:%s, Type:%d, Size:%d, Hash:%s", f.Name, f.Ftype, f.Size, f.Hash)
}

func NewFile() *File {
	return &File{
		Name:  "",
		Ftype: 0,
		Size:  0,
		Hash:  "",
	}
}

type FileStat struct {
	Hash           string
	Size           int64
	Blocks         int64
	Ftype          string
	CumulativeSize int64
}

func (f *FileStat) String() string {
	return fmt.Sprintf("Hash:%s, Size:%d, Blocks:%d, Type:%s, CumulativeSize:%d", f.Hash, f.Size, f.Blocks, f.Ftype, f.CumulativeSize)
}

func NewFileStat() *FileStat {
	return &FileStat{
		Hash:           "",
		Size:           0,
		Blocks:         0,
		Ftype:          "",
		CumulativeSize: 0,
	}
}

func IpfsFileCopy(source string, dest string, parents bool) error {
	var opts = NewIpfsOptions("files/cp")
	opts.PushArg(NewArgs("arg", source))
	opts.PushArg(NewArgs("arg", dest))
	opts.PushArg(NewArgs("parents", fmt.Sprintf("%t", parents)))
	result, err := ipfsCommandResult(opts)
	if err != nil {
		return err
	}
	if result.statusCode != 200 {
		return errors.New(string(result.respBody))
	}
	return nil
}

func IpfsFileStat(path string) (*FileStat, error) {
	var opts = NewIpfsOptions("files/stat")
	opts.PushArg(NewArgs("arg", path))
	result, err := ipfsCommandResult(opts)
	if err != nil {
		return nil, err
	}
	if result.statusCode != 200 {
		return nil, errors.New(string(result.respBody))
	}
	stat := NewFileStat()
	value := result.respBody
	if stat.Hash, err = jsonparser.GetString(value, "Hash"); err != nil {
		return nil, err
	}
	if stat.Ftype, err = jsonparser.GetString(value, "Type"); err != nil {
		return nil, err
	}
	if stat.Size, err = jsonparser.GetInt(value, "Size"); err != nil {
		return nil, err
	}
	if stat.Blocks, err = jsonparser.GetInt(value, "Blocks"); err != nil {
		return nil, err
	}
	if stat.CumulativeSize, err = jsonparser.GetInt(value, "CumulativeSize"); err != nil {
		return nil, err
	}
	return stat, nil
}

func IpfsFileList(path string) ([]*File, error) {
	var item []byte
	var jsonType jsonparser.ValueType
	var err error
	var opts = NewIpfsOptions("files/ls")
	if len(path) > 0 {
		opts.PushArg(NewArgs("arg", path))
	}
	result, err := ipfsCommandResult(opts)
	if err != nil {
		return nil, err
	}
	if result.statusCode != 200 {
		return nil, errors.New(string(result.respBody))
	}
	var resp = result.respBody
	if item, jsonType, _, err = jsonparser.Get(resp, "Entries"); err != nil {
		return nil, err
	}
	if jsonType == jsonparser.NotExist || jsonType != jsonparser.Array {
		return nil, errors.New("invalid response")
	}
	var rs = []*File{}
	jsonparser.ArrayEach(item, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if dataType == jsonparser.Object {
			file := NewFile()
			if file.Name, err = jsonparser.GetString(value, "Name"); err != nil {
				return
			}
			if file.Ftype, err = jsonparser.GetInt(value, "Type"); err != nil {
				return
			}
			if file.Size, err = jsonparser.GetInt(value, "Size"); err != nil {
				return
			}
			if file.Hash, err = jsonparser.GetString(value, "Hash"); err != nil {
				return
			}
			rs = append(rs, file)
		}
	})
	return rs, nil
}

type FileWriteOptions struct {
	File     string
	Offset   int64
	Create   bool
	Parents  bool
	Truncate bool
}

func NewFileWriteOptions(file string) *FileWriteOptions {
	return &FileWriteOptions{
		File:     file,
		Offset:   0,
		Create:   true,
		Parents:  false,
		Truncate: false,
	}
}

func IpfsFileWrite(wo *FileWriteOptions, buf []byte) (uint32, error) {
	var err error
	var writen uint32
	var crs *commandRs
	var opts = NewIpfsOptions("files/write")
	opts.PushArg(NewArgs("arg", wo.File))
	opts.PushArg(NewArgs("offset", fmt.Sprintf("%d", wo.Offset)))
	opts.PushArg(NewArgs("create", fmt.Sprintf("%t", wo.Create)))
	opts.PushArg(NewArgs("parents", fmt.Sprintf("%t", wo.Parents)))
	opts.PushArg(NewArgs("truncate", fmt.Sprintf("%t", wo.Truncate)))
	if crs, err = ipfsCommand(opts); err != nil {
		return 0, err
	}
	defer ipfs_close(crs.handle)
	if writen, err = writeBody(crs.handle, buf); err != nil {
		return 0, err
	}
	if _, err = readBodyAll(crs.handle); err != nil {
		return 0, err
	}
	return writen, nil
}

func IpfsFileRead(path string, offset uint64, buf []byte) (uint32, error) {
	var err error
	var readn uint32
	var crs *commandRs
	if cap(buf) == 0 {
		return 0, INVALID_PARAMETER
	}
	var opts = NewIpfsOptions("files/read")
	opts.PushArg(NewArgs("arg", path))
	opts.PushArg(NewArgs("offset", fmt.Sprintf("%d", offset)))
	opts.PushArg(NewArgs("count", fmt.Sprintf("%d", len(buf))))

	if crs, err = ipfsCommand(opts); err != nil {
		return 0, err
	}
	defer ipfs_close(crs.handle)
	if readn, err = readBody(crs.handle, buf); err != nil {
		return 0, err
	}
	return readn, nil
}
