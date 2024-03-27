package repo

import (
	"context"
	"test.com/project-user/internal/data/organization"
	"test.com/project-user/internal/database"
)

type OrganizationRepo interface {
	SaveOrganization(conn database.DbConn, ctx context.Context, org *organization.Organization) error
	FindOrganizationByMemId(ctx context.Context, memId int64) ([]*organization.Organization, error)
}
