package project_service_v1

import (
	"context"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"strconv"
	"test.com/project-common/encrypts"
	"test.com/project-common/errs"
	"test.com/project-common/tms"
	"test.com/project-grpc/project"
	"test.com/project-grpc/user/login"
	"test.com/project-project/internal/dao"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/database"
	"test.com/project-project/internal/database/tran"
	"test.com/project-project/internal/domain"
	"test.com/project-project/internal/repo"
	"test.com/project-project/internal/rpc"
	"test.com/project-project/pkg/model"
	"time"
)

type ProjectService struct {
	project.UnimplementedProjectServiceServer
	cache                  repo.Cache
	transaction            tran.Transaction
	menuRepo               repo.MenuRepo
	projectRepo            repo.ProjectRepo
	projectTemplateRepo    repo.ProjectTemplateRepo
	taskStagesTemplateRepo repo.TaskStagesTemplateRepo
	taskStagesRepo         repo.TaskStagesRepo
	projectLogRepo         repo.ProjectLogRepo
	taskRepo               repo.TaskRepo
	nodeDomain             *domain.ProjectNodeDomain
	taskDomain             *domain.TaskDomain
}

func New() *ProjectService {
	return &ProjectService{
		cache:                  dao.Rc,
		transaction:            dao.NewTransaction(),
		menuRepo:               dao.NewMenuDao(),
		projectRepo:            dao.NewProjectDao(),
		projectTemplateRepo:    dao.NewProjectTemplateDao(),
		taskStagesTemplateRepo: dao.NewTaskStagesTemplateDao(),
		taskStagesRepo:         dao.NewTaskStagesDao(),
		projectLogRepo:         dao.NewProjectLogDao(),
		taskRepo:               dao.NewTaskDao(),
		nodeDomain:             domain.NewProjectNodeDomain(),
		taskDomain:             domain.NewTaskDomain(),
	}
}

func (p *ProjectService) Index(context.Context, *project.IndexMessage) (*project.IndexResponse, error) {
	pms, err := p.menuRepo.FindMenus(context.Background())
	if err != nil {
		zap.L().Error("Index db FindMenus error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	childs := data.CovertChild(pms)
	var mms []*project.MenuMessage
	copier.Copy(&mms, childs)
	return &project.IndexResponse{Menus: mms}, nil
}

func (p *ProjectService) FindProjectByMemId(ctx context.Context, msg *project.ProjectRpcMessage) (*project.MyProjectResponse, error) {
	memberId := msg.MemberId
	page := msg.Page
	pageSize := msg.PageSize
	var pms []*data.ProjectAndMember
	var total int64
	var err error
	if msg.SelectBy == "" || msg.SelectBy == "my" {
		pms, total, err = p.projectRepo.FindProjectByMemId(ctx, memberId, "and deleted=0 ", page, pageSize)
	}
	if msg.SelectBy == "archive" {
		pms, total, err = p.projectRepo.FindProjectByMemId(ctx, memberId, "and archive=1 ", page, pageSize)
	}
	if msg.SelectBy == "deleted" {
		pms, total, err = p.projectRepo.FindProjectByMemId(ctx, memberId, "and deleted=1 ", page, pageSize)
	}
	if msg.SelectBy == "collect" {
		pms, total, err = p.projectRepo.FindCollectProjectByMemId(ctx, memberId, page, pageSize)
		for _, v := range pms {
			v.Collected = model.Collected
		}
	} else {
		collectPms, _, err := p.projectRepo.FindCollectProjectByMemId(ctx, memberId, page, pageSize)
		if err != nil {
			zap.L().Error("project FindProjectByMemId::FindCollectProjectByMemId error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
		var cMap = make(map[int64]*data.ProjectAndMember)
		for _, v := range collectPms {
			cMap[v.Id] = v
		}
		for _, v := range pms {
			if cMap[v.ProjectCode] != nil {
				v.Collected = model.Collected
			}
		}
	}
	if err != nil {
		zap.L().Error("project FindProjectByMemId error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if pms == nil {
		return &project.MyProjectResponse{Pm: []*project.ProjectMessage{}, Total: total}, nil
	}

	var pmm []*project.ProjectMessage
	copier.Copy(&pmm, pms)
	for _, v := range pmm {
		v.Code, _ = encrypts.EncryptInt64(v.ProjectCode, model.AESKey)
		pam := data.ToMap(pms)[v.Id]
		v.AccessControlType = pam.GetAccessControlType()
		v.OrganizationCode, _ = encrypts.EncryptInt64(pam.OrganizationCode, model.AESKey)
		v.JoinTime = tms.FormatByMill(pam.JoinTime)
		v.OwnerName = msg.MemberName
		v.Order = int32(pam.Sort)
		v.CreateTime = tms.FormatByMill(pam.CreateTime)
	}
	return &project.MyProjectResponse{Pm: pmm, Total: total}, nil
}

func (ps *ProjectService) FindProjectTemplate(ctx context.Context, msg *project.ProjectRpcMessage) (*project.ProjectTemplateResponse, error) {
	//1.根据viewType去查询项目模板表 得到list
	organizationCodeStr, _ := encrypts.Decrypt(msg.OrganizationCode, model.AESKey)
	organizationCode, _ := strconv.ParseInt(organizationCodeStr, 10, 64)
	page := msg.Page
	pageSize := msg.PageSize
	var pts []data.ProjectTemplate
	var total int64
	var err error
	if msg.ViewType == -1 {
		pts, total, err = ps.projectTemplateRepo.FindProjectTemplateAll(ctx, organizationCode, page, pageSize)
	}
	if msg.ViewType == 0 {
		pts, total, err = ps.projectTemplateRepo.FindProjectTemplateCustom(ctx, msg.MemberId, organizationCode, page, pageSize)
	}
	if msg.ViewType == 1 {
		pts, total, err = ps.projectTemplateRepo.FindProjectTemplateSystem(ctx, page, pageSize)
	}
	if err != nil {
		zap.L().Error("project FindProjectTemplate FindProjectTemplateSystem error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	//2.模型转换，拿到模板id列表 去 任务步骤模板表 去进行查询
	tsts, err := ps.taskStagesTemplateRepo.FindInProTemIds(ctx, data.ToProjectTemplateIds(pts))
	if err != nil {
		zap.L().Error("project FindProjectTemplate FindInProTemIds error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	var ptas []*data.ProjectTemplateAll
	for _, v := range pts {
		//写代码 该谁做的事情一定要交出去
		ptas = append(ptas, v.Convert(data.CovertProjectMap(tsts)[v.Id]))
	}
	//3.组装数据
	var pmMsgs []*project.ProjectTemplateMessage
	copier.Copy(&pmMsgs, ptas)
	return &project.ProjectTemplateResponse{Ptm: pmMsgs, Total: total}, nil
}

func (ps *ProjectService) SaveProject(ctxs context.Context, msg *project.ProjectRpcMessage) (*project.SaveProjectMessage, error) {
	organizationCodeStr, _ := encrypts.Decrypt(msg.OrganizationCode, model.AESKey)
	organizationCode, _ := strconv.ParseInt(organizationCodeStr, 10, 64)
	templateCodeStr, _ := encrypts.Decrypt(msg.TemplateCode, model.AESKey)
	templateCode, _ := strconv.ParseInt(templateCodeStr, 10, 64)
	//获取模板信息
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	stageTemplateList, err := ps.taskStagesTemplateRepo.FindByProjectTemplateId(ctx, int(templateCode))
	if err != nil {
		zap.L().Error("project SaveProject taskStagesTemplateRepo.FindByProjectTemplateId error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	//1. 保存项目表
	pr := &data.Project{
		Name:              msg.Name,
		Description:       msg.Description,
		TemplateCode:      int(templateCode),
		CreateTime:        time.Now().UnixMilli(),
		Cover:             "https://img2.baidu.com/it/u=792555388,2449797505&fm=253&fmt=auto&app=138&f=JPEG?w=667&h=500",
		Deleted:           model.NoDeleted,
		Archive:           model.NoArchive,
		OrganizationCode:  organizationCode,
		AccessControlType: model.Open,
		TaskBoardTheme:    model.Simple,
	}
	err = ps.transaction.Action(func(conn database.DbConn) error {
		err := ps.projectRepo.SaveProject(conn, ctx, pr)
		if err != nil {
			zap.L().Error("project SaveProject SaveProject error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		pm := &data.ProjectMember{
			ProjectCode: pr.Id,
			MemberCode:  msg.MemberId,
			JoinTime:    time.Now().UnixMilli(),
			IsOwner:     msg.MemberId,
			Authorize:   "",
		}
		//2. 保存项目和成员的关联表
		err = ps.projectRepo.SaveProjectMember(conn, ctx, pm)
		if err != nil {
			zap.L().Error("project SaveProject SaveProjectMember error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		//3. 生成任务的步骤
		for index, v := range stageTemplateList {
			taskStage := &data.TaskStages{
				ProjectCode: pr.Id,
				Name:        v.Name,
				Sort:        index + 1,
				Description: "",
				CreateTime:  time.Now().UnixMilli(),
				Deleted:     model.NoDeleted,
			}
			err := ps.taskStagesRepo.SaveTaskStages(ctx, conn, taskStage)
			if err != nil {
				zap.L().Error("project SaveProject taskStagesRepo.SaveTaskStages error", zap.Error(err))
				return errs.GrpcError(model.DBError)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	code, _ := encrypts.EncryptInt64(pr.Id, model.AESKey)
	rsp := &project.SaveProjectMessage{
		Id:               pr.Id,
		Code:             code,
		OrganizationCode: organizationCodeStr,
		Name:             pr.Name,
		Cover:            pr.Cover,
		CreateTime:       tms.FormatByMill(pr.CreateTime),
		TaskBoardTheme:   pr.TaskBoardTheme,
	}
	return rsp, nil
}

// 1. 查项目表
// 2. 项目和成员的关联表 查到项目的拥有者 去member表查名字
// 3. 查收藏表 判断收藏状态
func (ps *ProjectService) FindProjectDetail(ctx context.Context, msg *project.ProjectRpcMessage) (*project.ProjectDetailMessage, error) {
	projectCodeStr, _ := encrypts.Decrypt(msg.ProjectCode, model.AESKey)
	projectCode, _ := strconv.ParseInt(projectCodeStr, 10, 64)
	memberId := msg.MemberId
	zap.S().Info("project_code = ", projectCode)
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	projectAndMember, err := ps.projectRepo.FindProjectByPIdAndMemId(c, projectCode, memberId)
	if err != nil {
		zap.L().Error("project FindProjectDetail FindProjectByPIdAndMemId error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if projectAndMember == nil {
		return nil, errs.GrpcError(model.ParamsError)
	}
	ownerId := projectAndMember.IsOwner
	member, err := rpc.LoginServiceClient.FindMemInfoById(c, &login.UserMessage{MemId: ownerId})
	if err != nil {
		zap.L().Error("project rpc FindProjectDetail FindMemInfoById error", zap.Error(err))
		return nil, err
	}
	//去user模块去找了
	//TODO 优化 收藏的时候 可以放入redis
	isCollect, err := ps.projectRepo.FindCollectByPidAndMemId(c, projectCode, memberId)
	if err != nil {
		zap.L().Error("project FindProjectDetail FindCollectByPidAndMemId error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if isCollect {
		projectAndMember.Collected = model.Collected
	}
	var detailMsg = &project.ProjectDetailMessage{}
	copier.Copy(detailMsg, projectAndMember)
	detailMsg.OwnerAvatar = member.Avatar
	detailMsg.OwnerName = member.Name
	detailMsg.Code, _ = encrypts.EncryptInt64(projectAndMember.Id, model.AESKey)
	detailMsg.AccessControlType = projectAndMember.GetAccessControlType()
	detailMsg.OrganizationCode, _ = encrypts.EncryptInt64(projectAndMember.OrganizationCode, model.AESKey)
	detailMsg.Order = int32(projectAndMember.Sort)
	detailMsg.CreateTime = tms.FormatByMill(projectAndMember.CreateTime)
	return detailMsg, nil
}
func (ps *ProjectService) UpdateDeletedProject(ctx context.Context, msg *project.ProjectRpcMessage) (*project.DeletedProjectResponse, error) {
	projectCodeStr, _ := encrypts.Decrypt(msg.ProjectCode, model.AESKey)
	projectCode, _ := strconv.ParseInt(projectCodeStr, 10, 64)
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := ps.projectRepo.UpdateDeletedProject(c, projectCode, msg.Deleted)
	if err != nil {
		zap.L().Error("project RecycleProject DeleteProject error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &project.DeletedProjectResponse{}, nil
}
func (ps *ProjectService) UpdateProject(ctx context.Context, msg *project.UpdateProjectMessage) (*project.UpdateProjectResponse, error) {
	projectCodeStr, _ := encrypts.Decrypt(msg.ProjectCode, model.AESKey)
	projectCode, _ := strconv.ParseInt(projectCodeStr, 10, 64)
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	proj := &data.Project{
		Id:                 projectCode,
		Name:               msg.Name,
		Description:        msg.Description,
		Cover:              msg.Cover,
		TaskBoardTheme:     msg.TaskBoardTheme,
		Prefix:             msg.Prefix,
		Private:            int(msg.Private),
		OpenPrefix:         int(msg.OpenPrefix),
		OpenBeginTime:      int(msg.OpenBeginTime),
		OpenTaskPrivate:    int(msg.OpenTaskPrivate),
		Schedule:           msg.Schedule,
		AutoUpdateSchedule: int(msg.AutoUpdateSchedule),
	}
	err := ps.projectRepo.UpdateProject(c, proj)
	if err != nil {
		zap.L().Error("project UpdateProject::UpdateProject error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &project.UpdateProjectResponse{}, nil
}

func (ps *ProjectService) GetLogBySelfProject(ctx context.Context, msg *project.ProjectRpcMessage) (*project.ProjectLogResponse, error) {
	//根据用户id查询当前的用户的日志表

	projectLogs, total, err := ps.projectLogRepo.FindLogByMemberCode(context.Background(), msg.MemberId, msg.Page, msg.PageSize)
	if err != nil {
		zap.L().Error("project ProjectService::GetLogBySelfProject projectLogRepo.FindLogByMemberCode error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	//查询项目信息
	pIdList := make([]int64, len(projectLogs))
	mIdList := make([]int64, len(projectLogs))
	taskIdList := make([]int64, len(projectLogs))
	for _, v := range projectLogs {
		pIdList = append(pIdList, v.ProjectCode)
		mIdList = append(mIdList, v.MemberCode)
		taskIdList = append(taskIdList, v.SourceCode)
	}
	projects, err := ps.projectRepo.FindProjectByIds(context.Background(), pIdList)
	if err != nil {
		zap.L().Error("project ProjectService::GetLogBySelfProject projectLogRepo.FindProjectByIds error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	pMap := make(map[int64]*data.Project)
	for _, v := range projects {
		pMap[v.Id] = v
	}
	messageList, _ := rpc.LoginServiceClient.FindMemInfoByIds(context.Background(), &login.UserMessage{MIds: mIdList})
	mMap := make(map[int64]*login.MemberMessage)
	for _, v := range messageList.List {
		mMap[v.Id] = v
	}
	tasks, err := ps.taskRepo.FindTaskByIds(context.Background(), taskIdList)
	if err != nil {
		zap.L().Error("project ProjectService::GetLogBySelfProject projectLogRepo.FindTaskByIds error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	tMap := make(map[int64]*data.Task)
	for _, v := range tasks {
		tMap[v.Id] = v
	}
	var list []*data.IndexProjectLogDisplay
	for _, v := range projectLogs {
		display := v.ToIndexDisplay()
		display.ProjectName = pMap[v.ProjectCode].Name
		display.MemberAvatar = mMap[v.MemberCode].Avatar
		display.MemberName = mMap[v.MemberCode].Name
		display.TaskName = tMap[v.SourceCode].Name
		list = append(list, display)
	}
	var msgList []*project.ProjectLogMessage
	copier.Copy(&msgList, list)
	return &project.ProjectLogResponse{List: msgList, Total: total}, nil
}

func (ps *ProjectService) FindProjectByMemberId(ctx context.Context, msg *project.ProjectRpcMessage) (*project.FindProjectByMemberIdResponse, error) {
	isProjectCode := false
	var projectId int64
	if msg.ProjectCode != "" {
		projectId = encrypts.DecryptNoErr(msg.ProjectCode)
		isProjectCode = true
	}
	isTaskCode := false
	var taskId int64
	if msg.TaskCode != "" {
		taskId = encrypts.DecryptNoErr(msg.TaskCode)
		isTaskCode = true
	}
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if !isProjectCode && isTaskCode {
		projectCode, ok, bError := ps.taskDomain.FindProjectIdByTaskId(taskId)
		if bError != nil {
			return nil, bError
		}
		if !ok {
			return &project.FindProjectByMemberIdResponse{
				Project:  nil,
				IsOwner:  false,
				IsMember: false,
			}, nil
		}
		projectId = projectCode
		isProjectCode = true
	}
	if isProjectCode {
		//根据projectid和memberid查询
		pm, err := ps.projectRepo.FindProjectByPIdAndMemId(c, projectId, msg.MemberId)
		if err != nil {
			return nil, model.DBError
		}
		if pm == nil {
			return &project.FindProjectByMemberIdResponse{
				Project:  nil,
				IsOwner:  false,
				IsMember: false,
			}, nil
		}
		projectMessage := &project.ProjectMessage{}
		copier.Copy(projectMessage, pm)
		isOwner := false
		if pm.IsOwner == 1 {
			isOwner = true
		}
		return &project.FindProjectByMemberIdResponse{
			Project:  projectMessage,
			IsOwner:  isOwner,
			IsMember: true,
		}, nil
	}
	return &project.FindProjectByMemberIdResponse{}, nil
}
