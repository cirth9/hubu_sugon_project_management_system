package task_service_v1

import (
	"context"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"test.com/project-common/encrypts"
	"test.com/project-common/errs"
	"test.com/project-common/tms"
	"test.com/project-grpc/task"
	"test.com/project-grpc/user/login"
	"test.com/project-project/config"
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

type TaskService struct {
	task.UnimplementedTaskServiceServer
	cache                  repo.Cache
	transaction            tran.Transaction
	projectRepo            repo.ProjectRepo
	projectTemplateRepo    repo.ProjectTemplateRepo
	taskStagesTemplateRepo repo.TaskStagesTemplateRepo
	taskStagesRepo         repo.TaskStagesRepo
	taskRepo               repo.TaskRepo
	projectLogRepo         repo.ProjectLogRepo
	taskWorkTimeRepo       repo.TaskWorkTimeRepo
	fileRepo               repo.FileRepo
	sourceLinkRepo         repo.SourceLinkRepo
	taskWorkTimeDomain     *domain.TaskWorkTimeDomain
}

func New() *TaskService {
	return &TaskService{
		cache:                  dao.Rc,
		transaction:            dao.NewTransaction(),
		projectRepo:            dao.NewProjectDao(),
		projectTemplateRepo:    dao.NewProjectTemplateDao(),
		taskStagesTemplateRepo: dao.NewTaskStagesTemplateDao(),
		taskStagesRepo:         dao.NewTaskStagesDao(),
		taskRepo:               dao.NewTaskDao(),
		projectLogRepo:         dao.NewProjectLogDao(),
		taskWorkTimeRepo:       dao.NewTaskWorkTimeDao(),
		fileRepo:               dao.NewFileDao(),
		sourceLinkRepo:         dao.NewSourceLinkDao(),
		taskWorkTimeDomain:     domain.NewTaskWorkTimeDomain(),
	}
}

func (t *TaskService) TaskStages(co context.Context, msg *task.TaskReqMessage) (*task.TaskStagesResponse, error) {
	projectCode := encrypts.DecryptNoErr(msg.ProjectCode)
	page := msg.Page
	pageSize := msg.PageSize
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	stages, total, err := t.taskStagesRepo.FindStagesByProjectId(ctx, projectCode, page, pageSize)
	if err != nil {
		zap.L().Error("project SaveProject taskStagesRepo.FindStagesByProjectId error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}

	var tsMessages []*task.TaskStagesMessage
	copier.Copy(&tsMessages, stages)
	if tsMessages == nil {
		return &task.TaskStagesResponse{List: tsMessages, Total: 0}, nil
	}
	stagesMap := data.ToTaskStagesMap(stages)
	for _, v := range tsMessages {
		taskStages := stagesMap[int(v.Id)]
		v.Code = encrypts.EncryptNoErr(int64(v.Id))
		v.CreateTime = tms.FormatByMill(taskStages.CreateTime)
		v.ProjectCode = msg.ProjectCode
	}
	return &task.TaskStagesResponse{List: tsMessages, Total: total}, nil
}

func (t *TaskService) MemberProjectList(co context.Context, msg *task.TaskReqMessage) (*task.MemberProjectResponse, error) {
	//1. 去 project_member表 去查询 用户id列表
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	projectCode := encrypts.DecryptNoErr(msg.ProjectCode)
	projectMembers, total, err := t.projectRepo.FindProjectMemberByPid(ctx, projectCode)
	if err != nil {
		zap.L().Error("project MemberProjectList projectRepo.FindProjectMemberByPid error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	//2.拿上用户id列表 去请求用户信息
	if projectMembers == nil || len(projectMembers) <= 0 {
		return &task.MemberProjectResponse{List: nil, Total: 0}, nil
	}
	var mIds []int64
	pmMap := make(map[int64]*data.ProjectMember)
	for _, v := range projectMembers {
		mIds = append(mIds, v.MemberCode)
		pmMap[v.MemberCode] = v
	}
	//请求用户信息
	userMsg := &login.UserMessage{
		MIds: mIds,
	}
	memberMessageList, err := rpc.LoginServiceClient.FindMemInfoByIds(ctx, userMsg)
	if err != nil {
		zap.L().Error("project MemberProjectList LoginServiceClient.FindMemInfoByIds error", zap.Error(err))
		return nil, err
	}
	var list []*task.MemberProjectMessage
	for _, v := range memberMessageList.List {
		owner := pmMap[v.Id].IsOwner
		mpm := &task.MemberProjectMessage{
			MemberCode: v.Id,
			Name:       v.Name,
			Avatar:     v.Avatar,
			Email:      v.Email,
			Code:       v.Code,
		}
		if v.Id == owner {
			mpm.IsOwner = 1
		}
		list = append(list, mpm)
	}
	return &task.MemberProjectResponse{List: list, Total: total}, nil
}
func (t *TaskService) TaskList(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskListResponse, error) {
	stageCode := encrypts.DecryptNoErr(msg.StageCode)
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	taskList, err := t.taskRepo.FindTaskByStageCode(c, int(stageCode))
	if err != nil {
		zap.L().Error("project task TaskList FindTaskByStageCode error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	var taskDisplayList []*data.TaskDisplay
	var mIds []int64
	for _, v := range taskList {
		display := v.ToTaskDisplay()
		if v.Private == 1 {
			//代表隐私模式
			taskMember, err := t.taskRepo.FindTaskMemberByTaskId(ctx, v.Id, msg.MemberId)
			if err != nil {
				zap.L().Error("project task TaskList taskRepo.FindTaskMemberByTaskId error", zap.Error(err))
				return nil, errs.GrpcError(model.DBError)
			}
			if taskMember != nil {
				display.CanRead = model.CanRead
			} else {
				display.CanRead = model.NoCanRead
			}
		}
		taskDisplayList = append(taskDisplayList, display)
		mIds = append(mIds, v.AssignTo)
	}
	if mIds == nil || len(mIds) <= 0 {
		return &task.TaskListResponse{List: nil}, nil
	}
	// in ()
	messageList, err := rpc.LoginServiceClient.FindMemInfoByIds(ctx, &login.UserMessage{MIds: mIds})
	if err != nil {
		zap.L().Error("project task TaskList LoginServiceClient.FindMemInfoByIds error", zap.Error(err))
		return nil, err
	}
	memberMap := make(map[int64]*login.MemberMessage)
	for _, v := range messageList.List {
		memberMap[v.Id] = v
	}
	for _, v := range taskDisplayList {
		message := memberMap[encrypts.DecryptNoErr(v.AssignTo)]
		e := data.Executor{
			Name:   message.Name,
			Avatar: message.Avatar,
		}
		v.Executor = e
	}
	var taskMessageList []*task.TaskMessage
	copier.Copy(&taskMessageList, taskDisplayList)
	return &task.TaskListResponse{List: taskMessageList}, nil
}

func (t *TaskService) SaveTask(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskMessage, error) {
	//1. 检查业务逻辑
	if msg.Name == "" {
		return nil, errs.GrpcError(model.TaskNameNotNull)
	}
	stageCode := encrypts.DecryptNoErr(msg.StageCode)
	taskStages, err := t.taskStagesRepo.FindById(ctx, int(stageCode))
	if err != nil {
		zap.L().Error("project task SaveTask taskStagesRepo.FindById error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if taskStages == nil {
		return nil, errs.GrpcError(model.TaskStagesNotNull)
	}
	projectCode := encrypts.DecryptNoErr(msg.ProjectCode)
	project, err := t.projectRepo.FindProjectById(ctx, projectCode)
	if err != nil {
		zap.L().Error("project task SaveTask projectRepo.FindProjectById error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if project == nil || project.Deleted == model.Deleted {
		return nil, errs.GrpcError(model.ProjectAlreadyDeleted)
	}
	maxIdNum, err := t.taskRepo.FindTaskMaxIdNum(ctx, projectCode)
	if err != nil {
		zap.L().Error("project task SaveTask taskRepo.FindTaskMaxIdNum error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if maxIdNum == nil {
		a := 0
		maxIdNum = &a
	}
	maxSort, err := t.taskRepo.FindTaskSort(ctx, projectCode, stageCode)
	if err != nil {
		zap.L().Error("project task SaveTask taskRepo.FindTaskSort error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if maxSort == nil {
		a := 0
		maxSort = &a
	}
	assignTo := encrypts.DecryptNoErr(msg.AssignTo)
	ts := &data.Task{
		Name:        msg.Name,
		CreateTime:  time.Now().UnixMilli(),
		CreateBy:    msg.MemberId,
		AssignTo:    assignTo,
		ProjectCode: projectCode,
		StageCode:   int(stageCode),
		IdNum:       *maxIdNum + 1,
		Private:     project.OpenTaskPrivate,
		Sort:        *maxSort + 65536,
		BeginTime:   time.Now().UnixMilli(),
		EndTime:     time.Now().Add(2 * 24 * time.Hour).UnixMilli(),
	}
	err = t.transaction.Action(func(conn database.DbConn) error {
		err = t.taskRepo.SaveTask(ctx, conn, ts)
		if err != nil {
			zap.L().Error("project task SaveTask taskRepo.SaveTask error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}

		tm := &data.TaskMember{
			MemberCode: assignTo,
			TaskCode:   ts.Id,
			JoinTime:   time.Now().UnixMilli(),
			IsOwner:    model.Owner,
		}
		if assignTo == msg.MemberId {
			tm.IsExecutor = model.Executor
		}
		err = t.taskRepo.SaveTaskMember(ctx, conn, tm)
		if err != nil {
			zap.L().Error("project task SaveTask taskRepo.SaveTaskMember error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	display := ts.ToTaskDisplay()
	member, err := rpc.LoginServiceClient.FindMemInfoById(ctx, &login.UserMessage{MemId: assignTo})
	if err != nil {
		return nil, err
	}
	display.Executor = data.Executor{
		Name:   member.Name,
		Avatar: member.Avatar,
		Code:   member.Code,
	}
	//添加任务动态
	createProjectLog(t.projectLogRepo, ts.ProjectCode, ts.Id, ts.Name, ts.AssignTo, "create", "task")

	tm := &task.TaskMessage{}
	copier.Copy(tm, display)
	//发送kafka 缓存删除
	config.SendCache([]byte("task"))
	return tm, nil
}

func createProjectLog(
	logRepo repo.ProjectLogRepo,
	projectCode int64,
	taskCode int64,
	taskName string,
	toMemberCode int64,
	logType string,
	actionType string) {
	remark := ""
	if logType == "create" {
		remark = "创建了任务"
	}
	pl := &data.ProjectLog{
		MemberCode:  toMemberCode,
		SourceCode:  taskCode,
		Content:     taskName,
		Remark:      remark,
		ProjectCode: projectCode,
		CreateTime:  time.Now().UnixMilli(),
		Type:        logType,
		ActionType:  actionType,
		Icon:        "plus",
		IsComment:   0,
		IsRobot:     0,
	}
	logRepo.SaveProjectLog(pl)
}

func (t *TaskService) TaskSort(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskSortResponse, error) {
	//移动的任务id肯定有 preTaskCode
	preTaskCode := encrypts.DecryptNoErr(msg.PreTaskCode)
	toStageCode := encrypts.DecryptNoErr(msg.ToStageCode)
	if msg.PreTaskCode == msg.NextTaskCode {
		return &task.TaskSortResponse{}, nil
	}
	err := t.sortTask(preTaskCode, msg.NextTaskCode, toStageCode)
	if err != nil {
		return nil, err
	}
	return &task.TaskSortResponse{}, nil

}

func (t *TaskService) sortTask(preTaskCode int64, nextTaskCode string, toStageCode int64) error {
	//1. 从小到大排
	//2. 原有的顺序  比如 1 2 3 4 5 4排到2前面去 4的序号在1和2 之间 如果4是最后一个 保证 4比所有的序号都打 如果 排到第一位 直接置为0

	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	ts, err := t.taskRepo.FindTaskById(c, preTaskCode)
	if err != nil {
		zap.L().Error("project task TaskSort taskRepo.FindTaskById error", zap.Error(err))
		return errs.GrpcError(model.DBError)
	}
	err = t.transaction.Action(func(conn database.DbConn) error {
		//如果相等是不需要进行改变的
		ts.StageCode = int(toStageCode)
		if nextTaskCode != "" {
			//意味着要进行排序的替换
			nextTaskCode := encrypts.DecryptNoErr(nextTaskCode)
			next, err := t.taskRepo.FindTaskById(c, nextTaskCode)
			if err != nil {
				zap.L().Error("project task TaskSort taskRepo.FindTaskById error", zap.Error(err))
				return errs.GrpcError(model.DBError)
			}
			// next.Sort 要找到比它小的那个任务
			prepre, err := t.taskRepo.FindTaskByStageCodeLtSort(c, next.StageCode, next.Sort)
			if err != nil {
				zap.L().Error("project task TaskSort taskRepo.FindTaskByStageCodeLtSort error", zap.Error(err))
				return errs.GrpcError(model.DBError)
			}
			if prepre != nil {
				ts.Sort = (prepre.Sort + next.Sort) / 2
			}
			if prepre == nil {
				ts.Sort = 0
			}
			//sort := ts.Sort
			//ts.Sort = next.Sort
			//next.Sort = sort
			//err = t.taskRepo.UpdateTaskSort(c, conn, next)
			//if err != nil {
			//	zap.L().Error("project task TaskSort taskRepo.UpdateTaskSort error", zap.Error(err))
			//	return errs.GrpcError(model.DBError)
			//}
		} else {
			maxSort, err := t.taskRepo.FindTaskSort(c, ts.ProjectCode, int64(ts.StageCode))
			if err != nil {
				zap.L().Error("project task TaskSort taskRepo.FindTaskSort error", zap.Error(err))
				return errs.GrpcError(model.DBError)
			}
			if maxSort == nil {
				a := 0
				maxSort = &a
			}
			ts.Sort = *maxSort + 65536
		}
		if ts.Sort < 50 {
			//重置排序
			err = t.resetSort(toStageCode)
			if err != nil {
				zap.L().Error("project task TaskSort resetSort error", zap.Error(err))
				return errs.GrpcError(model.DBError)
			}
			return t.sortTask(preTaskCode, nextTaskCode, toStageCode)
		}
		err = t.taskRepo.UpdateTaskSort(c, conn, ts)
		if err != nil {
			zap.L().Error("project task TaskSort taskRepo.UpdateTaskSort error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		return nil
	})
	return err
}

func (t *TaskService) resetSort(stageCode int64) error {
	list, err := t.taskRepo.FindTaskByStageCode(context.Background(), int(stageCode))
	if err != nil {
		return err
	}
	return t.transaction.Action(func(conn database.DbConn) error {
		iSort := 65536
		for index, v := range list {
			v.Sort = (index + 1) * iSort
			return t.taskRepo.UpdateTaskSort(context.Background(), conn, v)
		}
		return nil
	})

}

func (t *TaskService) MyTaskList(ctx context.Context, msg *task.TaskReqMessage) (*task.MyTaskListResponse, error) {
	var tsList []*data.Task
	var err error
	var total int64
	if msg.TaskType == 1 {
		//我执行的
		tsList, total, err = t.taskRepo.FindTaskByAssignTo(ctx, msg.MemberId, int(msg.Type), msg.Page, msg.PageSize)
		if err != nil {
			zap.L().Error("project task MyTaskList taskRepo.FindTaskByAssignTo error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
	}
	if msg.TaskType == 2 {
		//我参与的
		tsList, total, err = t.taskRepo.FindTaskByMemberCode(ctx, msg.MemberId, int(msg.Type), msg.Page, msg.PageSize)
		if err != nil {
			zap.L().Error("project task MyTaskList taskRepo.FindTaskByMemberCode error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
	}
	if msg.TaskType == 3 {
		//我创建的
		tsList, total, err = t.taskRepo.FindTaskByCreateBy(ctx, msg.MemberId, int(msg.Type), msg.Page, msg.PageSize)
		if err != nil {
			zap.L().Error("project task MyTaskList taskRepo.FindTaskByCreateBy error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
	}
	if tsList == nil || len(tsList) <= 0 {
		return &task.MyTaskListResponse{List: nil, Total: 0}, nil
	}
	var pids []int64
	var mids []int64
	for _, v := range tsList {
		pids = append(pids, v.ProjectCode)
		mids = append(mids, v.AssignTo)
	}
	pListChan := make(chan []*data.Project)
	defer close(pListChan)
	mListChan := make(chan *login.MemberMessageList)
	defer close(mListChan)
	//1.
	go func() {
		pList, _ := t.projectRepo.FindProjectByIds(ctx, pids)
		pListChan <- pList
	}()
	go func() {
		mList, _ := rpc.LoginServiceClient.FindMemInfoByIds(ctx, &login.UserMessage{
			//2.  1,2这两个请求无关联性  go+channel
			MIds: mids,
		})
		mListChan <- mList
	}()
	pList := <-pListChan
	projectMap := data.ToProjectMap(pList)
	mList := <-mListChan
	mMap := make(map[int64]*login.MemberMessage)
	for _, v := range mList.List {
		mMap[v.Id] = v
	}
	var mtdList []*data.MyTaskDisplay
	for _, v := range tsList {
		memberMessage := mMap[v.AssignTo]
		name := memberMessage.Name
		avatar := memberMessage.Avatar
		mtd := v.ToMyTaskDisplay(projectMap[v.ProjectCode], name, avatar)
		mtdList = append(mtdList, mtd)
	}
	var myMsgs []*task.MyTaskMessage
	copier.Copy(&myMsgs, mtdList)
	return &task.MyTaskListResponse{List: myMsgs, Total: total}, nil
}

func (t *TaskService) ReadTask(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskMessage, error) {
	//根据taskCode查询任务详情 根据任务查询项目详情 根据任务查询任务步骤详情 查询任务的执行者的成员详情
	taskCode := encrypts.DecryptNoErr(msg.TaskCode)
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	taskInfo, err := t.taskRepo.FindTaskById(c, taskCode)
	if err != nil {
		zap.L().Error("project task ReadTask taskRepo FindTaskById error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if taskInfo == nil {
		return &task.TaskMessage{}, nil
	}
	display := taskInfo.ToTaskDisplay()
	if taskInfo.Private == 1 {
		//代表隐私模式
		taskMember, err := t.taskRepo.FindTaskMemberByTaskId(ctx, taskInfo.Id, msg.MemberId)
		if err != nil {
			zap.L().Error("project task TaskList taskRepo.FindTaskMemberByTaskId error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
		if taskMember != nil {
			display.CanRead = model.CanRead
		} else {
			display.CanRead = model.NoCanRead
		}
	}
	pj, err := t.projectRepo.FindProjectById(c, taskInfo.ProjectCode)
	display.ProjectName = pj.Name
	taskStages, err := t.taskStagesRepo.FindById(c, taskInfo.StageCode)
	display.StageName = taskStages.Name
	memberMessage, err := rpc.LoginServiceClient.FindMemInfoById(ctx, &login.UserMessage{MemId: taskInfo.AssignTo})
	if err != nil {
		zap.L().Error("project task TaskList LoginServiceClient.FindMemInfoById error", zap.Error(err))
		return nil, err
	}
	e := data.Executor{
		Name:   memberMessage.Name,
		Avatar: memberMessage.Avatar,
	}
	display.Executor = e
	var taskMessage = &task.TaskMessage{}
	copier.Copy(taskMessage, display)
	return taskMessage, nil
}
func (t *TaskService) ListTaskMember(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskMemberList, error) {
	//查询 task member表 根据memberCode去查询用户信息
	taskCode := encrypts.DecryptNoErr(msg.TaskCode)
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	taskMemberPage, total, err := t.taskRepo.FindTaskMemberPage(c, taskCode, msg.Page, msg.PageSize)
	if err != nil {
		zap.L().Error("project task TaskList taskRepo.FindTaskMemberPage error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	var mids []int64
	for _, v := range taskMemberPage {
		mids = append(mids, v.MemberCode)
	}
	messageList, err := rpc.LoginServiceClient.FindMemInfoByIds(ctx, &login.UserMessage{MIds: mids})
	mMap := make(map[int64]*login.MemberMessage, len(messageList.List))
	for _, v := range messageList.List {
		mMap[v.Id] = v
	}
	var taskMemeberMemssages []*task.TaskMemberMessage
	for _, v := range taskMemberPage {
		tm := &task.TaskMemberMessage{}
		tm.Code = encrypts.EncryptNoErr(v.MemberCode)
		tm.Id = v.Id
		message := mMap[v.MemberCode]
		tm.Name = message.Name
		tm.Avatar = message.Avatar
		tm.IsExecutor = int32(v.IsExecutor)
		tm.IsOwner = int32(v.IsOwner)
		taskMemeberMemssages = append(taskMemeberMemssages, tm)
	}
	return &task.TaskMemberList{List: taskMemeberMemssages, Total: total}, nil
}
func (t *TaskService) TaskLog(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskLogList, error) {
	taskCode := encrypts.DecryptNoErr(msg.TaskCode)
	all := msg.All
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var list []*data.ProjectLog
	var total int64
	var err error
	if all == 1 {
		//显示全部
		list, total, err = t.projectLogRepo.FindLogByTaskCode(c, taskCode, int(msg.Comment))
	}
	if all == 0 {
		//分页
		list, total, err = t.projectLogRepo.FindLogByTaskCodePage(c, taskCode, int(msg.Comment), int(msg.Page), int(msg.PageSize))
	}
	if err != nil {
		zap.L().Error("project task TaskLog projectLogRepo.FindLogByTaskCodePage error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if total == 0 {
		return &task.TaskLogList{}, nil
	}
	var displayList []*data.ProjectLogDisplay
	var mIdList []int64
	for _, v := range list {
		mIdList = append(mIdList, v.MemberCode)
	}
	messageList, err := rpc.LoginServiceClient.FindMemInfoByIds(c, &login.UserMessage{MIds: mIdList})
	mMap := make(map[int64]*login.MemberMessage)
	for _, v := range messageList.List {
		mMap[v.Id] = v
	}
	for _, v := range list {
		display := v.ToDisplay()
		message := mMap[v.MemberCode]
		m := data.Member{}
		m.Name = message.Name
		m.Id = message.Id
		m.Avatar = message.Avatar
		m.Code = message.Code
		display.Member = m
		displayList = append(displayList, display)
	}
	var l []*task.TaskLog
	copier.Copy(&l, displayList)
	return &task.TaskLogList{List: l, Total: total}, nil
}

func (t *TaskService) TaskWorkTimeList(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskWorkTimeResponse, error) {
	taskCode := encrypts.DecryptNoErr(msg.TaskCode)
	list, err := t.taskWorkTimeDomain.TaskWorkTimeList(taskCode)
	if err != nil {
		return nil, errs.GrpcError(err)
	}
	var l []*task.TaskWorkTime
	copier.Copy(&l, list)
	return &task.TaskWorkTimeResponse{List: l, Total: int64(len(l))}, nil
}

func (t *TaskService) SaveTaskWorkTime(ctx context.Context, msg *task.TaskReqMessage) (*task.SaveTaskWorkTimeResponse, error) {
	tmt := &data.TaskWorkTime{}
	tmt.BeginTime = msg.BeginTime
	tmt.Num = int(msg.Num)
	tmt.Content = msg.Content
	tmt.TaskCode = encrypts.DecryptNoErr(msg.TaskCode)
	tmt.MemberCode = msg.MemberId
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := t.taskWorkTimeRepo.Save(c, tmt)
	if err != nil {
		zap.L().Error("project task SaveTaskWorkTime taskWorkTimeRepo.Save error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &task.SaveTaskWorkTimeResponse{}, nil
}

func (t *TaskService) SaveTaskFile(ctx context.Context, msg *task.TaskFileReqMessage) (*task.TaskFileResponse, error) {
	taskCode := encrypts.DecryptNoErr(msg.TaskCode)
	//存file表
	f := &data.File{
		PathName:         msg.PathName,
		Title:            msg.FileName,
		Extension:        msg.Extension,
		Size:             int(msg.Size),
		ObjectType:       "",
		OrganizationCode: encrypts.DecryptNoErr(msg.OrganizationCode),
		TaskCode:         encrypts.DecryptNoErr(msg.TaskCode),
		ProjectCode:      encrypts.DecryptNoErr(msg.ProjectCode),
		CreateBy:         msg.MemberId,
		CreateTime:       time.Now().UnixMilli(),
		Downloads:        0,
		Extra:            "",
		Deleted:          model.NoDeleted,
		FileType:         msg.FileType,
		FileUrl:          msg.FileUrl,
		DeletedTime:      0,
	}
	err := t.fileRepo.Save(context.Background(), f)
	if err != nil {
		zap.L().Error("project task SaveTaskFile fileRepo.Save error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	//存入source_link
	sl := &data.SourceLink{
		SourceType:       "file",
		SourceCode:       f.Id,
		LinkType:         "task",
		LinkCode:         taskCode,
		OrganizationCode: encrypts.DecryptNoErr(msg.OrganizationCode),
		CreateBy:         msg.MemberId,
		CreateTime:       time.Now().UnixMilli(),
		Sort:             0,
	}
	err = t.sourceLinkRepo.Save(context.Background(), sl)
	if err != nil {
		zap.L().Error("project task SaveTaskFile sourceLinkRepo.Save error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &task.TaskFileResponse{}, nil
}

func (t *TaskService) TaskSources(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskSourceResponse, error) {
	taskCode := encrypts.DecryptNoErr(msg.TaskCode)
	sourceLinks, err := t.sourceLinkRepo.FindByTaskCode(context.Background(), taskCode)
	if err != nil {
		zap.L().Error("project task SaveTaskFile sourceLinkRepo.FindByTaskCode error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if len(sourceLinks) == 0 {
		return &task.TaskSourceResponse{}, nil
	}
	var fIdList []int64
	for _, v := range sourceLinks {
		fIdList = append(fIdList, v.SourceCode)
	}
	files, err := t.fileRepo.FindByIds(context.Background(), fIdList)
	if err != nil {
		zap.L().Error("project task SaveTaskFile fileRepo.FindByIds error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	fMap := make(map[int64]*data.File)
	for _, v := range files {
		fMap[v.Id] = v
	}
	var list []*data.SourceLinkDisplay
	for _, v := range sourceLinks {
		list = append(list, v.ToDisplay(fMap[v.SourceCode]))
	}
	var slMsg []*task.TaskSourceMessage
	copier.Copy(&slMsg, list)
	return &task.TaskSourceResponse{List: slMsg}, nil
}

func (t *TaskService) CreateComment(ctx context.Context, msg *task.TaskReqMessage) (*task.CreateCommentResponse, error) {
	taskCode := encrypts.DecryptNoErr(msg.TaskCode)
	taskById, err := t.taskRepo.FindTaskById(context.Background(), taskCode)
	if err != nil {
		zap.L().Error("project task CreateComment fileRepo.FindTaskById error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	pl := &data.ProjectLog{
		MemberCode:   msg.MemberId,
		Content:      msg.CommentContent,
		Remark:       msg.CommentContent,
		Type:         "createComment",
		CreateTime:   time.Now().UnixMilli(),
		SourceCode:   taskCode,
		ActionType:   "task",
		ToMemberCode: 0,
		IsComment:    model.Comment,
		ProjectCode:  taskById.ProjectCode,
		Icon:         "plus",
		IsRobot:      0,
	}
	t.projectLogRepo.SaveProjectLog(pl)
	return &task.CreateCommentResponse{}, nil
}
