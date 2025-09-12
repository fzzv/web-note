import { NextFunction, Request, Response } from "express";

/**
 * HTML 的 <form> 标签默认只支持 GET 和 POST
 * 但 RESTful API 常常需要 PUT、PATCH、DELETE 等方法
 * 为了绕过这个限制，前端可以在表单里加一个隐藏字段 _method，把要真正使用的 HTTP 方法放进去。
 * example:
 * <form action="/users/1" method="POST">
 *   <input type="hidden" name="_method" value="DELETE">
 *   <button type="submit">Delete User</button>
 * </form>
 */
function methodOverride(req: Request, res: Response, next: NextFunction) {
  if (req.body && typeof req.body === 'object' && '_method' in req.body) {
    req.method = req.body._method.toUpperCase();
    delete req.body._method;
  }
  next();
}

export default methodOverride;
