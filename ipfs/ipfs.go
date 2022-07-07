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

func (opts *IpfsOptions) ToJsonString() string {
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
