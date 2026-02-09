function get_auth_token() {
    // 从会话存储中获取 token
    token = window.sessionStorage.getItem("auth_token");

    if (token == "" || token == null) {
        // 如果 token 为空，尝试从 Cookie 获取刷新令牌
        refresh_token = getCookie("refresh_token");
        if (refresh_token != null) {
            $.ajax({
                type: "POST",
                url: "/token",
                data: { "refresh_token": refresh_token },
                async: false, // 保证在返回 token 前代码阻塞等待
                success: function (result) {
                    token = result;
                    // 将新获取的 token 存入 sessionStorage
                    window.sessionStorage.setItem("auth_token", result);
                }
            }).fail(function (result, result1, result2) {
                // 处理获取失败逻辑
            });
        }
    }
    return token;
}

// 通用的获取 Cookie 函数
function getCookie(name) {
    var arr, reg = new RegExp("(^| )" + name + "=([^;]*)(;|$)");
    if (arr = document.cookie.match(reg))
        return arr[2];
    else
        return null;
}