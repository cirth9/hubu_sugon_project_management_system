package code_gen

import "testing"

func TestGenStruct(t *testing.T) {
	//GenStruct("ms_project", "Project")
	GenProtoMessage("ms_project", "ProjectMessage")
}
