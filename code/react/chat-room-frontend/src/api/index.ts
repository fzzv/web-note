import axios from "axios";
import type { RegisterUser } from "../pages/Register";
import type { UpdatePassword } from "../pages/UpdatePassword";
import type { UserInfo } from "../pages/UpdateInfo";
import { message } from "antd";

const axiosInstance = axios.create({
  baseURL: 'http://localhost:3000/',
  timeout: 3000
});

axiosInstance.interceptors.request.use(function (config) {
  const accessToken = localStorage.getItem('token');

  if (accessToken) {
    config.headers.authorization = 'Bearer ' + accessToken;
  }
  return config;
})

let isRefreshing = false;
let refreshSubscribers: ((token: string) => void)[] = [];

function onRefreshed(token: string) {
  refreshSubscribers.forEach((cb) => cb(token));
  refreshSubscribers = [];
}
axiosInstance.interceptors.response.use(
  (response) => {
    if (response.data.code === 200 || response.data.code === 201) {
      return response.data;
    }
  }, async (error) => {
    const originalRequest = error.config;

    // 如果没有响应，直接拒绝
    if (!error.response) {
      return Promise.reject(error);
    }

    if (error.response?.data.code === 500) {
      message.error(error.response?.data.message);
    }

    // 如果返回 401，说明 access_token 失效
    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        // 其他请求等待刷新完成
        return new Promise((resolve) => {
          refreshSubscribers.push((token) => {
            originalRequest.headers.Authorization = `Bearer ${token}`;
            resolve(axiosInstance(originalRequest));
          });
        });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        const refreshToken = localStorage.getItem('refresh_token');
        const { data } = await refreshTokenApi(refreshToken!);

        const newAccessToken = data.access_token;
        localStorage.setItem('token', newAccessToken);
        axiosInstance.defaults.headers.common.Authorization = `Bearer ${newAccessToken}`;

        onRefreshed(newAccessToken);
        return axiosInstance(originalRequest);
      } catch (refreshError) {
        console.error('Refresh token invalid', refreshError);
        // refresh_token 也失效 → 登出
        localStorage.removeItem('token');
        localStorage.removeItem('refresh_token');
        window.location.href = '/login';
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }

    return Promise.reject(error);
  }
)

export async function login(username: string, password: string) {
  return await axiosInstance.post('/user/login', {
    username, password
  });
}

export async function refreshTokenApi(refreshToken: string) {
  return await axiosInstance.post('/user/refresh-token', {
    refresh_token: refreshToken
  });
}

export async function registerCaptcha(email: string) {
  return await axiosInstance.get('/user/register-captcha', {
    params: {
      address: email
    }
  });
}

export async function register(registerUser: RegisterUser) {
  return await axiosInstance.post('/user/register', registerUser);
}

export async function updatePasswordCaptcha(email: string) {
  return await axiosInstance.get('/user/update_password/captcha', {
    params: {
      address: email
    }
  });
}

export async function updatePassword(data: UpdatePassword) {
  return await axiosInstance.post('/user/update_password', data);
}

export async function getUserInfo() {
  return await axiosInstance.get('/user/info');
}

export async function updateInfo(data: UserInfo) {
  return await axiosInstance.post('/user/update', data);
}

export async function updateUserInfoCaptcha() {
  return await axiosInstance.get('/user/update/captcha');
}
