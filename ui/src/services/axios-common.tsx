import axios from "axios";

const http = (dbkey: String) => {
  return axios.create({
    baseURL: "",
    headers: {
      "Content-type": "application/json",
      "Db-Key": String(dbkey),
    },
  });
};

export default http;
