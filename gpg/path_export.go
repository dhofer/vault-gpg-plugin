package gpg

import (
	"bytes"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

func pathExportKeys(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "export/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the key",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathExportKeyRead,
		},
		HelpSynopsis:    pathExportHelpSyn,
		HelpDescription: pathExportHelpDesc,
	}
}

func (b *backend) pathExportKeyRead(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)
	entry, err := b.key(req.Storage, name)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	if !entry.Exportable {
		return logical.ErrorResponse("key is not exportable"), nil
	}

	var buf bytes.Buffer
	w, err := armor.Encode(&buf, openpgp.PrivateKeyType, nil)
	if err != nil {
		return nil, err
	}
	w.Write(entry.SerializedKey)
	if w.Close() != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"name": name,
			"key":  buf.String(),
		},
	}, nil
}

const pathExportHelpSyn = "Export named GPG key"
const pathExportHelpDesc = "This path is used to export the keys that are configured as exportable."
