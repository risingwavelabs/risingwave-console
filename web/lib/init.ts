import axios from "axios";

import { OpenAPI } from "../api-gen";
import toast from "react-hot-toast";

const baseUrl = "http://localhost:8000/api/v1";

const tokenKey = "token";
const userInfoKey = "user_info";

const initService = () => {
  OpenAPI.BASE = baseUrl;
  const userToken = localStorage.getItem(tokenKey) || "";
  OpenAPI.TOKEN = userToken;

  axios.interceptors.response.use(
    function (response) {
      return response;
    },

    async function (error) {
      // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment, @typescript-eslint/no-unsafe-member-access
      const { status } = error.response;
      console.log(error);

      if (status === 401) {
        localStorage.removeItem(tokenKey);
        localStorage.removeItem(userInfoKey);
        window.location.href = "/login";
        return Promise.reject(error);
      } else if (status === 500) {
        return Promise.reject({
          message: "An error occurred. Please try again later.",
        });
      } else {
        // eslint-disable-next-line @typescript-eslint/no-unsafe-argument, @typescript-eslint/no-unsafe-member-access
        toast.error(error.response.data);
        return Promise.reject(error);
      }
    },
  );
};

export { initService };
