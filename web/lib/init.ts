import axios from "axios";

import { OpenAPI } from "../api-gen";
import { DefaultService } from "../api-gen";
import toast from "react-hot-toast";

const baseUrl = "http://localhost:8020/api/v1";

const tokenKey = "token";
const refreshTokenKey = "refresh_token";

const initService = () => {
  OpenAPI.BASE = baseUrl;
  const userToken = localStorage.getItem(tokenKey) || "";
  OpenAPI.TOKEN = userToken;

  axios.interceptors.response.use(
    function (response) {
      return response;
    },

    async function (error) {
      console.log(error);
      // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment, @typescript-eslint/no-unsafe-member-access
      const { status } = error;

      if (status === 401) {
        const refreshToken = localStorage.getItem(refreshTokenKey);
        if (refreshToken) {
          try {
            const response = await DefaultService.refreshToken({ refreshToken });
            localStorage.setItem(tokenKey, response.accessToken);
            localStorage.setItem(refreshTokenKey, response.refreshToken);
            OpenAPI.TOKEN = response.accessToken;
            
            // Retry the original request with the new token
            const config = error.config;
            config.headers.Authorization = `Bearer ${response.accessToken}`;
            return axios(config);
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
        if (error.status >= 500) {
          toast.error(error.response.data);
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
