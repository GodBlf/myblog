---
name: myblog-anime-frontend
description: 生成或改造 myblog 的前端页面（views/*.html）为现有暗色二次元风格（玻璃拟态、霓虹渐变、粒子光效）。当需要新建页面、统一样式、补齐交互按钮/表单/列表/详情布局，或复用 /js/particles.js 与现有视觉语言时使用。
---

# 目标

- 保持与 `views/` 现有页面一致的视觉风格，不做风格漂移。
- 产出可直接落地的完整 HTML（`<style>` + 必要 `<script>`）。
- 优先复用仓库已有资源：`/js/particles.js`、`/js/my.js`、jQuery。

# 强制风格规范

- 定义暗色 + 霓虹色变量（建议保留以下命名）：

```css
:root {
  --bg-a: #101128;
  --bg-b: #26265c;
  --card: rgba(255, 255, 255, 0.12);
  --line: rgba(255, 255, 255, 0.30);
  --text: #f8f8ff;
  --sub: #c7ccff;
  --pink: #ff79c6;
  --blue: #7ab8ff;
  --danger: #ff9ebd;
}
```

- `body` 使用深色底 + 双径向渐变光斑：
  - `background-color: var(--bg-a)`
  - `background-image: radial-gradient(...pink...), radial-gradient(...blue...)`
  - `position: relative; z-index: 1; min-height: 100vh`
- 主内容容器使用玻璃拟态：
  - 半透明背景、`border`、`backdrop-filter: blur(10px)`、柔和阴影、圆角 22~24px。
- 交互按钮统一胶囊样式：
  - 常态：半透明背景 + 细边框。
  - 悬停：`transform: translateY(-2px)` + 粉蓝渐变高亮。
- 标题统一带轻微霓虹光晕（`text-shadow`）。

# 粒子与光影效果

- 页面尾部必须包含：

```html
<script src="/js/particles.js"></script>
```

- 不要自行再创建第二套全屏背景 Canvas；使用现有粒子脚本即可。
- 内容层保持在粒子层上方（保留 `body` 的 `z-index: 1` 语义）。

# 结构模板选择

- 登录/小表单页：居中 `.card`，宽度 `min(92vw, 440px)`。
- 列表页：`.container` + `.list` + `.item`，支持 hover 高亮。
- 详情页：`.panel` + 标题 + `hr` + 正文块（`white-space: pre-wrap`）。
- 顶部动作区：`.actions` + 多个 `.btn`，保持可换行。

# 交互与脚本约定

- 若页面涉及鉴权接口，请使用 `/js/my.js` 的 `get_auth_token()`。
- 需要发送鉴权头时，遵循现有模式：

```js
beforeSend: function (request) {
  var auth_token = get_auth_token();
  request.setRequestHeader("auth_token", auth_token);
}
```

- 项目已有页面普遍使用 jQuery，新增交互默认沿用 jQuery 风格。

# 输出步骤

1. 先判断页面类型（登录、列表、详情、混合编辑）。
2. 选择对应骨架（见 `assets/base_page.html` 进行改写）。
3. 套用统一色板、按钮、容器、光影与粒子脚本。
4. 再填入业务结构（Go 模板变量、按钮事件、AJAX 请求）。
5. 最后检查移动端显示（`meta viewport`、按钮换行、文本可读性）。

# 输出检查清单

- 是否包含暗色 + 粉蓝渐变光斑背景。
- 是否包含玻璃拟态容器。
- 是否包含统一 `.btn` 悬停动效。
- 是否在页面底部引用 `/js/particles.js`。
- 若调用受保护接口，是否通过 `get_auth_token()` 注入 `auth_token` 头。

