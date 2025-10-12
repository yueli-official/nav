package main

import (
	"nav/conf"
	"nav/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// 获取所有导航数据
func getNavs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": resource,
	})
}

// ===== 分类管理 =====

// 创建分类
func createCategory(c *gin.Context) {
	var req models.Category
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 检查 CID 是否已存在
	for _, cat := range resource {
		if cat.Category.CID == req.CID {
			c.JSON(http.StatusConflict, gin.H{"code": 409, "message": "分类ID已存在"})
			return
		}
	}

	// 添加新分类
	newCategory := models.CategoryData{
		Category: req,
		Groups:   []models.GroupWithCollections{},
	}
	resource = append(resource, newCategory)

	if err := conf.SaveNavConfig(configPath, resource); err != nil {
		log.Error().Err(err).Msg("保存配置失败")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功", "data": newCategory})
}

// 更新分类
func updateCategory(c *gin.Context) {
	cid := c.Param("cid")
	var req models.Category
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 查找并更新分类
	found := false
	for i := range resource {
		if resource[i].Category.CID == cid {
			resource[i].Category = req
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "分类不存在"})
		return
	}

	if err := conf.SaveNavConfig(configPath, resource); err != nil {
		log.Error().Err(err).Msg("保存配置失败")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

// 删除分类
func deleteCategory(c *gin.Context) {
	cid := c.Param("cid")

	// 查找并删除分类
	found := false
	newResource := []models.CategoryData{}
	for _, cat := range resource {
		if cat.Category.CID == cid {
			found = true
			continue
		}
		newResource = append(newResource, cat)
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "分类不存在"})
		return
	}

	resource = newResource

	if err := conf.SaveNavConfig(configPath, resource); err != nil {
		log.Error().Err(err).Msg("保存配置失败")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

// ===== 分组管理 =====

// 创建分组
func createGroup(c *gin.Context) {
	cid := c.Param("cid")
	var req models.Group
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 查找分类并添加分组
	found := false
	for i := range resource {
		if resource[i].Category.CID == cid {
			// 检查分组ID是否已存在
			for _, g := range resource[i].Groups {
				if g.Group.GID == req.GID {
					c.JSON(http.StatusConflict, gin.H{"code": 409, "message": "分组ID已存在"})
					return
				}
			}

			newGroup := models.GroupWithCollections{
				Group:       req,
				Collections: []models.Collection{},
			}
			resource[i].Groups = append(resource[i].Groups, newGroup)
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "分类不存在"})
		return
	}

	if err := conf.SaveNavConfig(configPath, resource); err != nil {
		log.Error().Err(err).Msg("保存配置失败")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功"})
}

// 更新分组
func updateGroup(c *gin.Context) {
	cid := c.Param("cid")
	gid := c.Param("gid")
	var req models.Group
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 查找并更新分组
	found := false
	for i := range resource {
		if resource[i].Category.CID == cid {
			for j := range resource[i].Groups {
				if resource[i].Groups[j].Group.GID == gid {
					resource[i].Groups[j].Group = req
					found = true
					break
				}
			}
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "分类或分组不存在"})
		return
	}

	if err := conf.SaveNavConfig(configPath, resource); err != nil {
		log.Error().Err(err).Msg("保存配置失败")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

// 删除分组
func deleteGroup(c *gin.Context) {
	cid := c.Param("cid")
	gid := c.Param("gid")

	// 查找并删除分组
	found := false
	for i := range resource {
		if resource[i].Category.CID == cid {
			newGroups := []models.GroupWithCollections{}
			for _, group := range resource[i].Groups {
				if group.Group.GID == gid {
					found = true
					continue
				}
				newGroups = append(newGroups, group)
			}
			resource[i].Groups = newGroups
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "分类或分组不存在"})
		return
	}

	if err := conf.SaveNavConfig(configPath, resource); err != nil {
		log.Error().Err(err).Msg("保存配置失败")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

// ===== 收藏项管理 =====

// 创建收藏项
func createCollection(c *gin.Context) {
	cid := c.Param("cid")
	gid := c.Param("gid")
	var req models.Collection
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 查找分组并添加收藏项
	found := false
	for i := range resource {
		if resource[i].Category.CID == cid {
			for j := range resource[i].Groups {
				if resource[i].Groups[j].Group.GID == gid {
					resource[i].Groups[j].Collections = append(resource[i].Groups[j].Collections, req)
					found = true
					break
				}
			}
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "分类或分组不存在"})
		return
	}

	if err := conf.SaveNavConfig(configPath, resource); err != nil {
		log.Error().Err(err).Msg("保存配置失败")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功"})
}

// 更新收藏项（根据link匹配）
func updateCollection(c *gin.Context) {
	cid := c.Param("cid")
	gid := c.Param("gid")

	type UpdateCollectionReq struct {
		OldLink string            `json:"old_link"`
		Data    models.Collection `json:"data"`
	}

	var req UpdateCollectionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 查找并更新收藏项
	found := false
	for i := range resource {
		if resource[i].Category.CID == cid {
			for j := range resource[i].Groups {
				if resource[i].Groups[j].Group.GID == gid {
					for k := range resource[i].Groups[j].Collections {
						if resource[i].Groups[j].Collections[k].Link == req.OldLink {
							resource[i].Groups[j].Collections[k] = req.Data
							found = true
							break
						}
					}
					break
				}
			}
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "收藏项不存在"})
		return
	}

	if err := conf.SaveNavConfig(configPath, resource); err != nil {
		log.Error().Err(err).Msg("保存配置失败")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

// 删除收藏项（通过query参数传递link）
func deleteCollection(c *gin.Context) {
	cid := c.Param("cid")
	gid := c.Param("gid")
	link := c.Query("link")

	if link == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少link参数"})
		return
	}

	// 查找并删除收藏项
	found := false
	for i := range resource {
		if resource[i].Category.CID == cid {
			for j := range resource[i].Groups {
				if resource[i].Groups[j].Group.GID == gid {
					newCollections := []models.Collection{}
					for _, col := range resource[i].Groups[j].Collections {
						if col.Link == link {
							found = true
							continue
						}
						newCollections = append(newCollections, col)
					}
					resource[i].Groups[j].Collections = newCollections
					break
				}
			}
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "收藏项不存在"})
		return
	}

	if err := conf.SaveNavConfig(configPath, resource); err != nil {
		log.Error().Err(err).Msg("保存配置失败")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}
