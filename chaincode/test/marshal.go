package test

import "encoding/json"

// MustMarshal convierte a JSON el valor v suministrado y genera un panic si hay
// algún error en el proceso.
//
// ATENCION: SOLO PARA TESTING
func MustMarshal(v interface{}) string {
	bs, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(bs)
}
