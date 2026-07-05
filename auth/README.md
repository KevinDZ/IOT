* EMQX Dashboard 中配置 HTTP 鉴权

** Method: POST
** URL: 0.0.0.0:8888/api/v1/emqx/acl
** Body 示例
{
  "clientid": "${clientid}",
  "username": "${username}",
  "action": "${action}",
  "topic": "${topic}"
}