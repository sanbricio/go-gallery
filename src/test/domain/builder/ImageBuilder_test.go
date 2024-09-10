package builder_test

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/domain/entities/builder"
	"testing"
)

func TestImageBuilderEmptyFields(t *testing.T) {
	//Caso 'idImage' vacio
	_, err := builder.NewImageBuilder().
		From("", "", "", "", "").
		SetId("").
		Build()
	assertBuilderException(t, err, "idImage")
	//Caso 'name' vacio
	_, err = builder.NewImageBuilder().
		From("", "extension", "contentFile", "owner", "size").
		Build()
	assertBuilderException(t, err, "name")
	//Caso 'extension' vacio
	_, err = builder.NewImageBuilder().
		From("name", "", "contentFile", "owner", "size").
		Build()
	assertBuilderException(t, err, "extension")
	//Caso 'contentFile' vacio
	_, err = builder.NewImageBuilder().
		From("name", "extension", "", "owner", "size").
		Build()
	assertBuilderException(t, err, "contentFile")
	//Caso 'owner' vacio
	_, err = builder.NewImageBuilder().
		From("name", "extension", "contentFile", "", "size").
		Build()
	assertBuilderException(t, err, "owner")
	//Caso 'size' vacio
	_, err = builder.NewImageBuilder().
		From("name", "extension", "contentFile", "owner", "").
		Build()
	assertBuilderException(t, err, "size")
}

func assertBuilderException(t *testing.T, err *exception.BuilderException, field string) {
	if err == nil {
		t.Errorf("TestImageBuilderEmptyFrom: se esperaba un error al intentar crear una Image con el campo '%v' vacio", field)
		t.FailNow()
	}

	if err.Field != field {
		t.Errorf("TestImageBuilderEmptyFrom: se esperaba un error del campo %v", field)
		t.FailNow()
	}
}
