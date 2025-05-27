import axios, { AxiosError } from "axios";

import { OpenAPI } from "../api-gen";
import { OpenAPI as AnchorOpenAPI } from "../api-anchor";
import toast from "react-hot-toast";
import { DefaultService as AnchorService } from "../api-anchor";

let baseUrl = "http://localhost:30080/api/v1";
if (typeof window !== "undefined") {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  baseUrl = (window as any).APP_ENDPOINT || "http://localhost:8020/api/v1";
}

AnchorOpenAPI.BASE = baseUrl;
OpenAPI.BASE = baseUrl;

const tokenKey = "token";
const refreshTokenKey = "refresh_token";

const initService = () => {
  const userToken = localStorage.getItem(tokenKey) || "";
  OpenAPI.TOKEN = userToken;
  AnchorOpenAPI.TOKEN = userToken;

  axios.interceptors.response.use(
    function (response) {
      return response;
    },

    async function (error: AxiosError) {
      console.log(error);
       
      const { status } = error;
      if (error.request?.responseURL?.includes("/auth/sign-in")) {
        return Promise.reject(error);
      }

      if (status === 401) {
        const refreshToken = localStorage.getItem(refreshTokenKey);
        if (refreshToken) {
          try {
            const response = await AnchorService.refreshToken({ refreshToken });
            localStorage.setItem(tokenKey, response.accessToken);
            localStorage.setItem(refreshTokenKey, response.refreshToken);
            OpenAPI.TOKEN = response.accessToken;
            AnchorOpenAPI.TOKEN = response.accessToken;

            // Retry the original request with the new token
            const config = error.config;
            if (config) {
              config.headers.Authorization = `Bearer ${response.accessToken}`;
              return axios(config);
            }
            return Promise.reject(error);
          } catch (refreshError) {
            console.error("Refresh token failed", refreshError);
            // If refresh token fails, clear storage and redirect to login
            localStorage.removeItem(tokenKey);
            localStorage.removeItem(refreshTokenKey);
            window.location.href = "/login";
            return Promise.reject(error);
          }
        }

        localStorage.removeItem(tokenKey);
        localStorage.removeItem(refreshTokenKey);
        window.location.href = "/login";
        return Promise.reject(error);
      } else if (status === 500) {
        return Promise.reject({
          message: "An error occurred. Please try again later.",
        });
      } else {
        if (error.status && error.status >= 500) {
          toast.error(error.response?.data as string);
        }
        if (error.status === 403) {
          toast.error("You are not authorized to perform this action");
        }
        if (error.status === 401) {
          toast.error("Cluster is running, please stop it before deleting");
        }
        return Promise.reject(error);
      }
    },
  );
};

export { initService };
