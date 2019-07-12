package authorization

import (
	"github.com/pkg/errors"

	"github.com/lalloni/fabrikit/chaincode/context"
)

func MSPID(id string) Check {
	return func(ctx *context.Context) error {
		mspid, err := ctx.ClientMSPID()
		if err != nil {
			return err
		}
		if mspid != id {
			return errors.Errorf("Client MSPID %q not allowed", mspid)
		}
		return nil
	}
}

func SubjectCommonName(cn string) Check {
	return func(ctx *context.Context) error {
		cert, err := ctx.ClientCertificate()
		if err != nil {
			return err
		}
		v := cert.Subject.CommonName
		if cn != v {
			return errors.Errorf("Client certificate common name %q not allowed", v)
		}
		return nil
	}
}

func SubjectSerialNumber(sn string) Check {
	return func(ctx *context.Context) error {
		cert, err := ctx.ClientCertificate()
		if err != nil {
			return err
		}
		v := cert.Subject.CommonName
		if sn != v {
			return errors.Errorf("Client certificate serial number %q not allowed", v)
		}
		return nil
	}
}
