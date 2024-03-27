package project

import (
	"github.com/gin-gonic/gin"
	"log"
	"test.com/project-api/api/midd"
	"test.com/project-api/router"
)

type RouterProject struct {
}

func init() {
	log.Println("init project router")
	ru := &RouterProject{}
	router.Register(ru)
}

func (*RouterProject) Route(r *gin.Engine) {
	//初始化grpc的客户端连接
	InitRpcProjectClient()
	h := New()
	group := r.Group("/project")
	group.Use(midd.TokenVerify())
	group.Use(Auth())
	group.Use(ProjectAuth())
	group.POST("/index", h.index)
	group.POST("/project/selfList", h.myProjectList)
	group.POST("/project", h.myProjectList)
	group.POST("/project_template", h.projectTemplate)
	group.POST("/project_template/uploadCover", h.projectTemplateUploadCover) //todo unfinished ,do not know the api document
	group.POST("/project/save", h.projectSave)
	group.POST("/project/read", h.readProject)
	group.POST("/project/recycle", h.recycleProject)
	group.POST("/project/recovery", h.recoveryProject)
	group.POST("/project_collect/collect", h.collectProject)
	group.POST("/project/edit", h.editProject)
	group.POST("/project/getLogBySelfProject", h.getLogBySelfProject)
	group.POST("/project/analysis", h.getProjectAnalysis)
	group.POST("/node", h.nodeList)
	group.POST("/notify", h.notifyMessage)                                    //todo unfinished ,do not know the api document
	group.POST("/organization", h.getOrganizationInfo)                        //todo unfinished ,do not know the api document
	group.POST("/department_member/searchInviteMember", h.searchInviteMember) //todo unfinished ,do not know the api document
	group.POST("/file", h.getFile)                                            //todo unfinished ,do not know the api document

	t := NewTask()
	group.POST("/task_stages", t.taskStages)
	group.POST("/project_member/index", t.memberProjectList)
	group.POST("/task_stages/tasks", t.taskList)
	group.POST("/task/save", t.saveTask)
	group.POST("/task/sort", t.taskSort)
	group.POST("/task/selfList", t.myTaskList)
	group.POST("/task/read", t.readTask)
	group.POST("/task_member", t.listTaskMember)
	group.POST("/task/taskLog", t.taskLog)
	group.POST("/task/_taskWorkTimeList", t.taskWorkTimeList)
	group.POST("/task/saveTaskWorkTime", t.saveTaskWorkTime)
	group.POST("/file/uploadFiles", t.uploadFiles)
	group.POST("/task/taskSources", t.taskSources)
	group.POST("/task/createComment", t.createComment)

	a := NewAccount()
	group.POST("/account", a.account)

	d := NewDepartment()
	group.POST("/department", d.department)
	group.POST("/department/save", d.save)
	group.POST("/department/read", d.read)

	auth := NewAuth()
	group.POST("/auth", auth.authList)
	group.POST("/auth/apply", auth.apply)

	menu := NewMenu()
	group.POST("/menu/menu", menu.menuList)

	event := NewEvent()
	group.POST("/events/myList", event.myEventProjectList) //todo unfinished ,do not know the api document
}
