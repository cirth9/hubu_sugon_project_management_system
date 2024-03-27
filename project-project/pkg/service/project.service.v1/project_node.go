package project_service_v1

import (
	"context"
	"github.com/jinzhu/copier"
	"test.com/project-common/errs"
	"test.com/project-grpc/project"
)

func (ps *ProjectService) NodeList(ctx context.Context, msg *project.ProjectRpcMessage) (*project.ProjectNodeResponseMessage, error) {
	list, err := ps.nodeDomain.TreeList()
	if err != nil {
		return nil, errs.GrpcError(err)
	}
	var nodes []*project.ProjectNodeMessage
	copier.Copy(&nodes, list)
	return &project.ProjectNodeResponseMessage{Nodes: nodes}, nil
}
