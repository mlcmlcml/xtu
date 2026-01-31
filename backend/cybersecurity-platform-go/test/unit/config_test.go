package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestHelpers(t *testing.T) {
	// 测试 GetProjectRoot
	root, err := helpers.GetProjectRoot()
	assert.NoError(t, err)
	assert.NotEmpty(t, root)
	assert.Contains(t, root, "cybersecurity-platform-go")
	t.Logf("项目根目录: %s", root)
	
	// 测试 GetTestArticleContent
	content := helpers.GetTestArticleContent(1)
	assert.Contains(t, content, "网络安全")
	assert.Contains(t, content, "测试文章")
	t.Logf("测试文章内容长度: %d", len(content))
}