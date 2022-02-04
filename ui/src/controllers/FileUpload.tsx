import React, { FormEvent, useState } from "react";
import { useSnackbar, VariantType } from "notistack";
import http from "../services/axios-common";
import UploadForm from "../components/UploadForm";

const FileUpload = () => {
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();
  const [file, setFile] = useState<string | Blob | File>("");

  const [uploadProgress, updateUploadProgress] = useState(0);
  const [uploadStatus, setUploadStatus] = useState(false);
  const [uploading, setUploading] = useState(false);

  const handleFileUpload = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    setUploading(true);
    const formData = new FormData();
    formData.append("file", file);

    http
      .post("/api/v1/upload", formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
        onUploadProgress: (ev: ProgressEvent) => {
          const progress = (ev.loaded / ev.total) * 100;
          updateUploadProgress(Math.round(progress));
        },
      })
      .then((resp) => {
        // our mocked response will always return true
        // in practice, you would want to use the actual response object
        setUploadStatus(true);
        setUploading(false);

        enqueueSnackbar("File uploaded successfully", {
          key: "success",
          variant: "success",
          draggable: true,
          onClick: () => {
            closeSnackbar("success");
          },
        });
        console.log(resp);
      })
      .catch((err) => {
        enqueueSnackbar(err.response.data.message, {
          key: "error",
          variant: "error",
          draggable: true,
          onClick: () => {
            closeSnackbar("error");
          },
        });
      });
  };

  return (
    <UploadForm
      handleFileUpload={handleFileUpload}
      setFile={setFile}
      uploadProgress={uploadProgress}
      enqueueSnackbar={enqueueSnackbar}
      closeSnackbar={closeSnackbar}
      file={file}
      uploading={uploading}
    />
  );
};

export default FileUpload;
