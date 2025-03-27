

# SAT Oneprep内容渲染规则文档

## 1. 基本元素类型 (type)

每个元素是一个字典，必须包含 `type` 字段，支持的主要类型包括：

### 1.1 string
- 普通字符串，直接显示文本内容
- 示例:
```json
{
  "type": "string",
  "content": "这是一段文本"
}
```

### 1.2 formula
- LaTeX 格式的公式
- 需将 `content.content` 转换为 LaTeX 显示
- 示例:
```json
{
  "type": "formula",
  "content": {
    "type": "latex",
    "content": "F(x) = 3.50 - 0.32x + 0.02x^2"
  }
}
```

### 1.3 image
三种可能格式：
- **base64**: 
```json
{
  "type": "image",
  "content": {
    "type": "base64",
    "content": "data:image/png;base64,..."
  }
}
```
- **url**: 使用链接
```json
{
  "type": "image",
  "content": {
    "type": "url",
    "content": "https://example.com/image.png"
  }
}
```
- **svg**: 使用SVG内容

### 1.4 figure
- 图表格式，直接还原原始 HTML
- 示例:
```json
{
  "type": "figure",
  "content": "<div class='figure'>...</div>"
}
```

### 1.5 table
- 表格格式，直接还原原始 HTML
- 示例:
```json
{
  "type": "table",
  "content": "<table>...</table>"
}
```

### 1.6 line_break
- 换行标志，生成换行效果
- 示例:
```json
{
  "type": "line_break"
}
```

### 1.7 blank
- 用于填空题，显示为下划线
- 示例:
```json
{
  "type": "blank"
}
```

### 1.8 html
- HTML内容，直接渲染
- 示例:
```json
{
  "type": "html",
  "content": "<div>HTML内容</div>"
}
```

## 2. 元数据处理 (metadata)

### 2.1 content_type
区分素材和问题的类型标记：
- **stimulus**: 表示素材部分
- **stem**: 表示问题部分
- 示例:
```json
{
  "type": "string",
  "content": "这是题目描述",
  "metadata": {
    "content_type": "stem"
  }
}
```

### 2.2 style
样式标记，支持多种样式：

- **p**: 段落标记，表示该元素后需要添加换行，切记在元素前方不需要换行！
- **em/i**: 斜体样式
- **strong/b**: 加粗样式
- **sup**: 上标样式
- **html-\***: 原始HTML样式，直接应用HTML样式

示例:
```json
{
  "type": "string",
  "content": "这是强调文本",
  "metadata": {
    "style": ["strong", "p"]
  }
}
```

## 3. 特殊处理规则

### 3.1 根号转换
将带有 `old-root-radicand` 类的span转换为LaTeX格式的根号：

HTML格式:
```html
<span class="old-root-radicand">16</span>
```

转换为:
```json
{
  "type": "formula",
  "content": {
    "type": "latex",
    "content": "\\sqrt{16}"
  }
}
```

## 4. 内容处理逻辑

### 4.1 文本合并规则
- 只合并相同 `content_type` 的连续字符串
- 使用空格连接合并后的字符串
- 如果任一字符串带有 `p` 样式，则合并后的内容也带有 `p` 样式

### 4.2 段落处理规则
- `LINE_BREAK` 类型或带有 `p` 样式的元素开始新段落
- 段落通过 `ParagraphContainer` 类型进行分组
- 每个段落可包含多个内容项

### 4.3 显示规则
- 段落样式元素通过 `display: block` 显示，并有适当边距
- 非段落样式元素通过 `display: inline-flex` 显示
- 公式根据 `display` 属性决定使用 `BlockMath` 或 `InlineMath`


