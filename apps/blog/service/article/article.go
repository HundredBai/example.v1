package article

import (
	"github.com/shipengqi/example.v1/apps/blog/model"
)

type Article interface {
	GetTags(maps map[string]interface{}, page int) ([]model.Tag, error)
}

type article struct {

}
