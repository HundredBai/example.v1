package v1

import (
	"fmt"
	app2 "github.com/shipengqi/example.v1/apps/blog/pkg/app"
	e2 "github.com/shipengqi/example.v1/apps/blog/pkg/e"
	"github.com/shipengqi/example.v1/apps/blog/pkg/export"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

type AddTagRequest struct {
	Name      string `form:"name" valid:"Required;MaxSize(100)"`
	CreatedBy string `form:"created_by" valid:"Required;MaxSize(100)"`
	// State     int    `form:"state" valid:"Range(0,1)"`
}

type AddTagResponse struct {
	Name string `json:"name"`
}

type EditTagRequest struct {
	Name       string `form:"name" valid:"Required;MaxSize(100)"`
	ModifiedBy string `form:"modified_by" valid:"Required;MaxSize(100)"`
	ID         int    `form:"id" valid:"Required;Min(1)"`
}

type EditTagResponse struct {
	Name string `json:"name"`
}

type ExportTagRequest struct {
	Name string `form:"name" valid:"MaxSize(100)"`
}

type ExportTagResponse struct {
	ExportUrl     string `json:"export_url"`
	ExportSaveUrl string `json:"export_save_url"`
}

// @Summary Get multiple article tags
// @Produce application/json
// @Param name query string false "Name"
// @Success 200 {object} app.Response
// @Failure 200 {object} app.Response
// @Router /api/v1/tags [get]
func GetTags(c *gin.Context) {
	name := c.Query("name")

	maps := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}
	page := com.StrTo(c.Query("page")).MustInt()
	data, err := svc.TagSvc.GetTags(maps, page)
	if err != nil {
		app2.SendResponse(c, err, data)
		return
	}
	app2.SendResponse(c, e2.OK, data)
}

// @Summary Add article tag
// @Produce application/json
// @Param name body string true "Name"
// @Param created_by body int false "CreatedBy"
// @Success 200 {object} app.Response
// @Failure 200 {object} app.Response
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context) {
	var form AddTagRequest

	err := app2.BindAndValid(c, &form)
	if err != nil {
		app2.SendResponse(c, err, nil)
		return
	}
	err = svc.TagSvc.AddTag(form.Name, form.CreatedBy)
	if err != nil {
		app2.SendResponse(c, err, nil)
		return
	}
	app2.SendResponse(c, e2.OK, nil)
}

// @Summary Update article tag
// @Produce application/json
// @Param id path int true "ID"
// @Param name body string true "Name"
// @Param state body int false "State"
// @Param modified_by body string true "ModifiedBy"
// @Success 200 {object} app.Response
// @Failure 200 {object} app.Response
// @Router /api/v1/tags/{id} [put]
func EditTag(c *gin.Context) {
	form := EditTagRequest{ID: com.StrTo(c.Param("id")).MustInt()}
	err := app2.BindAndValid(c, &form)
	if err != nil {
		app2.SendResponse(c, err, nil)
		return
	}

	data, err := svc.TagSvc.EditTag(form.ID, form.Name, form.ModifiedBy)
	if err != nil {
		app2.SendResponse(c, err, data)
		return
	}

	app2.SendResponse(c, e2.OK, data)
}

// @Summary Delete article tag
// @Produce application/json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 200 {object} app.Response
// @Router /api/v1/tags/{id} [delete]
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID must greater than 0")

	if valid.HasErrors() {
		app2.MarkErrors(valid.Errors)
		app2.SendResponse(c, e2.ErrBadRequest, nil)
		return
	}

	err := svc.TagSvc.DeleteTag(id)
	if err != nil {
		app2.SendResponse(c, err, nil)
		return
	}
	app2.SendResponse(c, e2.OK, nil)
}

// @Summary Export tags
// @Produce  json
// @Param name body string false "Name"
// @Success 200 {object} app.Response
// @Failure 200 {object} app.Response
// @Router /api/v1/tags/export [post]
func ExportTag(c *gin.Context) {
	form := ExportTagRequest{}
	err := app2.BindAndValid(c, &form)
	if err != nil {
		app2.SendResponse(c, err, nil)
		return
	}
	filename, err := svc.TagSvc.Export(form.Name)
	if err != nil {
		app2.SendResponse(c, err, nil)
		return
	}

	app2.SendResponse(c, e2.OK, ExportTagResponse{
		ExportUrl:     export.GetExcelFullUrl(filename),
		ExportSaveUrl: fmt.Sprintf("%s/%s", export.GetExcelPath(), filename),
	})
}

// @Summary Import tags
// @Produce  json
// @Param file body file true "Excel File"
// @Success 200 {object} app.Response
// @Failure 200 {object} app.Response
// @Router /api/v1/tags/import [post]
func ImportTag(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		app2.SendResponse(c, e2.Wrap(e2.ErrMultiFormErr, err.Error()), nil)
		return
	}

	err = svc.TagSvc.Import(file)
	if err != nil {
		app2.SendResponse(c, err, nil)
		return
	}

	app2.SendResponse(c, e2.OK, nil)
}
