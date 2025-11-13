import { Button, Form, Input, message } from 'antd';
import './index.css';
import { login } from '../../api';
import { useNavigate } from 'react-router-dom';

interface LoginUser {
  username: string;
  password: string;
}

const layout1 = {
  labelCol: { span: 4 },
  wrapperCol: { span: 20 }
}

const layout2 = {
  labelCol: { span: 0 },
  wrapperCol: { span: 24 }
}

export function Login() {
  const navigate = useNavigate();
  const onFinish = async (values: LoginUser) => {
    try {
      const res = await login(values.username, values.password);
      message.success('登录成功');

      localStorage.setItem('token', res.data.token.access_token);
      localStorage.setItem('refresh_token', res.data.token.refresh_token);
      localStorage.setItem('userInfo', JSON.stringify(res.data.user));
      setTimeout(() => {
        navigate('/');
      }, 1000);
    } catch (e: unknown) {
      if (e instanceof Error) {
        console.error(e.message);
      } else {
        console.error('Unknown error:', e);
      }
    }
  };
  return <div id="login-container">
    <h1>聊天室</h1>
    <Form
      {...layout1}
      onFinish={onFinish}
      colon={false}
      autoComplete="off"
    >
      <Form.Item
        label="用户名"
        name="username"
        rules={[{ required: true, message: '请输入用户名!' }]}
      >
        <Input />
      </Form.Item>

      <Form.Item
        label="密码"
        name="password"
        rules={[{ required: true, message: '请输入密码!' }]}
      >
        <Input.Password />
      </Form.Item>

      <Form.Item
        {...layout2}
      >
        <div className='links'>
          <a href='/register'>创建账号</a>
          <a href='/update_password'>忘记密码</a>
        </div>
      </Form.Item>

      <Form.Item
        {...layout2}
      >
        <Button className='btn' type="primary" htmlType="submit">
          登录
        </Button>
      </Form.Item>
    </Form>
  </div>
}
